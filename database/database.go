package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBSet() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal("Error al conectar a MongoDB:", err)
		return nil
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("No se pudo conectar a MongoDB:", err)
		return nil
	}

	log.Println("Conexión a MongoDB establecida exitosamente")
	return client
}

var Client *mongo.Client = DBSet()

func UserData(client *mongo.Client, collection_name string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("ecommerce").Collection(collection_name)
	return collection
}

func ProductData(client *mongo.Client, collection_name string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("ecommerce").Collection(collection_name)
	return collection
}
