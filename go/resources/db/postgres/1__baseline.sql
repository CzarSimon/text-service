-- +migrate Up
CREATE TABLE `language` (
  `id`         VARCHAR(5) NOT NULL,
  `created_at` TIMESTAMP NOT NULL,
  PRIMARY KEY (`id`)
);

CREATE TABLE `translated_text` (
  `id`         SERIAL PRIMARY KEY,
  `key`        VARCHAR(100) NOT NULL,
  `language`   VARCHAR(50) NOT NULL,
  `value`      TEXT NOT NULL,
  `created_at` TIMESTAMP NOT NULL,
  `updated_at` TIMESTAMP NOT NULL,
  FOREIGN KEY (`language`) REFERENCES `language`(`id`),
  UNIQUE(`key`, `language`)
);

CREATE TABLE `text_group` (
  `id`           VARCHAR(100) PRIMARY KEY,
  `created_at`   TIMESTAMP NOT NULL
);

CREATE TABLE `text_group_membership` (
  `id`         SERIAL PRIMARY KEY,
  `text_key`   VARCHAR(100) NOT NULL,
  `group_id`   VARCHAR(100) NOT NULL,
  `created_at` TIMESTAMP NOT NULL,
  FOREIGN KEY (`group_id`) REFERENCES `text_group`(`id`)
);

-- +migrate Down
DROP TABLE IF EXISTS `text_group_membership`;
DROP TABLE IF EXISTS `text_group`;
DROP TABLE IF EXISTS `translated_text`;
DROP TABLE IF EXISTS `language`;