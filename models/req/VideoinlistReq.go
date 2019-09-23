package req

import (
	"encoding/json"
)

type VideoinlistReq struct {
	Gid  int   `json:"gid"`
	Vids []int `json:"vids"`
}

func (v *VideoinlistReq) Unmarshal(jsonBody string) bool {
	err := json.Unmarshal([]byte(jsonBody), v)
	if err != nil || v.Gid == 0 {
		return false
	}
	return true
}
