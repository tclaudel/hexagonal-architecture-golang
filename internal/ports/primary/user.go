package primary

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/tclaudel/hexagonal-architecture-golang/internal/domain/entity"
)

var (
	ErrCreatingUser      = errors.New("error creating user")
	ErrUserAlreadyExists = errors.New("error user already exists")
	ErrUpdatingUser      = errors.New("error updating user")
	ErrRemovingUser      = errors.New("error removing user")
	ErrRetrievingUser    = errors.New("error retrieving user")
	ErrRetrievingUsers   = errors.New("error retrieving users")
	ErrUserNotFound      = errors.New("error user not found")
)

type UserUseCase interface {
	CreateUser(ctx context.Context, user *entity.User) error
	UpdateUser(ctx context.Context, user *entity.User) error
	RemoveUser(ctx context.Context, id uuid.UUID) error
	RetrieveUser(ctx context.Context, id uuid.UUID) (*entity.User, error)
	RetrieveUsers(ctx context.Context) ([]*entity.User, error)
}
