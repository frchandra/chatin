package service

import (
	"github.com/frchandra/chatin/app/model"
	"github.com/frchandra/chatin/app/repository"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetOrInsertOne(user *model.User) model.User {
	userResult := s.userRepo.GetOrInsertOne(user)
	return userResult
}
