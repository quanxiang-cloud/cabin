/*
Copyright 2022 QuanxiangCloud Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
     http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package time

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

//go:generate stringer -type UTC
type UTC int

const (
	UTC_12 UTC = iota - 12
	UTC_11
	UTC_10
	UTC_9
	UTC_8
	UTC_7
	UTC_6
	UTC_5
	UTC_4
	UTC_3
	UTC_2
	UTC_1
	UTC0
	UTC1
	UTC2
	UTC3
	UTC4
	UTC5
	UTC6
	UTC7
	UTC8
	UTC9
	UTC10
	UTC11
	UTC12
)

var (
	ErrFormat = errors.New("error format")
)

func Tolerant(str string) (UTC, error) {
	var l int
	if strings.Contains(str, "-") {
		l = strings.LastIndex(str, "-")
	} else {
		l = strings.LastIndex(str, "+")
	}

	if l == -1 {
		return UTC0, ErrFormat
	}

	i, err := strconv.Atoi(str[l:])
	if err != nil {
		return UTC0, nil
	}
	return UTC(i), nil
}

// Revise correct the time to the specified time zone, only 0 time zone.
func Revise(ts string, timezone UTC) (string, error) {
	t, err := time.Parse(ISO8601, ts)
	if err != nil {
		return "", err
	}
	return t.Add(time.Hour * time.Duration(timezone)).Format(ISO8601), nil
}

// Regular modified to 0 time zone time.
func Regular(ts string, timezone UTC) (string, error) {
	return Revise(ts, -timezone)
}
