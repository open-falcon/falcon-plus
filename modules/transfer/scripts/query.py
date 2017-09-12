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

url = "http://localhost:3060/graph/history"
r = requests.post(url, data=json.dumps(d))
print r.text

#curl "http://query.falcon.miliao.srv:9966/graph/history/one?cf=AVERAGE&endpoint=`hostname`&start=`date -d '1 hours ago' +%s`&counter=load.1min" |python -m json.tool
