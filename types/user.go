package types

import (
	"fmt"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost         = 12
	minFirstNameLength = 2
	minLastNameLength  = 2
	minPasswordLength  = 7
)

type CreateUserParams struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (p CreateUserParams) Validate() []string {
	errors := []string{}
	if len(p.FirstName) < minFirstNameLength {
		errors = append(errors, fmt.Sprintf("first name must be at least %d characters long", minFirstNameLength))
	}
	if len(p.LastName) < minLastNameLength {
		errors = append(errors, fmt.Sprintf("last name must be at least %d characters long", minLastNameLength))
	}
	if len(p.Password) < minPasswordLength {
		errors = append(errors, fmt.Sprintf("password must be at least %d characters long", minPasswordLength))
	}
	if !isEmailValid(p.Email) {
		errors = append(errors, fmt.Sprintf("invalid email address"))
	}

	return errors
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}

type User struct {
	ID                string `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName         string `bson:"first_name" json:"first_name"`
	LastName          string `bson:"last_name" json:"last_name"`
	Email             string `bson:"email" json:"email"`
	EncryptedPassword string `bson:"encrypted_password" json:"-"`
}

func NewUserFromParams(params CreateUserParams) (*User, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcryptCost)
	if err != nil {
		return nil, err
	}
	return &User{
		FirstName:         params.FirstName,
		LastName:          params.LastName,
		Email:             params.Email,
		EncryptedPassword: string(encpw),
	}, nil

}
