package stdlib

import (
	"github.com/tengolang/tengo/v3"
)

func wrapError(err error) tengo.Object {
	if err == nil {
		return tengo.TrueValue
	}
	return &tengo.Error{Value: &tengo.String{Value: err.Error()}}
}
