package mongo

import (
	"github.com/tclaudel/hexagonal-architecture-golang/internal/domain/entity"
	"github.com/tclaudel/hexagonal-architecture-golang/internal/ports/secondary"
)

type User struct {
	ID       string `bson:"_id"`
	Username string `bson:"username"`
	Email    string `bson:"email"`
}

func UserEntityToDTO(user *entity.User) User {
	email := user.Email()

	return User{
		ID:       user.ID().String(),
		Username: user.Username(),
		Email:    email.String(),
	}
}

func UserDTOToEntity(userDTO User) (*entity.User, error) {
	userParam := entity.UserParams{
		ID:       userDTO.ID,
		Username: userDTO.Username,
		Email:    userDTO.Email,
	}

	user, err := entity.NewUser(userParam)
	if err != nil {
		return nil, secondary.ErrInvalidUserDTO
	}

	return user, nil
}
