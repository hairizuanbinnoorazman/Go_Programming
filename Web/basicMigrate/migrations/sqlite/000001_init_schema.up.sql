CREATE TABLE IF NOT EXISTS `users` (
  `id` INTEGER PRIMARY KEY AUTOINCREMENT,      
  `first_name` varchar(100) NOT NULL default '',
  `last_name` varchar(100) NOT NULL default ''
);