package cli

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"git.pride.improwised.dev/Onboarding-2025/Krupanshi-Vaishnav/go-api/config"
	"git.pride.improwised.dev/Onboarding-2025/Krupanshi-Vaishnav/go-api/database"
	"github.com/doug-martin/goqu/v9"
	_ "github.com/lib/pq" // for postgres dialect
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// GetSeedCommandDef initialize migration command
func GetSeedCommandDef(cfg config.AppConfig, logger *zap.Logger) cobra.Command {
	seedCmd := cobra.Command{
		Use:   "seed",
		Short: "To run db seed",
		Long:  `This command is used to run seeding for database.`,
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPostgresSeed(cfg, logger)
		},
	}

	return seedCmd
}

func runPostgresSeed(cfg config.AppConfig, logger *zap.Logger) error {
	db, err := database.Connect(cfg.DB)

	if err != nil {
		logger.Error("Database connection error", zap.Error(err))
		return err
	}

	return SeedAllCSVs(cfg, db, logger)
}

func SeedAllCSVs(cfg config.AppConfig, db *goqu.Database, logger *zap.Logger) error {
	tx, err := db.Begin()
	if err != nil {
		logger.Error("Error starting transaction", zap.Error(err))
		return err
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				logger.Error("Error rolling back transaction", zap.Error(rollbackErr))
			} else {
				logger.Info("Transaction rolled back successfully")
			}
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				logger.Error("Error committing transaction", zap.Error(commitErr))
			}
		}
	}()

	seeds := []struct {
		name, path string
		fn         func(*goqu.TxDatabase, string, *zap.Logger) error
	}{
		{"movies", cfg.Movies, SeedMoviesMetadata},
		{"ratings", cfg.Ratings, SeedRatings},
		{"credits", cfg.Credits, SeedCredits},
	}

	for _, seed := range seeds {
		logger.Info("Seeding " + seed.name)
		if err = seed.fn(tx, seed.path, logger); err != nil {
			logger.Error("Error seeding table "+seed.name, zap.Error(err))
			return err
		}
	}

	logger.Info("All seeding completed successfully!")
	return nil
}

func readCSVFile(filePath string) (*csv.Reader, *os.File, []string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error opening file %s: %w", filePath, err)
	}
	reader := csv.NewReader(file)
	headers, err := reader.Read()
	if err != nil {
		file.Close()
		return nil, nil, nil, fmt.Errorf("error reading headers in %s: %w", filePath, err)
	}
	return reader, file, headers, nil
}

func parseValue(val string) interface{} {
	if val == "" {
		return nil
	}
	val = strings.TrimSpace(val)

	switch strings.ToLower(val) {
	case "inf", "+inf", "-inf", "infinity":
		return val
	}
	if i, err := strconv.Atoi(val); err == nil {
		return i
	}
	if f, err := strconv.ParseFloat(val, 64); err == nil {
		return f
	}
	if t, err := time.Parse("2006-01-02", val); err == nil {
		return t
	}
	if strings.HasPrefix(val, "{") || strings.HasPrefix(val, "[") {
		val = strings.ReplaceAll(val, "'", "\"")
		var v interface{}
		_ = json.Unmarshal([]byte(val), &v)
		b, _ := json.Marshal(v)
		return string(b)
	}
	return val
}

func insertRecords(tx *goqu.TxDatabase, table string, rows []interface{}, logger *zap.Logger) error {
	if len(rows) == 0 {
		logger.Info("No records to insert", zap.String("table", table))
		return nil
	}

	start := time.Now()
	_, err := tx.Insert(table).Rows(rows...).OnConflict(goqu.DoNothing()).Executor().Exec()
	logger.Info("Inserted", zap.String("table", table), zap.Int("rows", len(rows)), zap.Duration("duration", time.Since(start)))
	return err
}

