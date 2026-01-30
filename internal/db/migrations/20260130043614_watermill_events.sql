-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS `watermill_events` (
    `offset` BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `uuid` VARCHAR(36) NOT NULL,
    `payload` JSON DEFAULT NULL,
    `metadata` JSON DEFAULT NULL,
    `acked` BOOLEAN NOT NULL DEFAULT FALSE,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd
---
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS `watermill_events`;
-- +goose StatementEnd