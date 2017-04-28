-- Copyright 2016 Xiaomi, Inc.
--
-- Licensed under the Apache License, Version 2.0 (the "License");
-- you may not use this file except in compliance with the License.
-- You may obtain a copy of the License at
--
--     http://www.apache.org/licenses/LICENSE-2.0
--
-- Unless required by applicable law or agreed to in writing, software
-- distributed under the License is distributed on an "AS IS" BASIS,
-- WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
-- See the License for the specific language governing permissions and
-- limitations under the License.

SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='TRADITIONAL,ALLOW_INVALID_DATES';

-- -----------------------------------------------------
-- Schema falcon
-- -----------------------------------------------------
CREATE SCHEMA IF NOT EXISTS `falcon` DEFAULT CHARACTER SET utf8 ;
USE `falcon` ;

-- -----------------------------------------------------
-- Table `falcon`.`host`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `falcon`.`host`;
CREATE TABLE `falcon`.`host` (
  `id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `uuid` VARCHAR(128) NOT NULL DEFAULT '' COMMENT 'external/global id'
  `name` VARCHAR(128) NOT NULL DEFAULT '',
  `type` VARCHAR(64) NOT NULL DEFAULT '',
  `status` VARCHAR(64) NOT NULL DEFAULT '',
  `loc` VARCHAR(128) NOT NULL DEFAULT '',
  `idc` VARCHAR(128) NOT NULL DEFAULT '',
  `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `index_status` (`status`),
  INDEX `index_type` (`type`),
  INDEX `index_name` (`name`))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_unicode_ci COMMENT = '机器';


-- -----------------------------------------------------
-- Table `falcon`.`system`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `falcon`.`system`;
CREATE TABLE `falcon`.`system` (
  `id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
  `name` VARCHAR(32) NOT NULL COMMENT '系统名',
  `cname` VARCHAR(64) NOT NULL COMMENT '中文名',
  `developers` VARCHAR(255) NOT NULL COMMENT '系统开发人员',
  `email` VARCHAR(128) NOT NULL DEFAULT '' COMMENT '组邮箱',
  `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `index_name` (`name`))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_unicode_ci COMMENT = '系统模块';


-- -----------------------------------------------------
-- Table `falcon`.`scope`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `falcon`.`scope`;
CREATE TABLE `falcon`.`scope` (
  `id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(32) NOT NULL,
  `system_id` INT(11) UNSIGNED NOT NULL,
  `cname` VARCHAR(64) NOT NULL,
  `note` VARCHAR(255) NOT NULL DEFAULT '',
  `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `index_name` (`name`),
  INDEX `index_system_id` (`system_id`),
  CONSTRAINT `scope_rel_ibfk_1`
    FOREIGN KEY (`system_id`)
    REFERENCES `falcon`.`system` (`id`)
    ON DELETE CASCADE)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_unicode_ci COMMENT = '权限点';


-- -----------------------------------------------------
-- Table `falcon`.`role`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `falcon`.`role`;
CREATE TABLE `falcon`.`role` (
  `id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(32) NOT NULL,
  `cname` VARCHAR(64) NOT NULL,
  `note` VARCHAR(255) NOT NULL DEFAULT '',
  `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `index_name` (`name`))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_unicode_ci COMMENT = '角色';

