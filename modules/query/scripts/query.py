import requests
import time
import json

end = int(time.time())
start = end - 3600

d = {
        "start": start,
        "end": end,
        "cf": "AVERAGE",
        "endpoint_counters": [
            {
                "endpoint": "lg-op-mon-onebox01.bj",
                "counter": "load.1min",
            },
            {
                "endpoint": "lg-op-mon-onebox01.bj",
                "counter": "load.5min",
            },
            {
                "endpoint": "lg-op-mon-onebox01.bj",
                "counter": "load.15min",
            },
        ],
}

url = "http://localhost:19966/graph/history"
r = requests.post(url, data=json.dumps(d))
print r.text

#curl "http://query.falcon.miliao.srv:9966/graph/history/one?cf=AVERAGE&endpoint=`hostname`&start=`date -d '1 hours ago' +%s`&counter=load.1min" |python -m json.tool
