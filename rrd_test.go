package rrd

import (
	"testing"
	"time"
)

func TestAll(t *testing.T) {
	// Create
	const (
		dbfile    = "/tmp/test.rrd"
		step      = 1
		heartbeat = 2 * step
	)
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

	// Graph
	g := NewGrapher()
	g.SetTitle("Test")
	g.SetVlabel("some variable")
	g.SetSize(800, 300)
	g.Def("v1", dbfile, "g", "AVERAGE")
	g.Def("v2", dbfile, "cnt", "AVERAGE")
	g.Line(1, "v1", "ff0000", "var 1")
	g.Line(1.5, "v2", "0000ff", "var 2")
	now := time.Now()
	_, err = g.SaveGraph("/tmp/test_rrd.png", now.Add(-20*time.Second), now)
	if err != nil {
		t.Fatal(err)
	}
}
