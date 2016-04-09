package model

import (
	"fmt"
)

// code == 0 => success
// code == 1 => bad request
type SimpleRpcResponse struct {
	Code int `json:"code"`
}

func (this *SimpleRpcResponse) String() string {
	return fmt.Sprintf("<Code: %d>", this.Code)
}

type NullRpcRequest struct {
}
