package service

import (
	"errors"
	"github.com/frchandra/chatin/app/model"
	"github.com/frchandra/chatin/app/repository"
	"github.com/frchandra/chatin/app/util"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo  *repository.UserRepository
	tokenUtil *util.TokenUtil
}

func NewUserService(userRepo *repository.UserRepository, tokenUtil *util.TokenUtil) *UserService {
	return &UserService{userRepo: userRepo, tokenUtil: tokenUtil}
}

func (s *UserService) GetOneById(id string) (model.User, error) {
	resultUser, err := s.userRepo.GetOneById(id)
	if err != nil {
		return resultUser, err
	}
	return resultUser, nil
}

func (s *UserService) GetOrInsertOne(user *model.User) (model.User, error) {
	hashedCred, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost) //hash the credential
	user.Password = string(hashedCred)
	userResult, err := s.userRepo.GetOrInsertOne(user) //store new user to the database
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

func (s *UserService) ValidateLogin(user *model.User) (model.User, error) {

	resultUser, err := s.userRepo.GetOneByNameOrEmail(user) //get the user by credential pairs (email/name & password)
	if err != nil {
		return resultUser, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(resultUser.Password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return resultUser, errors.New("credential authentication error")
	}
	return resultUser, nil
}

func (s *UserService) GenerateToken(user *model.User) (*util.TokenDetails, error) {
	tokenDetails, err := s.tokenUtil.CreateToken(user.Id.Hex()) //create token for the user
	if err != nil {
		return tokenDetails, errors.New("credential authentication error")
	}

	if err = s.tokenUtil.StoreAuthn(user.Id.Hex(), tokenDetails); err != nil { //store the token to redis
		return tokenDetails, errors.New("credential preparation error")
	}

	return tokenDetails, nil //return the new created token
}
