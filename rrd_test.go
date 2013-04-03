package rrd

import (
	"fmt"
	"io/ioutil"
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

	// Info
	inf, err := Info(dbfile)
	if err != nil {
		t.Fatal(err)
	}
	for k, v := range inf {
		fmt.Printf("%s (%T): %v\n", k, v, v)
	}

	// Graph
	g := NewGrapher()
	g.SetTitle("Test")
	g.SetVLabel("some variable")
	g.SetSize(800, 300)
	g.SetWatermark("some watermark")
	g.Def("v1", dbfile, "g", "AVERAGE")
	g.Def("v2", dbfile, "cnt", "AVERAGE")
	g.VDef("max1", "v1,MAXIMUM")
	g.VDef("avg2", "v2,AVERAGE")
	g.Line(1, "v1", "ff0000", "var 1")
	g.Area("v2", "0000ff", "var 2")
	g.GPrintT("max1", "max1 at %c")
	g.GPrint("avg2", "avg2=%lf")
	g.PrintT("max1", "max1 at %c")
	g.Print("avg2", "avg2=%lf")

	now := time.Now()

	i, err := g.SaveGraph("/tmp/test_rrd1.png", now.Add(-20*time.Second), now)
	fmt.Printf("%+v\n", i)
	if err != nil {
		t.Fatal(err)
	}
	i, buf, err := g.Graph(now.Add(-20*time.Second), now)
	fmt.Printf("%+v\n", i)
	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile("/tmp/test_rrd2.png", buf, 0666)
	if err != nil {
		t.Fatal(err)
	}

	// Fetch
	end := time.Unix(int64(inf["last_update"].(uint)), 0)
	start := end.Add(-20 * step * time.Second)
	fmt.Printf("Fetch Params:\n")
	fmt.Printf("Start: %s\n", start)
	fmt.Printf("End: %s\n", end)
	fmt.Printf("Step: %s\n", step * time.Second)
	res, err := Fetch(dbfile, "AVERAGE", start, end, step * time.Second)
	if err != nil {
		t.Fatal(err)
	}
	defer res.FreeValues()
	fmt.Printf("FetchResult:\n")
	fmt.Printf("Start: %s\n", res.Start)
	fmt.Printf("End: %s\n", res.End)
	fmt.Printf("Step: %s\n", res.Step)
	for _, dsName := range res.DsNames {
		fmt.Printf("\t%s", dsName)
	}
	fmt.Printf("\n")

	tm := res.Start
	for row := 0; row < res.RowLen; row++ {
		tm = tm.Add(res.Step)
		fmt.Printf("%s / %d", tm, tm.Unix())
		for i := 0; i < len(res.DsNames); i++ {
			v := res.ValueAt(i, row)
			fmt.Printf("\t%e", v)
		}
		fmt.Printf("\n")
	}
}
