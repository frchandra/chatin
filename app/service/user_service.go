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

func (s *UserService) GetOneById(id string) (model.User, error) {
	resultUser, err := s.userRepo.GetOneById(id)
	if err != nil {
		return resultUser, err
	}
	return resultUser, nil
}

func (s *UserService) GetOrInsertOne(user *model.User) (model.User, error) {
	userResult, err := s.userRepo.GetOrInsertOne(user)
	if err != nil {
		return userResult, err
	}
	return userResult, nil
}

func (s *UserService) InsertOne(user *model.User) (model.User, error) {
	userResult, err := s.userRepo.InsertOne(user)
	if err != nil {
		return userResult, err
	}
	return userResult, nil
}

func (s *UserService) ValidateCredential() {

}
