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
) handler.Handler {
	return &Handler{
		Base: handler.Base{
			Validator: validator,
		},

		subscriptionsSvc: subscriptionsSvc,
	}
}

func (h *Handler) Register(router fiber.Router) {
	router = router.Group("subscriptions")

	router.Use(h.errorsHandler)

	router.Post("", h.post)
	router.Get("", h.list)
	router.Get(":uuid", h.get)
	router.Delete(":uuid", h.delete)
}

//	@Summary		Create or replace a subscription
//	@Description	Creates a new webhook subscription or replaces an existing one with the same UUID
//	@Tags			subscriptions
//	@Accept			json
//	@Produce		json
//	@Param			request	body		Subscription	true	"Subscription data"
//	@Success		200		{object}	Subscription
//	@Failure		400		{object}	fiberfx.ErrorResponse	"Bad Request"
//	@Failure		500		{object}	fiberfx.ErrorResponse	"Internal Server Error"
//	@Router			/subscriptions [post]
//
// Create or replace a subscription.
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

//	@Summary		List subscriptions
//	@Description	Retrieves a list of subscriptions by their UUIDs
//	@Tags			subscriptions
//	@Accept			json
//	@Produce		json
//	@Param			uuid	query		string	true	"Comma-separated list of subscription UUIDs"
//	@Success		200		{array}		Subscription
//	@Failure		400		{object}	fiberfx.ErrorResponse	"Bad Request"
//	@Failure		500		{object}	fiberfx.ErrorResponse	"Internal Server Error"
//	@Router			/subscriptions [get]
//
// List subscriptions.
func (h *Handler) list(c *fiber.Ctx) error {
	if c.Query("uuid") == "" {
		return fiber.NewError(fiber.StatusBadRequest, "uuid is required")
	}

	uuid := lo.FilterMap(
		strings.Split(c.Query("uuid"), ","),
		func(s string, _ int) (string, bool) {
			s = strings.TrimSpace(s)
			return s, s != ""
		},
	)
	if len(uuid) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "uuid is required")
	}

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

//	@Summary		Get a subscription
//	@Description	Retrieves a single subscription by its UUID
//	@Tags			subscriptions
//	@Accept			json
//	@Produce		json
//	@Param			uuid	path		string	true	"Subscription UUID"
//	@Success		200		{object}	Subscription
//	@Failure		400		{object}	fiberfx.ErrorResponse	"Bad Request"
//	@Failure		404		{object}	fiberfx.ErrorResponse	"Not Found"
//	@Failure		500		{object}	fiberfx.ErrorResponse	"Internal Server Error"
//	@Router			/subscriptions/{uuid} [get]
//
// Get a subscription.
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

//	@Summary		Delete a subscription
//	@Description	Deletes a subscription by its UUID
//	@Tags			subscriptions
//	@Accept			json
//	@Produce		json
//	@Param			uuid	path	string	true	"Subscription UUID"
//	@Success		204		"No Content"
//	@Failure		400		{object}	fiberfx.ErrorResponse	"Bad Request"
//	@Failure		500		{object}	fiberfx.ErrorResponse	"Internal Server Error"
//	@Router			/subscriptions/{uuid} [delete]
//
// Delete a subscription.
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
