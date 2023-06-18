package data

type TunnelData struct {
	Id           int64
	TunnelId     string
	Tls          bool
	IsReplay     bool
	LocalPort    int16
	RequestData  []byte
	ResponseData []byte
}
