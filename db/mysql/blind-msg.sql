CREATE DATABASE blind_msg;

use blind_msg

CREATE TABLE `bac` (
  `id`         bigint(20) NOT NULL AUTO_INCREMENT,
  `service`    int NOT NULL,             -- 1: LINE
  `user_id`    varchar(128) NOT NULL,    -- user_id
  `bac`        varchar(32) NOT NULL,     -- blind access code
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_bac1` (`service`, `user_id`),
  UNIQUE KEY `uk_bac2` (`bac`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
