package base

import (
	"runtime/debug"
	"time"

	"github.com/niean/go-metrics-lite"
)

var (
	debugMetrics struct {
		GCStats struct {
			LastGC     metrics.Gauge
			NumGC      metrics.Gauge
			Pause      metrics.Histogram
			PauseTotal metrics.Gauge
		}
		ReadGCStats metrics.Histogram
	}
	gcStats = debug.GCStats{Pause: make([]time.Duration, 11)}
)

func RegisterAndCaptureDebugGCStats(r metrics.Registry, d time.Duration) {
	registerDebugGCStats(r)
	go captureDebugGCStats(r, d)
}

// Capture new values for the Go garbage collector statistics exported in
// debug.GCStats.  This is designed to be called as a goroutine.
func captureDebugGCStats(r metrics.Registry, d time.Duration) {
	for _ = range time.Tick(d) {
		captureDebugGCStatsOnce(r)
	}
}

// Capture new values for the Go garbage collector statistics exported in
// debug.GCStats.  This is designed to be called in a background goroutine.
// Giving a registry which has not been given to RegisterDebugGCStats will
// panic.
//
// Be careful (but much less so) with this because debug.ReadGCStats calls
// the C function runtime·lock(runtime·mheap) which, while not a stop-the-world
// operation, isn't something you want to be doing all the time.
func captureDebugGCStatsOnce(r metrics.Registry) {
	lastGC := gcStats.LastGC
	t := time.Now()
	debug.ReadGCStats(&gcStats)
	debugMetrics.ReadGCStats.Update(int64(time.Since(t)))

	debugMetrics.GCStats.LastGC.Update(int64(gcStats.LastGC.UnixNano()))
	debugMetrics.GCStats.NumGC.Update(int64(gcStats.NumGC))
	if lastGC != gcStats.LastGC && 0 < len(gcStats.Pause) {
		debugMetrics.GCStats.Pause.Update(int64(gcStats.Pause[0]))
	}
	//debugMetrics.GCStats.PauseQuantiles.Update(gcStats.PauseQuantiles)
	debugMetrics.GCStats.PauseTotal.Update(int64(gcStats.PauseTotal))
}

// Register metrics for the Go garbage collector statistics exported in
// debug.GCStats.  The metrics are named by their fully-qualified Go symbols,
// i.e. debug.GCStats.PauseTotal.
func registerDebugGCStats(r metrics.Registry) {
	debugMetrics.GCStats.LastGC = metrics.NewGauge()
	debugMetrics.GCStats.NumGC = metrics.NewGauge()
	debugMetrics.GCStats.Pause = metrics.NewHistogram(metrics.NewExpDecaySample(1028, 0.015))
	debugMetrics.GCStats.PauseTotal = metrics.NewGauge()
	debugMetrics.ReadGCStats = metrics.NewHistogram(metrics.NewExpDecaySample(1028, 0.015))

	r.Register("debug.GCStats.LastGC", debugMetrics.GCStats.LastGC)
	r.Register("debug.GCStats.NumGC", debugMetrics.GCStats.NumGC)
	r.Register("debug.GCStats.Pause", debugMetrics.GCStats.Pause)
	r.Register("debug.GCStats.PauseTotal", debugMetrics.GCStats.PauseTotal)
	r.Register("debug.ReadGCStats", debugMetrics.ReadGCStats)
}
