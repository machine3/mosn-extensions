package integrate

import (
	_ "github.com/machine3/mosn-extensions/pkg/filter/stream/iprestriction"
	"mosn.io/mosn/pkg/configmanager"
	_ "mosn.io/mosn/pkg/filter/network/proxy"
	"mosn.io/mosn/pkg/mosn"
	_ "mosn.io/mosn/pkg/stream/http"
	"mosn.io/mosn/test/util"
	"net/http"
	"testing"
	"time"
)

func Test(t *testing.T) {
	// start a http server
	server := &http.Server{Addr: ":8080", Handler: &util.HTTPHandler{}}
	go server.ListenAndServe()
	defer server.Close()

	mosnConfig := configmanager.Load("config.json")
	configmanager.Reset()
	configmanager.SetMosnConfig(mosnConfig)
	mosn := mosn.NewMosn(mosnConfig)
	go mosn.Start()
	defer mosn.Close()

	time.Sleep(2 * time.Second) // wait mosn start

	client := http.Client{Timeout: 5 * time.Second}
	req, _ := http.NewRequest("GET", "http://localhost:2046/", nil)
	req.Header.Set("X-Real-Ip", "1.1.1.2")
	resp, err := client.Do(req)
	if err != nil {
		t.Error("test ip restriction error: ", err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected http status code was success 200, but got %d", resp.StatusCode)
	}
	resp.Body.Close() // release
}
