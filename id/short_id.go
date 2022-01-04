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
	"crypto/rand"
)

const (
	// DefaultShortNameLen is defult length of short name
	DefaultShortNameLen = 8
)
const (
	alphaTab   = "abcdfghjklmnpqrstvwxzABCDFGHJKLMNPQRSTVWXZ012456789"
	upperTab   = "ABCDFGHJKLMNPQRSTVWXZ012456789"
	trimBits   = 0
	tabLen     = len(alphaTab)
	headTabLen = tabLen - 10 // first byte dont allow number character

)

// ShortID  generate a random string with length n.
// If n<=0 it use n as default length(8).
// NOTE: we suggest n>=8, too short id with higher risk of collision.
func ShortID(n int) string {
	s, err := ShortIDWithError(n)
	if err != nil {
		panic(err)
	}
	return s
}

// UpperShortID  generate a random string with length n.
// If n<=0 it use n as default length(8).
// NOTE: we suggest n>=8, too short id with higher risk of collision.
// NOTE: it with much more probability to duplicate compare to ShortID.
func UpperShortID(n int) string {
	s, err := ShortIDWithDic(n, upperTab, len(upperTab)-10)
	if err != nil {
		panic(err)
	}
	return s
}

// ShortIDWithError  generate a random string with length n
func ShortIDWithError(n int) (string, error) {
	return ShortIDWithDic(n, alphaTab, headTabLen)
}

func withDefaultLength(n int) int {
	if n <= 0 {
		n = DefaultShortNameLen
	}
	return n
}

// ShortIDWithDic generate a shortID with dictionary.
// headDicLen is the length of first byte use.
// NOTE: it will panic if dic is empty or headDicLen overflow dic.
func ShortIDWithDic(n int, dic string, headDicLen int) (string, error) {
	n = withDefaultLength(n)
	b := make([]byte, n)
	if nr, err := rand.Read(b); err != nil || nr != len(b) {
		return "", err
	}

	mod := headDicLen
	dicLen := len(dic)
	for i, v := range b {
		idx := (int(v>>trimBits) % mod)
		mod = dicLen

		b[i] = dic[idx]
	}
	return string(b), nil
}
