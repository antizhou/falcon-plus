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

package model

import (
	"fmt"
	"encoding/json"
)

type EsItem struct {
	Metric    string  `json:"metric"`
	Tags      string  `json:"tags"`
	Value     float64 `json:"value"`
	Timestamp int64   `json:"timestamp"`
	Endpoint  string  `json:"endpoint"`
}

func (this *EsItem) String() string {
	return fmt.Sprintf(
		"<Metric:%s, Tags:%v, Value:%v, TS:%d>",
		this.Metric,
		this.Tags,
		this.Value,
		this.Timestamp,
	)
}

func (this *EsItem) EsString() (s string) {
	bs, err := json.Marshal(this)
	if err != nil {
		return ""
	}

	return string(bs)
}
