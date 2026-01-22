-- +goose Up
-- +goose StatementBegin
CREATE TABLE `subscriptions` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `uuid` varchar(36) NOT NULL,
    `url` varchar(256) NOT NULL,
    `secret` varchar(256) DEFAULT NULL,
    `created_at` datetime NOT NULL DEFAULT current_timestamp(),
    `updated_at` datetime NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
    PRIMARY KEY (`id`),
    UNIQUE KEY `uuid` (`uuid`)
) ENGINE = InnoDB;
-- +goose StatementEnd
-- +goose StatementBegin
CREATE TABLE `subscription_events` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `subscription_id` bigint(20) NOT NULL,
    `event` varchar(64) NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `subscription_id` (`subscription_id`, `event`),
    KEY `event` (`event`),
    CONSTRAINT `subscription_events_ibfk_1` FOREIGN KEY (`subscription_id`) REFERENCES `subscriptions` (`id`) ON DELETE CASCADE
) ENGINE = InnoDB;
-- +goose StatementEnd
---
-- +goose Down
-- +goose StatementBegin
DROP TABLE `subscription_events`;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE `subscriptions`;
-- +goose StatementEnd