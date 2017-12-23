package models

type DeviceRequest struct {
	Action string `json:"action"`
	DeviceUid string `json:"deviceUid"`
	Uid string `json:"uid"`
	Data map[string]interface{} `json:"data"`
}
