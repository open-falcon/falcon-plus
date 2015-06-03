import requests
import json

d = [
            {
                "endpoint": "host1",
                "counter": "load.15min",
            },
            {
                "endpoint": "host1",
                "counter": "net.if.in.bytes/iface=eth0",
            },
            {
                "endpoint": "127.0.0.1:8080",
                "counter": "qps",
            },
]

url = "http://127.0.0.1:9966/graph/info"
r = requests.post(url, data=json.dumps(d))
print r.text
