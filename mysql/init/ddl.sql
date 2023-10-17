CREATE SCHEMA IF NOT EXISTS `database` DEFAULT CHARACTER SET utf8mb4;
USE `database`;

SET CHARSET utf8mb4;

CREATE TABLE IF NOT EXISTS `database`.`user_scores` (
  `id` VARCHAR(128) NOT NULL COMMENT 'ユーザID',
  `score` INT UNSIGNED NOT NULL COMMENT 'スコア',
  PRIMARY KEY (`id`),
  INDEX `idx_score` (`score` DESC))
ENGINE = InnoDB
COMMENT = 'ユーザスコア';
