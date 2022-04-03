#!/bin/bash

docker stop ratd_mariadb
docker rm ratd_mariadb

docker run --detach --network ratd_default -p 3308:3306 --name ratd_mariadb \
--env MARIADB_ROOT_PASSWORD=${RATD_MYSQL_ROOT_PASS} \
--env MARIADB_USER=${RATD_MYSQL_USER} \
--env MARIADB_PASSWORD=${RATD_MYSQL_PASS} \
mariadb:latest

#docker network connect bridge ratd_default
