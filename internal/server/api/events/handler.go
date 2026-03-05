package events

import (
	"fmt"

	"github.com/capcom6/hook-relay/internal/events"
	"github.com/go-core-fx/fiberfx/handler"
	"github.com/go-core-fx/validator"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type Handler struct {
	handler.Base

	publisher *events.Publisher

	logger *zap.Logger
}

func NewHandler(
	validator *validator.Validator,
	publisher *events.Publisher,
	logger *zap.Logger,
) handler.Handler {
	return &Handler{
		Base: handler.Base{
			Validator: validator,
		},

		publisher: publisher,

		logger: logger,
	}
}

func (h *Handler) Register(router fiber.Router) {
	router = router.Group("events")

	router.Use(h.errorsHandler)

	router.Post("", h.post)
}

//	@Summary		Publish an event
//	@Description	Publishes an event to the event bus for delivery to subscribed webhooks
//	@Tags			events
//	@Accept			json
//	@Produce		json
//	@Param			request	body	EventRequest	true	"Event data"
//	@Success		202		"Accepted"
//	@Failure		400		{object}	fiberfx.ErrorResponse	"Bad Request"
//	@Failure		500		{object}	fiberfx.ErrorResponse	"Internal Server Error"
//	@Router			/events [post]
//
// Publish an event.
func (h *Handler) post(c *fiber.Ctx) error {
	req := new(EventRequest)
	if err := h.BodyParserValidator(c, req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := h.publisher.Publish(req.ID, events.Event{
		Type:    req.Type,
		Payload: req.Payload,
	}); err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	// h.logger.Info("event", zap.String("event_type", req.Type), zap.Any("payload", req.Payload))

	return c.SendStatus(fiber.StatusAccepted)
}

func (h *Handler) errorsHandler(c *fiber.Ctx) error {
	err := c.Next()
	if err == nil {
		return nil
	}

	// TODO: handle errors

	return err //nolint:wrapcheck //already wrapped
}
