package db

import (
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"l.hilmy.dev/backend/helpers/errorhandler"
)

type database struct {
	DB *mongo.Database
}

var db *database
var dbClient *mongo.Client

func New(dbAddr, dbUser, dbPwd string) {
	if db != nil && dbClient != nil {
		log.Println("database and it's client is exist")
		return
	}

	log.Printf("connecting to database...")
	dbUri := "mongodb+srv://" + dbUser + ":" + dbPwd + "@" + dbAddr

	dbclient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dbUri))
	if err != nil {
		errorhandler.LogErrorThenPanic(&err)
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	if err := dbclient.Ping(ctx, readpref.Primary()); err != nil {
		errorhandler.LogErrorThenPanic(&err)
	}

	dbClient = dbclient
	db = &database{DB: dbClient.Database("auth")}
}

func GetDB() *database {
	if db == nil {
		err := errors.New("database not found")
		errorhandler.LogErrorThenPanic(&err)
	}

	return db
}

func GetClient() *mongo.Client {
	if dbClient == nil {
		err := errors.New("database client not found")
		errorhandler.LogErrorThenPanic(&err)
	}

	return dbClient
}

func Init() {
	log.Println("initializing database...")

	{
		expireAfter := 3 * 30 * 24 * (time.Hour / time.Second) // 3 months
		coll := db.DB.Collection("url")
		_, err := coll.Indexes().CreateOne(
			context.TODO(),
			mongo.IndexModel{
				Keys:    bson.D{{Key: "updatedAt", Value: 1}},
				Options: options.Index().SetExpireAfterSeconds(int32(expireAfter)),
			},
		)
		if err != nil {
			errorhandler.LogErrorThenPanic(&err)
		}
	}
}
