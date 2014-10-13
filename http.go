package nessusgo

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"runtime"
	"strconv"
	"time"
)

var (
	// Returned if the specified resource does not exist.
	ErrNotFound = errors.New("Not Found")

	// Returned if the caller attempts to make a call or modify a resource
	// for which the caller is not authorized.
	//
	// The request was a valid request, the caller's authentication credentials
	// succeeded but those credentials do not grant the caller permission to
	// access the resource.
	ErrForbidden = errors.New("Forbidden")

	// Returned if the call requires authentication and either the credentials
	// provided failed or no credentials were provided.
	ErrNotAuthorized = errors.New("Unauthorized")

	// Returned if the caller submits a badly formed request. For example,
	// the caller can receive this return if you forget a required parameter.
	ErrBadRequest = errors.New("Bad Request")
)

const (
	Version          = "v0.1"
	DefaultUserAgent = "Nessus.Go/" + Version + " (" + runtime.GOOS + "; " + runtime.GOARCH + ")"
)

// DefaultClient uses DefaultTransport, and is used internall to execute
// all http.Requests. This may be overriden for unit testing purposes.
//
// IMPORTANT: this is not thread safe and should not be touched with
// the exception overriding for mock unit testing.
var DefaultClient = http.DefaultClient

func (c *Client) do(method string, path string, params url.Values, values url.Values, v interface{}) error {
	// create the URI
	uri, err := url.Parse(c.ServiceUrl + path)
	if err != nil {
		return err
	}

	if params != nil && len(params) > 0 {
		uri.RawQuery = params.Encode()
	}

	// set the SSL validation for the request
	DefaultClient.Transport = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}

	// create the request
	req := &http.Request{
		URL:        uri,
		Method:     method,
		ProtoMajor: 1,
		ProtoMinor: 1,
		Close:      true,
	}
	req.Header = http.Header{}
	req.Header.Set("User-Agent", DefaultUserAgent)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Fix for truncated gzip response: https://stackoverflow.com/questions/21147562/unexpected-eof-using-go-http-client
	req.Header.Add("Accept-Encoding", "identity")

	if c.isAuthenticated() {
		cookie := &http.Cookie{Name: "token", Value: c.Token}
		req.AddCookie(cookie)
	}

	if values == nil {
		values = url.Values{}
	}

	rand.Seed(time.Now().UTC().UnixNano())
	seq := strconv.Itoa(rand.Intn(9999))

	values.Add("seq", seq)
	values.Add("json", "1")

	//fmt.Printf("Values: %v\n", values)

	if values != nil && len(values) > 0 {
		body := []byte(values.Encode())
		buf := bytes.NewBuffer(body)
		req.Body = ioutil.NopCloser(buf)
		req.ContentLength = int64(buf.Len())
	}

	// make the request using the default http client
	resp, err := DefaultClient.Do(req)
	if err != nil {
		return err
	}

	// Read the bytes from the body (make sure we defer close the body)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	//fmt.Printf("Response Code: %v\n", resp.StatusCode)
	if path == "/poo" {
		fmt.Printf("DEBUG body: %s\n", string(body))
	}

	// Check for an http error status (ie not 200 StatusOK)
	switch resp.StatusCode {
	case 404:
		return ErrNotFound
	case 403:
		return ErrForbidden
	case 401:
		return ErrNotAuthorized
	case 400:
		return ErrBadRequest
	}

	// Unmarshall the JSON response
	if v != nil {
		//TODO: See if we can't verify the Seq and check for status: OK
		return json.Unmarshal(body, v)
	}

	return nil
}
