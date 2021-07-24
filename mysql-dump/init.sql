
DROP TABLE IF EXISTS `db_store`.`product`;
CREATE TABLE `db_store`.`product` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `price` bigint(20) NOT NULL DEFAULT 0,
  `stock` bigint(20) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
INSERT INTO product
(name, price, stock, created_at, updated_at)
VALUES
('Permen Minyak', 2000, 50, current_timestamp(), current_timestamp()),
('Keyboard Wireless', 3000, 50, current_timestamp(), current_timestamp());


DROP TABLE IF EXISTS `db_store`.`register_payment`;
CREATE TABLE `db_store`.`register_payment` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) NOT NULL,
  `cart_id` bigint(20) NOT NULL,
  `product` text COLLATE utf8mb4_unicode_ci NOT NULL,
  `status` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'pending' COMMENT 'pending; done; failed;',
  `created_at` timestamp NULL DEFAULT current_timestamp(),
  `expired_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;