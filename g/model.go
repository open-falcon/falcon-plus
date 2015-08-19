package g

import (
	"fmt"
)

type TransferResp struct {
	Msg        string
	Total      int
	ErrInvalid int
	Latency    int64
}

func (t *TransferResp) String() string {
	s := fmt.Sprintf("TransferResp total=%d, err_invalid=%d, latency=%dms",
		t.Total, t.ErrInvalid, t.Latency)
	if t.Msg != "" {
		s = fmt.Sprintf("%s, msg=%s", s, t.Msg)
	}
	return s
}
