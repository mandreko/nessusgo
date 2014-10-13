package nessusgo

import ()

type ServerResource struct {
	client *Client
}

// Requests the current Nessus server load and platform type.
func (r *ServerResource) GetLoad() (*ServerLoadRecord, error) {
	if !r.client.isAuthenticated() {
		return nil, ErrNotAuthenticated
	}

	record := ServerLoadRecord{}
	url_path := "/server/load"

	if err := r.client.do("POST", url_path, nil, nil, &record); err != nil {
		return nil, err
	}

	return &record, nil
}

// Requests the Nessus server settings including proxy information, User-Agent, and custom update host.
func (r *ServerResource) GetSettings() (*ServerSettingsRecord, error) {
	if !r.client.isAuthenticated() {
		return nil, ErrNotAuthenticated
	}

	record := ServerSettingsRecord{}
	url_path := "/server/securesettings/list"

	if err := r.client.do("POST", url_path, nil, nil, &record); err != nil {
		return nil, err
	}

	return &record, nil
}

// Updated Nessus server settings including proxy information, User-Agent, and custom update host
// func (r *ServerResource) GetSettings() (*ServerSettingsRecord, error) {
// 	if !r.client.isAuthenticated() {
// 		return nil, ErrNotAuthenticated
// 	}

// 	record := ServerSettingsRecord{}
// 	url_path := "/server/securesettings"

// 	if err := r.client.do("POST", url_path, nil, nil, &record); err != nil {
// 		return nil, err
// 	}

// 	return &record, nil
// }

// Server.GetLoad Structs

type ServerLoadRecord struct {
	Reply ServerLoadReply `json:"reply"`
}

type ServerLoadReply struct {
	Contents ServerLoadContents `json:"contents"`
	Sequence int                `json:"seq,string"`
	Status   string             `json:"status"`
}

type ServerLoadContents struct {
	Platform string         `json:"platform"`
	Load     ServerLoadLoad `json:"load"`
}

type ServerLoadLoad struct {
	LoadAverage         float64 `json:"loadavg,string"`
	NumberOfHosts       int     `json:"num_hosts,string"`
	NumberOfScans       int     `json:"num_scans,string"`
	NumberOfSessions    int     `json:"num_sessions,string"`
	NumberOfTcpSessions int     `json:"num_tcp_sessions,string"`
}

// Server.GetSettings Structs

type ServerSettingsRecord struct {
	Reply ServerSettingsReply `json:"reply"`
}

type ServerSettingsReply struct {
	Contents ServerSettingsContents `json:"contents"`
	Sequence int                    `json:"seq,string"`
	Status   string                 `json:"status"`
}

type ServerSettingsContents struct {
	SecureSettings ServerSettingsSecureSettings `json:"securesettings"`
}

type ServerSettingsSecureSettings struct {
	ProxySettings ServerSettingsProxySettings `json:"proxysettings"`
}

type ServerSettingsProxySettings struct {
	Proxy         string `json:"proxy"`
	ProxyPort     string `json:"proxy_port"`
	ProxyUsername string `json:"proxy_username"`
	ProxyPassword string `json:"proxy_password"`
	UserAgent     string `json:"user_agent"`
	CustomHost    string `json:"custom_host"`
}
