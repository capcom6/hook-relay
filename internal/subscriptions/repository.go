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
func (r *repository) Replace(ctx context.Context, subscription *subscriptionModel) error {
	err := r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		_, err := tx.NewInsert().Replace().Model(subscription).Exec(ctx)
		if err != nil {
			return err //nolint:wrapcheck //wrapped outside
		}

		events := lo.Map(
			subscription.Events,
			func(item eventModel, _ int) eventModel {
				item.SubscriptionID = subscription.ID
				return item
			},
		)

		_, err = tx.NewInsert().Model(&events).Exec(ctx)
		if err != nil {
			return err //nolint:wrapcheck //wrapped outside
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to replace subscription: %w", err)
	}

	return nil
}

// SelectByUUID retrieves a subscription by its UUID.
func (r *repository) SelectByUUID(ctx context.Context, uuid ...string) ([]subscriptionModel, error) {
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
	return subscriptions, nil
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
