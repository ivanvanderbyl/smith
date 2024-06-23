package httpdebugutil

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"
)

// Transport is an http.RoundTripper that logs request and response information.
//
// Usage:
//
//	client := &http.Client{
//	    Transport: &httpdebugutil.Transport{},
//	}
type Transport struct{}

func (s *Transport) RoundTrip(r *http.Request) (*http.Response, error) {
	output := bytes.NewBuffer(nil)

	bytes, _ := httputil.DumpRequestOut(r, true)
	for _, line := range strings.Split(string(bytes), "\n") {
		output.WriteString("<< " + line + "\n")
	}

	resp, err := http.DefaultTransport.RoundTrip(r)
	// err is returned after dumping the response

	respBytes, _ := httputil.DumpResponse(resp, true)
	for _, line := range strings.Split(string(respBytes), "\n") {
		output.WriteString(">> " + line + "\n")
	}

	fmt.Println(output.String())

	return resp, err
}
