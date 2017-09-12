import requests
# Copyright 2017 Xiaomi, Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


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
