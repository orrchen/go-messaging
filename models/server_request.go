package models

type ServerRequest struct{
	DeviceRequest DeviceRequest `json:"deviceRequest"`
	ServerId string `json:"serverId"` // unique identifier for each server in case of having more than 1 server
	ClientId string `json:"clientId"`
}