func SeedMoviesMetadata(tx *goqu.TxDatabase, filePath string, logger *zap.Logger) error {
	reader, file, headers, err := readCSVFile(filePath)
	if err != nil {
		logger.Error("csv load error", zap.Error(err))
		return err
	}
	defer file.Close()

	var (
		movieInserts, movieGenres, movieLanguages []interface{}
		genreMap                                  = make(map[int]string)
		languageMap                               = make(map[string]string)
	)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil || len(record) != len(headers) {
			continue
		}

		row := make(map[string]interface{})
		var movieID interface{}
		var genreJSONRaw, spokenLangJSONRaw string

		for i, col := range headers {
			val := record[i]
			switch col {
			case "id":
				movieID = parseValue(val)
				row["id"] = movieID
			case "original_language":
				code := val
				if code == "" {
					code = "xx"
				}
				row["original_language"] = code
				if _, exists := languageMap[code]; !exists {
					languageMap[code] = "Unknown"
				}
			case "genres":
				genreJSONRaw = val
			case "spoken_languages":
				spokenLangJSONRaw = val
			default:
				row[col] = parseValue(val)
			}
		}

		if movieID == nil {
			logger.Warn("missing movie ID", zap.String("filePath", filePath))
			continue
		}

		movieInserts = append(movieInserts, row)
		extractGenres(genreJSONRaw, movieID, genreMap, &movieGenres)
		extractSpokenLanguages(spokenLangJSONRaw, movieID, languageMap, &movieLanguages)
	}

	var genreRecords, languageRecords []interface{}
	for id, name := range genreMap {
		genreRecords = append(genreRecords, map[string]interface{}{
			"id":   id,
			"name": name,
		})
	}
	for iso, name := range languageMap {
		languageRecords = append(languageRecords, map[string]interface{}{
			"iso_code": iso,
			"name":     name,
		})
	}

	if err := insertRecords(tx, "genres", genreRecords, logger); err != nil {
		return fmt.Errorf("inserting genres: %w", err)
	}

	if err := insertRecords(tx, "languages", languageRecords, logger); err != nil {
		return fmt.Errorf("inserting languages: %w", err)
	}

	if err := insertRecords(tx, "movies", movieInserts, logger); err != nil {
		return fmt.Errorf("inserting movies: %w", err)
	}
	if err := insertRecords(tx, "movie_genres", movieGenres, logger); err != nil {
		return fmt.Errorf("inserting movie_genres: %w", err)
	}
	if err := insertRecords(tx, "movie_languages", movieLanguages, logger); err != nil {
		return fmt.Errorf("inserting movie_languages: %w", err)
	}

	return nil
}

func extractGenres(raw string, movieID interface{}, genreMap map[int]string, out *[]interface{}) {
	jsonStr := cleanJSONValue(raw)
	var genres []map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &genres); err != nil {
		return
	}
	for _, g := range genres {
		id := parseValue(fmt.Sprintf("%.0f", g["id"])).(int)
		name := g["name"].(string)
		genreMap[id] = name
		*out = append(*out, map[string]interface{}{
			"movieid": movieID,
			"genreid": id,
		})
	}
}

func extractSpokenLanguages(raw string, movieID interface{}, languageMap map[string]string, out *[]interface{}) {
	jsonStr := cleanJSONValue(raw)
	var langs []map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &langs); err != nil {
		return
	}
	for _, l := range langs {
		iso := l["iso_639_1"].(string)
		name := l["name"].(string)
		languageMap[iso] = name
		*out = append(*out, map[string]interface{}{
			"movieid":       movieID,
			"language_code": iso,
		})
	}
}

func cleanJSONValue(raw string) string {
	raw = strings.ReplaceAll(raw, "None", "null")
	clean := strings.ReplaceAll(raw, "'", "\"")
	re := regexp.MustCompile(`\\x[0-9a-fA-F]{2}`)
	return re.ReplaceAllString(clean, "")
}

