package frontend

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/kabukky/httpscerts"
)

func TestLaunchServer(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	cert, key := "test-cert.pem", "test-key.pem"
	secure, unsecure := "127.0.0.1:8000", "127.0.0.1:8001"
	if err := httpscerts.Generate(cert, key, secure); err != nil {
		t.Error("failure generating certificate / key")
	}
	defer os.Remove(cert)
	defer os.Remove(key)

	testFS := new(FrontendServer)
	testFS.LaunchServer(secure, unsecure, cert, key)
}

func TestStart(t *testing.T) {}

func TestStop(t *testing.T) {}
