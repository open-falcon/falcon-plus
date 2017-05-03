#!/bin/bash

declare -A confs
confs=(
    [%%AGENT_HTTP%%]=0.0.0.0:1988
    [%%AGGREGATOR_HTTP%%]=0.0.0.0:6055
    [%%GRAPH_HTTP%%]=0.0.0.0:6071
    [%%GRAPH_RPC%%]=0.0.0.0:6070
    [%%HBS_HTTP%%]=0.0.0.0:6031
    [%%HBS_RPC%%]=0.0.0.0:6030
    [%%JUDGE_HTTP%%]=0.0.0.0:6081
    [%%JUDGE_RPC%%]=0.0.0.0:6080
    [%%NODATA_HTTP%%]=0.0.0.0:6090
    [%%TRANSFER_HTTP%%]=0.0.0.0:6060
    [%%TRANSFER_RPC%%]=0.0.0.0:8433
    [%%REDIS%%]=127.0.0.1:6379
    [%%MYSQL%%]="root:@tcp(127.0.0.1:3306)"
    [%%PLUS_API_DEFAULT_TOKEN%%]="default-token-used-in-server-side"
)

configurer() {
    for i in "${!confs[@]}"
    do
        search=$i
        replace=${confs[$i]}

        if [ "$UNAME" == "Darwin" ] ; then
            # Note the "" and -e  after -i, needed in OS X
            find ./out/*/config/*.json -type f -exec sed -i .tpl -e "s/${search}/${replace}/g" {} \;
        else
            find ./out/*/config/*.json -type f -exec sed -i "s/${search}/${replace}/g" {} \;
        fi
    done
}
configurer
