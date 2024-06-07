package models

import (
	"errors"
	"time"
)

var (
	ErrNoRecord = errors.New("models: no matching records found")
	// Add a new ErrInvalidCredentials error. We'll use this later if a user
	// tries to login with an incorrect email address or password.
	ErrInvalidCredentials = errors.New("models : invalid credentials")
	// Add a new ErrDuplicateEmail error. We'll use this later if a user
	// tries to signup with an email address that's already in use.
	ErrDuplicateEmail = errors.New("models : duplicate email")
)

type Todo struct {
	ID      int
	Name    string
	Created time.Time
	Expires time.Time
}

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type SpecialTask struct {
	ID   int
	Name string
}
