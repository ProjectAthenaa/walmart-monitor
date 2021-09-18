package monitor

import (
	"encoding/json"
	"fmt"
	"github.com/ProjectAthenaa/sonic-core/fasttls"
	"github.com/ProjectAthenaa/sonic-core/sonic/antibots/perimeterx"
	"github.com/google/uuid"
	"github.com/prometheus/common/log"
)

type PayloadOut struct {
	Payload string `json:"payload"`
	AppID   string `json:"appId"`
	Tag     string `json:"tag"`
	Uuid    string `json:"uuid"`
	Ft      string `json:"ft"`
	Seq     string `json:"seq"`
	En      string `json:"en"`
	Pc      string `json:"pc"`
	Sid     string `json:"sid,omitempty"`
	Vid     string `json:"vid,omitempty"`
	Cts     string `json:"cts,omitempty"`
	Rsc     string `json:"rsc"`
	Cs      string `json:"cs"`
	Ci      string `json:"ci"`
}

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
		`content-type`:              {`application/x-www-form-urlencoded; charset=UTF-8`},
	}
}

func (tk *Task) PXHoldCaptcha(blockedUrl string) {
	tk.pxUUID = uuid.NewString()

	payload, err := pxClient.ConstructPayload(tk.Ctx, &perimeterx.Payload{
		Site: perimeterx.SITE_WALMART,
		Type: perimeterx.PXType_PX2,
		RSC:  0,
		Uuid: tk.pxUUID,
	})
	if err != nil {
		log.Info("px error")
		return
	}
	var p2struct *PayloadOut
	json.Unmarshal(payload.Payload, &p2struct)

	req, err := tk.NewRequest("POST", "https://collector-pxu6b0qd2s.px-cloud.net/api/v2/bundle", []byte(fmt.Sprintf(`payload=%s&appId=%s&tag=%s&uuid=%s&ft=%s&seq=%s&en=%s&pc=%s&pxhd=%s&rsc=%s`, p2struct.Payload, "PXu6b0qd2S", p2struct.Tag, tk.pxUUID, p2struct.Ft, "0", p2struct.En, p2struct.Pc, string(tk.Client.Jar.PeekValue("pxhd")), "1")))
	if err != nil {
		log.Info("px error")
		return
	}
	req.Headers = tk.GenerateDefaultHeaders()

	res, err := tk.Do(req)
	if err != nil {
		log.Info("px error")
		return
	}

	payload, err = pxClient.ConstructPayload(tk.Ctx, &perimeterx.Payload{
		Site:           perimeterx.SITE_WALMART,
		Type:           perimeterx.PXType_PX34,
		ResponseObject: res.Body,
		RSC:            1,
		Uuid:           tk.pxUUID,
	})
	if err != nil {
		log.Info("px error")
		return
	}
	var p3struct *PayloadOut
	json.Unmarshal(payload.Payload, &p3struct)

	req, err = tk.NewRequest("POST", "https://collector-pxu6b0qd2s.px-cloud.net/api/v2/bundle", []byte(fmt.Sprintf(`payload=%s&appId=%s&tag=%s&uuid=%s&ft=%s&seq=%s&en=%s&cs=%s&pc=%s&sid=%s&pxhd=%s&cts=%s&rsc=%s`, p3struct.Payload, "PXu6b0qd2S", p3struct.Tag, p3struct.Uuid, p3struct.Ft, "1", p3struct.En, p3struct.Cs, p3struct.Pc, p3struct.Sid, string(tk.Client.Jar.PeekValue("pxhd")), p3struct.Cts, p3struct.Rsc)))
	if err != nil {
		log.Info("px error")
		return
	}
	req.Headers = tk.GenerateDefaultHeaders()

	res, err = tk.Do(req)
	if err != nil {
		log.Info("px error")
		return
	}

	payload, err = pxClient.ConstructPayload(tk.Ctx, &perimeterx.Payload{
		Site:           perimeterx.SITE_WALMART,
		Type:           perimeterx.PXType_HCAPHIGH,
		ResponseObject: res.Body,
		RSC:            2,
		Uuid:           tk.pxUUID,
	})
	if err != nil {
		log.Info("px error")
		return
	}
	var hcapstruct *PayloadOut
	json.Unmarshal(payload.Payload, &hcapstruct)

	req, err = tk.NewRequest("POST", "https://collector-pxu6b0qd2s.px-cloud.net/api/v2/bundle", []byte(fmt.Sprintf(`payload=%s&appId=%s&tag=%s&uuid=%s&ft=%s&seq=%s&en=%s&cs=%s&pc=%s&sid=%s󠄶󠄳󠄱󠄹󠄴󠄵󠄳󠄶󠄷󠄶󠄷󠄳&vid=%s&ci=%s&pxhd=%s&cts=%s&rsc=%s`, hcapstruct.Payload, "PXu6b0qd2S", hcapstruct.Tag, hcapstruct.Uuid, hcapstruct.Ft, "5", hcapstruct.En, hcapstruct.Cs, hcapstruct.Pc, hcapstruct.Sid, hcapstruct.Vid, hcapstruct.Ci, string(tk.Client.Jar.PeekValue("pxhd")), "4")))
	if err != nil {
		log.Info("px error")
		return
	}
	req.Headers = tk.GenerateDefaultHeaders()

	res, err = tk.Do(req)
	if err != nil {
		log.Info("px error")
		return
	}
	cookie, err := pxClient.GetCookie(tk.Ctx, &perimeterx.GetCookieRequest{PXResponse: res.Body})
	if err != nil {
		log.Info("px error")
		return
	}
	tk.Client.Jar.Set("_px3", cookie.Value)
}

func (tk *Task) Homepage(){
	req, err := tk.NewRequest("GET", `https://www.walmart.com/`, nil)
	if err != nil {
		return
	}
	req.Headers = tk.GenerateDefaultHeaders()
	//getpx

	res, err := tk.Do(req)
	if err != nil {
		return
	}

	if res.StatusCode == 307{
		tk.PXHoldCaptcha(res.Headers["Location"][0])
	}
}