-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS `slot`
(
    `id`          INT(11)     NOT NULL AUTO_INCREMENT,
    `description` VARCHAR(50) NOT NULL DEFAULT '',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8
  COLLATE = utf8_general_ci;
-- +goose StatementEnd
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS `banner`
(
    `id`          INT(11)     NOT NULL AUTO_INCREMENT,
    `description` VARCHAR(50) NOT NULL DEFAULT '',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8
  COLLATE = utf8_general_ci;
-- +goose StatementEnd
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS `group`
(
    `id`          INT(11)     NOT NULL AUTO_INCREMENT,
    `description` VARCHAR(50) NOT NULL DEFAULT '',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8
  COLLATE = utf8_general_ci;
-- +goose StatementEnd
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS `slot2banner`
(
    `slot_id`   INT(11) NOT NULL DEFAULT 0,
    `banner_id` INT(11) NOT NULL DEFAULT 0,
    PRIMARY KEY (`slot_id`, `banner_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8
  COLLATE = utf8_general_ci;
-- +goose StatementEnd
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS `statistics`
(
    `id`           INT(11)    NOT NULL AUTO_INCREMENT,
    `slot_id`      INT(11)    NOT NULL DEFAULT 0,
    `banner_id`    INT(11)    NOT NULL DEFAULT 0,
    `group_id`     INT(11)    NOT NULL DEFAULT 0,
    `event_type`   TINYINT(1) NOT NULL DEFAULT 0,
    `processed_at` TIMESTAMP  NULL     DEFAULT NULL,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8
  COLLATE = utf8_general_ci;
-- +goose StatementEnd
-- +goose StatementBegin
INSERT INTO `slot` (`description`)
VALUES ('слот 1'),
       ('слот 2'),
       ('слот 3');
-- +goose StatementEnd
-- +goose StatementBegin
INSERT INTO `banner` (`description`)
VALUES ('banner 1'),
       ('banner 2'),
       ('banner 3'),
       ('banner 4'),
       ('banner 5');
-- +goose StatementEnd
-- +goose StatementBegin
INSERT INTO `group` (`description`)
VALUES ('девушки 20-25'),
       ('дедушки 80+');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE `slot`;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE `banner`;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE `group`;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE `slot2banner`;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE `statistics`;
-- +goose StatementEnd
