#!/bin/sh

pwddir=$(pwd)
rm -f /var/run/docker.pid
dockerd &
sleep 3s

for a in $(ls images)
do
    docker image load -i images/$a
done

cd $CHALDIR

for a in $(seq 0 1 $(($(yq '.services | length' docker-compose.yml) - 1)))
do
    if [ "$(yq ".services[.services | keys[$a]].image" docker-compose.yml)" == "null" ]
    then
        echo "Please setup image name." 1>&2
        exit 1
    fi
    if ! docker images --format "{{.Repository}}:{{.Tag}}" | grep "$(yq ".services[.services | keys[$a]].image" docker-compose.yml | sed 's/\./\\./g')"
    then
        docker compose build "$(yq ".services | keys[$a]" docker-compose.yml)"
    fi
done

cd $pwddir

for a in $(docker images --format "{{.Repository}}:{{.Tag}}")
do
    docker image save > images/$(echo "$a" | awk -F: '{print $1}' | sed 's/\//_/g').tar
done

./instancerapi
