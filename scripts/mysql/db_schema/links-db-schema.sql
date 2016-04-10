CREATE DATABASE falcon_links
  DEFAULT CHARACTER SET utf8
  DEFAULT COLLATE utf8_general_ci;
USE falcon_links;
SET NAMES utf8;


DROP TABLE IF EXISTS alert;
CREATE TABLE alert
(
  id        INT UNSIGNED NOT NULL AUTO_INCREMENT,
  path      VARCHAR(16)  NOT NULL DEFAULT '',
  content   TEXT         NOT NULL,
  create_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY alert_path(path)
)
  ENGINE =InnoDB
  DEFAULT CHARSET =utf8
  COLLATE =utf8_unicode_ci;

