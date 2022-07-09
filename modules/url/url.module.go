package url

import "github.com/gofiber/fiber/v2"

type url struct {
	router fiber.Router
}

func New(router *fiber.Router) *url {
	return &url{
		router: *router,
	}
}

func (u *url) Run() {
	u.createShortURLController()
	u.modifyShortURLController()
	u.moveShortURLsToNewUserIDController()
	u.getLongURLController()
	u.getShortURLsController()
	u.deleteShortURLController()
}
