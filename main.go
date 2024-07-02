package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
)

var ctx = context.Background()

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func initializeRedisClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis server address
		Password: "password",
		DB:       0,
	})
	return rdb
}

func cachePersonData(rdb *redis.Client, person Person) {
	jsonData, err := json.Marshal(person)
	if err != nil {
		log.Fatalf("Error marshalling JSON: %v", err)
	}

	err = rdb.Set(ctx, "person:1", jsonData, 0).Err()
	if err != nil {
		log.Fatalf("Error setting data in Redis: %v", err)
	}
}

func getPersonData(rdb *redis.Client) (Person, error) {
	var person Person

	val, err := rdb.Get(ctx, "person:1").Result()
	if err != nil {
		return person, err
	}

	err = json.Unmarshal([]byte(val), &person)
	if err != nil {
		return person, err
	}

	return person, nil
}

func main() {
	rdb := initializeRedisClient()

	person := Person{
		Name: "John Doe",
		Age:  30,
	}

	cachePersonData(rdb, person)

	retrievedPerson, err := getPersonData(rdb)
	if err != nil {
		log.Fatalf("Error retrieving data from Redis: %v", err)
	}

	fmt.Printf("Retrieved Person: %+v\n", retrievedPerson)
}
