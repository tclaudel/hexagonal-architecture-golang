//go:generate mockgen -source=user.go -destination=mocks/user.gen.go -package=mocks

package secondary

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/tclaudel/hexagonal-architecture-golang/internal/domain/entity"
)

var (
	ErrSavingUser        = errors.New("error saving user")
	ErrUserAlreadyExists = errors.New("error user already exists")
	ErrUpdatingUser      = errors.New("error updating user")
	ErrDeletingUser      = errors.New("error deleting user")
	ErrRetrievingUser    = errors.New("error retrieving user")
	ErrRetrievingUsers   = errors.New("error retrieving users")
	ErrUserNotFound      = errors.New("error user not found")
	ErrInvalidUserDTO    = errors.New("error invalid user dto")
)

type UserRepository interface {
	SaveUser(ctx context.Context, user *entity.User) error
	UpdateUser(ctx context.Context, user *entity.User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	GetUser(ctx context.Context, id uuid.UUID) (*entity.User, error)
	GetUsers(ctx context.Context) ([]*entity.User, error)
}
