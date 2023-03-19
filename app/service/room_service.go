package service

import (
	"github.com/frchandra/chatin/app/model"
	"github.com/frchandra/chatin/app/repository"
)

type RoomService struct {
	roomRepo *repository.RoomRepository
}

func NewRoomService(roomRepo *repository.RoomRepository) *RoomService {
	return &RoomService{roomRepo: roomRepo}
}

func (s *RoomService) InsertOne(room *model.Room) (model.Room, error) {
	resultRoom, err := s.roomRepo.InsertOne(room)
	return resultRoom, err
}