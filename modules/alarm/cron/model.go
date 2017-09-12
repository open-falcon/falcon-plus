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

package cron

type MailDto struct {
	Priority int    `json:"priority"`
	Metric   string `json:"metric"`
	Subject  string `json:"subject"`
	Content  string `json:"content"`
	Email    string `json:"email"`
	Status   string `json:"status"`
}

type SmsDto struct {
	Priority int    `json:"priority"`
	Metric   string `json:"metric"`
	Content  string `json:"content"`
	Phone    string `json:"phone"`
	Status   string `json:"status"`
}

type ImDto struct {
	Priority int    `json:"priority"`
	Metric   string `json:"metric"`
	Content  string `json:"content"`
	IM       string `json:"im"`
	Status   string `json:"status"`
}
