CREATE TABLE IF NOT EXISTS `users` (
  `id` int(11) NOT NULL auto_increment,      
  `first_name` varchar(100) NOT NULL default '',
  `last_name` varchar(100) NOT NULL default '',
   PRIMARY KEY  (`id`)
);