-- -----------------------------------------------------
-- Table `falcon`.`tag_role_scope`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `falcon`.`tag_role_scope`;
CREATE TABLE `falcon`.`tag_role_scope` (
  `id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `tag_id` INT(11) UNSIGNED NOT NULL,
  `role_id` INT(11) UNSIGNED NOT NULL,
  `scope_id` INT(11) UNSIGNED NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `index_tag_id` (`tag_id`),
  INDEX `index_role_id` (`role_id`),
  INDEX `index_scope_id` (`scope_id`),
  CONSTRAINT `unique_id`
    UNIQUE (`tag_id`,`role_id`,`scope_id`),
  CONSTRAINT `tag_role_scope_ibfk_1`
    FOREIGN KEY (`tag_id`)
    REFERENCES `falcon`.`tag` (`id`)
    ON DELETE CASCADE,
  CONSTRAINT `tag_role_scope_ibfk_2`
    FOREIGN KEY (`role_id`)
    REFERENCES `falcon`.`role` (`id`)
    ON DELETE CASCADE,
  CONSTRAINT `tag_role_scope_ibfk_3`
    FOREIGN KEY (`scope_id`)
    REFERENCES `falcon`.`scope` (`id`)
    ON DELETE CASCADE)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_unicode_ci COMMENT = '节点角色与权限点';


-- -----------------------------------------------------
-- Table `falcon`.`session`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `falcon`.`session`;
CREATE TABLE `falcon`.`session` (
  `session_key` CHAR(64) NOT NULL,
  `session_data` BLOB NULL DEFAULT NULL,
  `session_expiry` INT(11) UNSIGNED NOT NULL,
  PRIMARY KEY (`session_key`))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_unicode_ci;


-- -----------------------------------------------------
-- Table `falcon`.`tag`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `falcon`.`tag`;
CREATE TABLE `falcon`.`tag` (
  `id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(255) NOT NULL DEFAULT '',
  `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `index_name` (`name`))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_unicode_ci;

DROP TABLE IF EXISTS `falcon`.`tag_rel`;
CREATE TABLE `falcon`.`tag_rel` (
  `id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `tag_id` INT(11) UNSIGNED NOT NULL DEFAULT '0',
  `sup_tag_id` INT(11) UNSIGNED NOT NULL DEFAULT '0' COMMENT 'Superior/Self tag id',
  PRIMARY KEY (`id`),
  INDEX `index_tag_id` (`tag_id`),
  INDEX `index_sup_tag_id` (`sup_tag_id`),
  CONSTRAINT `tag_rel_ibfk_1`
    FOREIGN KEY (`tag_id`)
    REFERENCES `falcon`.`tag` (`id`)
    ON DELETE CASCADE,
  CONSTRAINT `tag_rel_ibfk_2`
    FOREIGN KEY (`sup_tag_id`)
    REFERENCES `falcon`.`tag` (`id`)
    ON DELETE CASCADE)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_unicode_ci;


-- -----------------------------------------------------
-- Table `falcon`.`tag_host`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `falcon`.`tag_host`;
CREATE TABLE `falcon`.`tag_host` (
  `id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `tag_id` INT(11) UNSIGNED NOT NULL DEFAULT '0',
  `host_id` INT(11) UNSIGNED NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  INDEX `index_tag_id` (`tag_id`),
  INDEX `index_host_id` (`host_id`),
  CONSTRAINT `tag_host_rel_ibfk_1`
    FOREIGN KEY (`tag_id`)
    REFERENCES `falcon`.`tag` (`id`)
    ON DELETE CASCADE,
  CONSTRAINT `tag_host_rel_ibfk_2`
    FOREIGN KEY (`host_id`)
    REFERENCES `falcon`.`host` (`id`)
    ON DELETE CASCADE)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_unicode_ci;


-- -----------------------------------------------------
-- Table `falcon`.`user`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `falcon`.`user`;
CREATE TABLE `falcon`.`user` (
  `id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `uuid` VARCHAR(128) NOT NULL DEFAULT '' COMMENT 'external/global id'
  `name` VARCHAR(32) NOT NULL,
  `cname` VARCHAR(64) NOT NULL DEFAULT '',
  `email` VARCHAR(128) NOT NULL DEFAULT '',
  `phone` VARCHAR(16) NOT NULL DEFAULT '',
  `im` VARCHAR(32) NOT NULL DEFAULT '',
  `qq` VARCHAR(16) NOT NULL DEFAULT '',
  `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `index_name` (`name`))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_unicode_ci;


-- -----------------------------------------------------
-- Table `falcon`.`tag_role_user`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `falcon`.`tag_role_user`;
CREATE TABLE `falcon`.`tag_role_user` (
  `id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `tag_id` INT(11) UNSIGNED NOT NULL DEFAULT '0',
  `role_id` INT(11) UNSIGNED NOT NULL DEFAULT '0',
  `user_id` INT(11) UNSIGNED NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  INDEX `index_tag_id` (`tag_id`),
  INDEX `index_role_id` (`role_id`),
  INDEX `index_user_id` (`user_id`),
  CONSTRAINT `unique_id`
    UNIQUE (`tag_id`,`role_id`,`user_id`),
  CONSTRAINT `tag_role_user_ibfk_1`
    FOREIGN KEY (`tag_id`)
    REFERENCES `falcon`.`tag` (`id`)
    ON DELETE CASCADE,
  CONSTRAINT `tag_role_user_ibfk_2`
    FOREIGN KEY (`role_id`)
    REFERENCES `falcon`.`role` (`id`)
    ON DELETE CASCADE,
  CONSTRAINT `tag_role_user_ibfk_3`
    FOREIGN KEY (`user_id`)
    REFERENCES `falcon`.`user` (`id`)
    ON DELETE CASCADE)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_unicode_ci;

