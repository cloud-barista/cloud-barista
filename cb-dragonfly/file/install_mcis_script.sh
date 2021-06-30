
#!/bin/bash

echo "[MCIS-Agent: Start to prepare a VM evaluation]"

echo "[MCIS-Agent: Install sysbench]"
sudo apt-get -y update
sudo apt-get -y install sysbench

echo "[MCIS-Agent: Install Ping]"
sudo apt-get -y install iputils-ping

echo "[MCIS-Agent: Install debconf-utils]"
sudo apt-get -y install debconf-utils
sudo debconf-set-selections <<< 'mysql-server mysql-server/root_password password psetri1234ak'
sudo debconf-set-selections <<< 'mysql-server mysql-server/root_password_again password psetri1234ak'

echo "[MCIS-Agent: Install MySQL]"
sudo DEBIAN_FRONTEND=noninteractive apt-get -y install mysql-server

echo "[MCIS-Agent: Generate dump tables for evaluation]"

mysql -u root -ppsetri1234ak -e "CREATE DATABASE sysbench;"
mysql -u root -ppsetri1234ak -e "CREATE USER 'sysbench'@'localhost' IDENTIFIED BY 'psetri1234ak';"
mysql -u root -ppsetri1234ak -e "GRANT ALL PRIVILEGES ON *.* TO 'sysbench'@'localhost' IDENTIFIED  BY 'psetri1234ak';"

echo "[MCIS-Agent: Preparation is done]"





