package http

import (
	"encoding/json"
	"net/http"
	_ "net/http/pprof"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/open-falcon/falcon-plus/modules/falcon2p8s/g"
	"github.com/open-falcon/falcon-plus/modules/falcon2p8s/middlewares"
	"github.com/open-falcon/falcon-plus/modules/falcon2p8s/rpc"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	nsema "github.com/toolkits/concurrent/semaphore"
)

var router *gin.Engine

func RenderJson(w http.ResponseWriter, v interface{}) {
	bs, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(bs)
}

var counterMutex, gaugeMutex sync.Mutex

func dealGaugeMetric(p8sItem *rpc.P8sItem) {
	counter, ok := g.CollectorMap.Load(p8sItem.PK)
	// 启动之后首次接收到该指标
	if !ok {
		gaugeMutex.Lock()
		gaugeVec := prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: p8sItem.Metric,
			}, p8sItem.TagKeys)
		uncheckedCollector := g.UncheckedCollector{C: gaugeVec, MetricType: p8sItem.MetricType}
		counter, ok = g.CollectorMap.Load(p8sItem.PK)
		if !ok {
			if err := prometheus.Register(uncheckedCollector); err != nil {
				log.Errorf("Error Register gauge. %s", err.Error())
			} else {
				gv, err := gaugeVec.GetMetricWithLabelValues(p8sItem.TagValues...)
				if err == nil {
					gv.Set(p8sItem.Value)
					g.CollectorMap.Store(p8sItem.PK, uncheckedCollector)
				} else {
					log.Errorf("Error when GetMetricWithLabelValues first time. item: %v, err: %s\n", *p8sItem, err.Error())
				}
			}
		}
		gaugeMutex.Unlock()
	}
	// exporter启动之后非首次接收到该指标
	if counter != nil {
		uncheckedCollector, ok := counter.(g.UncheckedCollector)
		if ok {
			uncheckedCollector.C.(*prometheus.GaugeVec).WithLabelValues(p8sItem.TagValues...).Set(p8sItem.Value)
		} else {
			log.Errorf("Error assert to UncheckedCollector. item: %v", p8sItem)
		}
	}
	// 记录当前指标存储的时间
	g.LastUpdateTimeOfGauge.Store(p8sItem.PKWithTagValue, g.LastUpdateTimeItem{
		LastUpdateTime: time.Now(),
		PK:             p8sItem.PK,
		TagValues:      p8sItem.TagValues,
	})
}

func dealCounterMetric(p8sItem *rpc.P8sItem) {
	counter, ok := g.CollectorMap.Load(p8sItem.PK)
	// exporter启动之后首次接收到该指标
	if !ok {
		counterMutex.Lock()
		counterVec := prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: p8sItem.Metric,
			}, p8sItem.TagKeys)
		uncheckedCollector := g.UncheckedCollector{C: counterVec, MetricType: p8sItem.MetricType}
		counter, ok = g.CollectorMap.Load(p8sItem.PK)
		if !ok {
			if err := prometheus.Register(uncheckedCollector); err != nil {
				log.Errorf("Error Register counter. %s", err.Error())
			} else {
				c, err := counterVec.GetMetricWithLabelValues(p8sItem.TagValues...)
				if err == nil {
					if p8sItem.Value >= 0 {
						c.Add(p8sItem.Value)
						g.CollectorMap.Store(p8sItem.PK, uncheckedCollector)
					} else {
						log.Errorf("Error counter value less than zero: %v\n", *p8sItem)
					}
				} else {
					log.Errorf("Error when GetMetricWithLabelValues first time. item: %v, err: %s\n", *p8sItem, err.Error())
				}
			}
		}
		counterMutex.Unlock()
	}
	// exporter启动之后非首次接收到该指标
	if counter != nil {
		uncheckedCollector, ok := counter.(g.UncheckedCollector)
		if ok {
			// counter类型需要add相较上一次值的增量
			lastValue, ok := g.CounterCollectorValueMap.Load(p8sItem.PKWithTagValue)
			if ok {
				lv, ok := lastValue.(float64)
				if ok {
					if p8sItem.Value-lv >= 0 {
						uncheckedCollector.C.(*prometheus.CounterVec).WithLabelValues(p8sItem.TagValues...).Add(p8sItem.Value - lv)
					} else {
						log.Errorf("Error counter value cannot decrease. current value %f, last value: %f. item: %v\n", p8sItem.Value, lv, *p8sItem)
					}
				} else {
					log.Errorf("Error assert last value: %v\n", *p8sItem)
				}
			} else {
				if p8sItem.Value >= 0 {
					uncheckedCollector.C.(*prometheus.CounterVec).WithLabelValues(p8sItem.TagValues...).Add(p8sItem.Value)
				} else {
					log.Errorf("Error counter value less than zero: %v\n", *p8sItem)
				}
			}
		} else {
			log.Errorf("Error assert to UncheckedCollector. item: %v", p8sItem)
		}
	}
	g.CounterCollectorValueMap.Store(p8sItem.PKWithTagValue, p8sItem.Value)
	// 记录当前指标存储的时间
	g.LastUpdateTimeOfCounter.Store(p8sItem.PKWithTagValue, g.LastUpdateTimeItem{
		LastUpdateTime: time.Now(),
		PK:             p8sItem.PK,
		TagValues:      p8sItem.TagValues,
	})
}

func Start() {
	go func() {
		concurrent := g.Config().Concurrent
		if concurrent == 0 {
			concurrent = 100
		}
		sema := nsema.NewSemaphore(concurrent)
		for {
			sema.Acquire()
			// 如果prometheus开始抓取数据，就延迟一分钟之后再更新collector中的数据
			// 避免同一指标上一分钟的监控数据被覆盖，导致点位缺失
			if g.IsScraping {
				time.Sleep(time.Millisecond * 1)
				sema.Release()
				continue
			}
			item := g.P8sItemQueue.Shift()
			go func() {
				defer sema.Release()
				defer func() {
					if r := recover(); r != nil {
						log.Errorf("Exception is recover. %#v", r)
						return
					}
				}()
				if item == nil {
					time.Sleep(time.Millisecond * 1)
					return
				}
				p8sItem, ok := item.(*rpc.P8sItem)
				if !ok {
					log.Errorf("Error when assert item to p8sItem. %#v", item)
					return
				}
				// 检测是否存在相同metric，但type不同的监控项
				metricType, ok := g.MetricTypeMap.Load(p8sItem.Metric)
				if ok {
					if p8sItem.MetricType != metricType {
						log.Errorf("Error metric same, but metric type not coincide. item: %v\n", *p8sItem)
						return
					}
				} else {
					g.MetricTypeMap.Store(p8sItem.Metric, p8sItem.MetricType)
				}

				if p8sItem.MetricType == g.Counter {
					dealCounterMetric(p8sItem)
				} else {
					dealGaugeMetric(p8sItem)
				}
			}()
		}
	}()

	go getqueueSize()

	go cleanOutdatedGaugeMetrics()
	go cleanOutdatedCounterMetrics()

	router = gin.Default()

	configCommonRoutes()

	if g.Config().LogLevel != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}
	router.Use(middlewares.CheckIsScraping())
	router.GET("/metrics", prometheusHandler())

	addr := g.Config().Http.Listen
	if addr == "" {
		addr = "0.0.0.0:9090"
	}
	log.Println("http listening", addr)
	if err := router.Run(addr); err != nil {
		log.Fatal("Error gin run: ", err.Error())
		return
	}
}

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
