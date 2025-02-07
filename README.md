# Version

### Latest

# Install

### NOTE

##### CMD

- install : /agent/path/DHNCenter install

- remove : /agent/path/DHNCenter remove

- start : systemctl start DHNCenter

- stop : systemctl stop DHNCenter

- restart : systemctl restart DHNCenter

- status : systemctl status DHNCenter

##### DB

- mariadb
  
- my.cnf
```ini
#
# This group is read both by the client and the server
# use it for options that affect everything
#
[client-server]

#
# include *.cnf from the config directory
#

[client]
ssl=0
port	= 3306
socket 	= /mariadb/data/mysql.sock

[mysqld]
ssl=0
port	= 3306
datadir	= /mariadb/data
socket	= /mariadb/data/mysql.sock
event_scheduler = ON

bind-address = 0.0.0.0

symbolic-links=0

max_allowed_packet=1024M
open_files_limit=20000
net_read_timeout=3600
net_write_timeout=3600

join_buffer_size = 5M
sort_buffer_size = 5M
read_buffer_size = 2M
read_rnd_buffer_size = 16M
thread_cache_size = 300

slow-query-log=1
slow-query-log-file=/mariadb/data/slow_query.log
long_query_time=5

skip-external-locking

max_connections = 1432
binlog_format = mixed
server-id = 1

expire_logs_days = 3
innodb_buffer_pool_size = 64G
innodb_file_per_table=1
innodb_flush_log_at_trx_commit = 0
innodb_lock_wait_timeout = 300
innodb_log_file_size = 1024M
innodb_log_buffer_size = 16M
innodb_write_io_threads = 32
innodb_read_io_threads = 32
innodb_io_capacity = 2000
innodb_flush_neighbors = 1

innodb_buffer_pool_instances = 8
innodb_page_cleaners = 8
innodb_autoinc_lock_mod = 2

sql_mode = NO_ENGINE_SUBSTITUTION
transaction_isolation = READ-COMMITTED


[mysqldump]
quick
max_allowed_packet = 16M

[mysql]
no-auto-rehash

[myismchk]
key_buffer_size = 128M
sort_buffer_size = 128M
read_buffer = 2M
write_buffer = 2M
```

# TEST

### 알림톡 500건

- Center : 0.1s

- Server : 0.4 ~ 10.0+-2s