-- -----------------------------------------------------
-- Table `falcon`.`log`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `falcon`.`log`;
CREATE TABLE `falcon`.`log` (
  `id`        INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `module`    VARCHAR(64) NOT NULL DEFAULT '',
  `module_id` INT(11) UNSIGNED NOT NULL DEFAULT 0,
  `user_id`   INT(11) UNSIGNED NOT NULL DEFAULT 0,
  `action`    TINYINT(4) UNSIGNED NOT NULL DEFAULT 0,
  `data`      BLOB NULL DEFAULT NULL,
  `time`      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_unicode_ci;


DROP VIEW IF EXISTS `falcon`.`tag_role_user_scope` ;
DROP TABLE IF EXISTS `falcon`.`tag_role_user_scope`;
CREATE OR REPLACE 
VIEW `falcon`.`tag_role_user_scope` AS
    SELECT 
        `a`.`tag_id` AS `user_tag_id`,
        `b`.`tag_id` AS `scope_tag_id`,
        `a`.`role_id` AS `role_id`,
        `a`.`user_id` AS `user_id`,
        `b`.`scope_id` AS `scope_id`
    FROM
        (`falcon`.`tag_role_user` `a`
        JOIN `falcon`.`tag_role_scope` `b` ON ((`a`.`role_id` = `b`.`role_id`)));


DROP VIEW IF EXISTS `falcon`.`user_scope` ;
DROP TABLE IF EXISTS `falcon`.`user_scope`;
CREATE OR REPLACE 
VIEW `falcon`.`user_scope` AS
    SELECT 
        `a`.`user_id` AS `user_id`,
        `a`.`scope_id` AS `scope_id`,
        `a`.`scope_tag_id` AS `tag_id`,
        `a`.`role_id` AS `role_id`
    FROM
        (`falcon`.`tag_role_user_scope` `a`
        JOIN `falcon`.`tag_rel` `b` ON (((`a`.`user_tag_id` = `b`.`tag_id`)
            AND (`a`.`scope_tag_id` = `b`.`sup_tag_id`))))
    GROUP BY `a`.`scope_id`
    HAVING (`a`.`scope_tag_id` = MAX(`a`.`scope_tag_id`));




SET SQL_MODE=@OLD_SQL_MODE;
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;

INSERT INTO `falcon`.`tag` (`name`) VALUES
    ('cop=xiaomi'),
    ('cop=xiaomi,owt=inf'),
    ('cop=xiaomi,owt=inf,pdl=falcon'),
    ('cop=xiaomi,owt=miliao'),
    ('cop=xiaomi,owt=miliao,pdl=op'),
    ('cop=xiaomi,owt=miliao,pdl=op,service=dpdk'),
    ('cop=xiaomi,owt=miliao,pdl=op,service=uaq');

INSERT INTO `falcon`.`tag_rel` (`tag_id`, `sup_tag_id`) VALUES
    (2, 1),
    (2, 2),
    (3, 1),
    (3, 2),
    (3, 3),
    (4, 1),
    (4, 4),
    (5, 1),
    (5, 4),
    (5, 5),
    (6, 1),
    (6, 4),
    (6, 5),
    (6, 6),
    (7, 1),
    (7, 4),
    (7, 5),
    (7, 7);

-- 1
INSERT INTO `falcon`.`host` (`uuid`, `name`, `type`, `status`, `loc`, `idc`) VALUES
    ('1', 'c3-op-mon-graph01.bj', '21vianet', 'online', 'bj', 'c3'),
    ('2', 'c3-op-mon-graph02.bj', '21vianet', 'online', 'bj', 'c3'),
    ('3', 'c3-op-mon-graph03.bj', '21vianet', 'online', 'bj', 'c3'),
    ('4', 'c3-op-mon-dpdk01.bj', 'machine', 'online', 'bj', 'lg'),
    ('5', 'c3-op-mon-dpdk02.bj', 'machine', 'online', 'bj', 'lg'),
    ('6', 'c3-op-mon-dpdk03.bj', 'machine', 'online', 'bj', 'lg'),
    ('7', 'c3-op-mon-uaq01.bj', 'machine', 'online', 'bj', 'lg'),
    ('8', 'c3-op-mon-uaq02.bj', 'machine', 'online', 'bj', 'lg'),
    ('9', 'c3-op-mon-uaq03.bj', 'machine', 'online', 'bj', 'lg');

INSERT INTO `falcon`.`user` (`uuid`, `name`, `cname`, `email`, `phone`, `im`, `qq`) VALUES
    ('cn=yubo,ou=users,dc=yubo,dc=org@ldap', 'yubo', 'yubo', 'yubo@yubo.org', '110', 'x80386', '20507'),
    ('tom@xiaomi.com', 'tom', 'tom', 'tom@yubo.org', '1860000000', '1234', '20507');


INSERT INTO `falcon`.`role` (`name`, `cname`, `note`) VALUES
    ('admin', '超级管理员', '配置所有选项'),
    ('manager', '管理员', 'user'),
    ('sre', '工程师', 'Site Reliability Engineering'),
    ('user', '普通用户', 'user');


INSERT INTO `falcon`.`system` (`name`, `cname`, `developers`, `email`) VALUES
    ('service-norns', '机器管理', 'yubo', 'yubo@xiaomi.com');


INSERT INTO `falcon`.`scope` (`name`, `system_id`, `cname`, `note`) VALUES
    ('service-norns-tag-edit', 1, '节点修改', '允许添加删除节点'),
    ('service-norns-host-operate', 1, '机器操作', '重启，改名，关机'),
    ('service-norns-host-bind', 1, '机器挂载', '允许挂载，删除机器与节点的对应关系'),
    ('service-norns-tag-read', 1, '节点读', '查看节点及节点下相关内容');


INSERT INTO `falcon`.`tag_role_scope` (`tag_id`, `role_id`, `scope_id`) VALUES
    (1, 1, 1),
    (1, 1, 2),
    (1, 1, 3),
    (1, 1, 4),
    (1, 2, 1),
    (1, 2, 2),
    (1, 2, 3),
    (1, 2, 4),
    (1, 3, 2),
    (1, 3, 3),
    (1, 3, 4),
    (1, 4, 4);


INSERT INTO `falcon`.`tag_role_user` (`tag_id`, `role_id`, `user_id`) VALUES
    (3, 3, 1);

