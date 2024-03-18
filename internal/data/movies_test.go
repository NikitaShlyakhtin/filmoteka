package data

import (
	"filmoteka/internal/validator"
	"testing"
	"time"
)

func TestMockMovieDB_Insert(t *testing.T) {
	mockModel := MockMovieDB{
		Actors: make(map[int64]*Actor),
		Movies: make(map[int64]*Movie),
	}

	actor1 := &Actor{
		ID:        1,
		FullName:  "John Doe",
		Gender:    "male",
		BirthDate: time.Date(2021, 8, 12, 0, 0, 0, 0, time.UTC),
	}
	actor2 := &Actor{
		ID:        2,
		FullName:  "Jane Doe",
		Gender:    "female",
		BirthDate: time.Date(2021, 8, 12, 0, 0, 0, 0, time.UTC),
	}

	mockModel.Actors[actor1.ID] = actor1
	mockModel.Actors[actor2.ID] = actor2

	movie := &Movie{
		ID:          1,
		Title:       "Movie 1",
		Description: "Description 1",
		Rating:      8.0,
		ReleaseDate: time.Date(2021, 8, 12, 0, 0, 0, 0, time.UTC),
		Actors:      []int64{1, 2},
	}

	mockModel.Movies[movie.ID] = movie

	t.Run("Valid", func(t *testing.T) {
		movie := &Movie{
			Title:       "Movie 2",
			Description: "Description 2",
			Rating:      8.5,
			ReleaseDate: time.Date(2000, 8, 12, 0, 0, 0, 0, time.UTC),
			Actors:      []int64{1},
		}

		err := mockModel.Insert(movie)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if mockModel.Movies[movie.ID] == nil {
			t.Error("expected movie to be inserted")
		}
	})

	t.Run("DuplicateTitle", func(t *testing.T) {
		movie := &Movie{
			Title:       "Movie 1",
			Description: "Description 1",
			Rating:      8.0,
			ReleaseDate: time.Date(2021, 8, 12, 0, 0, 0, 0, time.UTC),
			Actors:      []int64{1, 2},
		}

		err := mockModel.Insert(movie)
		if err == nil {
			t.Error("expected error, but got nil")
		}
	})

	t.Run("InvalidActors", func(t *testing.T) {
		movie := &Movie{
			Title:       "Movie 3",
			Description: "Description 3",
			Rating:      7.8,
			ReleaseDate: time.Date(2015, 8, 12, 0, 0, 0, 0, time.UTC),
			Actors:      []int64{3},
		}

		err := mockModel.Insert(movie)
		if err == nil {
			t.Error("expected error, but got nil")
		}
	})
}

func TestMockMovieDB_Get(t *testing.T) {
	mockModel := MockMovieDB{
		Actors: make(map[int64]*Actor),
		Movies: make(map[int64]*Movie),
	}

	actor1 := &Actor{
		ID:        1,
		FullName:  "John Doe",
		Gender:    "male",
		BirthDate: time.Date(2021, 8, 12, 0, 0, 0, 0, time.UTC),
	}
	actor2 := &Actor{
		ID:        2,
		FullName:  "Jane Doe",
		Gender:    "female",
		BirthDate: time.Date(2021, 8, 12, 0, 0, 0, 0, time.UTC),
	}

	mockModel.Actors[actor1.ID] = actor1
	mockModel.Actors[actor2.ID] = actor2

	movie := &Movie{
		ID:          1,
		Title:       "Movie 1",
		Description: "Description 1",
		Rating:      8.0,
		ReleaseDate: time.Date(2021, 8, 12, 0, 0, 0, 0, time.UTC),
		Actors:      []int64{1, 2},
	}

	mockModel.Movies[movie.ID] = movie

	t.Run("Valid", func(t *testing.T) {
		m, err := mockModel.Get(movie.ID)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if m == nil {
			t.Error("expected movie, but got nil")
		}
	})

	t.Run("InvalidID", func(t *testing.T) {
		m, err := mockModel.Get(2)
		if err != ErrRecordNotFound {
			t.Error("expected ErrRecordNotFound, but got nil")
		}

		if m != nil {
			t.Error("expected nil, but got movie")
		}
	})
}

