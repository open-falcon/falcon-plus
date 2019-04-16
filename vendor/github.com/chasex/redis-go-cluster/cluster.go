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
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Options is used to initialize a new redis cluster.
type Options struct {
	StartNodes []string // Startup nodes

	ConnTimeout  time.Duration // Connection timeout
	ReadTimeout  time.Duration // Read timeout
	WriteTimeout time.Duration // Write timeout

	KeepAlive int           // Maximum keep alive connecion in each node
	AliveTime time.Duration // Keep alive timeout
}

// Cluster is a redis client that manage connections to redis nodes,
// cache and update cluster info, and execute all kinds of commands.
// Multiple goroutines may invoke methods on a cluster simutaneously.
type Cluster struct {
	slots [kClusterSlots]*redisNode
	nodes map[string]*redisNode

	connTimeout  time.Duration
	readTimeout  time.Duration
	writeTimeout time.Duration

	keepAlive int
	aliveTime time.Duration

	updateTime time.Time
	updateList chan updateMesg

	rwLock sync.RWMutex

	closed bool
}

type updateMesg struct {
	node      *redisNode
	movedTime time.Time
}

// NewCluster create a new redis cluster client with specified options.
func NewCluster(options *Options) (*Cluster, error) {
	cluster := &Cluster{
		nodes:        make(map[string]*redisNode),
		connTimeout:  options.ConnTimeout,
		readTimeout:  options.ReadTimeout,
		writeTimeout: options.WriteTimeout,
		keepAlive:    options.KeepAlive,
		aliveTime:    options.AliveTime,
		updateList:   make(chan updateMesg),
	}

	for i := range options.StartNodes {
		node := &redisNode{
			address:      options.StartNodes[i],
			connTimeout:  options.ConnTimeout,
			readTimeout:  options.ReadTimeout,
			writeTimeout: options.WriteTimeout,
			keepAlive:    options.KeepAlive,
			aliveTime:    options.AliveTime,
		}

		err := cluster.update(node)
		if err != nil {
			continue
		} else {
			go cluster.handleUpdate()
			return cluster, nil
		}
	}

	return nil, fmt.Errorf("NewCluster: no valid node in %v", options.StartNodes)
}

// Do excute a redis command with random number arguments. First argument will
// be used as key to hash to a slot, so it only supports a subset of redis
// commands.
///
// SUPPORTED: most commands of keys, strings, lists, sets, sorted sets, hashes.
// NOT SUPPORTED: scripts, transactions, clusters.
//
// Particularly, MSET/MSETNX/MGET are supported using result aggregation.
// To MSET/MSETNX, there's no atomicity gurantee that given keys are set at once.
// It's possible that some keys are set, while others not.
//
// See README.md for more details.
// See full redis command list: http://www.redis.io/commands
func (cluster *Cluster) Do(cmd string, args ...interface{}) (interface{}, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("Do: no key found in args")
	}

	if cmd == "MSET" || cmd == "MSETNX" {
		return cluster.multiSet(cmd, args...)
	}

	if cmd == "MGET" {
		return cluster.multiGet(cmd, args...)
	}

	node, err := cluster.getNodeByKey(args[0])
	if err != nil {
		return nil, fmt.Errorf("Do: %v", err)
	}

	reply, err := node.do(cmd, args...)
	if err != nil {
		return nil, fmt.Errorf("Do: %v", err)
	}

	resp := checkReply(reply)

	switch resp {
	case kRespOK, kRespError:
		return reply, nil
	case kRespMove:
		return cluster.handleMove(node, reply.(redisError).Error(), cmd, args)
	case kRespAsk:
		return cluster.handleAsk(node, reply.(redisError).Error(), cmd, args)
	case kRespConnTimeout:
		return cluster.handleConnTimeout(node, cmd, args)
	}

	panic("unreachable")
}

// Close cluster connection, any subsequent method call will fail.
func (cluster *Cluster) Close() {
	cluster.rwLock.Lock()
	defer cluster.rwLock.Unlock()

	for addr, node := range cluster.nodes {
		node.shutdown()
		delete(cluster.nodes, addr)
	}

	cluster.closed = true
}

func (cluster *Cluster) handleMove(node *redisNode, replyMsg, cmd string, args []interface{}) (interface{}, error) {
	fields := strings.Split(replyMsg, " ")
	if len(fields) != 3 {
		return nil, fmt.Errorf("handleMove: invalid response \"%s\"", replyMsg)
	}

	// cluster has changed, inform update routine
	cluster.inform(node)

	newNode, err := cluster.getNodeByAddr(fields[2])
	if err != nil {
		return nil, fmt.Errorf("handleMove: %v", err)
	}

	return newNode.do(cmd, args...)
}

