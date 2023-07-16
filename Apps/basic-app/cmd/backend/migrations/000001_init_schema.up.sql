USE `basic`;

CREATE TABLE IF NOT EXISTS `users` (
  `id` varchar(40) NOT NULL,
  `email` varchar(250) NOT NULL default '',
  `password` varchar(250) default '',
  `name` varchar(250) default '',
  `forget_password_token` varchar(40) default '',
  `forget_password_expiry_date` datetime NOT NULL default '1970-01-01',
  `activation_token` varchar(40) default '',
  `activation_expiry_date` datetime NOT NULL default '1970-01-01',
  `activated` boolean,
  `date_created` datetime NOT NULL default '1970-01-01',
  `date_modified` datetime NOT NULL default '1970-01-01',
   PRIMARY KEY (`id`)
);