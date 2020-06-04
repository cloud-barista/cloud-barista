#!/bin/bash

HOST=${1-"localhost"}
CMD=${2-""}



	if [ "${CMD}" == "install" ]
		then
			echo "[Call Install]"
			curl -sX GET http://${HOST}:1324/milkyway/install | json_pp || return 1
			echo "#-----------------------------"
		else
			echo "[Call without Install]"
	fi

	echo "####################################################################"

	curl -sX GET http://${HOST}:1324/milkyway/init | json_pp || return 1
	echo "#-----------------------------"

	curl -sX GET http://${HOST}:1324/milkyway/cpus | json_pp || return 1
	echo "#-----------------------------"
	curl -sX GET http://${HOST}:1324/milkyway/cpum | json_pp || return 1
	echo "#-----------------------------"
	curl -sX GET http://${HOST}:1324/milkyway/memR | json_pp || return 1
	echo "#-----------------------------"
	curl -sX GET http://${HOST}:1324/milkyway/memW | json_pp || return 1
	echo "#-----------------------------"
	curl -sX GET http://${HOST}:1324/milkyway/fioR | json_pp || return 1
	echo "#-----------------------------"
	curl -sX GET http://${HOST}:1324/milkyway/fioW | json_pp || return 1
	echo "#-----------------------------"
	curl -sX GET http://${HOST}:1324/milkyway/dbR | json_pp || return 1
	echo "#-----------------------------"
	curl -sX GET http://${HOST}:1324/milkyway/dbW | json_pp || return 1
	echo "#-----------------------------"

	curl -sX GET http://${HOST}:1324/milkyway/rtt -H 'Content-Type: application/json' -d '{ "host": "localhost"}' |json_pp || return 1
	echo "#-----------------------------"
	curl -sX GET http://${HOST}:1324/milkyway/mrtt -H 'Content-Type: application/json' -d '{ "multihost": [{"host":"localhost", "spec":"spec1"},{"host":"localhost", "spec":"spec2"}]}' |json_pp || return 1
	echo "#-----------------------------"

	curl -sX GET http://${HOST}:1324/milkyway/clean | json_pp || return 1
	echo "#-----------------------------"

	

