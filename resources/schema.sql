DROP TABLE `member` IF EXISTS;
DROP TABLE `household` IF EXISTS;

CREATE TABLE `household` (
                             `id` int unsigned NOT NULL AUTO_INCREMENT,
                             `type` varchar(20) NOT NULL,
                             PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `member` (
                          `id` int unsigned NOT NULL AUTO_INCREMENT,
                          `name` varchar(128) NOT NULL,
                          `gender` varchar(20) NOT NULL,
                          `marital_status` varchar(20) NOT NULL,
                          `spouse_id` int unsigned DEFAULT NULL,
                          `occupation_type` varchar(128) NOT NULL,
                          `annual_income` double unsigned NOT NULL,
                          `dob` datetime NOT NULL,
                          `household_id` int unsigned NOT NULL,
                          PRIMARY KEY (`id`),
                          KEY `member_spouse_fk_1` (`spouse_id`),
                          KEY `member_household_fk_2` (`household_id`),
                          CONSTRAINT `member_household_fk_2` FOREIGN KEY (`household_id`) REFERENCES `household` (`id`),
                          CONSTRAINT `member_spouse_fk_1` FOREIGN KEY (`spouse_id`) REFERENCES `member` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;