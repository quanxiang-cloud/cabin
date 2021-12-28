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

import "fmt"

// Error error
type Error struct {
	Code    int64  `json:"code"`
	Message string `json:"msg,omitempty"`
}

// New return error with code and error message
func New(code int64, params ...interface{}) Error {
	if len(params) > 0 {
		return Error{
			Code:    code,
			Message: fmt.Sprintf(Translation(code), params...),
		}
	}
	return Error{
		Code: code,
	}
}

// Error return error string
func (e Error) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return Translation(e.Code)
}

// NewErrorWithString return error with code and string message
func NewErrorWithString(code int64, message string) Error {
	return Error{
		Code:    code,
		Message: message,
	}
}
