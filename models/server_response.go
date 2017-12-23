package models

type ServerResponse struct {
	DeviceResponse DeviceResponse `json:"deviceResponse"`
	ServerId string `json:"serverId"` // unique identifier for each server in case of having more than 1 server
	ClientId string `json:"clientId"`
	DeviceUid string `json:"deviceUid"`
}