func TestMockMovieDB_GetAll(t *testing.T) {
	mockModel := MockMovieDB{
		Actors: make(map[int64]*Actor),
		Movies: make(map[int64]*Movie),
	}

	actor1 := &Actor{
		ID:        1,
		FullName:  "John Doe",
		Gender:    "male",
		BirthDate: time.Date(2021, 8, 12, 0, 0, 0, 0, time.UTC),
	}
	actor2 := &Actor{
		ID:        2,
		FullName:  "Jane Doe",
		Gender:    "female",
		BirthDate: time.Date(2021, 8, 12, 0, 0, 0, 0, time.UTC),
	}

	mockModel.Actors[actor1.ID] = actor1
	mockModel.Actors[actor2.ID] = actor2

	movie := &Movie{
		ID:          1,
		Title:       "Movie 1",
		Description: "Description 1",
		ReleaseDate: time.Date(2021, 8, 12, 0, 0, 0, 0, time.UTC),
		Rating:      8.0,
		Actors:      []int64{1, 2},
	}

	mockModel.Movies[movie.ID] = movie

	movie = &Movie{
		ID:          2,
		Title:       "Movie 2",
		Description: "Description 2",
		ReleaseDate: time.Date(2000, 8, 12, 0, 0, 0, 0, time.UTC),
		Rating:      8.5,
		Actors:      []int64{1, 2},
	}

	mockModel.Movies[movie.ID] = movie

	t.Run("ValidDefault", func(t *testing.T) {
		movies, err := mockModel.GetAll(Filters{})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if movies[0].Rating < movies[1].Rating {
			t.Error("expected movies to be sorted by rating in descending order")
		}
	})

	t.Run("ValidRatingDesc", func(t *testing.T) {
		movies, err := mockModel.GetAll(Filters{Sort: "rating"})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if movies[0].Rating > movies[1].Rating {
			t.Error("expected movies to be sorted by rating in ascending order")
		}
	})

	t.Run("ValidTitle", func(t *testing.T) {
		movies, err := mockModel.GetAll(Filters{Sort: "title"})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if movies[0].Title > movies[1].Title {
			t.Error("expected movies to be sorted by rating in ascending order")
		}
	})

	t.Run("ValidTitleDesc", func(t *testing.T) {
		movies, err := mockModel.GetAll(Filters{Sort: "-title"})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if movies[0].Title < movies[1].Title {
			t.Error("expected movies to be sorted by rating in descending order")
		}
	})

	t.Run("ValidReleaseDate", func(t *testing.T) {
		movies, err := mockModel.GetAll(Filters{Sort: "release_date"})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if movies[0].ReleaseDate.After(movies[1].ReleaseDate) {
			t.Error("expected movies to be sorted by rating in ascending order")
		}
	})

	t.Run("ValidReleaseDateDesc", func(t *testing.T) {
		movies, err := mockModel.GetAll(Filters{Sort: "-release_date"})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if movies[0].ReleaseDate.Before(movies[1].ReleaseDate) {
			t.Error("expected movies to be sorted by rating in descending order")
		}
	})
}

