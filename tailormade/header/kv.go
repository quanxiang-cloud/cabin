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

package header

import "context"

const (
	requestID = "Request-Id"

	timezone = "Timezone"
)

// MutateContext return context.Context,
// carry request id and timezone if value exists.
func MutateContext(c CTX) context.Context {
	var (
		_requestID interface{} = "Request-Id"
		_timezone  interface{} = "Timezone"
	)
	ctx := context.Background()
	ctx = context.WithValue(ctx, _requestID, c.GetHeader(requestID))
	ctx = context.WithValue(ctx, _timezone, c.GetHeader(timezone))
	return ctx
}

type KV []string

func (k KV) Wreck() (string, string) {
	switch len(k) {
	case 0:
		return "", ""
	case 1:
		return k[0], ""
	default:
		return k[0], k[1]
	}
}

// Fuzzy return key and value as []interface
func (k KV) Fuzzy() (result []interface{}) {
	for _, elem := range k {
		result = append(result, elem)
	}
	return
}

// GetRequestIDKV return request id
func GetRequestIDKV(ctx context.Context) KV {
	i := ctx.Value(requestID)
	rid, ok := i.(string)
	if ok {
		return KV{requestID, rid}
	}
	return KV{requestID, "unexpected type"}
}

// GetTimezone return timezone
func GetTimezone(ctx context.Context) KV {
	i := ctx.Value(timezone)
	tz, ok := i.(string)
	if ok {
		return KV{timezone, tz}
	}
	return KV{timezone, "unexpected type"}
}
