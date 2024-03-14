CREATE TYPE gender AS ENUM ('male', 'female');

CREATE TABLE Actors (
    actor_id SERIAL PRIMARY KEY,
    full_name VARCHAR(200) UNIQUE NOT NULL,
    gender gender NOT NULL,
    birth_date DATE NOT NULL
);

CREATE TABLE Movies (
    movie_id SERIAL PRIMARY KEY,
    title VARCHAR(150) UNIQUE NOT NULL,
    description VARCHAR(1000) NOT NULL,
    release_date DATE NOT NULL,
    rating DECIMAL(3,1) NOT NULL CHECK (rating >= 0 AND rating <= 10)
);

CREATE TABLE Movies_actors (
    movie_id INT REFERENCES movies(movie_id) ON DELETE CASCADE,
    actor_id INT REFERENCES actors(actor_id) ON DELETE CASCADE,
    PRIMARY KEY (movie_id, actor_id)
);

CREATE TYPE user_role AS ENUM ('user', 'admin');

CREATE TABLE Users (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(100) NOT NULL,
    role user_role NOT NULL
);
