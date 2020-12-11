
#!/bin/bash

echo "[CB-Milkyway: Start to Delete Milkyway]"

echo "[CB-Milkyway: UnInstall sysbench]"
sudo apt-get purge -y update
sudo apt-get purge -y install sysbench

echo "[CB-Milkyway: UnInstall Ping]"
sudo apt-get purge -y iputils-ping

echo "[CB-Milkyway: UnInstall debconf-utils]"
sudo apt-get purge -y debconf-utils
sudo debconf-set-selections <<< 'mysql-server mysql-server/root_password password psetri1234ak'
sudo debconf-set-selections <<< 'mysql-server mysql-server/root_password_again password psetri1234ak'

echo "[CB-Milkyway: UnInstall MySQL]"
sudo DEBIAN_FRONTEND=noninteractive apt-get -y purge mysql-server

echo "[CB-Milkyway: Generate dump tables for evaluation]"

mysql -u root -ppsetri1234ak -e "CREATE DATABASE sysbench;"
mysql -u root -ppsetri1234ak -e "CREATE USER 'sysbench'@'localhost' IDENTIFIED BY 'psetri1234ak';"
mysql -u root -ppsetri1234ak -e "GRANT ALL PRIVILEGES ON *.* TO 'sysbench'@'localhost' IDENTIFIED  BY 'psetri1234ak';"

echo "[CB-Milkyway: Deletion is done]"








