package nessusgo

import (
	"encoding/json"
	"fmt"
	"io"
)

type Record struct {
	Reply Reply `json:"reply"`
}

func decode(r io.Reader) (x *Record, err error) {
	x = new(Record)
	err = json.NewDecoder(r).Decode(x)
	return
}

func (r Record) String() string {

	//return fmt.Sprintf("%b", r)
	out, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s", out)
}
