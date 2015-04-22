package g

type BuiltinItem struct {
	Metric string
	Tags   string
}

type BuiltinItemResp struct {
	Items     []*BuiltinItem
	Checksum  string
	Timestamp int64
}
