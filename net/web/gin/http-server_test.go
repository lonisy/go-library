package gin

import (
	"net/http"
	"os"
	"testing"
)

func TestHttpServer_Start(t *testing.T) {
	WebServer := NewHttpServer("test")
	WebServer.AnyApiRouter()
	if os.Getenv("APP_ENV") == "develop" {
		WebServer.Start(&http.Server{Addr: "127.0.0.1:8086", Handler: WebServer})
	} else {
		WebServer.Start(&http.Server{Addr: ":8086", Handler: WebServer})
	}
}
