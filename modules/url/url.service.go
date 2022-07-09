package url

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"l.hilmy.dev/backend/helpers/errorhandler"
	"l.hilmy.dev/backend/helpers/random"
)

type service struct {
	code int
	data interface{}
}

func (u *url) createShortURLService(createShortUrlReqField *createShortURLReqField) *service {
	{
		isRandomShortURL := false
		if len(createShortUrlReqField.ShortURL) == 0 {
			isRandomShortURL = true
			createShortUrlReqField.ShortURL = random.String(5, &random.Options{WithLowerChar: true, WithUpperChar: true, WithNumber: true})
		}

		for {
			urlData, err := getShortURLEntityByShortURL(&createShortUrlReqField.ShortURL)
			if !isRandomShortURL {
				if err != nil && err != mongo.ErrNoDocuments {
					return &service{code: fiber.StatusInternalServerError, data: fiber.ErrInternalServerError.Error()}
				}
				if urlData != nil && urlData.CreatedByUserID != createShortUrlReqField.UserID {
					return &service{code: fiber.StatusBadRequest, data: "short url already exist"}
				}
				break
			} else {
				if err == mongo.ErrNoDocuments {
					break
				}
			}
			createShortUrlReqField.ShortURL += random.String(1, &random.Options{WithLowerChar: true, WithUpperChar: true, WithNumber: true})
		}
	}

	if err := createShortURLEntity(
		&createShortLinkEntityParam{
			longURL:  &createShortUrlReqField.LongURL,
			shortURL: &createShortUrlReqField.ShortURL,
			userID:   &createShortUrlReqField.UserID,
		},
	); err != nil {
		return &service{code: fiber.StatusInternalServerError, data: fiber.ErrInternalServerError.Message}
	}

	return &service{code: fiber.StatusCreated}
}

func (u *url) modifyShortURLService(shortURL *string, modifyShortUrlReqField *modifyShortURLReqField) *service {
	if len(*shortURL) < 3 {
		err := errors.New("short url length must be more than three characters")
		errorhandler.LogErrorThenContinue(&err)
		return &service{code: fiber.StatusBadRequest, data: err.Error()}
	}

	urlData, err := getShortURLEntityByShortURL(shortURL)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &service{code: fiber.StatusNotFound, data: fiber.ErrNotFound.Error()}
		}
		return &service{code: fiber.StatusInternalServerError, data: fiber.ErrInternalServerError.Error()}
	}

	if urlData.CreatedByUserID != modifyShortUrlReqField.UserID {
		err := errors.New(fiber.ErrForbidden.Error())
		errorhandler.LogErrorThenContinue(&err)
		return &service{code: fiber.StatusForbidden, data: err.Error()}
	}

	updatedAt := time.Now()
	numberAccessed := uint(0)
	if len(modifyShortUrlReqField.LongURL) > 0 {
		numberAccessed = urlData.NumberAccessed
	}

	if err := modifyShortURLEntityByShortURL(
		shortURL,
		urlData,
		&modifyShortURLEntityParam{
			updatedAt:      &updatedAt,
			longURL:        &modifyShortUrlReqField.LongURL,
			numberAccessed: &numberAccessed,
			isShow:         &modifyShortUrlReqField.IsShow,
		},
	); err != nil {
		return &service{code: fiber.StatusInternalServerError, data: err.Error()}
	}

	return &service{code: fiber.StatusCreated}
}

func (u *url) moveShortURLsToNewUserIDService(moveShortURLsToNewUserIDReqField *moveShortURLsToNewUserIDReqField) *service {
	if err := modifyUserIDOfShortURLsEntity(&moveShortURLsToNewUserIDReqField.OldUserID, &moveShortURLsToNewUserIDReqField.NewUserID); err != nil {
		if err == mongo.ErrNoDocuments {
			return &service{code: fiber.StatusNotFound, data: fiber.ErrNotFound.Error()}
		}
		return &service{code: fiber.StatusInternalServerError, data: fiber.ErrInternalServerError.Error()}
	}

	return &service{code: fiber.StatusOK}
}

func (u *url) getLongURLService(shortURL *string) *service {
	urlData, err := getShortURLEntityByShortURL(shortURL)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &service{code: fiber.StatusNotFound, data: fiber.ErrNotFound.Error()}
		}
		return &service{code: fiber.StatusInternalServerError, data: fiber.ErrInternalServerError.Error()}
	}

	go func() {
		updatedAt := time.Now()
		numberAccessed := urlData.NumberAccessed + 1
		modifyShortURLEntityByShortURL(shortURL, urlData, &modifyShortURLEntityParam{updatedAt: &updatedAt, numberAccessed: &numberAccessed})
	}()

	return &service{code: fiber.StatusOK, data: urlData.LongURL}
}

func (u *url) getShortURLsService(userID *string, isShow *bool) *service {
	filter := bson.E{}
	if *isShow {
		filter = bson.E{Key: "isShow", Value: *isShow}
	}
	urlsData, err := getShortURLsEntityByUserID(userID, filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &service{code: fiber.StatusNotFound, data: fiber.ErrNotFound.Error()}
		}
		return &service{code: fiber.StatusInternalServerError, data: fiber.ErrInternalServerError.Error()}
	}

	data := []getShortURLsResField{}

	for _, urlData := range *urlsData {
		data = append(
			data,
			getShortURLsResField{
				CreatedAt:      urlData.CreatedAt,
				UpdatedAt:      urlData.UpdatedAt,
				LongURL:        urlData.LongURL,
				ShortURL:       urlData.ShortURL,
				NumberAccessed: urlData.NumberAccessed,
			},
		)
	}

	return &service{code: fiber.StatusOK, data: &data}
}

func (u *url) deleteShortURLService(shortURL *string, deleteShortURLReqField *deleteShortURLReqField) *service {
	urlData, err := getShortURLEntityByShortURL(shortURL)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &service{code: fiber.StatusNotFound, data: fiber.ErrNotFound.Error()}
		}
		return &service{code: fiber.StatusInternalServerError, data: fiber.ErrInternalServerError.Error()}
	}

	if urlData.CreatedByUserID != deleteShortURLReqField.UserID {
		err := errors.New(fiber.ErrForbidden.Error())
		errorhandler.LogErrorThenContinue(&err)
		return &service{code: fiber.StatusForbidden, data: err.Error()}
	}

	if err := deleteShortURLEntityByShortURL(shortURL); err != nil {
		if err == mongo.ErrNoDocuments {
			return &service{code: fiber.StatusNotFound, data: fiber.ErrNotFound.Error()}
		}
		return &service{code: fiber.StatusInternalServerError, data: fiber.ErrInternalServerError.Error()}
	}
	return &service{code: fiber.StatusOK}
}
