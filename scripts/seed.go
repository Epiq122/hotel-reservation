package main

import (
	"context"
	"log"

	"github.com/epiq122/hotel-reservation/db"
	"github.com/epiq122/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client     *mongo.Client
	hotelStore db.HotelStore
	roomStore  db.RoomStore
	ctx        = context.Background()
)

func seedHotel(name string, location string, rating int) {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}

	rooms := []types.Room{
		{
			Size:  "small",
			Price: 100.00,
		},
		{
			Size:  "normal",
			Price: 200.00,
		},
		{
			Size:  "kingsize",
			Price: 300.00,
		},
	}

	insertedHotel, err := hotelStore.CreateHotel(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}

	for _, room := range rooms {
		room.HotelID = insertedHotel.ID
		_, err := roomStore.CreateRoom(ctx, &room)
		if err != nil {
			log.Fatal(err)
		}

	}

}

func main() {
	seedHotel("Hilton", "New York", 5)
	seedHotel("Marriott", "San Francisco", 2)
	seedHotel("Sheraton", "Los Angeles", 3)

}

func init() {
	var err error
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore) // Pass hotelStore
}
