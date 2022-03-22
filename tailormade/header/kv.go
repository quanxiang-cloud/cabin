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

// predefined header
const (
	RequestID = "Request-Id"

	Timezone = "Timezone"

	TenantID = "Tenant-Id"
)

// MutateContext return context.Context,
// carry request id and timezone if value exists.
func MutateContext(c CTX) context.Context {
	var (
		_requestID interface{} = "Request-Id"
		_timezone  interface{} = "Timezone"
		_tenantID  interface{} = "Tenant-Id"
	)
	ctx := context.Background()
	ctx = context.WithValue(ctx, _requestID, c.GetHeader(RequestID))
	ctx = context.WithValue(ctx, _timezone, c.GetHeader(Timezone))
	ctx = context.WithValue(ctx, _tenantID, c.GetHeader(TenantID))
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
	i := ctx.Value(RequestID)
	rid, ok := i.(string)
	if ok {
		return KV{RequestID, rid}
	}
	return KV{RequestID, "unexpected type"}
}

// GetTimezone return timezone
func GetTimezone(ctx context.Context) KV {
	i := ctx.Value(Timezone)
	tz, ok := i.(string)
	if ok {
		return KV{Timezone, tz}
	}
	return KV{Timezone, "unexpected type"}
}

// GetTenantID return tenantID
func GetTenantID(ctx context.Context) KV {
	i := ctx.Value(TenantID)
	tid, ok := i.(string)
	if ok {
		return KV{TenantID, tid}
	}
	return KV{TenantID, "unexpected type"}
}
