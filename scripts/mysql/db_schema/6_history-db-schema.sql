CREATE DATABASE history
  DEFAULT CHARACTER SET utf8mb4
  DEFAULT COLLATE utf8mb4_unicode_ci;
USE history;
SET NAMES utf8mb4;

CREATE TABLE IF NOT EXISTS history (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `endpoint` varchar(255) NOT NULL DEFAULT '',
  `step` int(11) not null default 60 comment 'in second',
  `counter_type` varchar(16) not null comment 'GAUGE|COUNTER|DERIVE',
  `metric`   VARCHAR(128) NOT NULL DEFAULT '',
  `value`   VARCHAR(1024) NOT NULL DEFAULT '',  
  `tags`     VARCHAR(1024) NOT NULL DEFAULT '',
  `timestamp` int(11) DEFAULT NULL,
  INDEX(`endpoint`),
  INDEX(`metric`),
  INDEX(`timestamp`),
 PRIMARY KEY (`id`)
) ENGINE =InnoDB DEFAULT CHARSET =utf8mb4;

