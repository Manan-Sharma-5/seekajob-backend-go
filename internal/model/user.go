package model

type User struct {
    ID       string `bson:"_id,omitempty" json:"id"`
    Name     string `bson:"name" json:"name"`
    Email    string `bson:"email" json:"email"`
    Password string `bson:"password" json:"-"`
    Resume  string `bson:"resume" json:"resume"`
	IsCandidate bool   `bson:"is_candidate" json:"is_candidate"`
}

type SignupRequest struct {
    Name     string `json:"name" binding:"required"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
    IsCandidate *bool   `json:"is_candidate" binding:"required"`
}

type LoginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
    IsCandidate *bool   `json:"is_candidate" binding:"required"`
}
