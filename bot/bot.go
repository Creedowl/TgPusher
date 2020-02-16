package bot

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Creedowl/TgPusher/db"
	"github.com/Creedowl/TgPusher/utils"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var PusherBot *tgbotapi.BotAPI

func init() {
	token := utils.GetConfig().BotToken
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalln(err)
	}
	bot.Debug = utils.GetConfig().Debug
	PusherBot = bot
}

func StartBot() {
	bot := PusherBot
	log.Println("starting TgPusherBot...")
	Db := db.TgDatabase.DB
	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	var (
		lastUpdate   db.Message
		lastUpdateId int
	)
	err := Db.Collection("bot").FindOne(ctx, bson.D{}, options.FindOne().
		SetSort(bson.D{{"_id", -1}})).Decode(&lastUpdate)
	if err != nil {
		log.Println(err)
		lastUpdateId = 0
	} else {
		lastUpdateId = lastUpdate.UpdateId
	}
	u := tgbotapi.NewUpdate(lastUpdateId + 1)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}
	log.Println("bot listening")
	for update := range updates {
		ctx, _ = context.WithTimeout(context.Background(), time.Second*10)
		go solveUpdate(ctx, &update, bot, Db)
	}
}

func solveUpdate(ctx context.Context, update *tgbotapi.Update, bot *tgbotapi.BotAPI, Db *mongo.Database) {
	message := update.Message
	if message == nil {
		return
	}
	var user db.User
	err := Db.Collection("user").FindOne(ctx, bson.D{{"user_id", message.From.ID}}).
		Decode(&user)
	if err == mongo.ErrNoDocuments {
		user = db.User{
			ID:     primitive.NewObjectID(),
			UserID: message.From.ID,
			ChatID: message.Chat.ID,
			Name:   message.From.UserName,
			Token:  "",
		}
		_, err := Db.Collection("user").InsertOne(ctx, &user)
		if err != nil {
			log.Println(err)
		}
	}
	userID := user.ID
	origin, _ := json.Marshal(message)
	save := db.Message{
		ID:         primitive.NewObjectID(),
		UserID:     userID,
		MessageID:  message.MessageID,
		UpdateId:   update.UpdateID,
		FromUserId: message.From.ID,
		Date:       message.Date,
		ChatID:     message.Chat.ID,
		Text:       message.Text,
		Origin:     string(origin),
	}
	_, err = Db.Collection("bot").InsertOne(ctx, &save)
	if err != nil {
		log.Println(err)
	}
	//log.Println(res.InsertedID)
	//log.Println(update.UpdateID)
	msg := tgbotapi.NewMessage(message.Chat.ID, "🌸Q")
	if message.IsCommand() {
		switch message.Command() {
		case "start":
			msg.Text = "Welcome to TgPusher bot"
		case "help":
			msg.ParseMode = "MarkdownV2"
			msg.Text = utils.GetHelp("zh-hans")
		case "token":
			token := user.Token
			u, err := uuid.FromString(token)
			if err != nil {
				u = uuid.NewV4()
				res, err := Db.Collection("user").UpdateOne(ctx, bson.D{{"_id", user.ID}},
					bson.D{{"$set", bson.D{{"token", u.String()}}}})
				if err != nil {
					log.Println(err)
				}
				log.Println(res)
			}
			msg.ParseMode = "MarkdownV2"
			msg.Text = fmt.Sprintf("your token is: `%s`", u.String())
		case "revoke":
			u := uuid.NewV4()
			res, err := Db.Collection("user").UpdateOne(ctx, bson.D{{"_id", user.ID}},
				bson.D{{"$set", bson.D{{"token", u.String()}}}})
			if err != nil {
				log.Println(err)
			}
			log.Println(res)
			msg.ParseMode = "MarkdownV2"
			msg.Text = fmt.Sprintf("your new token is: `%s`", u.String())
		default:
			msg.Text = "unknown command!"
		}
	}
	_, err = bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}
