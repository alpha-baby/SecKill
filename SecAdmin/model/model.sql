CREATE TABLE `product` (
  `id` int(11) NOT NULL auto_increment,
  `name` varchar(1024) NOT NULL,
  `total` int(11) DEFAULT 0,
  `status` int(11) DEFAULT 0,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

ALTER TABLE product ADD UNIQUE(`name`);
alter table activity add column buy_limit int default 1;
alter table activity add column sec_speed int default 100;
alter table activity add column buy_rate float default 1.0;


create table `activity` (
`id` int(11) not null auto_increment,
`name` varchar(1024) not null ,
`product_id` int(11) default 0,
`start_time` int(11) default 0,
`end_time` int(11) default 0,
`total` int(11) default 0,
`status` int(11) default 0,
primary key (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;