func (cluster *Cluster) handleAsk(node *redisNode, replyMsg, cmd string, args []interface{}) (interface{}, error) {
	fields := strings.Split(replyMsg, " ")
	if len(fields) != 3 {
		return nil, fmt.Errorf("handleAsk: invalid response \"%s\"", replyMsg)
	}

	newNode, err := cluster.getNodeByAddr(fields[2])
	if err != nil {
		return nil, fmt.Errorf("handleAsk: %v", err)
	}

	conn, err := newNode.getConn()
	if err != nil {
		return nil, fmt.Errorf("handleAsk: %v", err)
	}

	conn.send("ASKING")
	conn.send(cmd, args...)

	err = conn.flush()
	if err != nil {
		conn.shutdown()
		return nil, fmt.Errorf("handleAsk: %v", err)
	}

	re, err := String(conn.receive())
	if err != nil || re != "OK" {
		conn.shutdown()
		return nil, fmt.Errorf("handleAsk: %v", err)
	}

	reply, err := conn.receive()
	if err != nil {
		conn.shutdown()
		return nil, fmt.Errorf("handleAsk: %v", err)
	}

	newNode.releaseConn(conn)

	return reply, nil
}

func (cluster *Cluster) handleConnTimeout(node *redisNode, cmd string, args []interface{}) (interface{}, error) {
	var randomNode *redisNode

	// choose a random node other than previous one
	cluster.rwLock.RLock()
	for _, randomNode = range cluster.nodes {
		if randomNode.address != node.address {
			break
		}
	}
	cluster.rwLock.RUnlock()

	reply, err := randomNode.do(cmd, args...)
	if err != nil {
		return nil, fmt.Errorf("handleConnTimeout: %v", err)
	}

	if _, ok := reply.(redisError); !ok {
		// we happen to choose the right node, which means
		// that cluster has changed, so inform update routine.
		cluster.inform(randomNode)
		return reply, nil
	}

	// ignore replies other than MOVED
	errMsg := reply.(redisError).Error()
	if len(errMsg) < 5 || string(errMsg[:5]) != "MOVED" {
		return nil, errors.New(errMsg)
	}

	// When MOVED received, we check wether move adress equal to
	// previous one. If equal, then it's just an connection timeout
	// error, return error and carry on. If not, then the master may
	// down or unreachable, a new master has served the slot, request
	// new master and update cluster info.
	//
	// TODO: At worst case, it will request redis 3 times on a single
	// command, will this be a problem?
	fields := strings.Split(errMsg, " ")
	if len(fields) != 3 {
		return nil, fmt.Errorf("handleConnTimeout: invalid response \"%s\"", errMsg)
	}

	if fields[2] == node.address {
		return nil, fmt.Errorf("handleConnTimeout: %s connection timeout", node.address)
	}

	// cluster change, inform back routine to update
	cluster.inform(randomNode)

	newNode, err := cluster.getNodeByAddr(fields[2])
	if err != nil {
		return nil, fmt.Errorf("handleConnTimeout: %v", err)
	}

	return newNode.do(cmd, args...)
}

const (
	kClusterSlots = 16384

	kRespOK          = 0
	kRespMove        = 1
	kRespAsk         = 2
	kRespConnTimeout = 3
	kRespError       = 4
)

func checkReply(reply interface{}) int {
	if _, ok := reply.(redisError); !ok {
		return kRespOK
	}

	errMsg := reply.(redisError).Error()

	if len(errMsg) >= 3 && string(errMsg[:3]) == "ASK" {
		return kRespAsk
	}

	if len(errMsg) >= 5 && string(errMsg[:5]) == "MOVED" {
		return kRespMove
	}

	if len(errMsg) >= 12 && string(errMsg[:12]) == "ECONNTIMEOUT" {
		return kRespConnTimeout
	}

	return kRespError
}

