package httputil

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	error2 "github.com/quanxiang-cloud/cabin/error"
)

// BindBody bind gin body
func BindBody(c *gin.Context, d interface{}) error {
	b := binding.Default(c.Request.Method, c.ContentType())
	bb, ok := b.(binding.BindingBody)
	if !ok {
		return error2.NewErrorWithString(error2.ErrParams, "binding type error:"+c.ContentType())
	}
	if err := c.ShouldBindBodyWith(d, bb); err != nil {
		return error2.NewErrorWithString(error2.ErrParams, err.Error())
	}
	return nil
}

// IsQueryMethod check if http method is query
func IsQueryMethod(method string) bool {
	switch method {
	// NOTE: parameter is in query GET, DELETE, HEAD
	case http.MethodGet, http.MethodDelete, http.MethodHead:
		return true
	}
	return false
}

// GetRequestArgs get request args
func GetRequestArgs(c *gin.Context, d interface{}) error {
	if d == nil {
		d = &json.RawMessage{}
	}

	// get query only in GET requestion
	if IsQueryMethod(c.Request.Method) {
		q := c.Request.URL.Query()
		raw := QueryToBody(q, false)
		return json.Unmarshal([]byte(raw), d)
	}

	return BindBody(c, d)
}
