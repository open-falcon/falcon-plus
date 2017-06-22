CREATE DATABASE graph
  DEFAULT CHARACTER SET utf8
  DEFAULT COLLATE utf8_general_ci;
USE graph;
SET NAMES utf8;

DROP TABLE if exists `graph`.`endpoint`;
CREATE TABLE `graph`.`endpoint` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `endpoint` varchar(255) NOT NULL DEFAULT '',
  `ts` int(11) DEFAULT NULL,
  `t_create` DATETIME NOT NULL COMMENT 'create time',
  `t_modify` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'last modify time',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_endpoint` (`endpoint`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE if exists `graph`.`endpoint_counter`;
CREATE TABLE `graph`.`endpoint_counter` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `endpoint_id` int(10) unsigned NOT NULL,
  `counter` varchar(255) NOT NULL DEFAULT '',
  `step` int(11) not null default 60 comment 'in second',
  `type` varchar(16) not null comment 'GAUGE|COUNTER|DERIVE',
  `ts` int(11) DEFAULT NULL,
  `t_create` DATETIME NOT NULL COMMENT 'create time',
  `t_modify` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'last modify time',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_endpoint_id_counter` (`endpoint_id`, `counter`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE if exists `graph`.`tag_endpoint`;
CREATE TABLE `graph`.`tag_endpoint` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `tag` varchar(255) NOT NULL DEFAULT '' COMMENT 'srv=tv',
  `endpoint_id` int(10) unsigned NOT NULL,
  `ts` int(11) DEFAULT NULL,
  `t_create` DATETIME NOT NULL COMMENT 'create time',
  `t_modify` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'last modify time',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_tag_endpoint_id` (`tag`, `endpoint_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

