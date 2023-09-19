package entity

import (
	"errors"
	"fmt"
	"net/mail"

	"github.com/google/uuid"
)

var (
	// ErrInvalidUserID is returned when the user id is invalid
	ErrInvalidUserID = errors.New("invalid user id")
	// ErrInvalidEmail is returned when the email is invalid
	ErrInvalidEmail = errors.New("invalid email")
	// ErrInvalidUsername is returned when the username is invalid
	ErrInvalidUsername = errors.New("invalid username")
)

type UserParams struct {
	ID       string
	Username string
	Email    string
}

// User represents a user
type User struct {
	id       uuid.UUID
	username string
	email    *mail.Address
}

// NewUser creates a new user
func NewUser(params UserParams) (*User, error) {
	userID, err := uuid.Parse(params.ID)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidUserID, err.Error())
	}

	emailAddress, err := mail.ParseAddress(params.Email)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidEmail, err.Error())
	}

	if params.Username == "" {
		return nil, fmt.Errorf("%w: %s", ErrInvalidUsername, "username cannot be empty")
	}

	return &User{
		id:       userID,
		username: params.Username,
		email:    emailAddress,
	}, nil
}

func (u User) ID() uuid.UUID {
	return u.id
}

func (u User) Username() string {
	return u.username
}

func (u User) Email() *mail.Address {
	return u.email
}
