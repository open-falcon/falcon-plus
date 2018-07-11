/*
* 建立告警归档资料表, 主要存储各个告警的最后触发状况
*/
CREATE TABLE IF NOT EXISTS event_cases(
                id VARCHAR(50),
                endpoint VARCHAR(512) NOT NULL,
                metric VARCHAR(1024) NOT NULL,
                func VARCHAR(512),
                cond VARCHAR(200) NOT NULL,
                note VARCHAR(500),
                max_step int(10) unsigned,
                current_step int(10) unsigned,
                priority INT(6) NOT NULL,
                status VARCHAR(20) NOT NULL,
                timestamp Timestamp NOT NULL,
                update_at Timestamp NULL DEFAULT NULL,
                closed_at Timestamp NULL DEFAULT NULL,
                closed_note VARCHAR(250),
                user_modified int(10) unsigned,
                tpl_creator VARCHAR(64),
                expression_id int(10) unsigned,
                strategy_id int(10) unsigned,
                template_id int(10) unsigned,
                process_note MEDIUMINT,
                process_status VARCHAR(20) DEFAULT 'unresolved',
                PRIMARY KEY (id),
                INDEX (endpoint, strategy_id, template_id)
)
        ENGINE =InnoDB
        DEFAULT CHARSET =utf8;


/*
* 建立告警归档资料表, 存储各个告警触发状况的历史状态
*/
CREATE TABLE IF NOT EXISTS events (
                id int(10) NOT NULL AUTO_INCREMENT,
                event_caseId VARCHAR(50),
                step int(10) unsigned,
                cond VARCHAR(200) NOT NULL,
                status int(3) unsigned DEFAULT 0,
                timestamp Timestamp,
                PRIMARY KEY (id),
                INDEX(event_caseId),
                FOREIGN KEY (event_caseId) REFERENCES event_cases(id)
                        ON DELETE CASCADE
                        ON UPDATE CASCADE
)
        ENGINE =InnoDB
        DEFAULT CHARSET =utf8;


