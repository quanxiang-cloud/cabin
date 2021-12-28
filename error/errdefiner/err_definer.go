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

// Package errdefiner enable centralized definition of error code.
package errdefiner

import (
	"fmt"
	"strings"

	errman "github.com/quanxiang-cloud/cabin/error"
)

// NOTE: create only one errorDefiner within an application on initialize
var inst = errorDefiner{
	codeTable: make(map[int64]string),
	cacheErr:  make(map[ErrorCode]*errman.Error),
}

func init() {
	errman.CodeTable = inst.codeTable
}

// ErrorDefiner object
type errorDefiner struct {
	codeTable map[int64]string            //  error code list
	cacheErr  map[ErrorCode]*errman.Error //cache fix error
}

// MustReg regist an error with duplicate code check
func MustReg(code int64, msg string) ErrorCode {
	return inst.MustReg(code, msg)
}

// MustReg regist an error with duplicate code check
func (r *errorDefiner) MustReg(code int64, msg string) ErrorCode {
	if _, ok := r.codeTable[code]; ok {
		panic(fmt.Errorf("duplicate code %d", code))
	}
	r.codeTable[code] = msg
	c := ErrorCode(code)
	if !withFormat(msg) { // without format parameter
		_ = c.NewError() // generate cacheErr
	}
	return c
}

func (r *errorDefiner) newError(c ErrorCode) *errman.Error {
	if cache, ok := r.cacheErr[c]; ok {
		return cache
	}

	err := errman.New(c.Int64())
	r.cacheErr[c] = &err // cache fix error
	return &err
}

func (r *errorDefiner) msg(c ErrorCode, paras []interface{}) string {
	if m, ok := r.codeTable[c.Int64()]; ok {
		return fmt.Sprintf(m, paras...)
	}
	return fmt.Sprintf("<unknown error code %d>", c)
}

func withFormat(s string) bool {
	return strings.IndexByte(s, '%') >= 0
}

// ErrorCode of int
type ErrorCode int64

// Int64 convert code to int64
func (c ErrorCode) Int64() int64 {
	return int64(c)
}

// WithFormat check if error message contains format
func (c ErrorCode) WithFormat() bool {
	return withFormat(c.Msg())
}

// NewError create an error without format
func (c ErrorCode) NewError() *errman.Error {
	return inst.newError(c)
}

// Msg format an error message string
func (c ErrorCode) Msg(paras ...interface{}) string {
	return inst.msg(c, paras)
}

// FmtError create an error with format
func (c ErrorCode) FmtError(paras ...interface{}) error {
	err := errman.New(c.Int64(), paras...)
	return &err
}

//------------------------------------------------------------------------------

// exports
const (
	ErrParams = errman.ErrParams
	Internal  = errman.Internal
	Unknown   = errman.Unknown
	Success   = errman.Success
)

// Errorf format an standard error with parameters
func Errorf(format string, paras ...interface{}) error {
	return fmt.Errorf(format, paras...)
}

// NewErrorWithString return an error with message
func NewErrorWithString(code int64, msg string) *errman.Error {
	err := errman.NewErrorWithString(code, msg)
	return &err
}
