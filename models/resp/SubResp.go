package resp

type SubResp struct {
	Me     int    `json:"me_uid"`
	Up     int    `json:"up_uid"`
	Action string `json:"action"`
}
