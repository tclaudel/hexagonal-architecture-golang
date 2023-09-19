package services

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tclaudel/hexagonal-architecture-golang/internal/domain/entity/fixtures"
	"github.com/tclaudel/hexagonal-architecture-golang/internal/ports/primary"
	"github.com/tclaudel/hexagonal-architecture-golang/internal/ports/secondary"
	"github.com/tclaudel/hexagonal-architecture-golang/internal/ports/secondary/mocks"
	"go.uber.org/mock/gomock"
)

func TestUserCreate_CreateUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		userSave    func(ctrl *gomock.Controller) *mocks.MockUserRepository
		expectedErr error
	}{
		{
			name: "should create user",
			userSave: func(ctrl *gomock.Controller) *mocks.MockUserRepository {
				userSave := mocks.NewMockUserRepository(ctrl)
				userSave.EXPECT().
					SaveUser(gomock.Any(), gomock.Any()).
					Return(nil)
				return userSave
			},
			expectedErr: nil,
		},
		{
			name: "should return error when user saver returns error",
			userSave: func(ctrl *gomock.Controller) *mocks.MockUserRepository {
				userSave := mocks.NewMockUserRepository(ctrl)

				userSave.EXPECT().
					SaveUser(gomock.Any(), gomock.Any()).
					Return(secondary.ErrSavingUser)

				return userSave
			},
			expectedErr: primary.ErrCreatingUser,
		},
		{
			name: "should return error when user saver returns user already exists error",
			userSave: func(ctrl *gomock.Controller) *mocks.MockUserRepository {
				userSave := mocks.NewMockUserRepository(ctrl)

				userSave.EXPECT().
					SaveUser(gomock.Any(), gomock.Any()).
					Return(secondary.ErrUserAlreadyExists)

				return userSave
			},
			expectedErr: primary.ErrUserAlreadyExists,
		},
	}

	for _, tt := range tests {
		test := tt

		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userUseCase := NewUserUseCase(test.userSave(ctrl))

			err := userUseCase.CreateUser(context.Background(), fixtures.User1(t))
			assert.ErrorIs(t, err, test.expectedErr)
		})
	}
}
