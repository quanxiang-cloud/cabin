package resp

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	e "github.com/quanxiang-cloud/cabin/error"
)

type context interface {
	JSON(code int, obj interface{})
}

type Resp struct {
	e.Error

	Data interface{} `json:"data"`
}

func (r *Resp) Context(c context, code ...int) {
	status := http.StatusOK
	if r.Code == e.Unknown {
		status = http.StatusInternalServerError
	} else if len(code) != 0 {
		status = code[0]
	}
	c.JSON(status, r)
}

func Format(data interface{}, err error) (r *Resp) {
	r = new(Resp)
	if err == nil {
		r.Code = e.Success
		r.Data = data
		return
	}

	var fail = func(err error) (r *Resp) {
		r.Code = e.Unknown
		return
	}

	switch err := err.(type) {
	case e.Error:
		r.Code = err.Code
	case *e.Error:
		if err == nil {
			return fail(nil)
		}
		r.Code = err.Code
	case validator.ValidationErrors:
		r.Code = e.ErrParams
	default:
		return fail(err)
	}

	r.Message = err.Error()
	return
}
