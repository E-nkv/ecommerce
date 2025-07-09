package customErrors

import "fmt"

var (
	NotFound error = fmt.Errorf("not found error")
	Internal       = fmt.Errorf("internal server error")
)
