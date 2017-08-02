USE uic;

LOCK TABLES `user` WRITE;
/*!40000 ALTER TABLE `user` DISABLE KEYS */;
INSERT INTO `user` VALUES (1,'testuser99','48515098703404ade1e8563667c2beab','testuser99','testuser99@open.com','000-000-000-0000','000000000000','000000000000',0,0,'2017-07-17 01:28:24'),(2,'root','db513ac1c000f2303a2f16496d9c22e5','ç®¡ç†å“¡','root@open.com','000-000-000-0000','000000000000','000000000000',2,0,'2017-07-17 01:28:24'),(3,'testuser92','e9dbac3cbd28d9c6f8fb03562075a4b6','testuser92','testuser92@open.com','000-000-000-0000','000000000000','000000000000',0,0,'2017-07-17 01:28:24');
/*!40000 ALTER TABLE `user` ENABLE KEYS */;
UNLOCK TABLES;

LOCK TABLES `session` WRITE;
/*!40000 ALTER TABLE `session` DISABLE KEYS */;
INSERT INTO `session` VALUES (1,2,'4833cd596ad211e79fb3001500c6ca5a',1500542904),(2,3,'48380ba36ad211e79fb3001500c6ca5a',1500542904),(3,1,'282a97736ad711e7b9c5001500c6ca5a',1500544998);
/*!40000 ALTER TABLE `session` ENABLE KEYS */;
UNLOCK TABLES;

LOCK TABLES `team` WRITE;
/*!40000 ALTER TABLE `team` DISABLE KEYS */;
INSERT INTO `team` VALUES (1,'team_A','this is resumeA',3,'2017-08-02 07:33:45'),(2,'team_B','this is resumeB',3,'2017-08-02 07:33:45'),(3,'team_C','this is resumeC',3,'2017-08-02 07:33:45'),(4,'team_D','this is resumeD',3,'2017-08-02 07:33:46'),(5,'team_D1','this is resumeD1',3,'2017-08-02 07:33:46');
/*!40000 ALTER TABLE `team` ENABLE KEYS */;
UNLOCK TABLES;

LOCK TABLES `rel_team_user` WRITE;
/*!40000 ALTER TABLE `rel_team_user` DISABLE KEYS */;
INSERT INTO `rel_team_user` VALUES (1,1,1),(2,1,2),(3,1,3),(4,2,1),(5,3,2),(6,3,3),(7,5,1),(8,5,2);
/*!40000 ALTER TABLE `rel_team_user` ENABLE KEYS */;
UNLOCK TABLES;


