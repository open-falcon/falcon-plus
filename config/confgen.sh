#!/bin/bash

declare -A confs
confs=(
    [%%QUERY_HTTP%%]=0.0.0.0:9966
    [%%HBS_HTTP%%]=0.0.0.0:6031
    [%%TRANSFER_HTTP%%]=0.0.0.0:6060
    [%%GRAPH_HTTP%%]=0.0.0.0:6071
    [%%HBS_RPC%%]=0.0.0.0:6030
    [%%REDIS%%]=127.0.0.1:6379
    [%%GRAPH_RPC%%]=0.0.0.0:6070
    [%%TRANSFER_RPC%%]=0.0.0.0:8433
    [%%MYSQL%%]="root:password@tcp(127.0.0.1:3306)"
)

configurer() {
    for i in "${!confs[@]}"
    do
        search=$i
        replace=${confs[$i]}
        # Note the "" after -i, needed in OS X
        find ./out/*/config/*.json -type f -exec sed -i "s/${search}/${replace}/g" {} \;
    done
}
configurer
