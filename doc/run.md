# centos7部署golang

### centos7 开放端口

```
firewall-cmd --zone=public --add-port=22/tcp --permanent
systemctl restart firewalld.service

查看端口
netstat -ntlp
```


### 安装mysql

```
curl -LO http://dev.mysql.com/get/mysql57-community-release-el7-11.noarch.rpm
sudo yum localinstall mysql57-community-release-el7-11.noarch.rpm
sudo yum install mysql-community-server
sudo systemctl enable mysqld
sudo systemctl start mysqld
sudo systemctl status mysqld
grep 'temporary password' /var/log/mysqld.log//查看默认密码
mysql -u root -p
set global validate_password_policy=0;
set global validate_password_length=4;
ALTER USER 'root'@'localhost' IDENTIFIED BY '123456';//修改密码

sudo firewall-cmd --zone=public --add-port=3306/tcp --permanent
sudo firewall-cmd --reload

//添加允许远程连接账户
GRANT ALL PRIVILEGES ON *.* TO 'admin'@'%' IDENTIFIED BY 'secret' WITH GRANT OPTION;

//修改编码 vi /etc/my.cnf
[mysqld]
# 在myslqd下添加如下键值对
character_set_server=utf8
init_connect='SET NAMES utf8'
```

### 安装redis

```
sudo yum update
yum install epel-release
yum install redis
sudo systemctl start redis

redis-cli>
config set stop-writes-on-bgsave-error no

//允许远程访问
/etc/redis.config
bind 127.0.0.1 改为 bind 0.0.0.0;


//允许端口访问
firewall-cmd --add-port=6379/tcp --permanent
firewall-cmd --reload
```

### 编译脚本

```
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
scp ./Journal root@140.82.56.114:/root
scp -r ./conf root@140.82.56.114:/root
```



### supervisor管理golang程序


+ sudo yum install python-setuptools
+ sudo easy_install supervisor
+ sudo `echo_supervisord_conf` > /etc/supervisord.conf
+ /etc/supervisord.conf 增加配置

```
[program:golang-http-server]
command=/home/golang/simple_http_server  //golang程序
autostart=true
autorestart=true
startsecs=10
stdout_logfile=/var/log/simple_http_server.log
stdout_logfile_maxbytes=1MB
stdout_logfile_backups=10
stdout_capture_maxbytes=1MB
stderr_logfile=/var/log/simple_http_server.log
stderr_logfile_maxbytes=1MB
stderr_logfile_backups=10
stderr_capture_maxbytes=1MB
```
+ 启动supervisor 

```
sudo /usr/bin/supervisord -c /etc/supervisord.conf
```

如果出现什么问题，可以查看日志进行分析，日志文件路径/tmp/supervisord.log


tips：如果修改了配置文件，可以用kill -HUP重新加载配置文件

```
cat /tmp/supervisord.pid | xargs sudo kill -HUP
```

+ 查看supervisor运行状态

```
supervisorctl
```

+ 关闭管理进程

```
supervisorctl stop all
```

### nginx 部署