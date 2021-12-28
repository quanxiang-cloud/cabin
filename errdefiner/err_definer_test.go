package errdefiner_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/quanxiang-cloud/cabin/errdefiner"
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
