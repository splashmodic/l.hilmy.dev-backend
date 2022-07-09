package url

import "time"

type createShortURLReqField struct {
	UserID   string `json:"userId" xml:"userId" form:"userId" validate:"required"`
	LongURL  string `json:"longUrl" xml:"longUrl" form:"longUrl" validate:"required,url"`
	ShortURL string `json:"shortUrl" xml:"shortUrl" form:"shortUrl" validate:"omitempty,startsnotwith=http,gte=3"`
}

type modifyShortURLReqField struct {
	UserID  string `json:"userId" xml:"userId" form:"userId" validate:"required"`
	LongURL string `json:"longUrl,omitempty" xml:"longUrl,omitempty" form:"longUrl,omitempty" validate:"omitempty,url"`
	IsShow  bool   `json:"isShow,omitempty" xml:"isShow,omitempty" form:"isShow,omitempty" validate:"omitempty"`
}

type moveShortURLsToNewUserIDReqField struct {
	OldUserID string `json:"oldUserId" xml:"oldUserId" form:"oldUserId" validate:"required"`
	NewUserID string `json:"newUserId" xml:"newUserId" form:"newUserId" validate:"required"`
}

type deleteShortURLReqField struct {
	UserID string `json:"userId" xml:"userId" form:"userId" validate:"required"`
}

type resField struct {
	Payload interface{} `json:"payload"`
}

type getShortURLsResField struct {
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	LongURL        string    `json:"longUrl"`
	ShortURL       string    `json:"shortUrl"`
	NumberAccessed uint      `json:"numberAccessed"`
}
