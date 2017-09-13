CREATE DATABASE IF NOT EXISTS `broker`;

USE `broker`;

DROP TABLE IF EXISTS `clients`;
CREATE TABLE `clients` (
    `id` int NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `login` varchar(300) NOT NULL,
    `password` varchar(300) NOT NULL,
    `balance` float NOT NULL
);

INSERT INTO `clients` (`id`, `login`,  `password`, `balance`)
    VALUES (1, 'Vasily', '123456', 200000),
    (2, 'Ivan', 'qwerty', 200000),
    (3, 'Olga', '1qaz2wsx', 200000);

DROP TABLE IF EXISTS `positions`;
CREATE TABLE `positions` (
    `id` int NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `user_id` int NOT NULL,
    `ticker` varchar(300) NOT NULL,
    `amount` int NOT NULL,
    `price` float NOT NULL,
    `type` varchar(300) NOT NULL,
    `status` int NOT NULL,
    KEY user_id(user_id)
);

INSERT INTO `positions` (`user_id`, `ticker`, `amount`, `price`, `type`, `status`)
    VALUES (1, 'SIM7', 200000, 1000.0, 'buy', 0),
    (1, 'RIM7', 200000, 1000.0, 'sell', 0),
    (2, 'RIM7', 200000, 1000.0, 'buy', 0);

DROP TABLE IF EXISTS `orders_history`;
CREATE TABLE `orders_history` (
    `id` int NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `time` int NOT NULL,
    `user_id` int,
    `ticker` varchar(300) NOT NULL,
    `amount` int NOT NULL,
    `price` float NOT NULL,
    `type` varchar(300) NOT NULL,
    KEY user_id(user_id)
);
