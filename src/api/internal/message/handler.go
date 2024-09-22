package message

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tmlsergen/messaging-service-api/internal/app"
)

type Handler struct {
	service messageService
}

type messageService interface {
	GetMessages(c *fiber.Ctx, page, limit int) (GetMessagesResponse, error)
	HandleCronAction(*fiber.Ctx, string) error
}

func NewHandler(service messageService) *Handler {
	return &Handler{service}
}

// GetMessages   return processed messages
// @Summary      Fetch processed messages
// @Description  Fetch processed messages
// @Tags         Messages
// @Accept       json
// @Produce      json
// @Param        page  query     int  true  "Page number"
// @Param        limit query     int  true  "Number of items per page"
// @Success      200   {object}  GetMessagesResponse  "OK"
// @Failure      400   {object}  app.HTTPError  "Bad Request"
// @Failure      500   {object}  app.HTTPError  "Internal Server Error"
// @Router       /api/v1/messages [get]
func (h *Handler) GetMessages(c *fiber.Ctx) error {
	page, limit := c.QueryInt("page", 1), c.QueryInt("limit", 10)

	resp, err := h.service.GetMessages(c, page, limit)
	if err != nil {
		return err
	}

	return c.JSON(resp)
}

type CronActionRequest struct {
	Action string `json:"action" validate:"required,oneof=start stop"`
}

type CronActionResponse struct {
	Status string `json:"status"`
}

// HandleCronAction handles the cron action
// @Summary      Handle cron action
// @Description  Handle cron action
// @Tags         Messages
// @Accept       json
// @Produce      json
// @Param        action  body      CronActionRequest  true  "Cron action to be processed"
// @Success      200     {object}  CronActionResponse "OK"
// @Failure      400     {object}  app.HTTPError "Bad Request"
// @Failure      500     {object}  app.HTTPError "Internal Server Error"
// @Router       /api/v1/messages/cron [post]
func (h *Handler) HandleCronAction(c *fiber.Ctx) error {
	var cronAction CronActionRequest
	if err := app.BindAndValidate(c, &cronAction); err != nil {
		return err
	}

	err := h.service.HandleCronAction(c, cronAction.Action)
	if err != nil {
		return err
	}

	resp := CronActionResponse{
		Status: "ok",
	}

	return c.JSON(resp)
}
