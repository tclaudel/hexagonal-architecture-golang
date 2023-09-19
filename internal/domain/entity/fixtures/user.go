package fixtures

import (
	"testing"

	"github.com/tclaudel/hexagonal-architecture-golang/internal/domain/entity"
)

const (
	user1ID = "4e671dcd-f0c5-47ee-b00a-d1f74048212c"
	user2ID = "a26502d5-1095-4f0e-a9b4-ba53898a6650"
)

func WithName(name string) UserOpt {
	return func(params *entity.UserParams) {
		params.Username = name
	}
}

type UserOpt func(params *entity.UserParams)

func Users(t *testing.T) []*entity.User {
	t.Helper()

	return []*entity.User{
		User1(t),
		User2(t),
	}
}

func User1(t *testing.T, opts ...UserOpt) *entity.User {
	t.Helper()

	user1Params := &entity.UserParams{
		ID:       user1ID,
		Username: "test user 1",
		Email:    "user1@test.com",
	}

	for _, opt := range opts {
		opt(user1Params)
	}

	user, err := entity.NewUser(*user1Params)
	if err != nil {
		t.Fatalf("error creating user: %s", err.Error())
	}

	return user
}

func User2(t *testing.T, opts ...UserOpt) *entity.User {
	t.Helper()

	user2Params := &entity.UserParams{
		ID:       user2ID,
		Username: "test user 2",
		Email:    "user2@test.com",
	}

	for _, opt := range opts {
		opt(user2Params)
	}

	user, err := entity.NewUser(*user2Params)
	if err != nil {
		t.Fatalf("error creating user: %s", err.Error())
	}

	return user
}

func UsersEquals(t *testing.T, user1, user2 *entity.User) bool {
	t.Helper()

	return user1.ID() == user2.ID() && user1.Username() == user2.Username() && user1.Email() == user2.Email()
}
