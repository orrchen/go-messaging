package lib

type Callbacks struct {
	OnNewConnection func (clientUid string)
	OnConnectionTerminated func (clientUid string)
	OnDataReceived func (clientUid string, data []byte)
}
