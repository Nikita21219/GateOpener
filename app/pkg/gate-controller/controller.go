package gate_controller

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"main/pkg/retry"
	"main/pkg/utils"
)

const (
	totalRetries  = 8
	urlAMVideoApi = "https://lk.amvideo-msk.ru/api/api4.php"

	EntryGateId = "1500"
	ExitGateId  = "1501"
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

func (gc *GateController) OpenGate(ctx context.Context, gateId string) error {
	if utils.Debug() {
		log.Println("debug mode active, gate not opened")
		return nil
	}

	client := http.Client{}

	bodyStr := fmt.Sprintf("type=open&id_shlag=%s&relay=0&sid=%s", gateId, os.Getenv("SID"))
	body := strings.NewReader(bodyStr)
	request, err := http.NewRequest("POST", gc.urlAMVideo, body)
	if err != nil {
		return err
	}

	gc.addHeaders(request)

	var resp *http.Response
	err = retry.Retry(ctx, totalRetries, func(ctx context.Context) error {
		resp, err = client.Do(request)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var a Dto
	err = json.Unmarshal(respBody, &a)
	if err != nil {
		return err
	}

	if !a.Result {
		switch {
		case a.Message == "":
			return fmt.Errorf("неизвестная ошибка")
		default:
			return fmt.Errorf(a.Message)
		}
	}

	return nil
}

func (gc *GateController) OpenGateForTimePeriod(ctx context.Context, ch chan error, duration time.Duration) {
	ticker := time.NewTicker(duration)

	go func() {
		for {
			select {
			case <-ticker.C:
				log.Println("gate opening mode stopped")
				return
			case <-ctx.Done():
				log.Println("gate opening mode stopped")
				return
			default:
				wg := sync.WaitGroup{}
				log.Println("gate opening mode active...")

				for _, gateId := range []string{EntryGateId, ExitGateId} {
					wg.Add(1)

					go func(gateId string) {
						defer wg.Done()

						err := gc.OpenGate(ctx, gateId)
						if err != nil {
							log.Println("error to open gate to entry:", err)
							ch <- err
						}
					}(gateId)
				}

				wg.Wait()
				time.Sleep(1 * time.Second)
			}
		}
	}()
}

func NewController() *GateController {
	return &GateController{
		urlAMVideo: urlAMVideoApi,
	}
}
