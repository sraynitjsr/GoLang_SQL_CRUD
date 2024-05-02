package service

import (
	"github.com/sraynitjsr/model"
	"github.com/sraynitjsr/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo}
}

func (service *UserService) GetUsers() ([]model.User, error) {
	return service.repo.GetAll()
}

func (service *UserService) GetUser(id int) (model.User, error) {
	return service.repo.GetByID(id)
}

func (service *UserService) CreateUser(user model.User) (model.User, error) {
	return service.repo.Create(user)
}

func (service *UserService) UpdateUser(id int, user model.User) (model.User, error) {
	return service.repo.Update(id, user)
}

func (service *UserService) DeleteUser(id int) error {
	return service.repo.Delete(id)
}
