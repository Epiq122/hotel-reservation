package db

import (
	"context"
	"fmt"

	"github.com/epiq122/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const userCollection = "users"

type Dropper interface {
	Drop(context.Context) error
}

type UserStore interface {
	Dropper

	GetUserById(context.Context, string) (*types.User, error)
	GetUsers(context.Context) ([]*types.User, error)
	CreateUser(context.Context, *types.User) (*types.User, error)
	DeleteUser(context.Context, string) error
	UpdateUser(ctx context.Context, filter bson.M, params types.UpdateUserParams) error
	GetUserByEmail(context.Context, string) (*types.User, error)
}

type MongoUserStore struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoUserStore(client *mongo.Client) *MongoUserStore {
	return &MongoUserStore{
		client:     client,
		collection: client.Database(DBNAME).Collection(userCollection),
	}
}

func (s *MongoUserStore) Drop(ctx context.Context) error {
	fmt.Println("----- Dropping users collection -----")
	return s.collection.Drop(ctx)

}

func (s *MongoUserStore) GetUsers(ctx context.Context) ([]*types.User, error) {
	cur, err := s.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var users []*types.User
	if err := cur.All(ctx, &users); err != nil {
		return []*types.User{}, nil
	}
	return users, nil
}

func (s *MongoUserStore) GetUserById(ctx context.Context, id string) (*types.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var user types.User
	if err := s.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *MongoUserStore) CreateUser(ctx context.Context, user *types.User) (*types.User, error) {
	res, err := s.collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	user.ID = res.InsertedID.(primitive.ObjectID).Hex()
	return user, nil
}

func (s *MongoUserStore) DeleteUser(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	// TODO: Mabe check if we did not delete a user and log it
	_, err = s.collection.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}
	return nil
}

func (s *MongoUserStore) UpdateUser(ctx context.Context, filter bson.M, params types.UpdateUserParams) error {

	update := bson.D{{"$set", params.ToBSON()}}
	_, err := s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (s *MongoUserStore) GetUserByEmail(ctx context.Context, email string) (*types.User, error) {
	var user types.User
	if err := s.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}
