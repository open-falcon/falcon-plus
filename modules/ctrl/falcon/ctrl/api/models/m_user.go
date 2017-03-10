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
	"time"

	"github.com/astaxie/beego/orm"
)

type User struct {
	Id         int64     `json:"id"`
	Uuid       string    `json:"uuid"`
	Name       string    `json:"name"`
	Cname      string    `json:"cname"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	Im         string    `json:"im"`
	Qq         string    `json:"qq"`
	CreateTime time.Time `json:"ctime"`
}

type Operator struct {
	User  *User
	Token int
	O     orm.Ormer
}

type OperatorInfo struct {
	User     *User `json:"user"`
	Reader   bool  `json:"reader"`
	Operator bool  `json:"operator"`
	Admin    bool  `json:"admin"`
}

func (op *Operator) Info() *OperatorInfo {
	if op.User == nil {
		return &OperatorInfo{}
	}
	return &OperatorInfo{
		User:     op.User,
		Reader:   op.IsReader(),
		Operator: op.IsOperator(),
		Admin:    op.IsAdmin(),
	}
}
func (op *Operator) IsAdmin() bool {
	return (op.Token & SYS_F_A_TOKEN) != 0
}

func (op *Operator) IsOperator() bool {
	return (op.Token & SYS_F_O_TOKEN) != 0
}

func (op *Operator) IsReader() bool {
	return (op.Token & SYS_F_R_TOKEN) != 0
}

func (op *Operator) UserTokens() (token int) {
	var (
		tids []int64
	)

	if op.User == nil {
		return 0
	}

	_, err := op.O.Raw(`
   SELECT b1.token_id
    FROM (SELECT a1.tag_id AS user_tag_id,
                a2.tag_id AS token_tag_id,
                a1.tpl_id AS role_id,
                a1.sub_id AS user_id,
                a2.sub_id AS token_id
          FROM tpl_rel a1
          JOIN tpl_rel a2
          ON a1.type_id = ? AND a1.sub_id = ? AND a2.type_id = ?
              AND a2.sub_id in (?, ?, ?) AND a1.tpl_id = a2.tpl_id) b1
    JOIN tag_rel b2
    ON b1.user_tag_id = b2.tag_id AND b1.token_tag_id = b2.sup_tag_id
    GROUP BY b1.token_id`,
		TPL_REL_T_ACL_USER, op.User.Id, TPL_REL_T_ACL_TOKEN,
		SYS_IDX_R_TOKEN, SYS_IDX_O_TOKEN,
		SYS_IDX_A_TOKEN).QueryRows(&tids)
	if err != nil {
		return 0
	}
	for _, tid := range tids {
		switch tid {
		case SYS_IDX_R_TOKEN:
			token |= SYS_F_R_TOKEN
		case SYS_IDX_O_TOKEN:
			token |= SYS_F_O_TOKEN
		case SYS_IDX_A_TOKEN:
			token |= SYS_F_A_TOKEN
		}
	}

	if op.User.Id < 3 {
		token |= SYS_F_A_TOKEN
	}

	if token&SYS_F_A_TOKEN != 0 {
		token |= SYS_F_O_TOKEN
	}
	if token&SYS_F_O_TOKEN != 0 {
		token |= SYS_F_R_TOKEN
	}

	// for dev
	if op.User.Name == "test" {
		token = SYS_F_A_TOKEN | SYS_F_O_TOKEN | SYS_F_R_TOKEN
		//token = SYS_F_O_TOKEN | SYS_F_R_TOKEN
		//token = SYS_F_R_TOKEN
	}

	return token
}

func (op *Operator) AddUser(user *User) (*User, error) {
	user.Id = 0
	id, err := op.O.Insert(user)
	if err != nil {
		return nil, err
	}
	user.Id = id
	moduleCache[CTL_M_USER].set(id, user)

	DbLog(op.O, op.User.Id, CTL_M_USER, id, CTL_A_ADD, jsonStr(user))
	return user, nil
}

// just called from profileFilter()
func GetUser(id int64, o orm.Ormer) (*User, error) {
	if user, ok := moduleCache[CTL_M_USER].get(id).(*User); ok {
		return user, nil
	}
	user := &User{Id: id}
	err := o.Read(user, "Id")
	if err == nil {
		moduleCache[CTL_M_USER].set(id, user)
	}
	return user, err
}

func (op *Operator) GetUser(id int64) (*User, error) {
	return GetUser(id, op.O)
}

func (op *Operator) GetUserByUuid(uuid string) (user *User, err error) {
	user = &User{Uuid: uuid}
	err = op.O.Read(user, "Uuid")
	return user, err
}

func (op *Operator) QueryUsers(query string) orm.QuerySeter {
	qs := op.O.QueryTable(new(User))
	if query != "" {
		qs = qs.SetCond(orm.NewCondition().Or("Name__icontains", query).Or("Email__icontains", query))
	}
	return qs
}

func (op *Operator) GetUsersCnt(query string) (int64, error) {
	return op.QueryUsers(query).Count()
}

func (op *Operator) GetUsers(query string, limit, offset int) (users []*User, err error) {
	_, err = op.QueryUsers(query).Limit(limit, offset).All(&users)
	return
}

func (op *Operator) UpdateUser(id int64, _u *User) (user *User, err error) {
	if user, err = op.GetUser(id); err != nil {
		return nil, ErrNoUsr
	}

	if _u.Name != "" && user.Name == "" {
		user.Name = _u.Name
	}
	if _u.Cname != "" {
		user.Cname = _u.Cname
	}
	if _u.Email != "" {
		user.Email = _u.Email
	}
	if _u.Phone != "" {
		user.Phone = _u.Phone
	}
	if _u.Im != "" {
		user.Im = _u.Im
	}
	if _u.Qq != "" {
		user.Qq = _u.Qq
	}
	_, err = op.O.Update(user)
	moduleCache[CTL_M_USER].set(id, user)
	DbLog(op.O, op.User.Id, CTL_M_USER, id, CTL_A_SET, "")
	return user, err
}

func (op *Operator) DeleteUser(id int64) error {
	if n, err := op.O.Delete(&User{Id: id}); err != nil || n == 0 {
		return ErrNoExits
	}
	moduleCache[CTL_M_USER].del(id)
	DbLog(op.O, op.User.Id, CTL_M_USER, id, CTL_A_DEL, "")

	return nil
}
