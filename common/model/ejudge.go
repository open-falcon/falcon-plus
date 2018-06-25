package model

type EHistoryData struct {
	Timestamp int64                  `json:"timestamp"`
	Filters   map[string]interface{} `json:"filters"`
}
