package main

import (
	"context"
	"flag"
	"log"

	"github.com/epiq122/hotel-reservation/api"
	"github.com/epiq122/hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dburi = "mongodb://localhost:27017"

var config = fiber.Config{

	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.JSON(fiber.Map{
			"error": err.Error(),
		})
	},
}

func main() {
	listenAddr := flag.String("listen-addr", ":3000", "listen address for the server")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dburi))
	if err != nil {
		log.Fatal(err)
	}

	// handlers initialization
	userHandler := api.NewUserHandler(db.NewMongoUserStore(client))

	app := fiber.New(config)
	apiv1 := app.Group("/api/v1")

	apiv1.Post("/user", userHandler.HandleCreateUser)
	apiv1.Get("/users", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)

	app.Listen(*listenAddr)

}
