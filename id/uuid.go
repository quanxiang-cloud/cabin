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
	"encoding/base64"

	"github.com/google/uuid"
)

// BaseUUID generate uuid by base64 encoding
func BaseUUID() string {
	id := base64Coder.EncodeToString(BytesUUID())
	return id
}

// HexUUID generate uuid by hex encoding
func HexUUID(upper bool) string {
	id := encodeToHexString(BytesUUID(), upper)
	return id
}

// StringUUID generate google/uuid as string
func StringUUID() string {
	return uuid.New().String()
}

// BytesUUID generate google/uuid as string
func BytesUUID() []byte {
	b, err := uuid.New().MarshalBinary()
	if err != nil {
		panic(err)
	}
	return b
}

// WithPrefix add prefix to id if prefix is not empty
func WithPrefix(id string, prefix string) string {
	if prefix != "" {
		id = prefix + id
	}
	return id
}

//------------------------------------------------------------------------------

var base64Coder = base64.RawURLEncoding

const hextable = "0123456789abcdef"
const hextableU = "0123456789ABCDEF"

func encodeToHexString(src []byte, upper bool) string {
	dst := make([]byte, len(src)*2)
	encodeHex(dst, src, upper)
	return string(dst)
}

func encodeHex(dst, src []byte, upper bool) int {
	tb := hextable
	if upper {
		tb = hextableU
	}

	for i, v := range src {
		j := (i << 1)
		dst[j] = tb[v>>4]
		dst[j+1] = tb[v&0x0f]
	}
	return len(src) * 2
}
