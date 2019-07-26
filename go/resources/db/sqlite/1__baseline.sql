-- +migrate Up
CREATE TABLE `language` (
  `id`         VARCHAR(5) NOT NULL,
  `created_at` DATETIME NOT NULL,
  PRIMARY KEY (`id`)
);

CREATE TABLE `translated_text` (
  `id`         INTEGER PRIMARY KEY,
  `key`        VARCHAR(100) NOT NULL,
  `language`   VARCHAR(50) NOT NULL,
  `value`      TEXT NOT NULL,
  `created_at` DATETIME NOT NULL,
  `updated_at` DATETIME NOT NULL,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`language`) REFERENCES `language`(`id`)
);

CREATE TABLE `group` (
  `id`           VARCHAR(100) PRIMARY KEY,
  `created_at`   DATETIME NOT NULL
);

CREATE TABLE `text_group` (
  `id`         INTEGER PRIMARY KEY,
  `text_key`   VARCHAR(100) NOT NULL,
  `group_id`   VARCHAR(100) NOT NULL,
  `created_at` DATETIME NOT NULL,
  FOREIGN KEY (`group_id`) REFERENCES `group`(`id`)
);

-- +migrate Down
DROP TABLE IF EXISTS `text_group`;
DROP TABLE IF EXISTS `group`;
DROP TABLE IF EXISTS `translated_text`;
DROP TABLE IF EXISTS `language`;