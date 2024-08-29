CREATE DATABASE IF NOT EXISTS `stori` /*!40100 DEFAULT CHARACTER SET utf8 */;

use stori;

CREATE TABLE IF NOT EXISTS `migration` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `csv_path` varchar(255) NOT NULL,
    `status` varchar(255) NOT NULL,
    `total_lines` int(11) NOT NULL,
    `processed_lines` int(11) NOT NULL,
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `csv_path` (`csv_path`)
);

CREATE TABLE IF NOT EXISTS `transaction` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `user_id` int(11) NOT NULL,
    `amount` decimal(10, 2) NOT NULL,
    `date_time` datetime NOT NULL,
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `user_id` (`user_id`),
    KEY `date_time` (`date_time`)
);