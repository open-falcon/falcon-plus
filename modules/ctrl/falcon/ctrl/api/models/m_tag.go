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
package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type Tag struct {
	Id         int64     `json:"id"`
	Name       string    `json:"name"`
	CreateTime time.Time `json:"ctime"`
}

type Tag_rel struct {
	Id       int64
	TagId    int64
	SupTagId int64
	Offset   int64
}

type TagNode struct {
	Key  string
	Must bool
}

type TagSchema struct {
	data  string
	nodes []TagNode
}

//cop,owt,pdl,servicegroup;service,jobgroup;job,sbs;mod;srv;grp;cluster;loc;idc;status;
// ',' : must
// ';' : or
func NewTagSchema(tag string) (*TagSchema, error) {
	var (
		i, j int
	)

	if tag == "" {
		return nil, nil
	}

	ret := &TagSchema{data: tag}
	for i, j = 0, 0; j < len(tag); j++ {
		if tag[j] == ',' {
			if i >= j {
				return nil, ErrParam
			}
			ret.nodes = append(ret.nodes, TagNode{
				Key:  strings.TrimSpace(tag[i:j]),
				Must: true,
			})
			i = j + 1
		} else if tag[j] == ';' {
			if i >= j {
				return nil, ErrParam
			}
			ret.nodes = append(ret.nodes, TagNode{
				Key:  strings.TrimSpace(tag[i:j]),
				Must: false,
			})
			i = j + 1
		}
	}
	if i != j || i == 0 {
		return nil, ErrParam
	}

	return ret, nil
}

func tagMap(tag string) (map[string]string, error) {
	var (
		i, j int
		k, v string
		ret  = make(map[string]string)
	)

	for i, j = 0, 0; j < len(tag); j++ {
		if tag[j] == '=' {
			k = strings.TrimSpace(tag[i:j])
			i = j + 1
		} else if tag[j] == ',' {
			v = strings.TrimSpace(tag[i:j])
			if len(k) > 0 && len(v) > 0 {
				ret[k] = v
				k, v = "", ""
			} else {
				return ret, ErrParam
			}
			i = j + 1
		}
	}

	v = strings.TrimSpace(tag[i:])
	if len(k) > 0 && len(v) > 0 {
		ret[k] = v
		return ret, nil
	} else {
		return ret, ErrParam
	}
}

func (ts *TagSchema) Fmt(tag string, force bool) (string, error) {
	var (
		ret string
		n   int
	)

	if tag == "" {
		return "", nil
	}

	m, err := tagMap(tag)
	if err != nil {
		return "", err
	}
	for _, node := range ts.nodes {
		if v, ok := m[node.Key]; ok {
			ret += fmt.Sprintf("%s=%s,", node.Key, v)
			n++
		} else if !force && node.Must {
			return ret, ErrParam
		}

		// done
		if n == len(m) {
			return ret[:len(ret)-1], nil
		}
	}

	// some m.key miss match
	if force && len(ret) > 1 {
		return ret[0 : len(ret)-1], nil
	}

	return ret, ErrParam
}

func TagRelation(t string) (ret []string) {

	if t == "" {
		return []string{""}
	}

	tags := strings.Split(t, ",")
	if len(tags) < 1 {
		return []string{""}
	}
	ret = make([]string, len(tags)+1)

	for tag, i := "", 1; i < len(ret); i++ {
		tag += tags[i-1] + ","
		ret[i] = tag[:len(tag)-1]
	}
	return ret
}

func TagParents(t string) (ret []string) {

	tags := strings.Split(t, ",")
	if len(tags) < 1 {
		return nil
	}
	ret = make([]string, len(tags))

	for tag, i := "", 1; i < len(ret); i++ {
		tag += tags[i-1] + ","
		ret[i] = tag[:len(tag)-1]
	}
	return ret
}

func TagParent(t string) string {
	if i := strings.LastIndexAny(t, ","); i > 0 {
		return t[:i]
	} else {
		return ""
	}
}

