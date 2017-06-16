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

package db

// graph.endpoint
type GraphEndpoint struct {
	Id       int64
	Endpoint string
}

// graph.tag_endpoint
type GraphTagEndpoint struct {
	Id         int64
	Tag        string
	EndpointId int64
}

// graph.endpoint_counter
type GraphEndpointCounter struct {
	Id         int64
	EndpointId int64
	Counter    string
}
