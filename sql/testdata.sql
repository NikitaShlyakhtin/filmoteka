INSERT INTO Actors (full_name, gender, birth_date) VALUES
    ('Brad Pitt', 'male', '1963-12-18'),
    ('Angelina Jolie', 'female', '1975-06-04'),
    ('Robert Downey Jr.', 'male', '1965-04-04'),
    ('Scarlett Johansson', 'female', '1984-11-22'),
    ('Denzel Washington', 'male', '1954-12-28'),
    ('Kate Winslet', 'female', '1975-10-05');

INSERT INTO Movies (title, description, release_date, rating) VALUES
    ('The Dark Knight', 'When the menace known as the Joker wreaks havoc and chaos on the people of Gotham, Batman must accept one of the greatest psychological and physical tests of his ability to fight injustice.', '2008-07-18', 9.0),
    ('Titanic', 'A seventeen-year-old aristocrat falls in love with a kind but poor artist aboard the luxurious, ill-fated R.M.S. Titanic.', '1997-12-19', 7.8),
    ('Avengers: Endgame', 'After the devastating events of Avengers: Infinity War, the universe is in ruins. With the help of remaining allies, the Avengers assemble once more in order to reverse Thanos'' actions and restore balance to the universe.', '2019-04-26', 8.4),
    ('The Matrix', 'A computer hacker learns from mysterious rebels about the true nature of his reality and his role in the war against its controllers.', '1999-03-31', 8.7),
    ('The Shawshank Redemption', 'Two imprisoned men bond over a number of years, finding solace and eventual redemption through acts of common decency.', '1994-10-14', 9.3),
    ('Pulp Fiction', 'The lives of two mob hitmen, a boxer, a gangster and his wife, and a pair of diner bandits intertwine in four tales of violence and redemption.', '1994-10-14', 8.9),
    ('Inglourious Basterds', 'In Nazi-occupied France during World War II, a plan to assassinate Nazi leaders by a group of Jewish U.S. soldiers coincides with a theatre owner''s vengeful plans for the same.', '2009-08-21', 8.3),
    ('The Godfather', 'The aging patriarch of an organized crime dynasty transfers control of his clandestine empire to his reluctant son.', '1972-03-24', 9.2);

INSERT INTO Movies_actors (movie_id, actor_id) VALUES
    (1, 4), -- Scarlett Johansson in The Dark Knight
    (3, 3), -- Robert Downey Jr. in Avengers: Endgame
    (3, 4), -- Scarlett Johansson in Avengers: Endgame
    (5, 1), -- Brad Pitt in The Shawshank Redemption
    (5, 5), -- Denzel Washington in The Shawshank Redemption
    (6, 1), -- Brad Pitt in Pulp Fiction
    (6, 2), -- Angelina Jolie in Pulp Fiction
    (7, 1), -- Brad Pitt in Inglourious Basterds
    (7, 2); -- Angelina Jolie in Inglourious Basterds

INSERT INTO Users (username, password_hash, role) VALUES
    ('admin', '$2a$12$6EASj861izXc62eMuaQGXOAOG/eWGHHcAYZTEP8GSoNG0qEWbRpDm', 'admin'), -- password: password123
    ('user', '$2a$12$6EASj861izXc62eMuaQGXOAOG/eWGHHcAYZTEP8GSoNG0qEWbRpDm', 'user'); -- password: password123