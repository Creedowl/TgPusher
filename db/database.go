package db

import (
	"context"
	"fmt"
	"github.com/Creedowl/TgPusher/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

//var Client *mongo.Client

var TgDatabase *DataBase

type DataBase struct {
	Client *mongo.Client
	DB     *mongo.Database
}

func init() {
	log.Println("connecting to database...")
	config := utils.GetConfig()
	uri := fmt.Sprintf("mongodb://%s:%d", config.DB.Host, config.DB.Port)
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalln(err)
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatalln("can't connect to database, timeout")
	}
	TgDatabase = &DataBase{
		Client: client,
		DB:     client.Database(config.DB.Database),
	}
}
