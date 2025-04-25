package service

import (
	"fmt"

	"github.com/manan-sharma-5/seekajob-backend/internal/model"
	"github.com/manan-sharma-5/seekajob-backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

func SignupUser(req model.SignupRequest) (string, error) {
    taken, err := repository.IsEmailTaken(req.Email)
    if err != nil {
        return "", err
    }
    if taken {
        return "", ErrEmailAlreadyInUse
    }

    hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }

    user := &model.User{
        Name:     req.Name,
        Email:    req.Email,
        Password: string(hashed),
    }

    userCreated, err := repository.CreateUser(user)

    if err != nil {
        return "", err
    }
    if userCreated == nil {
        return "", fmt.Errorf("user not created")
    }

    return userCreated.ID, nil
}

func LoginUser(req model.LoginRequest) (string, error) {
    userVerified, err := repository.VerifyUser(req.Email, req.Password, req.IsCandidate)
    if err != nil {
        return "", err
    }
    if userVerified == nil {
        return "", ErrInvalidCredentials
    }
    return userVerified.ID, nil
}



var ErrEmailAlreadyInUse = fmt.Errorf("email already registered")
var ErrInvalidCredentials = fmt.Errorf("invalid credentials")
