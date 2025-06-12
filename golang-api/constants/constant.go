package constants

const (
	LoadMoviesError  = "Failed to load movies"
	LoadRatingsError = "Failed to load ratings of movies"
	LoadCreditsError = "Failed to load credits of movies"
)

const (
	AddMovieError     = "Failed to add movie"
	AddRatingError    = "Failed to add ratings"
	DeleteRatingError = "Failed to delete rating"
	DeleteMovieError  = "Failed to delete movie"
	UpdateMovieError  = "Failed to update movie"
	UpdateRatingError = "Failed to update rating"
	UpdateCrewError   = "Failed to update crew member details"
	UpdateCastError   = "Failed to update cast member details"
)

const (
	MovieId = "movieId"
	UserId  = "userId"
	CastId  = "castId"
	CrewId  = "crewId"
)

const (
	AddRatingSuccess    = "Ratings added successfully"
	AddMovieSuccess     = "Movie added successfully"
	DeleteRatingSuccess = "Ratings deleted successfully"
	DeleteMovieSuccess  = "Movie deleted successfully"
	UpdateMovieSuccess  = "Movie updated successfully"
	UpdateRatingSuccess = "Ratings updated successfully"
	UpdateCrewSuccess   = "Crew Member details updated successfully"
	UpdateCastSuccess   = "Cast Member details updated successfully"
)

const (
	InvalidPageOrLimitError = "Invalid page number or limit number"
	PaginationError         = "Pagination error"
	InvalidRequestBody      = "Failed to parse request body"
	ValidationFailed        = "Request body is not as required"
	MovieCheckError         = "Movie not found"
)
