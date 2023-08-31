package download

import "net/http"

// Re-usable client with custom transport
var client = http.Client{
	Transport: NewAddHeaderTransport(nil),
}

// Custom header things below
// Source: https://stackoverflow.com/questions/51628755/how-to-add-default-header-fields-from-http-client
// 		   https://go.dev/play/p/FbkpFlyFCm_F

type AddHeaderTransport struct {
	T http.RoundTripper
}

func (adt *AddHeaderTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", "go-sdio-seeder/0.0.1")
	return adt.T.RoundTrip(req)
}

func NewAddHeaderTransport(T http.RoundTripper) *AddHeaderTransport {
	if T == nil {
		T = http.DefaultTransport
	}
	return &AddHeaderTransport{T}
}
