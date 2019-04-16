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

type multiTask struct {
	node *redisNode
	slot uint16

	cmd  string
	args []interface{}

	reply   interface{}
	replies []interface{}
	err     error

	done chan int
}

func (cluster *Cluster) multiSet(cmd string, args ...interface{}) (interface{}, error) {
	if len(args)&1 != 0 {
		return nil, fmt.Errorf("multiSet: invalid args %v", args)
	}

	tasks := make([]*multiTask, 0)

	cluster.rwLock.RLock()
	for i := 0; i < len(args); i += 2 {
		key, err := key(args[i])
		if err != nil {
			cluster.rwLock.RUnlock()
			return nil, fmt.Errorf("multiSet: invalid key %v", args[i])
		}

		slot := hash(key)

		var j int
		for j = 0; j < len(tasks); j++ {
			if tasks[j].slot == slot {
				tasks[j].args = append(tasks[j].args, args[i])   // key
				tasks[j].args = append(tasks[j].args, args[i+1]) // value

				break
			}
		}

		if j == len(tasks) {
			node := cluster.slots[slot]
			if node == nil {
				cluster.rwLock.RUnlock()
				return nil, fmt.Errorf("multiSet: %s[%d] no node found", key, slot)
			}

			task := &multiTask{
				node: node,
				slot: slot,
				cmd:  cmd,
				args: []interface{}{args[i], args[i+1]},
				done: make(chan int),
			}
			tasks = append(tasks, task)
		}
	}
	cluster.rwLock.RUnlock()

	for i := range tasks {
		go handleSetTask(tasks[i])
	}

	for i := range tasks {
		<-tasks[i].done
	}

	for i := range tasks {
		_, err := String(tasks[i].reply, tasks[i].err)
		if err != nil {
			return nil, err
		}
	}

	return "OK", nil
}

func (cluster *Cluster) multiGet(cmd string, args ...interface{}) (interface{}, error) {
	tasks := make([]*multiTask, 0)
	index := make([]*multiTask, len(args))

	cluster.rwLock.RLock()
	for i := 0; i < len(args); i++ {
		key, err := key(args[i])
		if err != nil {
			cluster.rwLock.RUnlock()
			return nil, fmt.Errorf("multiGet: invalid key %v", args[i])
		}

		slot := hash(key)

		var j int
		for j = 0; j < len(tasks); j++ {
			if tasks[j].slot == slot {
				tasks[j].args = append(tasks[j].args, args[i]) // key
				index[i] = tasks[j]

				break
			}
		}

		if j == len(tasks) {
			node := cluster.slots[slot]
			if node == nil {
				cluster.rwLock.RUnlock()
				return nil, fmt.Errorf("multiGet: %s[%d] no node found", key, slot)
			}

			task := &multiTask{
				node: node,
				slot: slot,
				cmd:  cmd,
				args: []interface{}{args[i]},
				done: make(chan int),
			}
			tasks = append(tasks, task)
			index[i] = tasks[j]
		}
	}
	cluster.rwLock.RUnlock()

	for i := range tasks {
		go handleGetTask(tasks[i])
	}

	for i := range tasks {
		<-tasks[i].done
	}

	reply := make([]interface{}, len(args))
	for i := range reply {
		if index[i].err != nil {
			return nil, index[i].err
		}

		if len(index[i].replies) < 0 {
			panic("unreachable")
		}

		reply[i] = index[i].replies[0]
		index[i].replies = index[i].replies[1:]
	}

	return reply, nil
}

func handleSetTask(task *multiTask) {
	task.reply, task.err = task.node.do(task.cmd, task.args...)
	task.done <- 1
}

func handleGetTask(task *multiTask) {
	task.replies, task.err = Values(task.node.do(task.cmd, task.args...))
	task.done <- 1
}
