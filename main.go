package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/epiq122/hotel-reservation/api"
	"github.com/epiq122/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dburi = "mongodb://localhost:27017"
const dbname = "hotel-reservation"
const userCollection = "users"

func main() {

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dburi))
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	coll := client.Database(dbname).Collection(userCollection)
	user := types.User{
		FirstName: "Rob",
		LastName:  "Admin",
	}

	_, err = coll.InsertOne(ctx, user)
	if err != nil {
		log.Fatal(err)
	}

	var rob types.User
	if err := coll.FindOne(ctx, bson.M{}).Decode(&rob); err != nil {
		log.Fatal(err)
	}
	fmt.Println(rob)

	listenAddr := flag.String("listen-addr", ":3000", "server listen address")
	flag.Parse()

	app := fiber.New()
	apiv1 := app.Group("/api/v1")

	apiv1.Get("/user", api.HandleGetUsers)
	apiv1.Get("/user/:id", api.HandleGetUser)

	app.Listen(*listenAddr)

}
