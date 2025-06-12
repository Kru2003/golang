package utils

import "git.pride.improwised.dev/Onboarding-2025/Krupanshi-Vaishnav/go-api/models"

////////////////////
// --- MOVIES  ---//
////////////////////

// swagger:parameters GetMovieByID
type RequestGetMovieByID struct {
	// in: path
	// required: true
	MovieID int `json:"movieId"`
}

// swagger:response ResponseGetMovieByID
type ResponseGetMovieByID struct {
	// in: body
	Body struct {
		// enum: success
		Status string       `json:"status"`
		Data   models.Movie `json:"data"`
	} `json:"body"`
}

// swagger:parameters ListMovies
type RequestListMovies struct {
	// in: query
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`
	Name     string `json:"name"`
	Genre    string `json:"genre"`
	Language string `json:"language"`
}

// swagger:response ResponseListMovies
type ResponseListMovies struct {
	// in: body
	Body struct {
		// enum: success
		Status string         `json:"status"`
		Data   []models.Movie `json:"data"`
	} `json:"body"`
}

// swagger:parameters DeleteMovieById
type RequestDeleteMovieByID struct {
	// in: path
	// required: true
	MovieID int `json:"movieId"`
}

// swagger:parameters AddMovie
type RequestAddMovie struct {
	// in: body
	// required: true
	Body struct {
		models.MovieWithMetadata
	}
}

// swagger:parameters UpdateMovie
type RequestUpdateMovie struct {
	// in: path
	// required: true
	MovieID int `json:"movieId"`
	// in: body
	Body struct {
		models.MovieWithMetadata
	}
}

////////////////////
// --- RATINGS ---//
////////////////////

// swagger:parameters ListAllMovieRatings
type RequestListAllMovieRatings struct {
	// in: query
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

// swagger:response ResponseListAllMovieRatings
type ResponseListAllMovieRatings struct {
	// in: body
	Body struct {
		// enum: success
		Status string               `json:"status"`
		Data   []models.MovieRating `json:"data"`
	}
}

// swagger:parameters GetRatingsByMovieId
type RequestGetRatingsByMovieId struct {
	// in: path
	// required: true
	MovieID int `json:"movieId"`
}

// swagger:response ResponseGetRatingsByMovieId
type ResponseGetRatingsByMovieId struct {
	// in: body
	Body struct {
		// enum: success
		Status string             `json:"status"`
		Data   models.MovieRating `json:"data"`
	} `json:"body"`
}

// swagger:parameters AddRating
type RequestAddRating struct {
	// in: path
	// required: true
	UserID int `json:"userId"`
	// in: body
	// required: true
	Body struct {
		MovieID int     `json:"movieId"`
		Rating  float32 `json:"rating"`
	}
}

// swagger:parameters UpdateRating
type RequestUpdateRating struct {
	// in: path
	// required: true
	MovieID int `json:"movieId"`
	//in: path
	// required: true
	UserID int `json:"userId"`
	// in: body
	Body struct {
		Rating float32 `json:"rating" validate:"required,gte=0,lte=5"`
	}
}

// swagger:parameters DeleteRating
type RequestDeleteRating struct {
	// in: path
	// required: true
	MovieID int `json:"movieId"`
	// in: path
	// required: true
	UserID int `json:"userId"`
}

/////////////////////////
// --- MOVIE_CASTS ---//
///////////////////////

// swagger:parameters ListCastMembers
type RequestListCastMembers struct {
	// in: path
	// required: true
	MovieID int `json:"movieId"`
}

// swagger:response ResponseListCastMembers
type ResponseListCastMembers struct {
	// in: body
	Body struct {
		// enum: success
		Status string             `json:"status"`
		Data   []models.MovieCast `json:"data"`
	} `json:"body"`
}

// swagger:parameters ListMoviesByCastId
type RequestListMoviesByCastId struct {
	// in: path
	// required: true
	CastID int `json:"castId"`
}

// swagger:response ResponseListMoviesByCastId
type ResponseListMoviesByCastId struct {
	// in: body
	Body struct {
		// enum: success
		Status string                 `json:"status"`
		Data   models.ActorWithMovies `json:"data"`
	}
}

// swagger:parameters AddMovieCastMember
type RequestAddMovieCastMember struct {
	// in: path
	// required: true
	MovieID int `json:"movieId"`
	// in: path
	// required: true
	CreditID int `json:"creditId"`
	// in: body
	Body struct {
		Character string `json:"character"`
		Order     int    `json:"order"`
	}
}

////////////////////////
// --- MOVIE_CREW ---//
//////////////////////

// swagger:parameters ListCrewMembers
type RequestListCrewMembers struct {
	// in: path
	// required: true
	MovieID int `json:"movieId"`
}

// swagger:response ResponseListCrewMembers
type ResponseListCrewMembers struct {
	// in: body
	Body struct {
		// enum: success
		Status string             `json:"status"`
		Data   []models.MovieCrew `json:"data"`
	} `json:"body"`
}

// swagger:parameters AddMovieCrewMember
type RequestAddMovieCrewMember struct {
	// in: path
	// required: true
	MovieID int `json:"movieId"`
	// in: path
	// required: true
	CreditID int `json:"creditId"`
	// in: body
	Body struct {
		Department string `json:"department"`
		Job        string `json:"job"`
	}
}

////////////////////
// --- GENERIC ---//
////////////////////

// Response is okay
// swagger:response GenericResOk
type ResOK struct {
	// in:body
	Body struct {
		// enum:success
		Status string `json:"status"`
	}
}

// Fail due to user invalid input
// swagger:response GenericResFailBadRequest
type ResFailBadRequest struct {
	// in: body
	Body struct {
		// enum: fail
		Status string      `json:"status"`
		Data   interface{} `json:"data"`
	} `json:"body"`
}

// Fail due to resource not exists
// swagger:response GenericResFailNotFound
type ResFailNotFound struct {
	// in: body
	Body struct {
		// enum: fail
		Status string      `json:"status"`
		Data   interface{} `json:"data"`
	} `json:"body"`
}

// Unexpected error occurred
// swagger:response GenericResError
type ResError struct {
	// in: body
	Body struct {
		// enum: error
		Status  string      `json:"status"`
		Data    interface{} `json:"data"`
		Message string      `json:"message"`
	} `json:"body"`
}
