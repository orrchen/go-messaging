package models

type DeviceResponse struct {
	Action string `json:"action"`
	Uid string `json:"uid"`
	Data map[string]interface{} `json:"data"`
	Status string `json:"status"`
	ErrorCode string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}