func (cluster *Cluster) update(node *redisNode) error {
	info, err := Values(node.do("CLUSTER", "SLOTS"))
	if err != nil {
		return err
	}

	errFormat := fmt.Errorf("update: %s invalid response", node.address)

	var nslots int
	slots := make(map[string][]uint16)

	for _, i := range info {
		m, err := Values(i, err)
		if err != nil || len(m) < 3 {
			return errFormat
		}

		start, err := Int(m[0], err)
		if err != nil {
			return errFormat
		}

		end, err := Int(m[1], err)
		if err != nil {
			return errFormat
		}

		t, err := Values(m[2], err)
		if err != nil || len(t) < 2 {
			return errFormat
		}

		var ip string
		var port int

		_, err = Scan(t, &ip, &port)
		if err != nil {
			return errFormat
		}
		addr := fmt.Sprintf("%s:%d", ip, port)

		slot, ok := slots[addr]
		if !ok {
			slot = make([]uint16, 0, 2)
		}

		nslots += end - start + 1

		slot = append(slot, uint16(start))
		slot = append(slot, uint16(end))

		slots[addr] = slot
	}

	// TODO: Is full coverage really needed?
	if nslots != kClusterSlots {
		return fmt.Errorf("update: %s slots not full covered", node.address)
	}

	cluster.rwLock.Lock()
	defer cluster.rwLock.Unlock()

	t := time.Now()
	cluster.updateTime = t

	for addr, slot := range slots {
		node, ok := cluster.nodes[addr]
		if !ok {
			node = &redisNode{
				address:      addr,
				connTimeout:  cluster.connTimeout,
				readTimeout:  cluster.readTimeout,
				writeTimeout: cluster.writeTimeout,
				keepAlive:    cluster.keepAlive,
				aliveTime:    cluster.aliveTime,
			}
		}

		n := len(slot)
		for i := 0; i < n-1; i += 2 {
			start := slot[i]
			end := slot[i+1]

			for j := start; j <= end; j++ {
				cluster.slots[j] = node
			}
		}

		node.updateTime = t
		cluster.nodes[addr] = node
	}

	// shrink
	for addr, node := range cluster.nodes {
		if node.updateTime != t {
			node.shutdown()

			delete(cluster.nodes, addr)
		}
	}

	return nil
}

func (cluster *Cluster) handleUpdate() {
	for {
		msg := <-cluster.updateList

		// TODO: control update frequency by updateTime and movedTime?

		err := cluster.update(msg.node)
		if err != nil {
			log.Printf("handleUpdate: %v\n", err)
		}
	}
}

func (cluster *Cluster) inform(node *redisNode) {
	mesg := updateMesg{
		node:      node,
		movedTime: time.Now(),
	}

	select {
	case cluster.updateList <- mesg:
		// Push update message, no more to do.
	default:
		// Update channel full, just carry on.
	}
}

func (cluster *Cluster) getNodeByAddr(addr string) (*redisNode, error) {
	cluster.rwLock.RLock()
	defer cluster.rwLock.RUnlock()

	if cluster.closed {
		return nil, fmt.Errorf("getNodeByAddr: cluster has been closed")
	}

	node, ok := cluster.nodes[addr]
	if !ok {
		return nil, fmt.Errorf("getNodeByAddr: %s not found", addr)
	}

	return node, nil
}

func (cluster *Cluster) getNodeByKey(arg interface{}) (*redisNode, error) {
	key, err := key(arg)
	if err != nil {
		return nil, fmt.Errorf("getNodeByKey: invalid key %v", key)
	}

	slot := hash(key)

	cluster.rwLock.RLock()
	defer cluster.rwLock.RUnlock()

	if cluster.closed {
		return nil, fmt.Errorf("getNodeByKey: cluster has been closed")
	}

	node := cluster.slots[slot]
	if node == nil {
		return nil, fmt.Errorf("getNodeByKey: %s[%d] no node found", key, slot)
	}

	return node, nil
}

func key(arg interface{}) (string, error) {
	switch arg := arg.(type) {
	case int:
		return strconv.Itoa(arg), nil
	case int64:
		return strconv.Itoa(int(arg)), nil
	case float64:
		return strconv.FormatFloat(arg, 'g', -1, 64), nil
	case string:
		return arg, nil
	case []byte:
		return string(arg), nil
	default:
		return "", fmt.Errorf("key: unknown type %T", arg)
	}
}

func hash(key string) uint16 {
	var s, e int
	for s = 0; s < len(key); s++ {
		if key[s] == '{' {
			break
		}
	}

	if s == len(key) {
		return crc16(key) & (kClusterSlots - 1)
	}

	for e = s + 1; e < len(key); e++ {
		if key[e] == '}' {
			break
		}
	}

	if e == len(key) || e == s+1 {
		return crc16(key) & (kClusterSlots - 1)
	}

	return crc16(key[s+1:e]) & (kClusterSlots - 1)
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
