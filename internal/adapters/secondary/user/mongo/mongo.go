package mongo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tclaudel/hexagonal-architecture-golang/internal/domain/entity"
	"github.com/tclaudel/hexagonal-architecture-golang/internal/ports/secondary"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/multierr"
)

const (
	collName       = "users"
	connectTimeout = time.Second
)

var (
	ErrConnectingMongoDatabase = errors.New("error connecting to mongo database")
)

var _ secondary.UserRepository = (*UserRepository)(nil)

type UserRepository struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoRepository(ctx context.Context, dbURI, dbName string) (*UserRepository, error) {
	opt := options.Client()

	opt.ApplyURI(dbURI)
	opt.SetConnectTimeout(connectTimeout)

	client, err := mongo.Connect(ctx, opt)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrConnectingMongoDatabase, err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrConnectingMongoDatabase, err)
	}

	err = client.Database(dbName).CreateCollection(ctx, collName)

	return &UserRepository{
		client: client,
		coll:   client.Database(dbName).Collection(collName),
	}, nil
}

func (u UserRepository) SaveUser(ctx context.Context, user *entity.User) error {
	userDTO := UserEntityToDTO(user)

	_, err := u.coll.InsertOne(ctx, userDTO)
	if err != nil {
		var writeErr mongo.WriteException

		ok := errors.As(err, &writeErr)
		if ok {
			var saveErr error

			for _, writeErr := range writeErr.WriteErrors {
				if writeErr.Code == 11000 {
					saveErr = multierr.Append(saveErr, fmt.Errorf("%w: %s", secondary.ErrUserAlreadyExists, err.Error()))
				}

				saveErr = multierr.Append(saveErr, fmt.Errorf("%w: %s", secondary.ErrSavingUser, err.Error()))
			}

			return saveErr
		}

		return fmt.Errorf("%w: %s", secondary.ErrSavingUser, err.Error())
	}

	return nil
}

func (u UserRepository) UpdateUser(ctx context.Context, user *entity.User) error {
	userDTO := UserEntityToDTO(user)

	res, err := u.coll.UpdateOne(ctx,
		bson.M{"_id": userDTO.ID}, bson.M{"$set": userDTO}, options.Update().SetUpsert(false))
	if err != nil {
		return fmt.Errorf("%w: %s", secondary.ErrUpdatingUser, err.Error())
	}

	if res.MatchedCount == 0 {
		return secondary.ErrUserNotFound
	}

	return nil
}

func (u UserRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	res, err := u.coll.DeleteOne(ctx, bson.M{"_id": id.String()})
	if err != nil {
		return fmt.Errorf("%w: %s", secondary.ErrDeletingUser, err.Error())
	}

	if res.DeletedCount == 0 {
		return secondary.ErrUserNotFound
	}

	return nil
}

func (u UserRepository) GetUser(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	bytes, err := u.coll.FindOne(ctx, bson.M{"_id": id.String()}).DecodeBytes()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("%w: %s", secondary.ErrUserNotFound, err.Error())
		}

		return nil, fmt.Errorf("%w: %s", secondary.ErrRetrievingUser, err.Error())
	}

	var userDTO User

	err = bson.Unmarshal(bytes, &userDTO)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", secondary.ErrRetrievingUser, err.Error())
	}

	user, err := UserDTOToEntity(userDTO)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", secondary.ErrRetrievingUser, err.Error())
	}

	return user, nil
}

func (u UserRepository) GetUsers(ctx context.Context) ([]*entity.User, error) {
	cursor, err := u.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("%w: %s", secondary.ErrRetrievingUsers, err.Error())
	}

	var usersDTO []User

	err = cursor.All(ctx, &usersDTO)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", secondary.ErrRetrievingUsers, err.Error())
	}

	users := make([]*entity.User, len(usersDTO))
	for i := range usersDTO {
		var err error

		users[i], err = UserDTOToEntity(usersDTO[i])
		if err != nil {
			return nil, fmt.Errorf("%w: %s", secondary.ErrRetrievingUsers, err.Error())
		}
	}

	return users, nil
}
