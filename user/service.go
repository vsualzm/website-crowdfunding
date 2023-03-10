package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	RegisterUserInput(input RegisterUserInput) (User, error)
	Login(LoginInput) (User, error)
}

type service struct {
	repository Repository
}

// RegisterUserInput implements Service
func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) RegisterUserInput(input RegisterUserInput) (User, error) {
	user := User{}
	user.Name = input.Name
	user.Email = input.Email
	user.Occupation = input.Occupation
	PasswordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)

	if err != nil {
		return user, err
	}

	user.PasswordHash = string(PasswordHash)
	user.Role = "user"

	NewUser, err := s.repository.Save(user)
	if err != nil {
		return NewUser, err
	}

	return NewUser, nil
}

func (s *service) Login(input LoginInput) (User, error) {
	email := input.Email
	password := input.Password

	user, err := s.repository.FindByEmail(email)
	if err != nil {
		return user, err
	}

	if user.ID == 0 {
		return user, errors.New("No user found on thath email")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return user, err
	}

	return user, nil

}

// mapping strcut input ke struct user
// simpan struct User melalui tepository
