#!/bin/bash

yum install -y yum-utils device-mapper-persistent-data lvm2
yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
yum install docker* -y
systemctl start docker
systemctl enable docker
docker run -d --name mysql -e MYSQL_ROOT_PASSWORD=nxera1688. -p 33306:3306 mysql
yum install mysql -y

clear
echo "install successful"
echo 'You can use "mysql -h 127.0.0.1 -P 33306 -u root -pnxera1688." to login and config Mysql'
echo "configuration follow below"
echo "mysql>
          -> CREATE DATABASE kubeStone;
          -> CREATE TABLE server (
          -> id INT AUTO_INCREMENT PRIMARY KEY,
          -> name VARCHAR(255) NOT NULL,
          -> ip VARCHAR(15) NOT NULL UNIQUE,
          -> port VARCHAR(15) NOT NULL,
          -> user VARCHAR(255) NOT NULL,
          -> password VARCHAR(255) NOT NULL
          -> );"

