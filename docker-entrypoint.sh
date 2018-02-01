#!/bin/bash
set -e
declare -A ARGS;
ARGS=(
    ["AGENT_HTTP"]="0.0.0.0:${AGENT_HTTP_PORT:-1988}"
    ['AGGREGATOR_HTTP']="0.0.0.0:${AGGREGATOR_HTTP_PORT:-6055}"
    ['GRAPH_HTTP']="0.0.0.0:${GRAPH_HTTP_PORT:-6071}"
    ['GRAPH_RPC']="0.0.0.0:${GRAPH_RPC_PORT:-6070}"
    ['HBS_HTTP']="0.0.0.0:${HBS_HTTP_PORT:-6031}"
    ['HBS_RPC']="0.0.0.0:${HBS_RPC_PORT:-6030}"
    ['JUDGE_HTTP']="0.0.0.0:${JUDGE_HTTP_PORT:-6081}"
    ['JUDGE_RPC']="0.0.0.0:${JUDGE_RPC_PORT:-6080}"
    ['NODATA_HTTP']="0.0.0.0:${NODATA_HTTP_PORT:-6090}"
    ['TRANSFER_HTTP']="0.0.0.0:${TRANSFER_HTTP_PORT:-6060}"
    ['TRANSFER_RPC']="0.0.0.0:${TRANSFER_RPC_PORT:-8433}"
    ['REDIS']="${REDIS:-'redis://127.0.0.1:6379'}"
    ['MYSQL']="${MYSQL:-'root:@tcp(127.0.0.1:3306)'}"
    ['PLUS_API_DEFAULT_TOKEN']="${PLUS_API_DEFAULT_TOKEN:-'default-token-used-in-server-side'}"
    ['PLUS_API_HTTP']="0.0.0.0:${PLUS_API_HTTP_PORT:-8080}"
 )

configurer() {
    for key in ${!ARGS[*]}
    do  
        search="%%${key}%%"
        replace=${ARGS["$key"]}
        #echo "$search = $replace"
        sysname=$(uname)
        if [ "$sysname" == "Darwin" ] ; then
            # Note the "" and -e  after -i, needed in OS X
            find ./config/*.json -type f -exec sed -i .tpl -e "s#${search}#${replace}#g" {} \;
        else
            find ./config/*.json -type f -exec sed -i "s#${search}#${replace}#g" {} \;
        fi
        
    done
}

configurer


exec "$@"
