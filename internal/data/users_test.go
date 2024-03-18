package data

import (
	"filmoteka/internal/validator"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestGeneratePasswordHash(t *testing.T) {
	plaintextPassword := "password123"
	hashedPassword, err := GeneratePasswordHash(plaintextPassword)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(plaintextPassword))
	if err != nil {
		t.Errorf("password hash does not match: %v", err)
	}
}

func TestPasswordMatches(t *testing.T) {
	plaintextPassword := "password123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), bcrypt.DefaultCost)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	p := &password{
		hash: hashedPassword,
	}

	matches, err := p.Matches(plaintextPassword)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !matches {
		t.Error("passwords should match")
	}
}

func TestPasswordMatches_MismatchedHash(t *testing.T) {
	plaintextPassword := "password123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("differentpassword"), bcrypt.DefaultCost)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	p := &password{
		hash: hashedPassword,
	}

	matches, err := p.Matches(plaintextPassword)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if matches {
		t.Error("passwords should not match")
	}
}

func TestSet(t *testing.T) {
	plaintextPassword := "password123"
	p := &password{}

	err := p.Set(plaintextPassword)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	err = bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		t.Errorf("password hash does not match: %v", err)
	}

	if *p.plaintext != plaintextPassword {
		t.Errorf("plaintext password does not match: expected %s, got %s", plaintextPassword, *p.plaintext)
	}
}

func TestValidatePasswordPlaintext(t *testing.T) {
	v := validator.New()

	t.Run("Valid", func(t *testing.T) {
		ValidatePasswordPlaintext(v, "password123")
		if !v.Valid() {
			t.Errorf("unexpected error: %v", v.Errors)
		}
	})

	t.Run("Empty", func(t *testing.T) {
		ValidatePasswordPlaintext(v, "")
		if v.Valid() {
			t.Error("expected error, but got none")
		}
	})

	t.Run("TooShort", func(t *testing.T) {
		ValidatePasswordPlaintext(v, "pass")
		if v.Valid() {
			t.Error("expected error, but got none")
		}
	})

	t.Run("TooLong", func(t *testing.T) {
		ValidatePasswordPlaintext(v, "password123password123password123password123password123password123password123password123password123password123")
		if v.Valid() {
			t.Error("expected error, but got none")
		}
	})
}

func TestValidateUser(t *testing.T) {
	v := validator.New()

	t.Run("Valid", func(t *testing.T) {
		user := &User{
			Name: "John Doe",
			Role: "user",
		}

		err := user.Password.Set("password123")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		ValidateUser(v, user)
		if !v.Valid() {
			t.Errorf("unexpected error: %v", v.Errors)
		}
	})

	t.Run("MissingName", func(t *testing.T) {
		user := &User{
			Name: "",
			Role: "user",
		}

		err := user.Password.Set("password123")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		ValidateUser(v, user)
		if v.Valid() {
			t.Error("expected error, but got none")
		}
	})

	t.Run("InvalidNameLength", func(t *testing.T) {
		user := &User{
			Name: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed ut ultrices nunc, at tincidunt nisl. Nulla facilisi. Sed id nunc auctor, ultrices nunc id, aliquam nunc. Sed nec nunc ut nunc aliquet tincidunt. Sed id nunc auctor, ultrices nunc id, aliquam nunc. Sed nec nunc ut nunc aliquet tincidunt.",
			Role: "user",
		}

		err := user.Password.Set("password123")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		ValidateUser(v, user)
		if v.Valid() {
			t.Error("expected error, but got none")
		}
	})

	t.Run("InvalidRole", func(t *testing.T) {
		user := &User{
			Name: "John Doe",
			Role: "guest",
		}

		err := user.Password.Set("password123")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		ValidateUser(v, user)
		if v.Valid() {
			t.Error("expected error, but got none")
		}
	})

	t.Run("MissingPasswordHash", func(t *testing.T) {
		user := &User{
			Name: "John Doe",
			Role: "user",
		}

		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic, but got none")
			}
		}()

		ValidateUser(v, user)
	})
}

func TestUserDB_Insert(t *testing.T) {
	mockDB := MockUserDB{
		Users: map[string]*User{
			"John Doe": {Name: "John Doe", Role: "user"},
		},
	}

	t.Run("Valid", func(t *testing.T) {
		user := &User{
			Name: "Jane Doe",
			Role: "user",
		}

		err := mockDB.Insert(user)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if user.ID == 0 {
			t.Error("expected ID to be set")
		}

		if _, exists := mockDB.Users[user.Name]; !exists {
			t.Error("user was not inserted into the database")
		}
	})

	t.Run("DuplicateName", func(t *testing.T) {
		user := &User{
			Name: "John Doe",
			Role: "user",
		}

		lenBefore := len(mockDB.Users)

		err := mockDB.Insert(user)
		if err == nil {
			t.Error("expected error, but got none")
		}

		if err != ErrDuplicateName {
			t.Errorf("expected ErrDuplicateName, but got %v", err)
		}

		if user.ID != 0 {
			t.Error("expected ID to be unset")
		}

		if len(mockDB.Users) != lenBefore {
			t.Errorf("expected number of users to be %d, but got %d", lenBefore, len(mockDB.Users))
		}

		if _, exists := mockDB.Users[user.Name]; !exists {
			t.Error("expected user to remain in the database")
		}
	})
}

func TestMockUserDB_Get(t *testing.T) {
	mockDB := MockUserDB{
		Users: map[string]*User{
			"John Doe": {Name: "John Doe", Role: "user"},
		},
	}

	t.Run("ExistingUser", func(t *testing.T) {
		username := "John Doe"
		expectedUser := &User{Name: "John Doe", Role: "user"}

		user, err := mockDB.Get(username)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if user == nil {
			t.Error("expected user, but got nil")
		}

		if user.Name != expectedUser.Name || user.Role != expectedUser.Role {
			t.Errorf("expected user %v, but got %v", expectedUser, user)
		}
	})

	t.Run("NonExistingUser", func(t *testing.T) {
		username := "Jane Doe"

		user, err := mockDB.Get(username)
		if err != ErrRecordNotFound {
			t.Errorf("expected ErrRecordNotFound, but got %v", err)
		}

		if user != nil {
			t.Errorf("expected nil user, but got %v", user)
		}
	})
}
