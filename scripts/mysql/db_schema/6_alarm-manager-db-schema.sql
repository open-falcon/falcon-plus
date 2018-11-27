-- MySQL dump 10.13  Distrib 5.6.40, for Linux (x86_64)
--
-- Host: 127.0.0.1    Database: alarm_manager
-- ------------------------------------------------------
-- Server version	5.7.19

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

CREATE DATABASE alarm_manager
  DEFAULT CHARACTER SET utf8
  DEFAULT COLLATE utf8_general_ci;
USE alarm_manager;
SET NAMES utf8;

--
-- Table structure for table `event_cases`
--

DROP TABLE IF EXISTS `event_cases`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `event_cases` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `event_case_id` varchar(50) DEFAULT NULL,
  `endpoint` varchar(300) NOT NULL,
  `counter` varchar(200) NOT NULL,
  `func` varchar(50) DEFAULT NULL,
  `cond` varchar(200) NOT NULL,
  `note` varchar(500) DEFAULT NULL,
  `max_step` int(10) unsigned DEFAULT NULL,
  `current_step` int(10) unsigned DEFAULT NULL,
  `priority` int(6) NOT NULL,
  `status` varchar(20) NOT NULL,
  `event_ts` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `update_at` timestamp NULL DEFAULT NULL,
  `closed_at` timestamp NULL DEFAULT NULL,
  `closed_note` varchar(250) DEFAULT NULL,
  `user_modified` int(10) unsigned DEFAULT NULL,
  `tpl_creator` varchar(64) DEFAULT NULL,
  `action_config_id` int(10) unsigned DEFAULT NULL,
  `expression_id` int(10) unsigned DEFAULT NULL,
  `strategy_id` int(10) unsigned DEFAULT NULL,
  `template_id` int(10) unsigned DEFAULT NULL,
  `process_note` mediumint(9) DEFAULT NULL,
  `process_status` varchar(20) DEFAULT 'unresolved',
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `event_case_id` (`event_case_id`),
  KEY `action_config_id` (`action_config_id`),
  KEY `endpoint` (`endpoint`,`strategy_id`,`template_id`),
  KEY `priority` (`priority`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `event_receiver`
--

DROP TABLE IF EXISTS `event_receiver`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `event_receiver` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `event_id` bigint(20) unsigned NOT NULL,
  `action_config_id` int(10) unsigned DEFAULT NULL,
  `uic` varchar(50) NOT NULL,
  `user` varchar(50) NOT NULL,
  `event_ts` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `uic` (`uic`),
  KEY `user` (`user`),
  KEY `uic_2` (`uic`,`user`),
  KEY `event_id` (`event_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `events`
--

DROP TABLE IF EXISTS `events`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `events` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `event_case_id` varchar(50) DEFAULT NULL,
  `current_step` int(10) unsigned DEFAULT NULL,
  `cond` varchar(200) NOT NULL,
  `status` varchar(20) NOT NULL,
  `event_ts` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `event_case_id` (`event_case_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `fault`
--

DROP TABLE IF EXISTS `fault`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `fault` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `title` varchar(255) DEFAULT NULL,
  `note` varchar(255) DEFAULT NULL,
  `creator` varchar(255) DEFAULT NULL,
  `owner` varchar(255) DEFAULT NULL,
  `state` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_fault_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `fault_comment`
--

DROP TABLE IF EXISTS `fault_comment`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `fault_comment` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `fault_id` int(10) unsigned DEFAULT NULL,
  `creator` varchar(255) DEFAULT NULL,
  `comment` varchar(999) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_fault_comment_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `fault_event`
--

DROP TABLE IF EXISTS `fault_event`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `fault_event` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `fault_id` int(10) unsigned DEFAULT NULL,
  `event_id` int(10) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_faultid_eventid` (`fault_id`,`event_id`),
  KEY `idx_fault_event_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `fault_follower`
--

DROP TABLE IF EXISTS `fault_follower`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `fault_follower` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `fault_id` int(10) unsigned DEFAULT NULL,
  `follower` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_faultid_follower` (`fault_id`,`follower`),
  KEY `idx_fault_follower_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `fault_tag`
--

DROP TABLE IF EXISTS `fault_tag`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `fault_tag` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `fault_id` int(10) unsigned DEFAULT NULL,
  `tag` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_faultid_tag` (`fault_id`,`tag`),
  KEY `idx_fault_tag_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `state_change_log`
--

DROP TABLE IF EXISTS `state_change_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `state_change_log` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `refer_table` varchar(255) DEFAULT NULL,
  `refer_id` int(10) unsigned DEFAULT NULL,
  `refer_field` varchar(255) DEFAULT NULL,
  `from` varchar(255) DEFAULT NULL,
  `to` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_state_change_log_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2018-11-15 18:08:45
