package data

import (
	"testing"
	"time"
)

func TestActorDB_Get(t *testing.T) {
	mockActorModel := MockActorDB{
		Actors: map[int64]*Actor{
			1: {
				ID:        1,
				FullName:  "John Doe",
				Gender:    "male",
				BirthDate: time.Date(2021, 8, 12, 0, 0, 0, 0, time.UTC),
				Movies:    []int{1, 2},
			},
		},
	}

	t.Run("Valid", func(t *testing.T) {
		actorID := int64(1)

		actor, err := mockActorModel.Get(actorID)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if actor.ID != 1 {
			t.Errorf("expected ID to be 1, but got %d", actor.ID)
		}

		if actor.FullName != "John Doe" {
			t.Errorf("expected full name to be John Doe, but got %s", actor.FullName)
		}
	})

	t.Run("Invalid", func(t *testing.T) {
		actorID := int64(2)

		_, err := mockActorModel.Get(actorID)
		if err == nil {
			t.Error("expected ErrRecordNotFound, but got nil")
		}
	})
}

func TestActorDB_GetAll(t *testing.T) {
	mockActorModel := MockActorDB{
		Actors: map[int64]*Actor{
			1: {
				ID:        1,
				FullName:  "John Doe",
				Gender:    "male",
				BirthDate: time.Date(2021, 8, 12, 0, 0, 0, 0, time.UTC),
				Movies:    []int{1, 2},
			},
		},
	}

	t.Run("Valid", func(t *testing.T) {
		actors, err := mockActorModel.GetAll()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if len(actors) != 1 {
			t.Errorf("expected 1 actor, but got %d", len(actors))
		}
	})

	t.Run("Invalid", func(t *testing.T) {
		mockActorModel.Actors = nil

		actors, err := mockActorModel.GetAll()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if len(actors) != 0 {
			t.Errorf("expected 0 actors, but got %d", len(actors))
		}
	})
}

func TestActorDB_Insert(t *testing.T) {
	mockActorModel := MockActorDB{
		Actors: map[int64]*Actor{
			1: {
				ID:        1,
				FullName:  "John Doe",
				Gender:    "male",
				BirthDate: time.Date(2021, 8, 12, 0, 0, 0, 0, time.UTC),
				Movies:    []int{1, 2},
			},
		},
	}

	t.Run("Valid", func(t *testing.T) {
		actor := &Actor{
			FullName:  "Jane Doe",
			Gender:    "female",
			BirthDate: time.Date(2021, 8, 12, 0, 0, 0, 0, time.UTC),
		}

		err := mockActorModel.Insert(actor)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("Invalid", func(t *testing.T) {
		actor := &Actor{
			FullName:  "John Doe",
			Gender:    "male",
			BirthDate: time.Date(2021, 8, 12, 0, 0, 0, 0, time.UTC),
		}

		err := mockActorModel.Insert(actor)
		if err == nil {
			t.Error("expected ErrDuplicateName, but got nil")
		}
	})
}

func TestActorDB_Update(t *testing.T) {
	mockActorModel := MockActorDB{
		Actors: map[int64]*Actor{
			1: {
				ID:        1,
				FullName:  "John Doe",
				Gender:    "male",
				BirthDate: time.Date(2021, 8, 12, 0, 0, 0, 0, time.UTC),
				Movies:    []int{1, 2},
			},
			2: {
				ID:        2,
				FullName:  "Max Verstapen",
				Gender:    "male",
				BirthDate: time.Date(2000, 8, 12, 0, 0, 0, 0, time.UTC),
				Movies:    []int{1, 2},
			},
		},
	}

	t.Run("Valid", func(t *testing.T) {
		actor := &Actor{
			ID:       2,
			FullName: "Max Verstappen",
		}

		err := mockActorModel.Update(actor)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if mockActorModel.Actors[2].FullName != "Max Verstappen" {
			t.Errorf("expected full name to be Max Verstappen, but got %s", mockActorModel.Actors[1].FullName)
		}
	})

	t.Run("Invalid", func(t *testing.T) {
		actor := &Actor{
			ID:       2,
			FullName: "John Doe",
		}

		err := mockActorModel.Update(actor)
		if err == nil {
			t.Error("expected ErrDuplicateName, but got nil")
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		actor := &Actor{
			ID:       3,
			FullName: "John Doe",
		}

		err := mockActorModel.Update(actor)
		if err == nil {
			t.Error("expected ErrRecordNotFound, but got nil")
		}
	})
}

func TestActorDB_Delete(t *testing.T) {
	mockActorModel := MockActorDB{
		Actors: map[int64]*Actor{
			1: {
				ID:        1,
				FullName:  "John Doe",
				Gender:    "male",
				BirthDate: time.Date(2021, 8, 12, 0, 0, 0, 0, time.UTC),
				Movies:    []int{1, 2},
			},
		},
	}

	t.Run("Valid", func(t *testing.T) {
		actorID := int64(1)

		err := mockActorModel.Delete(actorID)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("Invalid", func(t *testing.T) {
		actorID := int64(1)

		err := mockActorModel.Delete(actorID)
		if err == nil {
			t.Error("expected ErrRecordNotFound, but got nil")
		}
	})

}
