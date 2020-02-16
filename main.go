package main

import (
	"github.com/Creedowl/TgPusher/api"
	"github.com/Creedowl/TgPusher/bot"
	"log"
)

func main() {
	go bot.StartBot()
	r := api.NewEngine()
	err := r.Run()
	if err == nil {
		log.Fatalln("can't start server")
	}
	//fmt.Println(utils.GetConfig().DB.Host)
	//client := db.Client
	//collection := client.Database("test").Collection("naive")
	//ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	//cur, _ := collection.Find(ctx, bson.D{})
	//defer cur.Close(ctx)
	//
	//for cur.Next(ctx) {
	//	var result bson.M
	//	err := cur.Decode(&result)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	log.Println(result)
	//}
}
