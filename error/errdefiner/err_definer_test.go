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

package errdefiner_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/quanxiang-cloud/cabin/error/errdefiner"
)

func TestErrorDefiner(t *testing.T) {
	var errs = []errdefiner.ErrorCode{
		errdefiner.MustReg(1, "err1"),
		errdefiner.MustReg(2, "err2: %v"),
	}
	for _, v := range errs {
		m := v.Msg()
		var err error
		if strings.Contains(m, "%") {
			err = v.FmtError("foo")
		} else {
			err = v.NewError()
		}
		fmt.Printf("%v\n", err)
	}
}
