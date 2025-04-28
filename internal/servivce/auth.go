package service

import (
	"fmt"

	"github.com/manan-sharma-5/seekajob-backend/internal/model"
	"github.com/manan-sharma-5/seekajob-backend/internal/repository"
)

func SignupUser(req model.SignupRequest) (*model.User, error) {
    taken, err := repository.IsEmailTaken(req.Email)
    if err != nil {
        return nil, err
    }
    if taken {
        return nil, ErrEmailAlreadyInUse
    }

    if err != nil {
        return nil, err
    }

    user := &model.User{
        Name:     req.Name,
        Email:    req.Email,
        Password: req.Password,
        IsCandidate: req.IsCandidate,
    }

    userCreated, err := repository.CreateUser(user)

    if err != nil {
        return nil, err
    }
    if userCreated == nil {
        return nil, fmt.Errorf("user not created")
    }

    return userCreated, nil
}

func LoginUser(req model.LoginRequest) (*model.User, error) {
    userVerified, err := repository.VerifyUser(req.Email, req.Password, req.IsCandidate)
    if err != nil {
        return nil, err
    }
    if userVerified == nil {
        return nil, ErrInvalidCredentials
    }
    return userVerified, nil
}

func GetUserByID(id string) (*model.UserWithApplicationsAndJobs, error) {
    user, err := repository.GetUserByID(id)
    if err != nil {
        return nil, err
    }
    if user == nil {
        return nil, fmt.Errorf("user not found")
    }
    return user, nil
}



var ErrEmailAlreadyInUse = fmt.Errorf("email already registered")
var ErrInvalidCredentials = fmt.Errorf("invalid credentials")
