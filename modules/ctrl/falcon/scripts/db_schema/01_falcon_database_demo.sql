-- CREATE DATABASE falcon
-- DEFAULT CHARACTER SET utf8
-- DEFAULT COLLATE utf8_general_ci;
INSERT INTO `mysql`.`user`(`Host`,`User`,`Password`) VALUES ("localhost","falcon",password("1234")) ON DUPLICATE KEY UPDATE `Password`=password("1234");
GRANT ALL PRIVILEGES ON `falcon` .* TO  `falcon`  @localhost identified by '1234';
FLUSH PRIVILEGES;
