package tests

import (
	"context"
	"gotify/internal/cache"
	"log"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestWithRedis(t *testing.T) {
	log.Println("Start test now")

	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "redis:latest",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
	}
	redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatalf("Could not start redis: %s", err)
	}
	defer func() {
		if err := redisC.Terminate(ctx); err != nil {
			log.Fatalf("Could not stop redis: %s", err)
		}
	}()

	endpoint, err := redisC.Endpoint(ctx, "")
	if err != nil {
		t.Error(err)
	}

	cache := cache.New(endpoint)

	type Data struct {
		Name string
		Age  int
	}

	data := Data{
		Name: "Alice",
		Age:  30,
	}

	cache.SetStruct("testkey", data)

	otherData := Data{}
	cache.GetStruct("testkey", &otherData)

	if data == otherData {
		log.Println("Sucess!")
	} else {
		log.Printf("In: %+v \n", data)
		log.Printf("Out: %+v \n", otherData)
	}
}
