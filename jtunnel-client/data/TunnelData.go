package data

type TunnelData struct {
	Id           int64
	TunnelId     string
	IsReplay     bool
	LocalPort    int16
	RequestData  []byte
	ResponseData []byte
}
