package g

type Plugin struct {
	Path string
}

func (this *Plugin) String() string {
	return this.Path
}

type Host struct {
	HostId   int
	HostName string
	Pause    int
	Uuid     string
}

type AgentReq struct {
	Host
	Checksum string
}

type AgentPluginsResp struct {
	Plugins   []*Plugin
	HostName  string
	Checksum  string
	Timestamp int64
}

type BuiltinItem struct {
	Metric string
	Tags   string
}

type BuiltinItemResp struct {
	Items     []*BuiltinItem
	Checksum  string
	Timestamp int64
}
