package types

import "golang.org/x/crypto/bcrypt"

const (
	bcryptCost = 12
)

type CreateUserParams struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
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
