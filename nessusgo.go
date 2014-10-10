package nessusgo

import (
	//"errors"
	"fmt"
	"log"
	//"os"
	"bytes"
	"crypto/tls"
	//"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
)

type Client struct {
	verify_ssl bool
	base_url   string
	login      string
	password   string
}

func NewClient(base_url string, login string, password string) *Client {
	return &Client{verify_ssl: false, base_url: base_url, login: login, password: password}
}

func (c *Client) Authenticate(login string, password string) bool {

	resource := "/login"

	data := make(map[string]string)
	data["login"] = login
	data["password"] = password

	headers := make(map[string]string)

	resp := c.Post(resource, data, headers)

	fmt.Printf("Response Body: %v", string(resp))

	return false
}

func (c *Client) LogOut() {}

func (c *Client) Get(url string, params []string, headers []string) {}

func (c *Client) Post(resource string, post_data map[string]string, headers map[string]string) []byte {

	u, _ := url.ParseRequestURI(c.base_url)
	u.Path = resource
	urlStr := fmt.Sprintf("%v", u)

	seq := strconv.Itoa(rand.Intn(9999))

	data := url.Values{}
	data.Set("seq", seq)
	data.Add("json", "1")
	for k, v := range post_data {
		data.Add(k, v)
	}

	req, _ := http.NewRequest("POST", urlStr, bytes.NewBufferString(data.Encode()))
	req.Header.Set("User-Agent", "Nessus.Go v0.1")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	for k, v := range headers {
		req.Header.Add(k, v)
	}

	// Use SSL verification options
	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: !c.verify_ssl}}
	client := &http.Client{Transport: tr}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Could not make request: %v", err)
		//panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	// TODO: Verify check sequence

	// TODO: Verify "OK" status message

	return body
}
