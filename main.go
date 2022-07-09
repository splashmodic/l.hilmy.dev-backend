package main

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"l.hilmy.dev/backend/db"
	"l.hilmy.dev/backend/env"
	"l.hilmy.dev/backend/helpers/errorhandler"
	"l.hilmy.dev/backend/helpers/exithandler"
)

func main() {
	ex, _ := os.Executable()
	path := filepath.Dir(ex)
	{
		if len(os.Getenv("APP_MODE")) == 0 {
			if err := godotenv.Load("./.env"); err != nil {
				if err := godotenv.Load(path + "/.env"); err != nil {
					errorhandler.LogErrorThenPanic(&err)
				}
			}
		}
	}

	appName := env.Get(env.AppName).(string)
	appMode := env.Get(env.AppMode).(string)
	appAddr := env.Get(env.AppAddr).(string)
	appWebAddr := env.Get(env.AppWebAddr).(string)

	dbAddr := env.Get(env.DBAddr).(string)
	dbUser := env.Get(env.DBUser).(string)
	dbPwd := env.Get(env.DBPwd).(string)

	db.New(dbAddr, dbUser, dbPwd)
	disconnectDB := func() {
		log.Println("disconnecting database...")
		if err := db.GetClient().Disconnect(context.TODO()); err != nil {
			errorhandler.LogErrorThenContinue(&err)
		}
	}
	defer disconnectDB()
	exitHandler := exithandler.New(disconnectDB)
	db.Init()

	app := fiber.New(fiber.Config{
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
	})

	app.Use(recover.New(recover.Config{
		EnableStackTrace: func() bool {
			return appMode != "RELEASE"
		}(),
	}))

	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	app.Use(cors.New(cors.Config{
		AllowOrigins: func() string {
			if appMode == "RELEASE" && len(appWebAddr) > 0 {
				return appWebAddr
			}
			return "*"
		}(),
	}))

	if appMode != "RELEASE" {
		file, err := os.OpenFile(path+"/app.log", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			file.Close()
			errorhandler.LogErrorThenPanic(&err)
		}
		close := func() {
			log.Println("closing file...")
			if err := file.Close(); err != nil {
				errorhandler.LogErrorThenContinue(&err)
			}
		}
		defer close()
		exitHandler = exitHandler.Add(close)

		app.Use(logger.New(logger.Config{
			Format: "[${time}] ${ip}:${port} ${status} - ${latency} ${method} ${path}\n",
			Output: file,
		}))
	}

	exitHandler.Watch()

	module := module{appName: &appName, app: app}
	module.run()

	if err := app.Listen(appAddr); err != nil {
		errorhandler.LogErrorThenPanic(&err)
	}
}
