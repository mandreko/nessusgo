package nessusgo

import (
	//"errors"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
)

type Client struct {
	Verify_SSL bool
	base_url   string
	login      string
	password   string
}

func NewClient(base_url string) *Client {
	return &Client{Verify_SSL: false, base_url: base_url}
}

func (c *Client) Authenticate(login string, password string) (token string, response *Record) {

	resource := "/login"

	c.login = login
	c.password = password

	data := make(map[string]string)
	data["login"] = login
	data["password"] = password

	headers := make(map[string]string)

	var r = c.Post(resource, data, headers)
	response, _ = decode(bytes.NewReader(r))
	token = response.Reply.Contents.Token
	return
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
	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: !c.Verify_SSL}}
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

// Internal methods
func decode(r io.Reader) (x *Record, err error) {
	x = new(Record)
	err = json.NewDecoder(r).Decode(x)
	return
}

// JSON Parsing structs
func (r Record) String() string {

	//return fmt.Sprintf("%b", r)
	out, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s", out)
}

type Record struct {
	Reply Reply `json:"reply"`
}

type Reply struct {
	Contents Contents `json:"contents"`
	Sequence int      `json:"seq"`
	Status   string   `json:"status"`
}

type Contents struct {
	IdleTimeout     int    `json:"idle_timeout"`
	LoadedPluginSet int    `json:"loaded_plugin_set"`
	MSP             bool   `json:"msp"`
	PluginSet       int    `json:"plugin_set"`
	ScannerBootTime int    `json:"scanner_boot_time"`
	ServerUUID      string `json:"server_uuid"`
	Token           string `json:"token"`
	User            User   `json:"user"`
}

type User struct {
	Admin bool   `json:"admin"`
	Name  string `json:"name"`
}
