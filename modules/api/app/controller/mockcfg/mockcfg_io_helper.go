// Copyright 2017 Xiaomi, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
