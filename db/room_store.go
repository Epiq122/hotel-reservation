package db

import (
	"context"

	"github.com/epiq122/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const roomCollection = "rooms"

type RoomStore interface {
	CreateHotel(ctx context.Context, Room *types.Room) (*types.Room, error)
}

type MongoRoomStore struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoRoomStore(client *mongo.Client, dbname string) *MongoRoomStore {
	return &MongoRoomStore{
		client:     client,
		collection: client.Database(dbname).Collection(roomCollection),
	}
}

func (s *MongoRoomStore) CreateRoom(ctx context.Context, room *types.Room) (*types.Room, error) {
	res, err := s.collection.InsertOne(ctx, room)
	if err != nil {
		return nil, err
	}
	room.ID = res.InsertedID.(primitive.ObjectID)

	// update ther hotel with the new room id
	filter := bson.M{"_id": room.HotelID}
	update := bson.M{"$push": bson.M{"rooms": room.ID}}
	if err := s.client.Database("hotel").Collection("hotels").FindOneAndUpdate(ctx, filter, update).Err(); err != nil {
		return nil, err
	}
	return room, nil
}
