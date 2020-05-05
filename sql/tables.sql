create table tbl_app_info (
  app_id int auto_increment primary key,
  app_name varchar(1024) not null, 
  app_type varchar(64) not null, 
  create_time TIMESTAMP default current_timestamp,
  develop_path varchar(256) not null
) engine=innodb default charset=utf8 auto_increment=1;

create table tbl_app_ip (
  app_id int, 
  ip varchar(64),
  Key app_id_ip_index (app_id, ip)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 AUTO_INCREMENT=1;


create table tbl_log_info (
  log_id int auto_increment primary key,
  app_id varchar(1024) not null, 
  log_path varchar(64) not null, 
  topic varchar(1024) not null,
  create_time TIMESTAMP default current_timestamp,
  status tinyint default 1
) engine=innodb default charset=utf8 auto_increment=1;
