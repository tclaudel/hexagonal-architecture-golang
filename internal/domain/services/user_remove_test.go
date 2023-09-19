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

func TestUserRemove_RemoveUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		user        *entity.User
		userUseCase func(ctrl *gomock.Controller) *mocks.MockUserRepository
		expectedErr error
	}{
		{
			name: "should remove user",
			user: fixtures.User1(t),
			userUseCase: func(ctrl *gomock.Controller) *mocks.MockUserRepository {
				userDelete := mocks.NewMockUserRepository(ctrl)

				userDelete.EXPECT().
					DeleteUser(gomock.Any(), gomock.Any()).
					Return(nil)

				return userDelete
			},
			expectedErr: nil,
		},
		{
			name: "should return error when user deleter returns error",
			user: fixtures.User1(t),
			userUseCase: func(ctrl *gomock.Controller) *mocks.MockUserRepository {
				userDelete := mocks.NewMockUserRepository(ctrl)

				userDelete.EXPECT().
					DeleteUser(gomock.Any(), gomock.Any()).
					Return(secondary.ErrDeletingUser)

				return userDelete
			},
			expectedErr: primary.ErrRemovingUser,
		},
	}

	for _, tt := range tests {
		test := tt

		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userUseCase := NewUserUseCase(test.userUseCase(ctrl))

			err := userUseCase.RemoveUser(context.Background(), test.user.ID())
			assert.Equal(t, test.expectedErr, err)
		})
	}
}
