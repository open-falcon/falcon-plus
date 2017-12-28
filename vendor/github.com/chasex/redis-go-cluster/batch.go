// Copyright 2015 Joel Wu
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package redis

import (
    "fmt"
)

// Batch pack multiple commands, which should be supported by Do method.
type Batch struct {
    cluster	*Cluster
    batches	[]nodeBatch
    index	[]int
}

type nodeBatch struct {
    node    *redisNode
    cmds    []nodeCommand

    err	    error
    done    chan int
}

type nodeCommand struct {
    cmd	    string
    args    []interface{}
    reply   interface{}
    err	    error
}

// NewBatch create a new batch to pack mutiple commands.
func (cluster *Cluster) NewBatch() *Batch {
    return &Batch{
	cluster: cluster,
	batches: make([]nodeBatch, 0),
	index: make([]int, 0),
    }
}

// Put add a redis command to batch, DO NOT put MGET/MSET/MSETNX.
func (batch *Batch) Put(cmd string, args ...interface{}) error {
    if len(args) < 1 {
	return fmt.Errorf("Put: no key found in args")
    }

    if cmd == "MGET" || cmd == "MSET" || cmd == "MSETNX" {
	return fmt.Errorf("Put: %s not supported", cmd)
    }

    node, err := batch.cluster.getNodeByKey(args[0])
    if err != nil {
	return fmt.Errorf("Put: %v", err)
    }

    var i int
    for i = 0; i < len(batch.batches); i++ {
	if batch.batches[i].node == node {
	    batch.batches[i].cmds = append(batch.batches[i].cmds,
		nodeCommand{cmd: cmd, args: args})

	    batch.index = append(batch.index, i)
	    break
	}
    }

    if i == len(batch.batches) {
	batch.batches = append(batch.batches,
	    nodeBatch{
		node: node,
		cmds: []nodeCommand{{cmd: cmd, args: args}},
		done: make(chan int)})
	batch.index = append(batch.index, i)
    }

    return nil
}

// RunBatch execute commands in batch simutaneously. If multiple commands are 
// directed to the same node, they will be merged and sent at once using pipeling.
func (cluster *Cluster) RunBatch(bat *Batch) ([]interface{}, error) {
    for i := range bat.batches {
	go doBatch(&bat.batches[i])
    }

    for i := range bat.batches {
	<-bat.batches[i].done
    }

    var replies []interface{}
    for _, i := range bat.index {
	if bat.batches[i].err != nil {
	    return nil, bat.batches[i].err
	}

	replies = append(replies, bat.batches[i].cmds[0].reply)
	bat.batches[i].cmds = bat.batches[i].cmds[1:]
    }

    return replies, nil
}

func doBatch(batch *nodeBatch) {
    conn, err := batch.node.getConn()
    if err != nil {
	batch.err = err
	batch.done <- 1
	return
    }

    for i := range batch.cmds {
	conn.send(batch.cmds[i].cmd, batch.cmds[i].args...)
    }

    err = conn.flush()
    if err != nil {
	batch.err = err
	conn.shutdown()
	batch.done <- 1
	return
    }

    for i := range batch.cmds {
	reply, err := conn.receive()
	if err != nil {
	    batch.err = err
	    conn.shutdown()
	    batch.done <- 1
	    return
	}

	batch.cmds[i].reply, batch.cmds[i].err = reply, err
    }

    batch.node.releaseConn(conn)
    batch.done <- 1
}
