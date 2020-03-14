#!/bin/sh

DOCKER_DIR=/open-falcon
of_bin=$DOCKER_DIR/open-falcon
DOCKER_HOST_IP=$(route -n | awk '/UG[ \t]/{print $2}')

# Search $1 and replace with $2 or $3(defualt)
replace() {
    target=$2
    if [ -z "$target" ]; then
      target=$3
    fi
    find $DOCKER_DIR/*/config/*.json -type f -exec sed -i "s/$1/$2/g" {} \;
}

replace "%%MYSQL%%" "$MYSQL_PORT" "$DOCKER_HOST_IP:3306"
replace "%%REDIS%%" "$REDIS_PORT" "$DOCKER_HOST_IP:6379"
replace "%%AGGREGATOR_HTTP%%" "$AGGREGATOR_HTTP" "0.0.0.0:6055"
replace "%%GRAPH_HTTP%%" "$GRAPH_HTTP" "0.0.0.0:6071"
replace "%%GRAPH_RPC%%" "$GRAPH_RPC" "0.0.0.0:6070"
replace "%%HBS_HTTP%%" "$HBS_HTTP" "0.0.0.0:6031"
replace "%%HBS_RPC%%" "$HBS_RPC" "0.0.0.0:6030"
replace "%%JUDGE_HTTP%%" "$JUDGE_HTTP" "0.0.0.0:6081"
replace "%%JUDGE_RPC%%" "$JUDGE_RPC" "0.0.0.0:6080"
replace "%%NODATA_HTTP%%" "$NODATA_HTTP" "0.0.0.0:6090"
replace "%%TRANSFER_HTTP%%" "$TRANSFER_HTTP" "0.0.0.0:6060"
replace "%%TRANSFER_RPC%%" "$TRANSFER_RPC" "0.0.0.0:8433"
replace "%%PLUS_API_HTTP%%" "$PLUS_API_HTTP" "0.0.0.0:8080"
replace "%%AGENT_HOSTNAME%%" "$AGENT_HOSTNAME" ""

#use absolute path of metric_list_file in docker
TAB=$'\t'; sed -i "s|.*metric_list_file.*|${TAB}\"metric_list_file\": \"$DOCKER_DIR/api/data/metric\",|g" $DOCKER_DIR/api/config/*.json;

action=$1
module_name=$2
case $action in
 run)
        $DOCKER_DIR/"$module_name"/bin/falcon-"$module_name" -c /open-falcon/"$module_name"/config/cfg.json
        ;;
 *)
        supervisorctl $*
        ;;
esac
