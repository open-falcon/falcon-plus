// Copyright 2017 Xiaomi, Inc.
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

package cron

import (
	"log"
	"math/rand"
	"time"

	"github.com/open-falcon/falcon-plus/modules/aggregator/g"
)

type Worker struct {
	Ticker      *time.Ticker
	ClusterItem *g.Cluster
	Quit        chan struct{}
}

func NewWorker(ci *g.Cluster) Worker {
	w := Worker{}
	//	w.Ticker = time.NewTicker(time.Duration(ci.Step) * time.Second)
	w.Quit = make(chan struct{})
	w.ClusterItem = ci
	return w
}

func (this Worker) Start() {
	go func() {
		s1 := rand.NewSource(time.Now().UnixNano() * this.ClusterItem.Id)
		r1 := rand.New(s1)
		// 60s, step usually
		delay := r1.Int63n(60000)

		if g.Config().Debug {
			log.Printf("[I] after %5d ms, start worker %d", delay, this.ClusterItem.Id)
		}

		time.Sleep(time.Duration(delay) * time.Millisecond)
		this.Ticker = time.NewTicker(time.Duration(this.ClusterItem.Step) * time.Second)
		for {
			select {
			case <-this.Ticker.C:
				WorkerRun(this.ClusterItem)
			case <-this.Quit:
				if g.Config().Debug {
					log.Println("[I] drop worker", this.ClusterItem)
				}
				this.Ticker.Stop()
				return
			}
		}
	}()
}

func (this Worker) Drop() {
	close(this.Quit)
}

var Workers = make(map[string]Worker)

func deleteNoUseWorker(m map[string]*g.Cluster) {
	del := []string{}
	for key, worker := range Workers {
		if _, ok := m[key]; !ok {
			worker.Drop()
			del = append(del, key)
		}
	}

	for _, key := range del {
		delete(Workers, key)
	}
}

func createWorkerIfNeed(m map[string]*g.Cluster) {
	for key, item := range m {
		if _, ok := Workers[key]; !ok {
			if item.Step <= 0 {
				log.Println("[W] invalid cluster(step <= 0):", item)
				continue
			}
			worker := NewWorker(item)
			Workers[key] = worker
			worker.Start()
		}
	}
}
