package domain

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestRegisterUser(t *testing.T) {
	name, email, password := "expected name", "email@mail.com", "password"

	t.Run("register a user", func(t *testing.T) {
		got, _ := UserRegister(name, email, password)

		if got.ID == "" {
			t.Error("expected an ID, but got nothing")
		}

		if got.CreatedAt == nil {
			t.Error("expected a created date, but got a nil")
		}

		if got.UpdatedAt != nil {
			t.Error("expected a nil, but got a date")
		}

		if got.Name != name {
			t.Errorf("expected %q, but got %q", name, got.Name)
		}

		if got.Email != email {
			t.Errorf("expected %q, but got %q", email, got.Email)
		}

		if err := bcrypt.CompareHashAndPassword([]byte(got.Password), []byte(password)); err != nil {
			t.Error("password is not valid")
		}

		if got.Profile != PROFILE_CLIENT {
			t.Errorf("expected %q, but got %q", PROFILE_CLIENT, got.Profile)
		}
	})

	t.Run("empty name", func(t *testing.T) {
		_, err := UserRegister("", email, password)
		if err == nil {
			t.Error("expected an error, but got a nil")
		}
	})

	t.Run("empty email", func(t *testing.T) {
		_, err := UserRegister(name, "", password)
		if err == nil {
			t.Error("expected an error, but got a nil")
		}
	})

	t.Run("invalid email", func(t *testing.T) {
		_, err := UserRegister(name, "mail.com", password)
		if err == nil {
			t.Error("expected an error, but got a nil")
		}
	})

	t.Run("valid email .br", func(t *testing.T) {
		emailBr := "email@mail.com.br"
		got, err := UserRegister(name, emailBr, password)
		if err != nil {
			t.Errorf("expected a nil, but got an error %v", err)
		}
		if got.Email != emailBr {
			t.Errorf("expected %q, but got %q", emailBr, email)
		}
	})

	t.Run("empty password", func(t *testing.T) {
		_, err := UserRegister(name, email, "")
		if err == nil {
			t.Error("expected an error, but got a nil")
		}
	})

	t.Run("invalid password", func(t *testing.T) {
		_, err := UserRegister(name, email, "123")
		if err == nil {
			t.Error("expected an error, but got nil")
		}
	})
}

func TestIsCorrectPassword(t *testing.T) {
	user, _ := UserRegister("name", "email@mail.com", "password")

	t.Run("valid password", func(t *testing.T) {
		if !user.IsCorrectPassword("password") {
			t.Error("expected a valid password, but got an invalid")
		}
	})

	t.Run("invalid password", func(t *testing.T) {
		if user.IsCorrectPassword("password1") {
			t.Error("expected an invalid password, but got a valid")
		}
	})
}
