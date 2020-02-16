package api

import (
	"context"
	"fmt"
	"github.com/Creedowl/TgPusher/bot"
	"github.com/Creedowl/TgPusher/db"
	"github.com/Creedowl/TgPusher/utils"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"time"
)

func rootHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "welcome to TgPusher",
	})
}

func msgRootHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "Usage: TODO",
	})
}

func msgPostHandler(ctx *gin.Context) {
	token := ctx.Param("token")
	u, err := uuid.FromString(token)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "Wrong uuid",
		})
		return
	}
	var user db.User
	c, _ := context.WithTimeout(context.Background(), time.Second*10)
	collection := db.TgDatabase.DB.Collection("user")
	err = collection.FindOne(c, bson.D{{"token", u.String()}}).Decode(&user)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg": "User not found",
		})
		return
	}
	var message db.PushMessage
	if ctx.ShouldBind(&message) != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "Invalid message, please check your attributes",
		})
		return
	}
	message.ID = primitive.NewObjectID()
	message.UserID = user.ID
	message.Status = "pending"
	res, err := db.TgDatabase.DB.Collection("message").InsertOne(c, &message)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "Internal server error",
		})
		return
	}
	// push message to user
	go func() {
		m := tgbotapi.NewMessage(user.ChatID, "")
		msg := ""
		switch message.Type {
		case "markdown":
			m.ParseMode = "MarkdownV2"
			if message.Title != "" {
				msg = fmt.Sprintf("*%s*\n", message.Title)
			}
			msg += message.Content
		case "text":
			fallthrough
		default:
			if message.Title != "" {
				msg = fmt.Sprintf("Title: %s\nContent: ", message.Title)
			}
			msg += message.Content
		}
		m.Text = msg
		_, err := bot.PusherBot.Send(m)
		if err != nil {
			log.Println(err)
			_, err = db.TgDatabase.DB.Collection("message").UpdateOne(c, bson.D{{"_id", res.InsertedID}},
				bson.D{{"$set", bson.D{{"status", "failed"}}}})
			if err != nil {
				log.Println(err)
			}
			return
		}
		_, err = db.TgDatabase.DB.Collection("message").UpdateOne(c, bson.D{{"_id", res.InsertedID}},
			bson.D{{"$set", bson.D{{"status", "succeeded"}}}})
		if err != nil {
			log.Println(err)
		}
	}()
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "ok",
	})
}

func NewEngine() *gin.Engine {
	if utils.GetConfig().Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/", rootHandler)
	pushGroup := r.Group("/msg")
	{
		pushGroup.GET("/", msgRootHandler)
		pushGroup.POST("/:token", msgPostHandler)
	}

	return r
}
