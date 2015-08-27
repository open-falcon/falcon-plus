import requests
import json

d = [
            {
                "endpoint": "hh-op-mon-tran01.bj",
                "counter": "load.15min",
            },
            {
                "endpoint": "hh-op-mon-tran01.bj",
                "counter": "net.if.in.bytes/iface=eth0",
            },
            {
                "endpoint": "10.202.31.14:7934",
                "counter": "p2-com.xiaomi.miui.mibi.service.MibiService-method-createTradeV1",
            },
]

url = "http://query.falcon.miliao.srv:9966/graph/info"
r = requests.post(url, data=json.dumps(d))
print r.text

#curl "localhost:9966/graph/info/one?endpoint=`hostname`&counter=load.1min" |python -m json.tool
