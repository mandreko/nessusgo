package nessusgo

type Contents struct {
	IdleTimeout     int    `json:"idle_timeout,string"`
	LoadedPluginSet int    `json:"loaded_plugin_set,string"`
	MSP             string `json:"msp"` // Temporarily a string, due to uppercase "TRUE" Json parsing issues
	PluginSet       int    `json:"plugin_set,string"`
	ScannerBootTime int    `json:"scanner_boottime,string"`
	ServerUUID      string `json:"server_uuid"`
	Token           string `json:"token"`
	User            User   `json:"user"`
	Scans           Scans  `json:"scans"`
}
