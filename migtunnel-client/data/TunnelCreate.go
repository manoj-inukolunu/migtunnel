package data

type TunnelCreateRequest struct {
	HostName        string
	TunnelName      string
	LocalServerPort int16
	Tls             bool
	TlsServerFQDN   string
}