func movieExists(tx *goqu.TxDatabase, movieID int) (bool, error) {
	var count int
	// Query the movies table to check if the movieId exists
	ds := tx.From("movies").
		Where(goqu.C("id").Eq(movieID)).Select(goqu.COUNT("*")) // Assuming "id" is the primary key of the movies table

	_, err := ds.ScanVal(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func SeedRatings(tx *goqu.TxDatabase, filePath string, logger *zap.Logger) error {
	reader, file, headers, err := readCSVFile(filePath)
	if err != nil {
		logger.Error("csv load error", zap.Error(err))
		return err
	}
	defer file.Close()

	var rows []interface{}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil || len(record) != len(headers) {
			continue
		}

		movieID := parseValue(record[slices.Index(headers, "movieId")]).(int)
		if exists, err := movieExists(tx, movieID); err != nil || !exists {
			continue
		}

		row := make(map[string]interface{}, len(headers))
		for i, col := range headers {
			val := record[i]
			if col == "movieId" {
				row["movie_id"] = parseValue(val)
			} else if col == "userId" {
				row["user_id"] = parseValue(val)

			} else if col == "timestamp" {
				ts, _ := strconv.ParseInt(val, 10, 64)
				row[col] = time.Unix(ts, 0)
			} else {
				row[col] = parseValue(val)
			}
		}
		rows = append(rows, row)
	}

	if err := insertRecords(tx, "ratings", rows, logger); err != nil {
		return fmt.Errorf("inserting ratings: %w", err)
	}
	return nil
}

func SeedCredits(tx *goqu.TxDatabase, filePath string, logger *zap.Logger) error {
	reader, file, headers, err := readCSVFile(filePath)
	if err != nil {
		logger.Error("csv load error", zap.Error(err))
		return err
	}
	defer file.Close()

	creditsMap := make(map[string]goqu.Record)
	castMap := make(map[string]goqu.Record)
	crewMap := make(map[string]goqu.Record)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil || len(record) != len(headers) {
			logger.Warn("Skipping bad record", zap.Error(err))
			continue
		}

		movieID := parseValue(record[slices.Index(headers, "id")]).(int)
		if exists, err := movieExists(tx, movieID); err != nil || !exists {
			continue
		}

		var castRaw, crewRaw string
		row := make(map[string]interface{}, len(headers))

		for i, col := range headers {
			val := record[i]
			switch col {
			case "cast":
				castRaw = record[i]
			case "crew":
				crewRaw = record[i]
			default:
				row[col] = parseValue(val)
			}
		}

		extractCreditsRecords(movieID, castRaw, "cast", creditsMap, castMap)
		extractCreditsRecords(movieID, crewRaw, "crew", creditsMap, crewMap)
	}

	// Insert into tables
	if err := insertRecords(tx, "credits", valuesFromMap(creditsMap), logger); err != nil {
		return fmt.Errorf("inserting people: %w", err)
	}
	if err := insertRecords(tx, "movie_casts", valuesFromMap(castMap), logger); err != nil {
		return fmt.Errorf("inserting movie_casts: %w", err)
	}
	if err := insertRecords(tx, "movie_crew", valuesFromMap(crewMap), logger); err != nil {
		return fmt.Errorf("inserting movie_crew: %w", err)
	}

	return nil
}

func valuesFromMap(m map[string]goqu.Record) []interface{} {
	out := make([]interface{}, 0, len(m))
	for _, v := range m {
		out = append(out, v)
	}
	return out
}

func extractCreditsRecords(movieID interface{}, rawJSON string, role string, creditsMap map[string]goqu.Record, roleMap map[string]goqu.Record) {
	cleanedJSON := cleanJSONValue(rawJSON)
	var creditsList []map[string]interface{}

	if err := json.Unmarshal([]byte(cleanedJSON), &creditsList); err != nil {
		return
	}

	for _, person := range creditsList {
		rawID, hasID := person["id"].(float64)
		if !hasID {
			continue
		}
		personID := fmt.Sprintf("%.0f", rawID)
		creditID, _ := person["credit_id"].(string)

		// Add person to peopleMap if not already present
		if _, exists := creditsMap[personID]; !exists {
			creditsMap[personID] = goqu.Record{
				"id":           personID,
				"name":         person["name"],
				"gender":       int(person["gender"].(float64)),
				"profile_path": person["profile_path"],
			}
		}

		// Add role-specific entry to roleMap if not already present
		if _, exists := roleMap[creditID]; !exists {
			roleRecord := goqu.Record{
				"movie_id":  movieID,
				"person_id": personID,
				"credit_id": creditID,
			}

			if role == "cast" {
				roleRecord["cast_id"] = int(person["cast_id"].(float64))
				roleRecord["character"] = person["character"]
				roleRecord["cast_order"] = int(person["order"].(float64))
			} else if role == "crew" {
				roleRecord["department"] = person["department"]
				roleRecord["job"] = person["job"]
			}

			roleMap[creditID] = roleRecord
		}
	}
}
