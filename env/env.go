package env

import (
	"os"
)

const (
	AppName    = iota
	AppMode    = iota
	AppAddr    = iota
	AppWebAddr = iota

	DBAddr = iota
	DBUser = iota
	DBPwd  = iota

	MailSMTPHost     = iota
	MailSMTPPort     = iota
	MailSenderName   = iota
	MailSenderMail   = iota
	MailAuthEmail    = iota
	MailAuthPassword = iota
)

func Get(enum int) interface{} {
	switch enum {
	case AppName:
		return os.Getenv("APP_NAME")
	case AppMode:
		mode := os.Getenv("APP_MODE")
		if mode != "RELEASE" {
			mode = "DEBUG"
		}
		return mode
	case AppAddr:
		return os.Getenv("APP_ADDR")
	case AppWebAddr:
		return os.Getenv("APP_WEB_ADDR")

	case DBAddr:
		return os.Getenv("MONGODB_ADDR")
	case DBUser:
		return os.Getenv("MONGO_INITDB_ROOT_USERNAME")
	case DBPwd:
		return os.Getenv("MONGO_INITDB_ROOT_PASSWORD")
	case MailSMTPHost:
		return os.Getenv("MAIL_SMTP_HOST")
	case MailSMTPPort:
		return os.Getenv("MAIL_SMTP_PORT")
	case MailSenderName:
		return os.Getenv("MAIL_SENDER_NAME")
	case MailSenderMail:
		return os.Getenv("MAIL_SENDER_EMAIL")
	case MailAuthEmail:
		return os.Getenv("MAIL_AUTH_EMAIL")
	case MailAuthPassword:
		return os.Getenv("MAIL_AUTH_PASSWORD")
	}

	panic("unknown env variable")
}
