package mongo

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/tclaudel/hexagonal-architecture-golang/internal/domain/entity"
	"github.com/tclaudel/hexagonal-architecture-golang/internal/domain/entity/fixtures"
	"github.com/tclaudel/hexagonal-architecture-golang/internal/ports/secondary"
)

const dbName = "database"

type TestSuite struct {
	suite.Suite
	mongoDB *UserRepository
}

func mongoURI(t *testing.T) string {
	const (
		MongoUriEnv     = "MONGO_URI"
		defaultMongoURI = "mongodb://localhost:27017"
	)

	value := os.Getenv(MongoUriEnv)
	if len(value) == 0 {
		return defaultMongoURI
	}

	return value
}

func (suite *TestSuite) SetupSuite() {
	mongoRepo, err := NewMongoRepository(context.Background(), mongoURI(suite.T()), dbName)
	if err != nil {
		suite.T().Fatalf("failed to create MongoDB: %v, mongo URI : %s", err, mongoURI(suite.T()))
	}

	suite.mongoDB = mongoRepo
}

func (suite *TestSuite) TearDownTest() {
	err := suite.mongoDB.coll.Drop(context.Background())
	if err != nil {
		suite.T().Fatalf("failed to drop collection: %v", err)
	}
}

func (suite *TestSuite) TearDownSubTest() {
	err := suite.mongoDB.coll.Drop(context.Background())
	if err != nil {
		suite.T().Fatalf("failed to drop collection: %v", err)
	}
}

func TestMongoRepository(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (suite *TestSuite) TestSaveUser() {
	tests := []struct {
		name          string
		users         []*entity.User
		expectedErrs  []error
		expectedUsers []*entity.User
	}{
		{
			name:          "should save user",
			users:         []*entity.User{fixtures.User1(suite.T())},
			expectedUsers: []*entity.User{fixtures.User1(suite.T())},
		},
		{
			name:          "should save multiple users",
			users:         []*entity.User{fixtures.User1(suite.T()), fixtures.User2(suite.T())},
			expectedUsers: []*entity.User{fixtures.User1(suite.T()), fixtures.User2(suite.T())},
		},
		{
			name:         "should return an error when saving a user with an existing ID",
			users:        []*entity.User{fixtures.User1(suite.T()), fixtures.User1(suite.T())},
			expectedErrs: []error{nil, secondary.ErrUserAlreadyExists},
		},
	}

	for _, tt := range tests {
		test := tt

		suite.T().Run(test.name, func(t *testing.T) {
			ctx := context.Background()

			for i := range test.users {

				err := suite.mongoDB.SaveUser(ctx, test.users[i])

				if test.expectedErrs != nil {
					assert.ErrorIs(suite.T(), err, test.expectedErrs[i])
				} else {
					assert.NoError(suite.T(), err)
					assert.Equal(suite.T(), test.expectedUsers[i].ID().String(), test.users[i].ID().String())
					assert.Equal(suite.T(), test.expectedUsers[i].Username(), test.users[i].Username())
					assert.Equal(suite.T(), test.expectedUsers[i].Email().String(), test.users[i].Email().String())
				}
			}

			suite.TearDownSubTest()
		})
	}
}

func (suite *TestSuite) TestUpdateUser() {
	tests := []struct {
		name          string
		users         []*entity.User
		setup         func()
		expectedErrs  []error
		expectedUsers []*entity.User
	}{
		{
			name: "should update user",
			setup: func() {
				err := suite.mongoDB.SaveUser(context.Background(), fixtures.User1(suite.T()))
				assert.NoError(suite.T(), err)
			},
			users:         []*entity.User{fixtures.User1(suite.T())},
			expectedUsers: []*entity.User{fixtures.User1(suite.T())},
		},
		{
			name: "should update multiple users",
			setup: func() {
				err := suite.mongoDB.SaveUser(context.Background(), fixtures.User1(suite.T()))
				assert.NoError(suite.T(), err)
				err = suite.mongoDB.SaveUser(context.Background(), fixtures.User2(suite.T()))
				assert.NoError(suite.T(), err)
			},
			users:         []*entity.User{fixtures.User1(suite.T()), fixtures.User2(suite.T())},
			expectedUsers: []*entity.User{fixtures.User1(suite.T()), fixtures.User2(suite.T())},
		},
		{
			name:  "should return an error when updating a user which does not exist",
			users: []*entity.User{fixtures.User1(suite.T())},
			expectedErrs: []error{
				secondary.ErrUserNotFound,
			},
		},
	}

	for _, tt := range tests {
		test := tt

		suite.T().Run(test.name, func(t *testing.T) {
			ctx := context.Background()

			if test.setup != nil {
				test.setup()
			}

			for i := range test.users {
				err := suite.mongoDB.UpdateUser(ctx, test.users[i])

				if test.expectedErrs != nil {
					assert.ErrorIs(suite.T(), err, test.expectedErrs[i])
				} else {
					assert.NoError(suite.T(), err)
					assert.Equal(suite.T(), test.expectedUsers[i].ID().String(), test.users[i].ID().String())
					assert.Equal(suite.T(), test.expectedUsers[i].Username(), test.users[i].Username())
					assert.Equal(suite.T(), test.expectedUsers[i].Email().String(), test.users[i].Email().String())
				}
			}

			suite.TearDownSubTest()
		})
	}
}
