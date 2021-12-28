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

package error

const (
	// Unknown unknown error
	Unknown int64 = -1
	// Internal internal server error
	Internal int64 = -2

	// Success success
	Success int64 = 0
	// ErrParams parameter error
	ErrParams int64 = 1
)

// Table map[code]opywriting
type Table map[int64]string

// CodeTable error value comparison table
var CodeTable Table

var baseCode = map[int64]string{
	Unknown:   "unknown err",
	Internal:  "internal server error",
	Success:   "success",
	ErrParams: "parameter error",
}

// Translation translation code to message
func Translation(code int64) string {
	if CodeTable != nil {
		if text, ok := CodeTable[code]; ok {
			return text
		}
	}
	if text, ok := baseCode[code]; ok {
		return text
	}
	return "unknown code."
}
