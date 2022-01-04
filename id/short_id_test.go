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

package id

import (
	"fmt"
	"testing"
)

const (
	testNameLen = 8
	maxTest     = 10000_00
)

/*
Testlength=5
ShortID:      len=5 1864/1000000  duplicate 0.19%
UpperShortID: len=5 30827/1000000 duplicate 3.08%
String:       len=5 34200/1000000 duplicate 3.42%

Testlength=6
ShortID:      len=6 40/1000000   duplicate 0.00%
UpperShortID: len=6 1093/1000000 duplicate 0.11%
String:       len=6 1278/1000000 duplicate 0.13%

Testlength=8
String:       len=8 2/1000000    duplicate 0.00%
*/
func testRandID(name string, t *testing.T, fn func(int) string) {
	m := map[string]struct{}{}

	errCnt := 0
	const reportRound = maxTest / 100
	for i := 0; i < maxTest; i++ {
		id := fn(testNameLen)
		if _, ok := m[id]; ok {
			if errCnt%reportRound == 0 {
				fmt.Printf("%s: duplicate short-id %s at %d/%d\n", name, id, i+1, maxTest)
			}
			errCnt++
		}
		m[id] = struct{}{}
		if i == 0 {
			fmt.Println(name, ":", id)
		}
	}
	if errCnt > 0 {
		fmt.Printf("************%s: len=%d %d/%d duplicate %.2f%%************\n",
			name, testNameLen, errCnt, maxTest, float64(errCnt)/float64(maxTest)*100)
	}
}

func TestRandID(t *testing.T) {
	testRandID("ShortID", t, ShortID)
	testRandID("UpperShortID", t, UpperShortID)
	testRandID("String", t, String)
}
