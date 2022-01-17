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

package resp

import (
	"net/http"

	"github.com/quanxiang-cloud/cabin/logger"

	"github.com/go-playground/validator/v10"
	e "github.com/quanxiang-cloud/cabin/error"
)

type context interface {
	JSON(code int, obj interface{})
}

// Resp unified encapsulation of business return values.
type Resp struct {
	e.Error

	Data interface{} `json:"data"`
}

// Context write http status
func (r *Resp) Context(c context, code ...int) {
	status := http.StatusOK
	if r.Code == e.Unknown {
		status = http.StatusInternalServerError
	} else if len(code) != 0 {
		status = code[0]
	}
	c.JSON(status, r)
}

// Format formatting return value
func Format(data interface{}, err error) (r *Resp) {
	r = new(Resp)
	if err == nil {
		r.Code = e.Success
		r.Data = data
		return
	}

	var fail = func(r *Resp, err error) *Resp {
		r.Code = e.Unknown
		logger.Logger.PutError(err, "resp-known-error")
		return r
	}

	switch err := err.(type) {
	case e.Error:
		r.Code = err.Code
	case *e.Error:
		if err == nil {
			return fail(r, nil)
		}
		r.Code = err.Code
	case validator.ValidationErrors:
		r.Code = e.ErrParams
	default:
		return fail(r, err)
	}

	r.Message = err.Error()
	return
}
