USE falcon_portal;
LOCK TABLES `tpl` WRITE;
/*!40000 ALTER TABLE `tpl` DISABLE KEYS */;
INSERT INTO `tpl` VALUES (1,'mytpl1',0,0,'testuser99','2017-07-27 10:31:55'),(2,'mytpl2',0,0,'testuser99','2017-07-27 10:31:55'),(3,'mytpl3',0,0,'testuser99','2017-07-27 10:31:55');
/*!40000 ALTER TABLE `tpl` ENABLE KEYS */;
UNLOCK TABLES;
