package constants

// params
const (
	ParamMid = "movieId"
	CastId   = "castId"
	UserId   = "userId"
	CreditId = "creditId"
)

// Success messages
const (
	DeleteMovieSuccess  = "movie deleted successfully"
	DeleteRatingSuccess = "ratings deleted successfully"
	AddMovieSuccess     = "movie added successfully"
	AddRatingSuccess    = "rating added successfully"
	AddMovieCrewSuccess = "movie crew added successfully"
	AddMovieCastSuccess = "movie cast added successfully"
	UpdateMovieSuccess  = "movie updated successfully"
	UpdateRatingSuccess = "ratings updated successfully"
)

// Fail messages
const (
	MovieNotExist      = "movie does not exists"
	CastsNotExist      = "casts does not exists for given movie"
	ActorNotExist      = "actor does not exists"
	RatingNotExist     = "rating does not exists for given movie and user"
	InvalidPageOrLimit = "invalid page or limit value"
	InvalidRequestBody = "invalid request values"
	ValidationFailed   = "invalid input"
)

// Error messages
const (
	ErrHealthCheckDb = "error while checking health of database"
	ErrGetMovie      = "error while get movie"
	ErrGetRatings    = "error while get ratings"
	ErrGetCasts      = "error while get movie casts"
	ErrGetCrew       = "error while get movie crew"
	ErrAddMovie      = "error while adding movie"
	ErrAddRating     = "error while adding movie ratings"
	ErrAddMovieCrew  = "error while adding movie crew member"
	ErrAddMovieCast  = "error while adding movie cast"
	UpdateMovieError = "error while updating movie"
	ErrUpdateRating  = "error while updating rating"
	ErrDeleteRating  = "error while delete rating"
	ErrDeleteMovie   = "error while deleting move"
)
