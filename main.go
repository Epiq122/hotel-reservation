package main

import (
	"context"
	"flag"
	"log"

	"github.com/epiq122/hotel-reservation/api"
	"github.com/epiq122/hotel-reservation/api/middleware"
	"github.com/epiq122/hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	var (
		// handlers initialization
		hotelStore = db.NewMongoHotelStore(client)
		roomStore  = db.NewMongoRoomStore(client, hotelStore)
		userStore  = db.NewMongoUserStore(client)
		store      = &db.Store{
			Hotel: hotelStore,
			Room:  roomStore,
			User:  userStore,
		}
		userHandler  = api.NewUserHandler(userStore)
		hotelHandler = api.NewHotelHandler(store)
		authHandler  = api.NewAuthHandler(userStore)
		app          = fiber.New(config)
		auth         = app.Group("/api")
		apiv1        = app.Group("/api/v1", middleware.JWTAuthentication)
	)

	// auth routes
	auth.Post("/auth", authHandler.HandleAuthenticate)

	// Version API routes
	// users routes
	apiv1.Post("/user", userHandler.HandleCreateUser)
	apiv1.Put("/user/:id", userHandler.HandleUpdateUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiv1.Get("/users", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)

	// hotels routes
	apiv1.Get("/hotels", hotelHandler.HandleGetHotels)
	apiv1.Get("/hotel/:id", hotelHandler.HandleGetHotel)

	apiv1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)

	app.Listen(*listenAddr)

}
