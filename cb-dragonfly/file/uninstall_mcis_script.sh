
#!/bin/bash

echo "[MCIS-Agent: Start to Delete Milkyway]"

echo "[MCIS-Agent: UnInstall sysbench]"
sudo apt-get purge -y update
sudo apt-get purge -y sysbench

echo "[MCIS-Agent: UnInstall Ping]"
sudo apt-get purge -y iputils-ping

echo "[MCIS-Agent: UnInstall debconf-utils]"
sudo apt-get purge -y debconf-utils
sudo debconf-set-selections <<< 'mysql-server mysql-server/root_password password psetri1234ak'
sudo debconf-set-selections <<< 'mysql-server mysql-server/root_password_again password psetri1234ak'

echo "[MCIS-Agent: UnInstall MySQL]"
sudo DEBIAN_FRONTEND=noninteractive apt-get -y purge mysql-server

echo "[MCIS-Agent: Generate dump tables for evaluation]"

mysql -u root -ppsetri1234ak -e "CREATE DATABASE sysbench;"
mysql -u root -ppsetri1234ak -e "CREATE USER 'sysbench'@'localhost' IDENTIFIED BY 'psetri1234ak';"
mysql -u root -ppsetri1234ak -e "GRANT ALL PRIVILEGES ON *.* TO 'sysbench'@'localhost' IDENTIFIED  BY 'psetri1234ak';"

echo "[MCIS-Agent: Deletion is done]"








