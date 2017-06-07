package metrics

import (
	"sync"
	"time"
)

// Meters count events to produce exponentially-weighted moving average rates
// at one-, five-, and fifteen-minutes and a mean rate.
type Meter interface {
	Count() int64
	Mark(int64)
	Rate1() float64
	Rate5() float64
	Rate15() float64
	RateMean() float64
	RateStep() float64
	Snapshot() Meter
}

// GetOrRegisterMeter returns an existing Meter or constructs and registers a
// new StandardMeter.
func GetOrRegisterMeter(name string, r Registry) Meter {
	if nil == r {
		r = DefaultRegistry
	}
	return r.GetOrRegister(name, NewMeter).(Meter)
}

// NewMeter constructs a new StandardMeter and launches a goroutine.
func NewMeter() Meter {
	if UseNilMetrics {
		return NilMeter{}
	}
	m := newStandardMeter()
	arbiter.Lock()
	defer arbiter.Unlock()
	arbiter.meters = append(arbiter.meters, m)
	if !arbiter.started {
		arbiter.started = true
		go arbiter.tick()
	}
	return m
}

// NewMeter constructs and registers a new StandardMeter and launches a
// goroutine.
func NewRegisteredMeter(name string, r Registry) Meter {
	c := NewMeter()
	if nil == r {
		r = DefaultRegistry
	}
	r.Register(name, c)
	return c
}

// MeterSnapshot is a read-only copy of another Meter.
type MeterSnapshot struct {
	count, _lastCount                        int64
	rate1, rate5, rate15, rateMean, rateStep float64
	_lastTime                                time.Time
}

// Count returns the count of events at the time the snapshot was taken.
func (m *MeterSnapshot) Count() int64 { return m.count }

// Mark panics.
func (*MeterSnapshot) Mark(n int64) {
	panic("Mark called on a MeterSnapshot")
}

// Rate1 returns the one-minute moving average rate of events per second at the
// time the snapshot was taken.
func (m *MeterSnapshot) Rate1() float64 { return m.rate1 }

// Rate5 returns the five-minute moving average rate of events per second at
// the time the snapshot was taken.
func (m *MeterSnapshot) Rate5() float64 { return m.rate5 }

// Rate15 returns the fifteen-minute moving average rate of events per second
// at the time the snapshot was taken.
func (m *MeterSnapshot) Rate15() float64 { return m.rate15 }

// RateMean returns the meter's mean rate of events per second at the time the
// snapshot was taken.
func (m *MeterSnapshot) RateMean() float64 { return m.rateMean }
func (m *MeterSnapshot) RateStep() float64 { return m.rateStep }

// Snapshot returns the snapshot.
func (m *MeterSnapshot) Snapshot() Meter { return m }

// NilMeter is a no-op Meter.
type NilMeter struct{}

// Count is a no-op.
func (NilMeter) Count() int64 { return 0 }

// Mark is a no-op.
func (NilMeter) Mark(n int64) {}

// Rate1 is a no-op.
func (NilMeter) Rate1() float64 { return 0.0 }

// Rate5 is a no-op.
func (NilMeter) Rate5() float64 { return 0.0 }

// Rate15is a no-op.
func (NilMeter) Rate15() float64 { return 0.0 }

// RateMean is a no-op.
func (NilMeter) RateMean() float64 { return 0.0 }
func (NilMeter) RateStep() float64 { return 0.0 }

// Snapshot is a no-op.
func (NilMeter) Snapshot() Meter { return NilMeter{} }

// StandardMeter is the standard implementation of a Meter.
type StandardMeter struct {
	lock        sync.RWMutex
	snapshot    *MeterSnapshot
	a1, a5, a15 EWMA
	startTime   time.Time
}

func newStandardMeter() *StandardMeter {
	nw := time.Now()
	return &StandardMeter{
		snapshot:  &MeterSnapshot{_lastTime: nw},
		a1:        NewEWMA1(),
		a5:        NewEWMA5(),
		a15:       NewEWMA15(),
		startTime: nw,
	}
}

// Count returns the number of events recorded.
func (m *StandardMeter) Count() int64 {
	m.lock.RLock()
	count := m.snapshot.count
	m.lock.RUnlock()
	return count
}

// Mark records the occurance of n events.
func (m *StandardMeter) Mark(n int64) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.snapshot.count += n
	m.a1.Update(n)
	m.a5.Update(n)
	m.a15.Update(n)
	m.updateSnapshot()
}

// Rate1 returns the one-minute moving average rate of events per second.
func (m *StandardMeter) Rate1() float64 {
	m.lock.RLock()
	rate1 := m.snapshot.rate1
	m.lock.RUnlock()
	return rate1
}

// Rate5 returns the five-minute moving average rate of events per second.
func (m *StandardMeter) Rate5() float64 {
	m.lock.RLock()
	rate5 := m.snapshot.rate5
	m.lock.RUnlock()
	return rate5
}

// Rate15 returns the fifteen-minute moving average rate of events per second.
func (m *StandardMeter) Rate15() float64 {
	m.lock.RLock()
	rate15 := m.snapshot.rate15
	m.lock.RUnlock()
	return rate15
}

// RateMean returns the meter's mean rate of events per second.
func (m *StandardMeter) RateMean() float64 {
	m.lock.RLock()
	rateMean := m.snapshot.rateMean
	m.lock.RUnlock()
	return rateMean
}

// RateStep returns the meter's step rate of events per second
func (m *StandardMeter) RateStep() float64 {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.updateSnapshotOnStep()
	return m.snapshot.rateStep
}

// Snapshot updates rate.step and returns a read-only copy of the meter.
func (m *StandardMeter) Snapshot() Meter {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.updateSnapshotOnStep()
	snapshot := *m.snapshot
	return &snapshot
}

func (m *StandardMeter) updateSnapshot() {
	// should run with write lock held on m.lock
	snapshot := m.snapshot
	snapshot.rate1 = m.a1.Rate()
	snapshot.rate5 = m.a5.Rate()
	snapshot.rate15 = m.a15.Rate()
	snapshot.rateMean = float64(snapshot.count) / time.Since(m.startTime).Seconds()
}

func (m *StandardMeter) updateSnapshotOnStep() {
	// should run with write lock held on m.lock
	snapshot := m.snapshot
	snapshot.rate1 = m.a1.Rate()
	snapshot.rate5 = m.a5.Rate()
	snapshot.rate15 = m.a15.Rate()

	nw := time.Now()
	sub := nw.Sub(m.startTime).Seconds()
	step := nw.Sub(snapshot._lastTime).Seconds()

	if sub > 0 {
		snapshot.rateMean = float64(snapshot.count) / sub
	}

	if step > 0 {
		snapshot.rateStep = float64(snapshot.count-snapshot._lastCount) / step
	}
	snapshot._lastCount = snapshot.count
	snapshot._lastTime = nw
}

func (m *StandardMeter) tick(now time.Time) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.a1.Tick()
	m.a5.Tick()
	m.a15.Tick()
	m.updateSnapshot()
}

type meterArbiter struct {
	sync.RWMutex
	started bool
	meters  []*StandardMeter
	ticker  *time.Ticker
}

var arbiter = meterArbiter{ticker: time.NewTicker(5e9)}

// Ticks meters on the scheduled interval
func (ma *meterArbiter) tick() {
	for {
		select {
		case nw := <-ma.ticker.C:
			ma.tickMeters(nw)
		}
	}
}

func (ma *meterArbiter) tickMeters(nw time.Time) {
	ma.RLock()
	defer ma.RUnlock()
	for _, meter := range ma.meters {
		meter.tick(nw)
	}
}
