package g

import (
	"log"
	"runtime"
)

// TODO
// change graph.store cache struct(key: md5->uuid)
// flush when query happens seems unreasonable
// shrink packages

// CHANGE LOGS
// 0.4.8 fix filename bug emporarily, fix dirty-index-cache bug of query,
//		 add filter for debug
// 0.4.9 mv db back to g, add rpc.last
// 0.5.0 rm trace, add history&last api
// 0.5.1 add http interface v2, using form args
// 0.5.2 add last_raw
// 0.5.3 fix bug of last&last_raw
// 0.5.4 fix bug of Query.merge
// 0.5.5 use commom(rm model), fix sync disk
// 0.5.7 set xff to 0 from 0.5, in order to support irregular step counter
// 0.5.8 clean GraphItems/historyCache Cache at regular intervals
// 0.5.9 add flush style(flush by number of every counter's monitoring data)

const (
	VERSION         = "0.5.9"
	GAUGE           = "GAUGE"
	DERIVE          = "DERIVE"
	COUNTER         = "COUNTER"
	DEFAULT_STEP    = 60      //s
	MIN_STEP        = 30      //s
	CLEAN_CACHE     = 86400   //s the step that clean GraphItems/historyCache Cache
	CACHE_DELAY     = 1800    //s
	CACHE_TIME      = 1800000 //ms
	FLUSH_DISK_STEP = 1000    //ms
	FLUSH_MIN_COUNT = 6       //  flush counter to disk when its number of monitoring data greater than FLUSH_MIN_COUNT
	FLUSH_MAX_WAIT  = 86400   //s flush counter to disk if it not be flushed within FLUSH_MAX_WAIT seconds
)

const (
	GRAPH_F_MISS uint32 = 1 << iota
	GRAPH_F_ERR
	GRAPH_F_SENDING
	GRAPH_F_FETCHING
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}
