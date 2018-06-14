#coding:utf8
import sys
import random
import time
import json
import socket

import requests


args= sys.argv[1:]
if args:
    host = args[0]
else:
    host = "127.0.0.1:1988"

errors = (
        "warn out of memeory",
        "Error 1086",
        )

ts = int(time.time())
payload = [
        {
            "endpoint": socket.gethostname(),
            "metric": "str.match",
            "timestamp": ts,
            "step": 30,
            "value": random.choice(errors),
            "counterType": "STRMATCH",
            "tags": "project=test,platform=linux",
            },
        ]

url = "http://{host}/v1/push".format(host=host)

r = requests.post(url=url, data=json.dumps(payload))
resp_code = r.status_code
body = r.content

print "host", host
print "payload", payload
print "resp_code", resp_code
print "body", body

