CREATE DATABASE `amper` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci */ /*!80016 DEFAULT ENCRYPTION='N' */;

CREATE TABLE `object_sys` (
  `id` int NOT NULL AUTO_INCREMENT,
  `apiName` varchar(250) DEFAULT NULL,
  `title` varchar(500) DEFAULT NULL,
  `titlePlural` varchar(500) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `field_sys` (
  `id` int NOT NULL AUTO_INCREMENT,
  `apiName` varchar(256) DEFAULT NULL,
  `label` varchar(256) DEFAULT NULL,
  `type` varchar(45) DEFAULT NULL,
  `status` tinyint DEFAULT NULL,
  `required` tinyint DEFAULT NULL,
  `entityId` int DEFAULT NULL,
  `createdBy` int DEFAULT NULL,
  `textLength` int DEFAULT NULL,
  `objectReference` int DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `entityId_idx` (`entityId`),
  CONSTRAINT `entityId` FOREIGN KEY (`entityId`) REFERENCES `object_sys` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `object_type_sys` (
  `id` int NOT NULL AUTO_INCREMENT,
  `object_id` int NOT NULL,
  `apiName` varchar(512) NOT NULL,
  `label` varchar(512) NOT NULL,
  `extends` int DEFAULT NULL,
  `created_by` int NOT NULL,
  `created_date` datetime NOT NULL,
  PRIMARY KEY (`id`),
  KEY `object_id_idx` (`object_id`),
  KEY `extends_idx` (`extends`),
  CONSTRAINT `extends` FOREIGN KEY (`extends`) REFERENCES `object_type_sys` (`id`) ON DELETE RESTRICT,
  CONSTRAINT `object_id` FOREIGN KEY (`object_id`) REFERENCES `object_sys` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=98 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `object_type_field_sys` (
  `id` int NOT NULL AUTO_INCREMENT,
  `object_type_id` int NOT NULL,
  `field_id` int NOT NULL,
  `created_by` int NOT NULL,
  `created_date` datetime NOT NULL,
  PRIMARY KEY (`id`),
  KEY `object_type_id_idx` (`object_type_id`),
  KEY `field_d_idx` (`field_id`),
  CONSTRAINT `field_d` FOREIGN KEY (`field_id`) REFERENCES `field_sys` (`id`),
  CONSTRAINT `object_type_id` FOREIGN KEY (`object_type_id`) REFERENCES `object_type_sys` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `profile_sys` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(256) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `users_sys` (
  `id` int NOT NULL AUTO_INCREMENT,
  `firstName` varchar(256) NOT NULL,
  `lastName` varchar(256) NOT NULL,
  `middleName` varchar(256) NOT NULL DEFAULT '',
  `username` varchar(256) NOT NULL,
  `password` varchar(256) NOT NULL,
  `photo` longtext,
  `profile` int NOT NULL,
  `email` varchar(256) NOT NULL,
  `active` int NOT NULL DEFAULT '0',
  `activationCode` varchar(45) NOT NULL DEFAULT '',
  `amperId` int NOT NULL,
  `state` int NOT NULL,
  `config` longtext,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id_UNIQUE` (`id`),
  UNIQUE KEY `username_UNIQUE` (`username`),
  KEY `profile_idx` (`profile`),
  CONSTRAINT `profile` FOREIGN KEY (`profile`) REFERENCES `profile_sys` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=155 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `users_detail_sys`(
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` int NOT NULL,
  `info` longtext,
  `about_me` longtext,
  `responsibilities` longtext,
  `skills` varchar(1000) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id_UNIQUE` (`id`),
  KEY `user_idx` (`user_id`),
  CONSTRAINT `user_id` FOREIGN KEY (`user_id`) REFERENCES `users_sys` (`id`)
  ) ENGINE=InnoDB AUTO_INCREMENT=155 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `users_relationship_sys`(
  `id` int NOT NULL AUTO_INCREMENT,
  `employee_id` int NOT NULL,
  `manager_id` int NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id_UNIQUE` (`id`),
  KEY `employee_idx` (`employee_id`),
  KEY `manager_idx` (`manager_id`),
  CONSTRAINT `employee_id` FOREIGN KEY (`employee_id`) REFERENCES `users_sys` (`id`),
  CONSTRAINT `manager_id` FOREIGN KEY (`manager_id`) REFERENCES `users_sys` (`id`)
  ) ENGINE=InnoDB AUTO_INCREMENT=155 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


CREATE TABLE `dashboard_sys` (
  `id` int NOT NULL AUTO_INCREMENT,
  `label` varchar(16) NOT NULL,
  `description` varchar(150) NOT NULL,
  `configuration` varchar(10240) NOT NULL,
  `created_date` datetime NOT NULL,
  `created_by` int NOT NULL,
  PRIMARY KEY (`id`),
  KEY `userId_idx` (`created_by`),
  CONSTRAINT `userId` FOREIGN KEY (`created_by`) REFERENCES `users_sys` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `widget_sys` (
  `id` int NOT NULL AUTO_INCREMENT,
  `dashboard_id` int NOT NULL,
  `label` varchar(16) NOT NULL,
  `description` varchar(150) NOT NULL,
  `configuration` varchar(10240) NOT NULL,
  `created_date` datetime NOT NULL,
  `created_by` int NOT NULL,
  PRIMARY KEY (`id`),
  KEY `userId_idx` (`created_by`),
  KEY `dashboardId_idx` (`dashboard_id`),
  CONSTRAINT `dashboardId` FOREIGN KEY (`dashboard_id`) REFERENCES `dashboard_sys` (`id`),
  CONSTRAINT `wUserId` FOREIGN KEY (`created_by`) REFERENCES `users_sys` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `amper_sys` (
  `id` int NOT NULL AUTO_INCREMENT,
  `identifier` varchar(250) DEFAULT NULL,
  `name` varchar(250) DEFAULT NULL,
  `type` varchar(100) DEFAULT NULL,
  `address` varchar(500) DEFAULT NULL,
  `port` varchar(10) DEFAULT NULL,
  `state` int NOT NULL DEFAULT '0',
  `state_update_date` datetime NOT NULL,
  `usage` BIGINT NOT NULL DEFAULT '0',
  `limitation` BIGINT NOT NULL DEFAULT '0',
  `directory` varchar(250) DEFAULT NULL,
  `key` varchar(100) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `chat_channel_group_sys` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(250) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `chat_channel_sys` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(250) DEFAULT NULL,
  `group_id` int NOT NULL,
  `amper_id` int NOT NULL,
  `user_ids` longtext,
  `batch_ids` varchar(5000) DEFAULT NULL,
  PRIMARY KEY (`id`),
  CONSTRAINT `groupId` FOREIGN KEY (`group_id`) REFERENCES `chat_channel_group_sys` (`id`),
  CONSTRAINT `amperId` FOREIGN KEY (`amper_id`) REFERENCES `amper_sys` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

/* Make sure the 'profile' id references to a valid profile with a fregin key 
   Make sure the inserted user has id as 1, as below, since it is used as a system user id */
INSERT INTO `amper`.`users_sys`
(`id`, `firstName`, `lastName`, `middleName`, `username`, `password`, `photo`, `profile`, `email`, `active`, `activationCode`, `amperId`, `state`) 
VALUES 
(1, 'Amper', 'Cloud', '', 'amper', '7be4c8a38c3efc72627a7a0f7cf3999218f1490da051795c0ec5fc309534bee078c864fc', '', 1, 'your@email.address', 1, 'code', 1, 1);
