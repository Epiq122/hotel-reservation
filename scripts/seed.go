package main

import (
	"context"
	"fmt"
	"log"

	"github.com/epiq122/hotel-reservation/db"
	"github.com/epiq122/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	hotelStore := db.NewMongoHotelStore(client, db.DBNAME)
	roomStore := db.NewMongoRoomStore(client, db.DBNAME)

	hotel := types.Hotel{
		Name:     "Hotel California",
		Location: "USA",
	}

	rooms := []types.Room{
		{
			Type:      types.Single,
			BasePrice: 100.00,
		},
		{
			Type:      types.Double,
			BasePrice: 200.00,
		},
		{
			Type:      types.Seaside,
			BasePrice: 300.00,
		},
	}
	insertedHotel, err := hotelStore.CreateHotel(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}

	for _, room := range rooms {
		room.HotelID = insertedHotel.ID
		insertedRoom, err := roomStore.CreateRoom(ctx, &room)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(insertedRoom)

	}

}
