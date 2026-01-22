package subscriptions

import (
	"errors"
	"fmt"
	"strings"

	"github.com/capcom6/hook-relay/internal/subscriptions"
	"github.com/go-core-fx/fiberfx/handler"
	"github.com/go-core-fx/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

type Handler struct {
	handler.Base

	subscriptionsSvc *subscriptions.Service
}

func NewHandler(
	subscriptionsSvc *subscriptions.Service,
	validator *validator.Validator,
) *Handler {
	return &Handler{
		Base: handler.Base{
			Validator: validator,
		},

		subscriptionsSvc: subscriptionsSvc,
	}
}

func (h *Handler) Register(router fiber.Router) {
	router.Use(h.errorsHandler)

	router.Post("/", h.post)
	router.Get("/", h.list)
	router.Get("/:uuid", h.get)
	router.Delete("/:uuid", h.delete)
}

func (h *Handler) post(c *fiber.Ctx) error {
	req := new(Subscription)
	if err := h.BodyParserValidator(c, req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	res, err := h.subscriptionsSvc.Replace(c.Context(), subscriptions.SubscriptionIn{
		UUID:   req.UUID,
		URL:    req.URL,
		Secret: string(req.Secret),
		Events: req.Events,
	})
	if err != nil {
		return fmt.Errorf("failed to create subscription: %w", err)
	}

	return c.Status(fiber.StatusOK).JSON(newSubscription(res))
}

func (h *Handler) list(c *fiber.Ctx) error {
	if c.Query("uuid") == "" {
		return fiber.NewError(fiber.StatusBadRequest, "uuid is required")
	}

	uuid := strings.Split(c.Query("uuid"), ",")
	res, err := h.subscriptionsSvc.SelectByUUID(c.Context(), uuid...)
	if err != nil {
		return fmt.Errorf("failed to get subscriptions: %w", err)
	}

	return c.Status(fiber.StatusOK).JSON(
		lo.Map(
			res,
			func(s subscriptions.Subscription, _ int) Subscription {
				return *newSubscription(&s)
			},
		),
	)
}

func (h *Handler) get(c *fiber.Ctx) error {
	id := c.Params("uuid")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "uuid is required")
	}

	res, err := h.subscriptionsSvc.GetByUUID(c.Context(), id)
	if err != nil {
		return fmt.Errorf("failed to get subscription: %w", err)
	}

	return c.Status(fiber.StatusOK).JSON(newSubscription(res))
}

func (h *Handler) delete(c *fiber.Ctx) error {
	id := c.Params("uuid")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "uuid is required")
	}

	err := h.subscriptionsSvc.DeleteByUUID(c.Context(), id)
	if err != nil {
		return fmt.Errorf("failed to delete subscription: %w", err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) errorsHandler(c *fiber.Ctx) error {
	err := c.Next()
	if err == nil {
		return nil
	}

	if errors.Is(err, subscriptions.ErrNotFound) {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return err //nolint:wrapcheck //already wrapped
}
