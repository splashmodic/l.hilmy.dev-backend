package url

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"l.hilmy.dev/backend/db"
	"l.hilmy.dev/backend/helpers/errorhandler"
)

type urlEntity struct {
	ID              string    `bson:"_id"`
	CreatedAt       time.Time `bson:"createdAt"`
	UpdatedAt       time.Time `bson:"updatedAt"`
	CreatedByUserID string    `bson:"createdByUserId"`
	LongURL         string    `bson:"longUrl"`
	ShortURL        string    `bson:"shortUrl"`
	NumberAccessed  uint      `bson:"numberAccessed"`
	IsShow          bool      `bson:"isShow"`
}

type createShortLinkEntityParam struct {
	longURL  *string
	shortURL *string
	userID   *string
}

type modifyShortURLEntityParam struct {
	updatedAt      *time.Time
	longURL        *string
	numberAccessed *uint
	isShow         *bool
}

func createShortURLEntity(param *createShortLinkEntityParam) error {
	db := db.GetDB()
	coll := db.DB.Collection("url")

	docStruct := &urlEntity{
		ID:              *param.shortURL,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		CreatedByUserID: *param.userID,
		LongURL:         *param.longURL,
		ShortURL:        *param.shortURL,
		NumberAccessed:  0,
		IsShow:          true,
	}

	_, err := coll.ReplaceOne(context.TODO(), bson.M{"_id": docStruct.ID}, docStruct, options.Replace().SetUpsert(true))
	if err != nil {
		errorhandler.LogErrorThenContinue(&err)
		return err
	}

	return nil
}

func modifyShortURLEntityByShortURL(shortURL *string, data *urlEntity, param *modifyShortURLEntityParam) error {
	db := db.GetDB()
	coll := db.DB.Collection("url")

	docStruct := data

	if param.updatedAt != nil {
		docStruct.UpdatedAt = *param.updatedAt
	}
	if param.longURL != nil {
		docStruct.LongURL = *param.longURL
	}
	if param.numberAccessed != nil {
		docStruct.NumberAccessed = *param.numberAccessed
	}
	if param.isShow != nil {
		docStruct.IsShow = *param.isShow
	}

	_, err := coll.ReplaceOne(context.TODO(), bson.M{"_id": docStruct.ID}, docStruct)
	if err != nil {
		errorhandler.LogErrorThenContinue(&err)
		return err
	}

	return nil
}

func modifyUserIDOfShortURLsEntity(oldUserID, newUserID *string) error {
	db := db.GetDB()
	coll := db.DB.Collection("url")

	dataResult, err := coll.UpdateMany(
		context.TODO(),
		bson.M{"createdByUserId": oldUserID},
		bson.D{{Key: "$set", Value: bson.D{{Key: "createdByUserId", Value: newUserID}}}},
	)
	if err != nil {
		errorhandler.LogErrorThenContinue(&err)
		return err
	}

	if dataResult.MatchedCount == 0 {
		err := mongo.ErrNoDocuments
		errorhandler.LogErrorThenContinue(&err)
		return err
	}

	return nil
}

func getShortURLEntityByShortURL(shortURL *string) (*urlEntity, error) {
	db := db.GetDB()
	coll := db.DB.Collection("url")

	docStruct := &urlEntity{}

	if err := coll.FindOne(context.TODO(), bson.M{"_id": *shortURL}).Decode(docStruct); err != nil {
		errorhandler.LogErrorThenContinue(&err)
		return nil, err
	}

	return docStruct, nil
}

func getShortURLsEntityByUserID(userID *string, filter ...bson.E) (*[]urlEntity, error) {
	db := db.GetDB()
	coll := db.DB.Collection("url")

	docArrayOfStruct := &[]urlEntity{}

	filterFind := bson.D{{Key: "createdByUserId", Value: *userID}}
	if filter != nil {
		filterFind = append(filterFind, filter...)
	}

	urlsDataCursor, err := coll.Find(context.TODO(), filterFind)
	if err != nil {
		errorhandler.LogErrorThenContinue(&err)
		return nil, err
	}

	if err := urlsDataCursor.All(context.TODO(), docArrayOfStruct); err != nil {
		errorhandler.LogErrorThenContinue(&err)
		return nil, err
	}

	return docArrayOfStruct, nil
}

func deleteShortURLEntityByShortURL(shortURL *string) error {
	db := db.GetDB()
	coll := db.DB.Collection("url")

	dataResult, err := coll.DeleteOne(context.TODO(), bson.M{"_id": *shortURL})
	if err != nil {
		errorhandler.LogErrorThenContinue(&err)
		return err
	}
	if dataResult.DeletedCount == 0 {
		err := mongo.ErrNoDocuments
		errorhandler.LogErrorThenContinue(&err)
		return err
	}

	return nil
}
