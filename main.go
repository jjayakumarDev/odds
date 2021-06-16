package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// delay boolean flag
var delay bool

// retrieve sports data
func getsportsResponse() map[string]interface{} {
	apiKey := os.Getenv("apiKey")
	sportsResponse, err := http.Get("https://api.the-odds-api.com/v3/sports?api_key=" + apiKey + "")
	if err != nil {
		log.Fatal(err)
	}
	sportsResponseData, err := ioutil.ReadAll(sportsResponse.Body)
	if err != nil {
		log.Fatal(err)
	}
	var sportsData map[string]interface{}
	json.Unmarshal(sportsResponseData, &sportsData)
	return sportsData
}

// retrieve odds data
func getOddsResponse(delay bool) map[string]interface{} {
	if delay == true {
		//Delaying 5 minutes if delay flag is true
		time.Sleep(5 * time.Minute)
	}
	apiKey := os.Getenv("apiKey")
	sportKey := "upcoming"
	region := "uk"
	mkt := "h2h"
	oddsResponse, err := http.Get("https://api.the-odds-api.com/v3/odds?api_key=" + apiKey + "&sport=" + sportKey + "&region=" + region + "&mkt=" + mkt + "")
	if err != nil {
		log.Fatal(err)
	}
	oddsResponseData, err := ioutil.ReadAll(oddsResponse.Body)
	if err != nil {
		log.Fatal(err)
	}
	var oddsData map[string]interface{}
	json.Unmarshal(oddsResponseData, &oddsData)
	delay = true
	return oddsData
}

func main() {
	//Setting environment variable, can be removed if the variable aready been set
	os.Setenv("conenctionString", "mongodb+srv://user001:user001@cluster0.mu1mb.mongodb.net/myFirstDatabase?authSource=admin&replicaSet=atlas-14uvvk-shard-0&w=majority&readPreference=primary&retryWrites=true&ssl=true")
	os.Setenv("apiKey", "cc480abae2fce0ecf14b62c4866026fc")
	conenctionString := os.Getenv("conenctionString")

	client, err := mongo.NewClient(options.Client().ApplyURI(conenctionString))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	//Created database 'betme' in mongodb in cloud
	betmeDatabase := client.Database("betme")

	//Created collection for sports
	sports := betmeDatabase.Collection("sports")

	//Created collections for odds
	odds := betmeDatabase.Collection("odds")

	//key can be randomly generated number or we can define the struct and insert data without key
	sports.InsertOne(ctx, bson.D{
		{Key: "200521", Value: getsportsResponse()},
	})
	odds.InsertOne(ctx, bson.D{
		{Key: "200521", Value: getOddsResponse(delay)},
	})

	// Looping and inserting data once an hour
	for {
		time.Sleep(1 * time.Hour)
		odds.InsertOne(ctx, bson.D{
			{Key: "200521", Value: getOddsResponse(delay)},
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}
