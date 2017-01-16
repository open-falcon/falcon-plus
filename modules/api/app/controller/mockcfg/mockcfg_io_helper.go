package mockcfg

import "errors"

type APICreateNoDataInputs struct {
	Name string `json:"name" binding:"required"`
	Obj  string `json:"obj" binding:"required"`
	//group, host, other
	ObjType string  `json:"obj_type" binding:"required"`
	Metric  string  `json:"metric" binding:"required"`
	Tags    string  `json:"tags" binding:"exists"`
	DsType  string  `json:"dstype" binding:"required"`
	Step    int     `json:"step" binding:"required"`
	Mock    float64 `json:"mock" binding:"exists"`
}

func (this APICreateNoDataInputs) CheckFormat() (err error) {
	switch {
	case this.ObjType != "group" && this.ObjType != "host" && this.ObjType != "other":
		err = errors.New("obj_type only accpect \"group, host, other\"")
	}
	return
}

type APIUpdateNoDataInputs struct {
	ID  int64  `json:"id" binding:"required"`
	Obj string `json:"obj" binding:"required"`
	//group, host, other
	ObjType string  `json:"obj_type" binding:"required"`
	Metric  string  `json:"metric" binding:"required"`
	Tags    string  `json:"tags" binding:"exists"`
	DsType  string  `json:"dstype" binding:"required"`
	Step    int     `json:"step" binding:"required"`
	Mock    float64 `json:"mock" binding:"exists"`
}

func (this APIUpdateNoDataInputs) CheckFormat() (err error) {
	switch {
	case this.ObjType != "group" && this.ObjType != "host" && this.ObjType != "other":
		err = errors.New("obj_type only accpect \"group, host, other\"")
	}
	return
}
