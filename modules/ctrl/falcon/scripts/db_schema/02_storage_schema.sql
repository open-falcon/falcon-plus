USE falcon;
SET NAMES utf8;

DROP TABLE if exists `falcon`.`host`;
CREATE TABLE `falcon`.`host` (
	`id` int(10) unsigned NOT NULL AUTO_INCREMENT,
	`host` varchar(255) NOT NULL DEFAULT '',
	`ts` int(11) DEFAULT NULL,
	`t_create` DATETIME NOT NULL COMMENT 'create time',
	`t_modify` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'last modify time',
	PRIMARY KEY (`id`),
	UNIQUE KEY `idx_host` (`host`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE if exists `falcon`.`counter`;
CREATE TABLE `falcon`.`counter` (
	`id` int(10) unsigned NOT NULL AUTO_INCREMENT,
	`counter` varchar(255) NOT NULL DEFAULT '',
	`host_id` int(10) unsigned NOT NULL,
	`step` int(11) not null default 60 comment 'in second',
	`type` varchar(16) not null comment 'GAUGE|COUNTER|DERIVE',
	`ts` int(11) DEFAULT NULL,
	`t_create` DATETIME NOT NULL COMMENT 'create time',
	`t_modify` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'last modify time',
	PRIMARY KEY (`id`),
	UNIQUE KEY `idx_host_id_counter` (`counter`, `host_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE if exists `falcon`.`tag`;
CREATE TABLE `falcon`.`tag` (
	`id` int(10) unsigned NOT NULL AUTO_INCREMENT,
	`tag` varchar(255) NOT NULL DEFAULT '' COMMENT 'srv=tv',
	`host_id` int(10) unsigned NOT NULL,
	`ts` int(11) DEFAULT NULL,
	`t_create` DATETIME NOT NULL COMMENT 'create time',
	`t_modify` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'last modify time',
	PRIMARY KEY (`id`),
	UNIQUE KEY `idx_tag_host_id` (`tag`, `host_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
