package structs

import "git.pride.improwised.dev/Onboarding-2025/Krupanshi-Vaishnav/go-api/models"

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
		Status string          `json:"status"`
		Data   []models.Movies `json:"data"`
	} `json:"body"`
}

// swagger:parameters GetMovieByID
type RequestGetMovieByID struct {
	// in: path
	// required: true
	MovieID string `json:"movieId"`
}

// swagger:response ResponseGetMovieByID
type ResponseGetMovieByID struct {
	// in: body
	Body struct {
		// enum: success
		Status string        `json:"status"`
		Data   models.Movies `json:"data"`
	} `json:"body"`
}

// swagger:parameters AddMovie
type RequestAddMovie struct {
	// in: body
	// required: true
	Body struct {
		models.Movies
	}
}

// swagger:parameters UpdateMovie
type RequestUpdateMovie struct {
	// in: path
	// required: true
	MovieID string `json:"movieId"`
	// in: body
	Body struct {
		models.Movies
	}
}

// swagger:parameters DeleteMovieById
type RequestDeleteMovieByID struct {
	// in: path
	// required: true
	MovieID string `json:"movieId"`
}

// swagger:parameters ListAllMovieRatings
type RequestListAllMovieRatings struct {
	// in: query
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

// swagger:response ResponseListMovies
type ResponseListAllMovieRatings struct {
	// in: body
	Body struct {
		// enum: success
		Status string           `json:"status"`
		Data   []models.Ratings `json:"data"`
	} `json:"body"`
}

// swagger:parameters GetRatingsByMovieId
type RequestGetRatingsByMovieId struct {
	// in: path
	// required: true
	MovieID string `json:"movieId"`
}

// swagger:response ResponseGetRatingsByMovieId
type ResponseGetRatingsByMovieId struct {
	// in: body
	Body struct {
		// enum: success
		Status string         `json:"status"`
		Data   models.Ratings `json:"data"`
	} `json:"body"`
}

// swagger:parameters AddRating
type RequestAddRAting struct {
	// in: body
	// required: true
	Body struct {
		models.Movies
	}
}

// swagger:parameters DeleteRating
type RequestDeleteRating struct {
	// in: path
	// required: true
	MovieID string `json:"movieId"`
	// in: path
	// required: true
	UserID string `json:"userId"`
}

// swagger:parameters UpdateRating
type RequestUpdateRating struct {
	// in: path
	// required: true
	MovieID string `json:"movieId"`
	//in: path
	// required: true
	UserID string `json:"userId"`
	// in: body
	Body struct {
		models.Ratings
	}
}

// swagger:parameters ListCrewMembers
type RequestListCrewMembers struct {
	// in: path
	// required: true
	MovieID string `json:"movieId"`
}

// swagger:response ResponseListCrewMembers
type ResponseListCrewMembers struct {
	// in: body
	Body struct {
		// enum: success
		Status string              `json:"status"`
		Data   []models.CrewMember `json:"data"`
	} `json:"body"`
}

// swagger:parameters UpdateCrewMember
type RequestUpdateCrewMember struct {
	// in: path
	// required: true
	MovieID string `json:"movieId"`
	// in: path
	// required: true
	CrewID string `json:"crewId"`
	// in: body
	Body struct {
		models.CrewMember
	}
}

// swagger:parameters ListCastMembers
type RequestListCastMembers struct {
	// in: path
	// required: true
	MovieID string `json:"movieId"`
}

// swagger:response ResponseListCastMembers
type ResponseListCastMembers struct {
	// in: body
	Body struct {
		// enum: success
		Status string              `json:"status"`
		Data   []models.CastMember `json:"data"`
	} `json:"body"`
}

// swagger:parameters ListMoviesByCastId
type RequestListMoviesByCastId struct {
	// in: path
	// required: true
	CastID string `json:"castId"`
}

// swagger:response ResponseListMoviesByCastId
type ResponseListMoviesByCastId struct {
	// in: body
	Body struct {
		// enum: success
		Status string   `json:"status"`
		Data   []string `json:"data"`
	} `json:"body"`
}

// swagger:parameters UpdateCastMember
type RequestUpdateCastMember struct {
	// in: path
	// required: true
	MovieID string `json:"movieId"`
	// in: path
	// required: true
	CastId string `json:"castId"`
	// in: body
	Body struct {
		models.CastMember
	}
}

// swagger:response GenericSuccessResponse
type GenericSuccessResponse struct {
	// in: body
	Body struct {
		// enum: success
		Status  int    `json:"status"`
		Message string `json:"message"`
	} `json:"body"`
}

// swagger:response ValidationErrorResponse
type ValidationErrorResponse struct {
	// in : body
	Body struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	} `json:"body"`
}

// swagger:response GenericErrorResponse
type GenericErrorResponse struct {
	// in : body
	Body struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	} `json:"body"`
}
