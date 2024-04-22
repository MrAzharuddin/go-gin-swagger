package daos

import (
	"errors"
	"github.com/MrAzharuddin/swagger-test/user-service/pkg/rest/server/models"
)

var users = make(map[int64]*models.User)

type UserDao struct {
}

func NewUserDao() (*UserDao, error) {
	return &UserDao{}, nil
}

func (userDao *UserDao) CreateUser(user *models.User) (*models.User, error) {
	users[user.Id] = user

	return user, nil
}

func (userDao *UserDao) GetUser(id int64) (*models.User, error) {
	if user, ok := users[id]; ok {
		return user, nil
	}

	return &models.User{}, errors.New("user not found")
}

func (userDao *UserDao) UpdateUser(id int64, user *models.User) (*models.User, error) {
	if id != user.Id {
		return nil, errors.New("id and payload don't match")
	}
	users[user.Id] = user

	return user, nil
}

func (userDao *UserDao) DeleteUser(id int64) error {
	if _, ok := users[id]; ok {
		delete(users, id)
		return nil
	}

	return errors.New("user not found")
}
