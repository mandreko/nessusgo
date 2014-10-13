package nessusgo

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
)

var (
	ErrNotAuthenticated = errors.New("You must authenticate before using this function.")
)

// New creates an instance of the Nessus Client
func New(nessusUrl string) *Client {
	c := &Client{}
	c.ServiceUrl = nessusUrl

	c.Server = &ServerResource{c}
	c.Users = &UserResource{c}
	return c
}

type Client struct {
	ServiceUrl string
	Token      string

	Server *ServerResource
	Users  *UserResource
}

func (c *Client) isAuthenticated() bool {
	return c.Token != ""
}

func (c *Client) Login(username, password string) error {
	record := Record{}
	path := "/login"
	values := url.Values{}

	values.Add("login", username)
	values.Add("password", password)

	if err := c.do("POST", path, nil, values, &record); err != nil {
		return err
	}

	c.Token = record.Reply.Contents.Token

	return nil
}

func (c *Client) Logout() error {
	if !c.isAuthenticated() {
		return ErrNotAuthenticated
	}

	path := "/logout"

	if err := c.do("POST", path, nil, nil, nil); err != nil {
		return err
	}

	c.Token = ""

	return nil
}

func (c *Client) Feed() (*Record, error) {
	if !c.isAuthenticated() {
		return nil, ErrNotAuthenticated
	}

	record := Record{}
	path := "/feed"

	if err := c.do("POST", path, nil, nil, &record); err != nil {
		return nil, err
	}

	return &record, nil
}

func (c *Client) Uuid() (*string, error) {
	uuid := ""
	if !c.isAuthenticated() {
		return &uuid, ErrNotAuthenticated
	}

	record := Record{}
	path := "/uuid"

	if err := c.do("POST", path, nil, nil, &record); err != nil {
		return &uuid, err
	}

	uuid = record.Reply.Contents.UUID
	return &uuid, nil
}

// func (c *Client) GetCert() (string, error) {
// 	record := ""
// 	path := "/getcert"

// 	if err := c.do("POST", path, nil, nil, &record); err != nil {
// 		return "", err
// 	}
// 	fmt.Printf("DEBUG body2: %v\n", record)
// 	return record, nil
// }

//TODO: Move all this to other files

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
	Sequence int      `json:"seq,string"`
	Status   string   `json:"status"`
}

type Contents struct {
	IdleTimeout      int                     `json:"idle_timeout,string"`
	LoadedPluginSet  int                     `json:"loaded_plugin_set,string"`
	MSP              string                  `json:"msp"` // Temporarily a string, due to uppercase "TRUE" Json parsing issues
	PluginSet        int                     `json:"plugin_set,string"`
	ScannerBootTime  int                     `json:"scanner_boottime,string"`
	ServerUUID       string                  `json:"server_uuid"`
	Token            string                  `json:"token"`
	User             user                    `json:"user"`
	Scans            Scans                   `json:"scans"`
	Expiration       int                     `json:"expiration,string"`
	ExpirationTime   int                     `json:"expiration_time,string"`
	Feed             string                  `json:"feed"`
	NessusType       string                  `json:"nessus_type"`
	NessusUIVersion  string                  `json:"nessus_ui_version"`
	ServerVersion    string                  `json:"server_version"`
	WebServerVersion string                  `json:"web_server_version"`
	Notifications    map[string]Notification `json:"notifications"`
	UUID             string                  `json:"uuid"`
	Users            users                   `json:"users"`
}

type Notification struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

type Scans struct {
	ScanList []string `json:"scanlist"`
}

type users struct {
	Users []*user `json:"user"`
}

type user struct {
	Name      string `json:"name"`
	IsAdmin   string `json:"admin"`
	LastLogin int    `json:"lastlogin"`
}
