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

package utils

import (
	"fmt"
	"testing"
)

func origPK(endpoint, metric string, tags map[string]string) string {
	if tags == nil || len(tags) == 0 {
		return fmt.Sprintf("%s/%s", endpoint, metric)
	}
	return fmt.Sprintf("%s/%s/%s", endpoint, metric, SortedTags(tags))
}

func origPK2(endpoint, counter string) string {
	return fmt.Sprintf("%s/%s", endpoint, counter)
}

func origUUID(endpoint, metric string, tags map[string]string, dstype string, step int) string {
	if tags == nil || len(tags) == 0 {
		return fmt.Sprintf("%s/%s/%s/%d", endpoint, metric, dstype, step)
	}
	return fmt.Sprintf("%s/%s/%s/%s/%d", endpoint, metric, SortedTags(tags), dstype, step)
}

var pkCase = []struct {
	endpoint string
	metric   string
	tags     map[string]string
	except   string
}{
	{"endpoint1", "metric1", nil, "endpoint1/metric1"},
	{"endpoint1", "metric1", map[string]string{}, "endpoint1/metric1"},
	{"endpoint1", "metric1", map[string]string{"k1": "v1", "k2": "v2"}, "endpoint1/metric1/k1=v1,k2=v2"},
	{"endpoint1", "metric1", map[string]string{"k2": "v2", "k1": "v1"}, "endpoint1/metric1/k1=v1,k2=v2"},
}

var pk2Case = []struct {
	endpoint string
	counter  string
	except   string
}{
	{"endpoint1", "counter1", "endpoint1/counter1"},
}

var uuidCase = []struct {
	endpoint string
	metric   string
	tags     map[string]string
	dstype   string
	step     int
	except   string
}{
	{"endpoint1", "metric1", nil, "ds", 10, "endpoint1/metric1/ds/10"},
	{"endpoint1", "metric1", map[string]string{}, "ds", 10, "endpoint1/metric1/ds/10"},
	{"endpoint1", "metric1", map[string]string{"k1": "v1", "k2": "v2"}, "ds", 10, "endpoint1/metric1/k1=v1,k2=v2/ds/10"},
	{"endpoint1", "metric1", map[string]string{"k2": "v2", "k1": "v1"}, "ds", 10, "endpoint1/metric1/k1=v1,k2=v2/ds/10"},
}

func Test_PK(t *testing.T) {
	for _, pk := range pkCase {
		if PK(pk.endpoint, pk.metric, pk.tags) != origPK(pk.endpoint, pk.metric, pk.tags) || PK(pk.endpoint, pk.metric, pk.tags) != pk.except {
			t.Error("not except")
		}
	}
}

func Test_PK2(t *testing.T) {
	for _, pk2 := range pk2Case {
		if PK2(pk2.endpoint, pk2.counter) != origPK2(pk2.endpoint, pk2.counter) || PK2(pk2.endpoint, pk2.counter) != pk2.except {
			t.Error("not except")
		}
	}
}

func Test_UUID(t *testing.T) {
	for _, uuid := range uuidCase {
		if UUID(uuid.endpoint, uuid.metric, uuid.tags, uuid.dstype, uuid.step) != origUUID(uuid.endpoint, uuid.metric, uuid.tags, uuid.dstype, uuid.step) ||
			UUID(uuid.endpoint, uuid.metric, uuid.tags, uuid.dstype, uuid.step) != uuid.except {
			t.Error("not except")
		}
	}
}

var (
	testTags = map[string]string{"k1": "v1", "k2": "v2"}
)

func Benchmark_PK(b *testing.B) {
	for i := 0; i < b.N; i++ {
		PK("endpoint1", "metric1", testTags)
	}
}

func Benchmark_PK_orig(b *testing.B) {
	for i := 0; i < b.N; i++ {
		origPK("endpoint1", "metric1", testTags)
	}
}

func Benchmark_PK2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		PK2("endpoint1", "counter1")
	}
}

func Benchmark_PK2_orig(b *testing.B) {
	for i := 0; i < b.N; i++ {
		origPK2("endpoint1", "counter1")
	}
}

func Benchmark_UUID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		UUID("endpoint1", "metric1", testTags, "dt", 10)
	}
}

func Benchmark_UUID_orig(b *testing.B) {
	for i := 0; i < b.N; i++ {
		origUUID("endpoint1", "metric1", testTags, "dt", 10)
	}
}
