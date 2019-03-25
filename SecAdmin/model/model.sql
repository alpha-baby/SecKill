CREATE TABLE `product` (
  `id` int(11) NOT NULL auto_increment,
  `name` varchar(1024) NOT NULL,
  `total` int(11) DEFAULT 0,
  `status` int(11) DEFAULT 0,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8

