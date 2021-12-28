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
	"time"
)

const (
	// ISO8601 ISO8601
	ISO8601 = "2006-01-02T15:04:05.999Z"
)

// Time time now utc
func Time() time.Time {
	return time.Now().UTC()
}

// NowUnix current UTC timestamp.
func NowUnix() int64 {
	return time.Now().UTC().UnixNano() / 1e6
}

// Now current time in ISO8601 format.
func Now() string {
	return time.Now().UTC().Format(ISO8601)
}

// Format timestamp to time string.
func Format(t int64) string {
	return time.Unix(0, t*1e6).UTC().Format(ISO8601)
}

// Unix time string to timestamp.
func Unix(ts string) (int64, error) {
	t, err := time.Parse(ISO8601, ts)
	if err != nil {
		return 0, err
	}

	return t.UTC().UnixNano() / 1e6, nil
}
