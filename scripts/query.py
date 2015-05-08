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
                "endpoint": "host1",
                "counter": "load.1min",
            },
            {
                "endpoint": "host1",
                "counter": "load.5min",
            },
            {
                "endpoint": "host1",
                "counter": "load.15min",
            },
        ],
}

url = "http://127.0.0.1:9966/graph/history"
r = requests.post(url, data=json.dumps(d))
print r.text
