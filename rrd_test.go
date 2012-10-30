package rrd

import (
	"testing"
	"time"
)

func TestAll(t *testing.T) {
	const dbfile = "/tmp/test.rrd"

	// Create
	c := NewCreater(dbfile, time.Now(), 1)
	c.RRA("AVERAGE", 0.5, 1, 100)
	c.RRA("AVERAGE", 0.5, 5, 100)
	c.DS("cnt", "COUNTER", 10, 0, 100)
	c.DS("g", "GAUGE", 10, 0, 60)
	err := c.Create(true)
	if err != nil {
		t.Fatal(err)
	}

	// Update
	u := NewUpdater(dbfile)
	for i := 0; i < 10; i++ {
		time.Sleep(time.Second)
		err := u.Update(time.Now(), i, 1.5*float64(i))
		if err != nil {
			t.Fatal(err)
		}
	}

	// Update with cache
	for i := 10; i < 20; i++ {
		time.Sleep(time.Second)
		u.Cache(time.Now(), i, 1.5*float64(i))
	}
	err = u.Update()
	if err != nil {
		t.Fatal(err)
	}
	err = u.Update()
	if err != nil {
		t.Fatal(err)
	}
}
