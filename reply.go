package nessusgo

type Reply struct {
	Contents Contents `json:"contents"`
	Sequence int      `json:"seq,string"`
	Status   string   `json:"status"`
}