func TestMockMovieDB_Update(t *testing.T) {
	mockModel := MockMovieDB{
		Actors: make(map[int64]*Actor),
		Movies: make(map[int64]*Movie),
	}

	actor1 := &Actor{
		ID:        1,
		FullName:  "John Doe",
		Gender:    "male",
		BirthDate: time.Date(2021, 8, 12, 0, 0, 0, 0, time.UTC),
	}
	actor2 := &Actor{
		ID:        2,
		FullName:  "Jane Doe",
		Gender:    "female",
		BirthDate: time.Date(2021, 8, 12, 0, 0, 0, 0, time.UTC),
	}

	mockModel.Actors[actor1.ID] = actor1
	mockModel.Actors[actor2.ID] = actor2

	movie := &Movie{
		ID:          1,
		Title:       "Movie 1",
		Description: "Description 1",
		ReleaseDate: time.Date(2021, 8, 12, 0, 0, 0, 0, time.UTC),
		Rating:      8.0,
		Actors:      []int64{1, 2},
	}

	mockModel.Movies[movie.ID] = movie

	movie = &Movie{
		ID:          2,
		Title:       "Movie 2",
		Description: "Description 2",
		ReleaseDate: time.Date(2000, 8, 12, 0, 0, 0, 0, time.UTC),
		Rating:      8.5,
		Actors:      []int64{1, 2},
	}

	mockModel.Movies[movie.ID] = movie

	t.Run("Valid", func(t *testing.T) {
		movie := &Movie{
			ID:     1,
			Rating: 9.0,
		}

		err := mockModel.Update(*movie)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if mockModel.Movies[movie.ID].Rating != 9.0 {
			t.Error("expected movie to be updated")
		}
	})

	t.Run("InvalidID", func(t *testing.T) {
		movie := &Movie{
			ID:     3,
			Rating: 7.0,
		}

		err := mockModel.Update(*movie)
		if err == nil {
			t.Error("expected error, but got nil")
		}

		if err != ErrRecordNotFound {
			t.Errorf("expected ErrRecordNotFound, but got %v", err)
		}
	})

	t.Run("InvalidActor", func(t *testing.T) {
		movie := &Movie{
			ID:     1,
			Actors: []int64{1, 10},
		}

		err := mockModel.Update(*movie)
		if err == nil {
			t.Error("expected error, but got nil")
		}

		if err != ErrActorsNotFound {
			t.Errorf("expected ErrInvalidActor, but got %v", err)
		}
	})
}

