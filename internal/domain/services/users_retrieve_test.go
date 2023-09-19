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

func TestUsersRetrieve_UsersRetrieve(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		userUseCase   func(ctrl *gomock.Controller) *mocks.MockUserRepository
		expectedUsers []*entity.User
		expectedErr   error
	}{
		{
			name: "should retrieve users",
			userUseCase: func(ctrl *gomock.Controller) *mocks.MockUserRepository {
				mock := mocks.NewMockUserRepository(ctrl)

				mock.EXPECT().
					GetUsers(gomock.Any()).
					Return(fixtures.Users(t), nil)

				return mock
			},
			expectedUsers: fixtures.Users(t),
			expectedErr:   nil,
		},
		{
			name: "should return error when users getter returns error",
			userUseCase: func(ctrl *gomock.Controller) *mocks.MockUserRepository {
				mock := mocks.NewMockUserRepository(ctrl)

				mock.EXPECT().
					GetUsers(gomock.Any()).
					Return(nil, secondary.ErrRetrievingUsers)

				return mock
			},
			expectedUsers: nil,
			expectedErr:   primary.ErrRetrievingUsers,
		},
	}

	for _, tt := range tests {
		test := tt

		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userUseCase := NewUserUseCase(tt.userUseCase(ctrl))

			users, err := userUseCase.RetrieveUsers(context.Background())
			if test.expectedUsers != nil {
				for i, user := range users {
					assert.True(t, fixtures.UsersEquals(t, test.expectedUsers[i], user))
				}
			}

			assert.ErrorIs(t, err, test.expectedErr)
		})
	}
}
