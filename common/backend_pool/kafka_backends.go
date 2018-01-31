// Copyright 2018 CMCC IOT, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package backend_pool

import (
	"errors"
	"fmt"
	"github.com/Shopify/sarama"
	"log"
	"time"
)

// kafka
type KafkaConnPoolHelper struct {
	p            sarama.SyncProducer
	maxConns     int
	connTimeout  int
	writeTimeout int
	address      string
}

func newKafkaConnPool(address string, maxConns int, connTimeout int, writeTimeout int) sarama.SyncProducer {
	kafka_config := sarama.NewConfig()
	kafka_config.Net.MaxOpenRequests = maxConns
	kafka_config.Net.DialTimeout = time.Duration(connTimeout) * time.Millisecond
	kafka_config.Net.WriteTimeout = time.Duration(writeTimeout) * time.Millisecond
	kafka_config.Producer.RequiredAcks = sarama.WaitForAll // ACK
	//    kafka_config.Producer.Partitioner = sarama.NewRandomPartitioner  // 随机分区
	kafka_config.Producer.Partitioner = sarama.NewManualPartitioner
	kafka_config.Producer.Return.Successes = true // 返回true

	kafka_client, err := sarama.NewSyncProducer([]string{address}, kafka_config)
	if err != nil {
		log.Fatalln("producer close, err:", err)
		return kafka_client
	}

	//    defer kafka_client.Close()

	return kafka_client
}

func NewKafkaConnPoolHelper(address string, maxConns, connTimeout, writeTimeout int) *KafkaConnPoolHelper {
	return &KafkaConnPoolHelper{
		p:            newKafkaConnPool(address, maxConns, connTimeout, writeTimeout),
		maxConns:     maxConns,
		connTimeout:  connTimeout,
		writeTimeout: writeTimeout,
		address:      address,
	}
}

func (this *KafkaConnPoolHelper) Send(value []byte, topic string, partition int32) (err error) {
	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	msg.Partition = 0
	msg.Value = sarama.StringEncoder(value)
	partition, offset, err := this.p.SendMessage(msg)

	if err == nil {
		return nil
	}
	time.Sleep(100 * time.Millisecond)
	if err != nil {
		return errors.New("Failed to produce message to kafka cluster.")
	}
	if partition != 0 {
		return errors.New("Message should have been produced to partition 0!")
	}
	return errors.New(fmt.Sprintf("Produced message to partition %d with offset %d", partition, offset))
}
