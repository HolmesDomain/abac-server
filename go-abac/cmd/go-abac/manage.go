package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Request string

func Log(msg string) string {
	return time.Now().Format("2006-01-02 15:04:05") + msg
}

func postPolicy(pType string, V0 string, V1 string, V2 string) bool {
	clientOptions := options.Client().ApplyURI(Mongo)
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB")
	collection := client.Database("PolicyDatabase").Collection("casbin_rule")

	post := CasbinModel{pType, V0, V1, V2}
	insertResult, err := collection.InsertOne(context.TODO(), post)

	if err != nil {
		log.Fatal(err)
		return false
	} else {
		fmt.Println("Inserted post with ID:", insertResult.InsertedID)
		return true
	}
}

func removePolicy(id string) bool {
	// Connection check
	clientOptions := options.Client().ApplyURI(Mongo)
	client, err := mongo.NewClient(clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB")

	// Get collection and findOne example of object
	collection := client.Database("PolicyDatabase").Collection("casbin_rule")

	// convert hex mongodb _id string
	var objID, _ = primitive.ObjectIDFromHex(id)

	result, err := collection.DeleteOne(ctx, bson.M{"_id": objID})

	if err != nil {
		log.Fatal(err)
		return false
	} else {
		fmt.Printf("DeleteOne removed %v document(s)\n", result.DeletedCount)
		return true
	}
}

func getSinglePolicy(sub, obj, act string) []primitive.M {
	clientOptions := options.Client().ApplyURI(Mongo)
	client, err := mongo.NewClient(clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB")

	collection := client.Database("PolicyDatabase").Collection("casbin_rule")

	cursor, err := collection.Find(
		context.Background(),
		bson.D{{"v0", "r.sub.Label == " + "'" + sub + "'"},
			{"v1", "/" + obj},
			{"v2", act}})

	var policyQuery []bson.M

	if err = cursor.All(ctx, &policyQuery); err != nil {
		log.Fatal(err)
	}

	return policyQuery
}

func getMatchingPolicy(sub string) []primitive.M {
	clientOptions := options.Client().ApplyURI(Mongo)
	client, err := mongo.NewClient(clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB")

	collection := client.Database("PolicyDatabase").Collection("casbin_rule")

	cursor, err := collection.Find(context.Background(), bson.D{{"v0", "r.sub.Label == " + "'" + sub + "'"}})

	var policyQuery []bson.M

	if err = cursor.All(ctx, &policyQuery); err != nil {
		log.Fatal(err)
	}

	return policyQuery
}

func getThePolicyV2() []primitive.M {
	clientOptions := options.Client().ApplyURI(Mongo)
	client, err := mongo.NewClient(clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB")

	collection := client.Database("PolicyDatabase").Collection("casbin_rule")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	var policies []bson.M

	if err = cursor.All(ctx, &policies); err != nil {
		log.Fatal(err)
	}

	return policies
}

func getAccessLog(user, resource, action string, result bool) string {
	Request = " [Request: " + user + " " + resource + " " + action + "] " + fmt.Sprint(result)

	return Log(Request)
}

func getLog() string {
	return Log("")
}

type CasbinModeler struct {
	PType string `json:PType,omitempty”`
	V0    string `json:V0,omitempty”`
	V1    string `json:V1,omitempty”`
	V2    string `json:V2,omitempty”`
}
