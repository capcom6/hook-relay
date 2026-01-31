package subscriptions

import (
	"context"
	"fmt"

	"github.com/samber/lo"
	"github.com/uptrace/bun"
)

type repository struct {
	db *bun.DB
}

func newRepository(db *bun.DB) *repository {
	return &repository{
		db: db,
	}
}

// Replace replaces an existing subscription.
func (r *repository) Replace(ctx context.Context, subscription SubscriptionIn) (*Subscription, error) {
	model := newSubscriptionModel(subscription)

	err := r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		_, err := tx.NewInsert().Replace().Model(model).Returning("*").Exec(ctx)
		if err != nil {
			return err //nolint:wrapcheck //wrapped outside
		}

		events := lo.Map(
			subscription.Events,
			func(item string, _ int) eventModel {
				return eventModel{
					BaseModel:      bun.BaseModel{},
					ID:             0,
					SubscriptionID: model.ID,
					Event:          item,
				}
			},
		)

		_, err = tx.NewInsert().Model(&events).Exec(ctx)
		if err != nil {
			return err //nolint:wrapcheck //wrapped outside
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to replace subscription: %w", err)
	}

	return model.toDomain(), nil
}

// SelectByUUID retrieves a subscription by its UUID.
func (r *repository) SelectByUUID(ctx context.Context, uuid ...string) ([]Subscription, error) {
	if len(uuid) == 0 {
		return nil, nil
	}

	subscriptions := []subscriptionModel{}
	err := r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		err := tx.NewSelect().
			Model(&subscriptions).
			Relation("Events").
			Where("uuid IN (?)", bun.In(uuid)).
			Scan(ctx)
		if err != nil {
			return err //nolint:wrapcheck //wrapped outside
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get subscriptions: %w", err)
	}
	return lo.Map(subscriptions, func(item subscriptionModel, _ int) Subscription { return *item.toDomain() }), nil
}

func (r *repository) SelectByEvent(ctx context.Context, event string) ([]Subscription, error) {
	subscriptions := []subscriptionModel{}
	err := r.db.NewSelect().
		Model(&subscriptions).
		Relation("Events").
		Where("id IN (?)", r.db.NewSelect().
			Model((*eventModel)(nil)).
			Column("subscription_id").
			Where("event = ?", event)).
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscriptions: %w", err)
	}
	return lo.Map(subscriptions, func(item subscriptionModel, _ int) Subscription { return *item.toDomain() }), nil
}

// DeleteByUUID deletes a subscription.
func (r *repository) DeleteByUUID(ctx context.Context, uuid string) error {
	_, err := r.db.NewDelete().
		Model((*subscriptionModel)(nil)).
		Where("uuid = ?", uuid).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete subscription: %w", err)
	}
	return nil
}
