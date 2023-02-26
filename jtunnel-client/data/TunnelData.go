package data

type TunnelData struct {
	Id           int64
	TunnelId     string
	IsReplay     bool
	RequestData  []byte
	ResponseData []byte
}
