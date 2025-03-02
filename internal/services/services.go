package service

type SvcLayer struct {
	HTTPServer HTTPServer
}

type HTTPServer interface {
	GenServer(spec []byte) ([]byte, error)
}
