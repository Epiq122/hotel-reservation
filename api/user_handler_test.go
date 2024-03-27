package api

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http/httptest"
	"testing"

	"github.com/epiq122/hotel-reservation/db"
	"github.com/epiq122/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	testdbUri = "mongodb://localhost:27017"
	dbname    = "hotel-reservation-test"
)

type testdb struct {
	db.UserStore
}

func (t *testdb) teardown(*testing.T) {
	if err := t.UserStore.Drop(context.TODO()); err != nil {
		log.Fatal(err)

	}

}

func setup(t *testing.T) *testdb {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(testdbUri))
	if err != nil {
		log.Fatal(err)
	}
	return &testdb{db.NewMongoUserStore(client, dbname)}
}

func TestCreateUser(t *testing.T) {
	testdb := setup(t)
	defer testdb.teardown(t)

	app := fiber.New()
	userHandler := NewUserHandler(testdb.UserStore)
	app.Post("/", userHandler.HandleCreateUser)

	params := types.CreateUserParams{
		Email:     "epiqtest@test.com",
		FirstName: "Epiq",
		LastName:  "Test",
		Password:  "password123",
	}

	b, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}

	var user types.User
	json.NewDecoder(resp.Body).Decode(&user)

	if len(user.ID) == 0 {
		t.Error("User ID is empty")
	}
	if len(user.EncryptedPassword) > 0 {
		t.Error("Encrypted password is not to be included in the response")
	}
	if user.Email != params.Email || user.FirstName != params.FirstName || user.LastName != params.LastName {
		t.Errorf("User does not match expected. Got: %#v, Expected Email: %s, First Name: %s, Last Name: %s", user, params.Email, params.FirstName, params.LastName)
	}

}
