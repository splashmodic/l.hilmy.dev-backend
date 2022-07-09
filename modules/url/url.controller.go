package url

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"l.hilmy.dev/backend/helpers/errorhandler"
	"l.hilmy.dev/backend/helpers/validator"
)

func (u *url) createShortURLController() {
	u.router.Post("", func(c *fiber.Ctx) error {
		req := new(createShortURLReqField)

		if err := c.BodyParser(req); err != nil {
			errorhandler.LogErrorThenContinue(&err)
			return c.Status(fiber.StatusBadRequest).JSON(&resField{Payload: fiber.ErrBadRequest.Error()})
		}

		if err := validator.Struct(req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(&resField{Payload: fiber.ErrBadRequest.Error()})
		}

		service := u.createShortURLService(req)

		return c.Status(service.code).JSON(&resField{Payload: service.data})
	})
}

func (u *url) modifyShortURLController() {
	u.router.Patch(":shortUrl", func(c *fiber.Ctx) error {
		shortUrl := c.Params("shortUrl")

		if len(shortUrl) == 0 {
			err := errors.New("blank short url")
			errorhandler.LogErrorThenContinue(&err)
			return c.Status(fiber.StatusBadRequest).JSON(&resField{Payload: err.Error()})
		}

		req := new(modifyShortURLReqField)

		if err := c.BodyParser(req); err != nil {
			errorhandler.LogErrorThenContinue(&err)
			return c.Status(fiber.StatusBadRequest).JSON(&resField{Payload: fiber.ErrBadRequest.Error()})
		}

		if err := validator.Struct(req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(&resField{Payload: fiber.ErrBadRequest.Error()})
		}

		service := u.modifyShortURLService(&shortUrl, req)

		return c.Status(service.code).JSON(&resField{Payload: service.data})
	})
}

func (u *url) moveShortURLsToNewUserIDController() {
	u.router.Patch("m/move", func(c *fiber.Ctx) error {
		req := new(moveShortURLsToNewUserIDReqField)

		if err := c.BodyParser(req); err != nil {
			errorhandler.LogErrorThenContinue(&err)
			return c.Status(fiber.StatusBadRequest).JSON(&resField{Payload: fiber.ErrBadRequest.Error()})
		}

		if err := validator.Struct(req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(&resField{Payload: fiber.ErrBadRequest.Error()})
		}

		service := u.moveShortURLsToNewUserIDService(req)

		return c.Status(service.code).JSON(&resField{Payload: service.data})
	})
}

func (u *url) getLongURLController() {
	u.router.Get(":shortUrl", func(c *fiber.Ctx) error {
		shortURL := c.Params("shortUrl")

		if len(shortURL) == 0 {
			err := errors.New("short url can't be empty")
			errorhandler.LogErrorThenContinue(&err)
			return c.Status(fiber.StatusBadRequest).JSON(&resField{Payload: err.Error()})
		}

		service := u.getLongURLService(&shortURL)

		return c.Status(service.code).JSON(&resField{Payload: service.data})
	})
}

func (u *url) getShortURLsController() {
	u.router.Get("u/:userId", func(c *fiber.Ctx) error {
		userID := c.Params(("userId"))
		isShow, err := strconv.ParseBool(c.Query("isShow"))
		if err != nil {
			errorhandler.LogErrorThenContinue(&err)
			isShow = true
		}

		if len(userID) == 0 {
			err := errors.New("user id can't be empty")
			errorhandler.LogErrorThenContinue(&err)
			return c.Status(fiber.StatusBadRequest).JSON(&resField{Payload: err.Error()})
		}

		service := u.getShortURLsService(&userID, &isShow)

		return c.Status(service.code).JSON(&resField{Payload: service.data})
	})
}

func (u *url) deleteShortURLController() {
	u.router.Delete(":shortUrl", func(c *fiber.Ctx) error {
		shortURL := c.Params("shortUrl")

		if len(shortURL) == 0 {
			err := errors.New("short url can't be empty")
			errorhandler.LogErrorThenContinue(&err)
			return c.Status(fiber.StatusBadRequest).JSON(&resField{Payload: err.Error()})
		}

		req := new(deleteShortURLReqField)

		if err := c.BodyParser(req); err != nil {
			errorhandler.LogErrorThenContinue(&err)
			return c.Status(fiber.StatusBadRequest).JSON(&resField{Payload: fiber.ErrBadRequest.Error()})
		}

		if err := validator.Struct(req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(&resField{Payload: fiber.ErrBadRequest.Error()})
		}

		service := u.deleteShortURLService(&shortURL, req)

		return c.Status(service.code).JSON(&resField{Payload: service.data})

	})
}
