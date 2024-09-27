package server

import (
	"context"
	"go_ant_work/internal/database"
)

type UserResponse struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Avatar    string `json:"avatar"`
	Email     string `json:"email"`
	AccountId string `json:"accountId"`
}

type UserUpdateRequest struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Email  string `json:"email"`
}

func (s *Server) GetUserById(ctx context.Context, id string) (*UserResponse, error) {
	user, err := s.db.FindUserById(id)

	if err != nil {
		return nil, err
	}

	userResponse := &UserResponse{
		Id:        user.Id,
		Name:      user.Name,
		Avatar:    user.Avatar,
		Email:     user.Email,
		AccountId: user.AccountId,
	}

	return userResponse, nil
}

func (s *Server) UpdateUser(ctx context.Context, userUpdateRequest UserUpdateRequest) (string, error) {
	user := &database.User{
		Id:     userUpdateRequest.Id,
		Name:   userUpdateRequest.Name,
		Avatar: userUpdateRequest.Avatar,
		Email:  userUpdateRequest.Email,
	}

	_, err := s.db.UpdateUser(user)

	if err != nil {
		return "", err
	}

	return userUpdateRequest.Id, nil
}
