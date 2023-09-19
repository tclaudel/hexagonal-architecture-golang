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

func TestUserRetrieve_RetrieveUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		user         *entity.User
		userUseCase  func(ctrl *gomock.Controller) *mocks.MockUserRepository
		expectedUser *entity.User
		expectedErr  error
	}{
		{
			name: "should retrieve user",
			user: fixtures.User1(t),
			userUseCase: func(ctrl *gomock.Controller) *mocks.MockUserRepository {
				mock := mocks.NewMockUserRepository(ctrl)

				mock.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Return(fixtures.User1(t), nil)

				return mock
			},
			expectedUser: fixtures.User1(t),
			expectedErr:  nil,
		},
		{
			name: "should return error when user getter returns error",
			user: fixtures.User1(t),
			userUseCase: func(ctrl *gomock.Controller) *mocks.MockUserRepository {
				mock := mocks.NewMockUserRepository(ctrl)

				mock.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Return(nil, secondary.ErrRetrievingUser)

				return mock
			},
			expectedUser: nil,
			expectedErr:  primary.ErrRetrievingUser,
		},
		{
			name: "should return error when user getter returns user not found error",
			user: fixtures.User1(t),
			userUseCase: func(ctrl *gomock.Controller) *mocks.MockUserRepository {
				mock := mocks.NewMockUserRepository(ctrl)

				mock.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Return(nil, secondary.ErrUserNotFound)

				return mock
			},
			expectedUser: nil,
			expectedErr:  primary.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		test := tt

		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userUseCase := NewUserUseCase(test.userUseCase(ctrl))

			user, err := userUseCase.RetrieveUser(context.Background(), test.user.ID())
			if test.expectedUser != nil {
				assert.True(t, fixtures.UsersEquals(t, test.expectedUser, user))
			}

			assert.ErrorIs(t, err, test.expectedErr)
		})
	}
}
