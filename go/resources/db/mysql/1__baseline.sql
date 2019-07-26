-- +migrate Up
CREATE TABLE `language` (
  `id`         VARCHAR(5) NOT NULL,
  `created_at` TIMESTAMP NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;;

CREATE TABLE `translated_text` (
  `id`         INT AUTO_INCREMENT PRIMARY KEY,
  `key`        VARCHAR(100) NOT NULL,
  `language`   VARCHAR(50) NOT NULL,
  `value`      TEXT NOT NULL,
  `created_at` TIMESTAMP NOT NULL,
  `updated_at` TIMESTAMP NOT NULL,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`language`) REFERENCES `language`(`id`),
  UNIQUE(`key`, `language`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;;

CREATE TABLE `group` (
  `id`           VARCHAR(100) PRIMARY KEY,
  `created_at`   TIMESTAMP NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;;

CREATE TABLE `text_group` (
  `id`         INT AUTO_INCREMENT PRIMARY KEY,
  `text_key`   VARCHAR(100) NOT NULL,
  `group_id`   VARCHAR(100) NOT NULL,
  `created_at` TIMESTAMP NOT NULL,
  FOREIGN KEY (`group_id`) REFERENCES `group`(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;;

-- +migrate Down
DROP TABLE IF EXISTS `text_group`;
DROP TABLE IF EXISTS `group`;
DROP TABLE IF EXISTS `translated_text`;
DROP TABLE IF EXISTS `language`;