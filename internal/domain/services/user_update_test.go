package services

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tclaudel/hexagonal-architecture-golang/internal/domain/entity"
	"github.com/tclaudel/hexagonal-architecture-golang/internal/domain/entity/fixtures"
	"github.com/tclaudel/hexagonal-architecture-golang/internal/ports/primary"
	"github.com/tclaudel/hexagonal-architecture-golang/internal/ports/secondary"
	"github.com/tclaudel/hexagonal-architecture-golang/internal/ports/secondary/mocks"
	"go.uber.org/mock/gomock"
)

func TestUserUpdate_UpdateUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		user        *entity.User
		userUseCase func(ctrl *gomock.Controller) *mocks.MockUserRepository
		expectedErr error
	}{
		{
			name: "should update user",
			user: fixtures.User1(t),
			userUseCase: func(ctrl *gomock.Controller) *mocks.MockUserRepository {
				userUpdate := mocks.NewMockUserRepository(ctrl)

				userUpdate.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Return(nil)

				return userUpdate
			},
			expectedErr: nil,
		},
		{
			name: "should return error when user updater returns error",
			user: fixtures.User1(t),
			userUseCase: func(ctrl *gomock.Controller) *mocks.MockUserRepository {
				userUpdate := mocks.NewMockUserRepository(ctrl)

				userUpdate.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Return(secondary.ErrUpdatingUser)

				return userUpdate
			},
			expectedErr: primary.ErrUpdatingUser,
		},
		{
			name: "should return error when user updater returns user not found error",
			user: fixtures.User1(t),
			userUseCase: func(ctrl *gomock.Controller) *mocks.MockUserRepository {
				userUpdate := mocks.NewMockUserRepository(ctrl)

				userUpdate.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Return(secondary.ErrUserNotFound)

				return userUpdate
			},
			expectedErr: primary.ErrUserNotFound,
		},
	}

	for _, test := range tests {
		tt := test

		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userUseCase := NewUserUseCase(tt.userUseCase(ctrl))

			err := userUseCase.UpdateUser(context.Background(), tt.user)
			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
			}
		})
	}
}
