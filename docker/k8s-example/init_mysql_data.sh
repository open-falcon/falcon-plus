#!/bin/sh

mysql_pod=$(kubectl get pods | grep mysql | awk '{print $1}')
cd /tmp && \
	git clone --depth=1 https://github.com/open-falcon/falcon-plus && \
	cd /tmp/falcon-plus/ && \
	for x in `ls ./scripts/mysql/db_schema/*.sql`; do
	    echo init mysql table $x ...;
	    kubectl exec -it $mysql_pod -- mysql -uroot -p123456 < $x;
	done

rm -rf /tmp/falcon-plus/