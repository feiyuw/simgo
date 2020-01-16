package protocols

type HTTPClient struct {
	addr string // connected service addr, eg. http://127.0.0.1:2345
}

func NewHTTPClient(addr string) (*HTTPClient, error) {
	return &HTTPClient{addr: addr}, nil
}
