USE falcon_portal;

LOCK TABLES `grp` WRITE;
/*!40000 ALTER TABLE `grp` DISABLE KEYS */;
INSERT INTO `grp` VALUES (1,'testhg1','root','2017-08-22 03:09:08',1);
INSERT INTO `grp` VALUES (2,'testhg2','root','2017-08-22 03:09:09',1);
INSERT INTO `grp` VALUES (3,'testhg3','root','2017-08-22 03:09:10',1);
/*!40000 ALTER TABLE `grp` ENABLE KEYS */;
UNLOCK TABLES;
