# Movies API

This is a **Movies API** built with Golang that provides information about movies, ratings, cast, and crew. It allows users to **fetch**, **add**, **update**, and **delete** movie-related data.

## Features

- **Movies API** – Fetch, search, add, update, and delete movies.
- **Movie Ratings API** – Fetch, add, update, and remove ratings.
- **Cast API** – Fetch cast details and update cast members for movies.
- **Crew API** – Fetch and update crew members for movies.
- **Swagger** – For documentation and testing

---

## Setup and Installation

### **1. Prerequisites**

Ensure you have the following installed:

- **Go** (1.18 or later)

---

### **2. Clone the Repository**

git clone git.pride.improwised.dev/Onboarding-2025/Krupanshi-Vaishnav/go-api

cd golang-api

---

### **3.  Create the .env File**

Create a .env file in the root directory and add the following environment variables:
```
###Application Configuration
APP_PORT=127.0.0.1:4000
APP_ENV=local

###CSV File Paths
MOVIES=./data/movies_metadata.csv
CREDITS=./data/credits.csv
RATINGS=./data/ratings_small.csv
```
**Modify the paths as per your system.**

---

### **4. Install Dependencies**

Run the following command to install the dependencies:
go mod tidy

---

### **5. Run the Application**

go run app.go api
or use if you have Makefile:
make start-api

---

### **6. Endpoints**

**Movies API**

- GET /movies – List all movies with their ID and name.
- GET /movies/:movieId – Get details of a specific movie.
- GET /movies?name=moviename – Search movies by name (supports partial matches).
- GET /movies?genre=genre – Get movies filtered by genre.
- GET /movies?language=language – Get movies filtered by language.
- POST /movies – Add a new movie.
- PUT /movies/{id} – Update specific movie details.
- DELETE /movies/{id} – Delete a specific movie.

**Movie Ratings API**

- GET /ratings – List all movies with their ratings.
- GET /ratings/movie/:movieId/ratings – Get the overall rating for a particular movie.
- POST /ratings – Add a rating for a movie (user ID in body).
- PUT ratings/movies/:movieId/user/:userId/ratings – Edit a user's rating for a movie.
- DELETE /ratings/movies/:movieId/user/userId/ratings – Remove a user's rating for a movie.

**Cast API**

- GET /movies/:movieId/casts – List cast members of a particular movie.
- GET /casts/:castId/movies – List movies in which an actor played a role.
- PUT /movies/:movieId/casts/:castId – Add or update cast members for a particular movie.

**Crew API**

- GET /movies/:movieId/crew – List crew members of a particular movie.
- PUT /movies/:movieId/crew/:crewId – Add or update crew members for a particular movie.

---

### **7. Testing the API**

You can use cURL, Postman, Swagger, or any API client to test the API.

Example to get all movies:
curl -X GET http://127.0.0.1:4000/movies

To use swagger:
http://127.0.0.1:4000/docs
