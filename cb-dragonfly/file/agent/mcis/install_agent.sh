#!/usr/bin/env bash

#######################################
#   Monitoring Agent Install Script   #
#                                     #
#   [ innogrid ]                      #
#   Monitoring Agent Install Script   #
#                                     #
#######################################

# API CONFIG
AGENT_REPO="http://{{api_server}}"
AGENT_CONF_PATH="/etc/telegraf"

# VM INFORMATION
os_type=$(hostnamectl | grep 'Operating System' | awk '{print $3}')
os_type=$(echo $os_type | tr 'a-z' 'A-Z')
arch=$(arch)

# VARIABLE
declare -a dpkg=`which dpkg`
declare -a rpm=`which rpm`
declare -a yum=`which yum`
declare -a wget=`which wget`
declare -a awk=`which awk`
declare -a sed=`which sed`
declare -a date=`which date`
declare -a systemctl=`which systemctl`
declare -a echo=`which echo`
declare -a clear=`which clear`
declare -a cat=`which cat`
declare -a grep=`which grep`
declare -a cut=`which cut`
declare -a ip=`which ip`
declare -a wc=`which wc`
declare -a head=`which head`
declare -a rev=`which rev`
declare -a sudo=`which sudo`
declare -a expr=`which expr`
declare -a mkdir=`which mkdir`
declare -a rm=`which rm`

setup_agent()
{
        echo "start setup_agent()"

        # install agent
        if [ "$os_type" == "UBUNTU" ]; then
                $sudo $wget -O cb-agent.deb "$AGENT_REPO/mon/file/agent/pkg?osType=$os_type&arch=$arch"
                $sudo $dpkg -i cb-agent.deb > /dev/null 2>&1
                $sudo $rm cb-agent.deb
        elif [ "$os_type" == "CENTOS" ]; then
                $sudo $wget -O cb-agent.rpm "$AGENT_REPO/mon/file/agent/pkg?osType=$os_type&arch=$arch"
                $sudo $rpm -ivh cb-agent.rpm > /dev/null 2>&1
                $sudo $rm cb-agent.rpm
        fi

        # install iostat tools
        if [ "$os_type" == "UBUNTU" ]; then
                $sudo apt-get install sysstat -y
        elif [ "$os_type" == "CENTOS" ]; then
                $sudo yum install sysstat -y
        fi
}

setup_config()
{
        echo "start setup_config()"

        $sudo $rm $AGENT_CONF_PATH/telegraf.conf
        $sudo $wget -O $AGENT_CONF_PATH/"telegraf.conf" "$AGENT_REPO/mon/file/agent/conf?mcis_id={{mcis_id}}&vm_id={{vm_id}}"
}

start_agent()
{
        echo "start start_agent()"

        $sudo $systemctl enable telegraf
        $sudo $systemctl start telegraf
}

step()
{
        setup_agent
        setup_config
        start_agent
}

step

