-- +migrate Up
CREATE TABLE IF NOT EXISTS movie_languages (
    movieid int references movies (id) on delete cascade,
    language_code varchar(2) references languages (iso_code) on delete cascade,
    PRIMARY KEY (movieid, language_code)
);
