package monitor

import (
	"github.com/ProjectAthenaa/sonic-core/fasttls"
	"github.com/ProjectAthenaa/sonic-core/sonic/antibots/perimeterx"
	"github.com/prometheus/common/log"
	"time"
)

func (tk *Task) GenerateDefaultHeaders() fasttls.Headers {
	return fasttls.Headers{
		"sec-ch-ua":                 {`"Google Chrome";v="93", " Not;A Brand";v="99", "Chromium";v="93"`},
		"sec-ch-ua-mobile":          {"?0"},
		"upgrade-insecure-requests": {"1"},
		"user-agent":                {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.63 Safari/537.36"},
		"accept":                    {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"},
		"sec-fetch-site":            {"none"},
		"sec-fetch-mode":            {"navigate"},
		"sec-fetch-user":            {"?1"},
		"sec-fetch-dest":            {"document"},
		"accept-encoding":           {"gzip, deflate, br"},
		"accept-language":           {"en-US,en;q=0.9"},
	}
}

func (tk *Task) PXLoop(){
	for range time.Tick(500 * time.Millisecond){
		payload, err := pxClient.ConstructPayload(tk.Ctx, &perimeterx.Payload{
			Site:           perimeterx.SITE_WALMART,
			Type:           perimeterx.PXType_PX34,
			Cookie:         "",
			ResponseObject: nil,
			Token:          "",
			RSC:            0,
		})

		if err != nil{
			log.Error(err)
			continue
		}

		req, err := tk.NewRequest("POST", "https://collector-pxu6b0qd2s.px-cloud.net/api/v2/collector", payload.Payload)
		if err != nil {
			log.Error(err)
			continue
		}
		req.Headers = tk.GenerateDefaultHeaders()

		res, err := tk.Do(req)
		if err != nil {
			log.Error(err)
			continue
		}

		cookie, err := pxClient.GetCookie(tk.Ctx, &perimeterx.GetCookieRequest{PXResponse: res.Body})
		if err != nil {
			log.Error(err)
			continue
		}
		tk.Client.Jar.Set("_px3", cookie.Value)
	}
}