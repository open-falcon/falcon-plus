package rrdlite

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"testing"
	"time"
)

const (
	dbfile    = "/tmp/test.rrd"
	step      = 1
	heartbeat = 2 * step
	b_size    = 100000
	work_size = 10
)

var now time.Time
var wg sync.WaitGroup

func init() {
	now = time.Now()
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func testAll(t *testing.T) {
	// Create

	c := NewCreator(dbfile, time.Now(), step)
	c.RRA("AVERAGE", 0.5, 1, 100)
	c.RRA("AVERAGE", 0.5, 5, 100)
	c.DS("cnt", "COUNTER", heartbeat, 0, 100)
	c.DS("g", "GAUGE", heartbeat, 0, 60)
	err := c.Create(true)
	if err != nil {
		t.Fatal(err)
	}

	// Update
	u := NewUpdater(dbfile)
	for i := 0; i < 10; i++ {
		time.Sleep(step * time.Second)
		err := u.Update(time.Now(), i, 1.5*float64(i))
		if err != nil {
			t.Fatal(err)
		}
	}

	// Update with cache
	for i := 10; i < 20; i++ {
		time.Sleep(step * time.Second)
		u.Cache(time.Now(), i, 2*float64(i))
	}
	err = u.Update()
	if err != nil {
		t.Fatal(err)
	}

	// Info
	inf, err := Info(dbfile)
	if err != nil {
		t.Fatal(err)
	}
	for k, v := range inf {
		fmt.Printf("%s (%T): %v\n", k, v, v)
	}

	// Fetch
	end := time.Unix(int64(inf["last_update"].(uint)), 0)
	start := end.Add(-20 * step * time.Second)
	fmt.Printf("Fetch Params:\n")
	fmt.Printf("Start: %s\n", start)
	fmt.Printf("End: %s\n", end)
	fmt.Printf("Step: %s\n", step*time.Second)
	fetchRes, err := Fetch(dbfile, "AVERAGE", start, end, step*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	defer fetchRes.FreeValues()
	fmt.Printf("FetchResult:\n")
	fmt.Printf("Start: %s\n", fetchRes.Start)
	fmt.Printf("End: %s\n", fetchRes.End)
	fmt.Printf("Step: %s\n", fetchRes.Step)
	for _, dsName := range fetchRes.DsNames {
		fmt.Printf("\t%s", dsName)
	}
	fmt.Printf("\n")

	row := 0
	for ti := fetchRes.Start.Add(fetchRes.Step); ti.Before(end) || ti.Equal(end); ti = ti.Add(fetchRes.Step) {
		fmt.Printf("%s / %d", ti, ti.Unix())
		for i := 0; i < len(fetchRes.DsNames); i++ {
			v := fetchRes.ValueAt(i, row)
			fmt.Printf("\t%e", v)
		}
		fmt.Printf("\n")
		row++
	}
}

func add(b *testing.B, filename string) {
	c := NewCreator(filename, now, step)
	c.RRA("AVERAGE", 0.5, 1, 100)
	c.RRA("AVERAGE", 0.5, 5, 100)
	c.DS("g", "GAUGE", heartbeat, 0, 60)
	err := c.Create(true)
	if err != nil {
		b.Fatal(err)
	}
}

func update(b *testing.B, filename string) {
	u := NewUpdater(filename)
	err := u.Update(now, 1.5)
	if err != nil {
		b.Fatal(err)
	}
}

func fetch(b *testing.B, filename string, start, end time.Time) {
	if fetchRes, err := Fetch(filename, "AVERAGE", start, end, step*time.Second); err != nil {
		b.Fatal(err)
	} else {
		fetchRes.FreeValues()
	}
}

func BenchmarkAdd(b *testing.B) {
	b.StopTimer()
	if err := exec.Command("rm", "-rf", "/tmp/rrd").Run(); err != nil {
		b.Fatal(err)
	}
	if err := os.Mkdir("/tmp/rrd", 0755); err != nil {
		b.Fatal(err)
	}
	for i := 0; i < 256; i++ {
		if err := os.Mkdir(fmt.Sprintf("/tmp/rrd/%d", i), 0755); err != nil {
			b.Fatal(err)
		}
	}
	b.StartTimer()
	b.N = b_size
	n := b.N / work_size
	for j := 0; j < work_size; j++ {
		wg.Add(1)
		go func(j int) {
			defer wg.Done()
			for i := 0; i < n; i++ {
				filename := fmt.Sprintf("/tmp/rrd/%d/%d-%d.rrd", i%256, j, i)
				add(b, filename)
			}
		}(j)
	}
	wg.Wait()
}

func BenchmarkUpdate(b *testing.B) {
	b.N = b_size
	n := b.N / work_size

	for j := 0; j < work_size; j++ {
		wg.Add(1)
		go func(j int) {
			defer wg.Done()
			for i := 0; i < n; i++ {
				filename := fmt.Sprintf("/tmp/rrd/%d/%d-%d.rrd", i%256, j, i)
				update(b, filename)
			}
		}(j)
	}
	wg.Wait()
}

func BenchmarkFetch(b *testing.B) {
	b.N = b_size
	n := b.N / work_size
	start := time.Unix(now.Unix()-step, 0)
	end := start.Add(20 * step * time.Second)
	for j := 0; j < work_size; j++ {
		wg.Add(1)
		go func(j int) {
			defer wg.Done()
			for i := 0; i < n; i++ {
				filename := fmt.Sprintf("/tmp/rrd/%d/%d-%d.rrd", i%256, j, i)
				fetch(b, filename, start, end)
			}
		}(j)
	}
	wg.Wait()
}