func TestValidateMovie(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		v := validator.New()

		movie := &Movie{
			Title:       "Movie 1",
			Description: "Description 1",
			Rating:      8.0,
			ReleaseDate: time.Date(2021, 8, 12, 0, 0, 0, 0, time.UTC),
			Actors:      []int64{1, 2},
		}

		ValidateMovie(v, movie)
		if !v.Valid() {
			t.Errorf("unexpected error: %v", v.Errors)
		}
	})

	t.Run("TitleMustBeProvided", func(t *testing.T) {
		v := validator.New()

		movie := &Movie{
			Title:       "",
			Description: "Description 1",
			Rating:      8.0,
			ReleaseDate: time.Date(2021, 8, 12, 0, 0, 0, 0, time.UTC),
			Actors:      []int64{1, 2},
		}

		ValidateMovie(v, movie)
		if v.Valid() {
			t.Error(`expected error "title": "must be provided", got none`)
		}

		if v.Errors["title"] != "must be provided" {
			t.Error(`expected error "title": "must be provided", got none`)
		}
	})

	t.Run("TitleExceedsLimit", func(t *testing.T) {
		v := validator.New()

		movie := &Movie{
			Title:       "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nulla facilisis, urna non tincidunt ultrices, justo nisl fringilla ipsum, sed lacinia elit nunc id nunc.",
			Description: "Description 1",
			Rating:      8.0,
			ReleaseDate: time.Date(2021, 8, 12, 0, 0, 0, 0, time.UTC),
			Actors:      []int64{1, 2},
		}

		ValidateMovie(v, movie)
		if v.Valid() {
			t.Error(`expected error "title": "must be no more than 150 symbols", got none`)
		}

		if v.Errors["title"] != "must be no more than 150 symbols" {
			t.Error(`expected error "title": "must be no more than 150 symbols", got none`)
		}
	})

	t.Run("DescriptionMustBeProvided", func(t *testing.T) {
		v := validator.New()

		movie := &Movie{
			Title:       "Movie 1",
			Description: "",
			Rating:      8.0,
			ReleaseDate: time.Date(2021, 8, 12, 0, 0, 0, 0, time.UTC),
			Actors:      []int64{1, 2},
		}

		ValidateMovie(v, movie)
		if v.Valid() {
			t.Error(`expected error "description": "must be provided", got none`)
		}

		if v.Errors["description"] != "must be provided" {
			t.Error(`expected error "description": "must be provided", got none`)
		}
	})

	t.Run("DescriptionExceedsLimit", func(t *testing.T) {
		v := validator.New()

		movie := &Movie{
			Title:       "Movie 1",
			Description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Suspendisse vel elementum lectus. Etiam aliquam, dui ac scelerisque feugiat, lacus orci bibendum turpis, vitae vehicula erat leo scelerisque orci. Morbi mollis scelerisque erat eget gravida. Nullam nec imperdiet nisi, non consectetur sem. Nunc volutpat elit in ultricies feugiat. Suspendisse potenti. Donec facilisis diam tristique erat bibendum, ut maximus lorem malesuada. Nulla in ultrices est. Proin rutrum odio tortor, vel consequat magna venenatis in. Aenean volutpat ipsum nisi, ac venenatis massa hendrerit sit amet. Nam ac dapibus nisi. Etiam quis est iaculis, mollis eros at, luctus nunc. Phasellus ultrices elit vel fringilla lobortis. Nullam felis risus, semper vitae bibendum sit amet, scelerisque et orci. Fusce non viverra metus. Phasellus dignissim mattis convallis. Nulla vel tortor lectus. Curabitur arcu lorem, lacinia sed purus malesuada, tincidunt porta nibh. Cras ut risus at erat malesuada mattis. Ut lacinia dolor non nibh rhoncus bibendum. In a ex eget lectus commodo eleifend. Fusce arcu lorem, suscipit quis ornare in, elementum sit amet ex. Aliquam eget tortor ullamcorper, mattis neque eu, pharetra orci. Praesent posuere felis at dolor facilisis, eget semper neque mattis. Duis egestas faucibus euismod. Sed felis nunc, tincidunt vitae orci ac, finibus rhoncus sapien. In convallis cursus rutrum. Orci varius natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. In sit amet sagittis libero. Vivamus posuere quam rhoncus iaculis accumsan. Nullam magna mi, ornare in justo quis, dapibus sagittis urna. Aenean ac condimentum arcu, a sagittis magna. Nullam aliquam dolor vitae libero blandit, vel luctus lorem faucibus. Sed dapibus lorem sit amet gravida pellentesque. Donec id bibendum urna, vitae ornare augue. Phasellus elit purus, hendrerit sit amet interdum non, venenatis ut libero. Cras id sollicitudin risus. Sed pulvinar lacus nec dapibus maximus.",
			Rating:      8.0,
			ReleaseDate: time.Date(2021, 8, 12, 0, 0, 0, 0, time.UTC),
			Actors:      []int64{1, 2},
		}

		ValidateMovie(v, movie)
		if v.Valid() {
			t.Error(`expected error "description": "must be no more than 1000 symbols", got none`)
		}

		if v.Errors["description"] != "must be no more than 1000 symbols" {
			t.Error(`expected error "description": "must be no more than 1000 symbols", got none`)
		}
	})

	t.Run("RatingInvalid", func(t *testing.T) {
		v := validator.New()

		movie := &Movie{
			Title:       "Movie 1",
			Description: "Description 1",
			Rating:      11.0,
			ReleaseDate: time.Date(2021, 8, 12, 0, 0, 0, 0, time.UTC),
			Actors:      []int64{1, 2},
		}

		ValidateMovie(v, movie)
		if v.Valid() {
			t.Error(`expected error "rating": "must be between 0 and 10", got none`)
		}

		if v.Errors["rating"] != "must be between 0 and 10" {
			t.Error(`expected error "rating": "must be between 0 and 10", got none`)
		}
	})

	t.Run("ReleaseDateMustBeProvided", func(t *testing.T) {
		v := validator.New()

		movie := &Movie{
			Title:       "Movie 1",
			Description: "Description 1",
			Rating:      8.0,
			Actors:      []int64{1, 2},
		}

		ValidateMovie(v, movie)
		if v.Valid() {
			t.Error(`expected error "release_date": "must be provided", got none`)
		}

		if v.Errors["release_date"] != "must be provided" {
			t.Error(`expected error "release_date": "must be provided", got none`)
		}
	})

	t.Run("ActorsMustBeProvided", func(t *testing.T) {
		v := validator.New()

		movie := &Movie{
			Title:       "Movie 1",
			Description: "Description 1",
			Rating:      8.0,
			ReleaseDate: time.Date(2021, 8, 12, 0, 0, 0, 0, time.UTC),
		}

		ValidateMovie(v, movie)
		if v.Valid() {
			t.Error(`expected error "actors": "must contain at least one actor", got none`)
		}

		if v.Errors["actors"] != "must contain at least one actor" {
			t.Error(`expected error "actors": "must contain at least one actor", got none`)
		}
	})

}
