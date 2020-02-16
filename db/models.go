package db

import "go.mongodb.org/mongo-driver/bson/primitive"

type Message struct {
	ID         primitive.ObjectID `bson:"_id"`
	UserID     primitive.ObjectID `bson:"user_id"`
	MessageID  int                `bson:"message_id"`
	UpdateId   int                `bson:"update_id"`
	FromUserId int                `bson:"from_user_id"`
	Date       int                `bson:"date"`
	ChatID     int64              `bson:"chat_id"`
	Text       string             `bson:"text"`
	Origin     string             `bson:"origin"`
}

type User struct {
	ID     primitive.ObjectID `bson:"_id"`
	UserID int                `bson:"user_id"`
	ChatID int64              `bson:"chat_id"`
	Name   string             `bson:"name"`
	Token  string             `bson:"token"`
}

type PushMessage struct {
	ID      primitive.ObjectID `bson:"_id"`
	UserID  primitive.ObjectID `bson:"user_id"`
	Title   string             `bson:"title" json:"title" form:"title"`
	Type    string             `bson:"type" json:"type" form:"type"`
	Content string             `bson:"content" json:"content" form:"content" binding:"required"`
	Status  string             `bson:"status"`
}