func TagLast(t string) string {
	return t[strings.LastIndexAny(t, ",")+1:]
}

func (op *Operator) addTag(t *Tag, schema *TagSchema) (id int64, err error) {
	if schema != nil {
		t.Name, err = schema.Fmt(t.Name, false)
		if err != nil {
			return
		}
	}

	// TODO: check parent exist/acl
	if _, err = op.Access(SYS_O_TOKEN,
		TagParent(t.Name), false); err != nil {
		return
	}

	t.Id = 0
	if id, err = op.O.Insert(t); err != nil {
		return
	}
	fmt.Printf("id: %d tag: %s\n", id, t.Name)

	if rels := TagRelation(t.Name); len(rels) > 0 {
		var (
			tags []*Tag
			arg  = make([]interface{}, len(rels))
		)

		for i, v := range rels {
			arg[i] = v
		}

		_, err = op.O.QueryTable(new(Tag)).
			Filter("Name__in", arg...).All(&tags)
		if err != nil {
			return
		}
		tag_rels := make([]Tag_rel, len(tags))
		for i, tag := range tags {
			tag_rels[i] = Tag_rel{
				TagId:    id,
				SupTagId: tag.Id,
				Offset:   int64(len(tags) - 1 - i)}
		}
		_, err = op.O.InsertMulti(10, tag_rels)
		if err != nil {
			return
		}
	}

	t.Id = id
	moduleCache[CTL_M_TAG].set(id, t)
	DbLog(op.O, op.User.Id, CTL_M_TAG, id, CTL_A_ADD, "")

	return id, err
}

func (op *Operator) AddTag(t *Tag) (id int64, err error) {
	return op.addTag(t, sysTagSchema)
}

func (op *Operator) GetTag(id int64) (*Tag, error) {
	if t, ok := moduleCache[CTL_M_TAG].get(id).(*Tag); ok {
		return t, nil
	}
	t := &Tag{Id: id}
	err := op.O.Read(t, "Id")
	if err == nil {
		moduleCache[CTL_M_TAG].set(id, t)
	}
	return t, err
}

func (op *Operator) GetTagByName(tag string) (t *Tag, err error) {
	t = &Tag{Name: tag}

	err = op.O.Read(t, "Name")
	return
}

func (op *Operator) QueryTags(query string) orm.QuerySeter {
	// TODO: acl filter
	qs := op.O.QueryTable(new(Tag))
	if query != "" {
		qs = qs.Filter("Name__icontains", query)
	}
	return qs
}

func (op *Operator) GetTagsCnt(query string) (int64, error) {
	return op.QueryTags(query).Count()
}

func (op *Operator) GetTags(query string, limit, offset int) (tags []*Tag, err error) {
	_, err = op.QueryTags(query).Limit(limit, offset).All(&tags)
	return
}

func (op *Operator) UpdateTag(id int64, _t *Tag) (t *Tag, err error) {
	if t, err = op.GetTag(id); err != nil {
		return
	}

	if _, err = op.Access(SYS_O_TOKEN,
		TagParent(t.Name), true); err != nil {
		return
	}

	if _t.Name != "" {
		t.Name = _t.Name
	}
	_, err = op.O.Update(t)
	moduleCache[CTL_M_TAG].set(id, t)
	DbLog(op.O, op.User.Id, CTL_M_TAG, id, CTL_A_SET, "")
	return t, err
}

func (op *Operator) DeleteTag(id int64) (err error) {
	var n int64
	var tag *Tag

	if tag, err = op.GetTag(id); err != nil {
		return
	}

	if _, err = op.Access(SYS_O_TOKEN,
		TagParent(tag.Name), false); err != nil {
		return
	}

	if n, err = op.O.Delete(&Tag{Id: id}); err != nil || n == 0 {
		return ErrNoExits
	}
	moduleCache[CTL_M_TAG].del(id)
	DbLog(op.O, op.User.Id, CTL_M_TAG, id, CTL_A_DEL, "")

	return nil
}
