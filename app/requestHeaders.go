package app

import "fmt"

type RequestHeaders struct {
	Method         string
	Path           string
	Host           string
	Agent          string
	ContentType    string
	ContentLength  string
	AcceptEncoding string
	Body           string
}

func (rh *RequestHeaders) Inspect() string {
	return fmt.Sprintf("%+v", rh)
}
