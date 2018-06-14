create table if not exists eexp (
	`id` int(10) unsigned NOT NULL AUTO_INCREMENT,
	`filters` varchar(1024) NOT NULL, 
	`conditions`  varchar(1024) NOT NULL, 
	`priority`    TINYINT(4)       NOT NULL DEFAULT 0,
	`note`        VARCHAR(1024)    NOT NULL DEFAULT '',		  
	`max_step`    INT(11)          NOT NULL DEFAULT 1,
	`create_user` varchar(64) NOT NULL DEFAULT '',
	`pause` tinyint(1) NOT NULL DEFAULT 0,
	PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
