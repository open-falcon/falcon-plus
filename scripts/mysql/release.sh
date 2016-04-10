#!/bin/bash

WORKSPACE=${GOPATH%%/}/src/github.com/open-falcon

output_dir=$WORKSPACE/output
tmp_dir=$WORKSPACE/tmp
rm -rf $output_dir $tmp_dir &> /dev/null
mkdir -p $WORKSPACE $output_dir $tmp_dir

echo "working at" $WORKSPACE
echo "output dir is" $output_dir
echo

gitorg="https://github.com/open-falcon"

pre_components=(
https://github.com/open-falcon/common,$GOPATH/src/github.com/open-falcon/common
https://github.com/open-falcon/rrdlite,$GOPATH/src/github.com/open-falcon/rrdlite
)
for c in ${pre_components[@]};do
    repo=`echo -n $c | awk -F ',' '{print $1}'`
    target_dir=`echo -n $c | awk -F',' '{print $2}'`

    if [ -d "$target_dir" ];then
        echo "pre pull $repo ..."
        ( cd $target_dir && git pull origin master &> /dev/null)
    else
        echo "pre clone $repo ..."
        git clone $repo $target_dir &> /dev/null
    fi
done


graph_components=("agent" "transfer" "graph" "query" "dashboard")
judge_components=("portal" "fe" "hbs" "judge" "sender" "alarm" "links")
other_components=("task" "gateway" "nodata" "aggregator")

all_components=""
for c in ${graph_components[@]};do
    all_components+="$c "
done
for c in ${judge_components[@]};do
    all_components+="$c "
done
for c in ${other_components[@]};do
    all_components+="$c "
done

for c in `echo -n $all_components | tr ' ' '\n'`;do
    repo=${gitorg}/${c}.git

    cd $WORKSPACE
    if [ -d "$c" ];then
        echo "pull $repo ..."
        cd $WORKSPACE/$c && git reset --hard HEAD; git pull origin master &> /dev/null
    else
        echo "clone $repo ..."
        cd $WORKSPACE && git clone $repo &> /dev/null
    fi
done

for c in `echo -n $all_components | tr ' ' '\n'`;do
    echo "======== packing" $c "..."
    apptar=falcon-${c}-*.tar.gz
    cd $WORKSPACE/$c && rm -rf *.tar.gz ; go get &> /dev/null; bash control pack &>/dev/null; mv $apptar $output_dir

done

OF_RELEASE_VERSION=v0.1.0
cd $output_dir && tar -zcf $WORKSPACE/of-release-$OF_RELEASE_VERSION.tar.gz *.tar.gz &> /dev/null
echo
echo "--> $WORKSPACE/of-release-$OF_RELEASE_VERSION.tar.gz"
echo
