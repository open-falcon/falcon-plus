# redis-go-cluster
redis-go-cluster is a golang implementation of redis client based on Gary Burd's
[Redigo](https://github.com/garyburd/redigo). It caches slot info at local and 
updates it automatically when cluster change. The client manages a connection pool 
for each node, uses goroutine to execute as concurrently as possible, which leads 
to its high efficiency and low lantency.

**Supported**:
* Most commands of keys, strings, lists, sets, sorted sets, hashes.
* MGET/MSET
* Pipelining

**NOT supported**:
* Cluster commands
* Pub/Sub
* Transaction
* Lua script

## Documentation
[API Reference](https://godoc.org/github.com/chasex/redis-go-cluster)

## Installation
Install redis-go-cluster with go tool:
```
    go get github.com/chasex/redis-go-cluster
```
    
## Usage
To use redis cluster, you need import the package and create a new cluster client
with an options:
```go
import "github.com/chasex/redis-go-cluster"

cluster, err := redis.NewCluster(
    &redis.Options{
	StartNodes: []string{"127.0.0.1:7000", "127.0.0.1:7001", "127.0.0.1:7002"},
	ConnTimeout: 50 * time.Millisecond,
	ReadTimeout: 50 * time.Millisecond,
	WriteTimeout: 50 * time.Millisecond,
	KeepAlive: 16,
	AliveTime: 60 * time.Second,
    })
```

### Basic
redis-go-cluster has compatible interface to [Redigo](https://github.com/garyburd/redigo), 
which uses a print-like API for all redis commands. When executing a command, it need a key 
to hash to a slot, then find the corresponding redis node. Do method will choose first
argument in args as the key, so commands which are independent from keys are not supported,
such as SYNC, BGSAVE, RANDOMKEY, etc. 

**RESTRICTION**: Please be sure the first argument in args is key.

See full redis commands: http://www.redis.io/commands

```go
cluster.Do("SET", "foo", "bar")
cluster.Do("INCR", "mycount", 1)
cluster.Do("LPUSH", "mylist", "foo", "bar")
cluster.Do("HMSET", "myhash", "f1", "foo", "f2", "bar")
```
You can use help functions to convert reply to int, float, string, etc.
```go
reply, err := Int(cluster.Do("INCR", "mycount", 1))
reply, err := String(cluster.Do("GET", "foo"))
reply, err := Strings(cluster.Do("LRANGE", "mylist", 0, -1))
reply, err := StringMap(cluster.Do("HGETALL", "myhash"))
```
Also, you can use Values and Scan to convert replies to multiple values with different types.
```go
_, err := cluster.Do("MSET", "key1", "foo", "key2", 1024, "key3", 3.14, "key4", "false")
reply, err := Values(cluster.Do("MGET", "key1", "key2", "key3", "key4"))
var val1 string
var val2 int
reply, err = Scan(reply, &val1, &val2)
var val3 float64
reply, err = Scan(reply, &val3)
var val4 bool
reply, err = Scan(reply, &val4)

```

### Multi-keys
Mutiple keys command - MGET/MSET are supported using result aggregation.
Processing steps are as follows:
- First, split the keys into multiple nodes according to their hash slot.
- Then, start a goroutine for each node to excute MGET/MSET commands and wait them finish.
- Last, collect and rerange all replies, return back to caller.

**NOTE**: Since the keys may spread across mutiple node, there's no atomicity gurantee that 
all keys will be set at once. It's possible that some keys are set while others are not.

### Pipelining
Pipelining is supported through the Batch interface. You can put multiple commands into a 
batch as long as it is supported by Do method. RunBatch will split these command to distinct
nodes and start a goroutine for each node. Commands hash to same nodes will be merged and sent 
using pipelining. After all commands done, it rearrange results as MGET/MSET do. Result is a 
slice of each command's reply, you can use Scan to convert them to other types.
```go
batch := cluster.NewBatch()
err = batch.Put("LPUSH", "country_list", "France")
err = batch.Put("LPUSH", "country_list", "Italy")
err = batch.Put("LPUSH", "country_list", "Germany")
err = batch.Put("INCRBY", "countries", 3)
err = batch.Put("LRANGE", "country_list", 0, -1)
reply, err = cluster.RunBatch(batch)

var resp int
for i := 0; i < 4; i++ {
    reply, err = redis.Scan(reply, &resp)    
}

countries, err := Strings(reply[0], nil)
```

## Contact
Bug reports and feature requests are welcome.
If you have any question, please email me wuxibin2012@gmail.com.

## License
redis-go-cluster is available under the [Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0.html).
