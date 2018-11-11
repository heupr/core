-- MySQL dump 10.16  Distrib 10.1.26-MariaDB, for debian-linux-gnu (x86_64)
--
-- Host: 127.0.0.1    Database: heupr
-- ------------------------------------------------------
-- Server version	5.5.8

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `github_event_assignees`
--

DROP TABLE IF EXISTS `github_event_assignees`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `github_event_assignees` (
  `id` bigint(11) NOT NULL AUTO_INCREMENT,
  `repo_id` int(11) DEFAULT NULL,
  `issues_id` int(11) DEFAULT NULL,
  `number` int(11) DEFAULT NULL,
  `is_closed` tinyint(1) DEFAULT NULL,
  `is_pull` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`id`)
) AUTO_INCREMENT=37372;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `github_event_assignees_lk`
--

DROP TABLE IF EXISTS `github_event_assignees_lk`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `github_event_assignees_lk` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `github_event_assignees_fk` bigint(20) DEFAULT NULL,
  `assignee` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL,
  PRIMARY KEY (`id`)
) AUTO_INCREMENT=37757;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `github_events`
--

DROP TABLE IF EXISTS `github_events`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `github_events` (
  `id` bigint(11) NOT NULL AUTO_INCREMENT,
  `repo_id` int(11) DEFAULT NULL,
  `issues_id` int(11) DEFAULT NULL,
  `number` int(11) DEFAULT NULL,
  `action` varchar(32) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL,
  `is_closed` tinyint(1) DEFAULT NULL,
  `closed_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `is_pull` tinyint(1) DEFAULT NULL,
  `payload` JSON COLLATE utf8_bin NOT NULL,
  PRIMARY KEY (`id`)
) AUTO_INCREMENT=46424;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `integrations`
--

DROP TABLE IF EXISTS `integrations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `integrations` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `repo_id` int(11) DEFAULT NULL,
  `app_id` int(11) DEFAULT NULL,
  `installation_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) AUTO_INCREMENT=360;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `integrations_settings`
--

DROP TABLE IF EXISTS `integrations_settings`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `integrations_settings` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `repo_id` int(11) DEFAULT NULL,
  `start_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `email` varchar(25) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL,
  `twitter` varchar(25) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL,
  `enable_triager` tinyint(1) NOT NULL,
  `enable_labeler` tinyint(1) NOT NULL,
  PRIMARY KEY (`id`)
) AUTO_INCREMENT=251;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `integrations_settings_ignorelabels_lk`
--

DROP TABLE IF EXISTS `integrations_settings_ignorelabels_lk`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `integrations_settings_ignorelabels_lk` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `integrations_settings_fk` bigint(20) DEFAULT NULL,
  `label` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL,
  PRIMARY KEY (`id`)
) AUTO_INCREMENT=19;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `integrations_settings_ignoreusers_lk`
--

DROP TABLE IF EXISTS `integrations_settings_ignoreusers_lk`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `integrations_settings_ignoreusers_lk` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `integrations_settings_fk` bigint(20) DEFAULT NULL,
  `user` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL,
  PRIMARY KEY (`id`)
) AUTO_INCREMENT=11;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `integrations_settings_labels_bif_lk`
--

DROP TABLE IF EXISTS `integrations_settings_labels_bif_lk`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `integrations_settings_labels_bif_lk` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `integrations_settings_fk` bigint(20) DEFAULT NULL,
  `repo_id` int(11) DEFAULT NULL,
  `bug` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL,
  `feature` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL,
  `improvement` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL,
  PRIMARY KEY (`id`)
) AUTO_INCREMENT=22;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2018-11-11 20:09:52
