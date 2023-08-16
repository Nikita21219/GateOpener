package gate_controller

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"main/internal/amvideo"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type GateController struct {
	urlAMVideo string
}

func (gc *GateController) addHeaders(r *http.Request) {
	r.Header.Add("Host", "lk.amvideo-msk.ru")
	r.Header.Add("content-type", "application/x-www-form-urlencoded; charset=UTF-8")
	r.Header.Add("accept", "application/json, text/javascript, */*; q=0.01")
	r.Header.Add("sec-fetch-site", "cross-site")
	r.Header.Add("accept-language", "ru")
	r.Header.Add("sec-fetch-mode", "cors")
	r.Header.Add("origin", "null")
	r.Header.Add("user-agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 16_5_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148")
	r.Header.Add("sec-fetch-dest", "empty")
	r.Header.Add("pragma", "no-cache")
	r.Header.Add("cache-control", "no-cache")
}

func (gc *GateController) OpenGate(entry bool) error {
	var gateId string

	client := http.Client{}
	switch entry {
	case true:
		gateId = "1500"
	default:
		gateId = "1501"
	}

	bodyStr := fmt.Sprintf("type=open&id_shlag=%s&relay=0&sid=%s", gateId, os.Getenv("SID"))
	body := strings.NewReader(bodyStr)
	request, err := http.NewRequest("POST", gc.urlAMVideo, body)
	if err != nil {
		return err
	}

	gc.addHeaders(request)

	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	a := amvideo.AMVideoDto{}
	err = json.Unmarshal(respBody, &a)
	if err != nil {
		return err
	}

	if !a.Result {
		return fmt.Errorf(a.Message)
	}

	return nil
}

func (gc *GateController) OpenGateAlways(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("gate opening mode stopped")
				return
			default:
				wg := sync.WaitGroup{}
				log.Println("gate opening mode active...")

				for _, val := range []bool{true, false} {
					wg.Add(1)
					val := val
					go func() {
						defer wg.Done()
						err := gc.OpenGate(val)
						if err != nil {
							log.Println("error to open gate to entry:", err)
						}
					}()
				}

				wg.Wait()
				time.Sleep(1 * time.Second)
			}
		}
	}()
}

func NewController(urlAMVideoApi string) *GateController {
	return &GateController{
		urlAMVideo: urlAMVideoApi,
	}
}
