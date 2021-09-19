package monitor

import (
	"encoding/json"
	"fmt"
	"github.com/ProjectAthenaa/sonic-core/sonic/antibots/perimeterx"
	"github.com/ProjectAthenaa/sonic-core/sonic/core"
	"github.com/google/uuid"
	"github.com/prometheus/common/log"
	"regexp"
	"time"
)

var (
	itemOfferRe = regexp.MustCompile(`offers":\["(\w+)"`)
)

func (tk *Task) iteration() error{
	count := core.Base.GetRedis("cache").PubSubNumSub(tk.Ctx, tk.Data.RedisChannel).Val()

	if v, ok := count[tk.Data.RedisChannel]; v == 0 || !ok || tk.pubcount != 1{
		time.Sleep(time.Second)
		return nil
	}

	tk.GetPX()

	req, err := tk.NewRequest("GET", fmt.Sprintf(`https://www.walmart.com/terra-firma/item/%s`, tk.sku), nil)
	if err != nil {
		return err
	}
	req.Headers = tk.GenerateDefaultHeaders()
	//getpx

	res, err := tk.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode == 307{
		tk.PXHoldCaptcha(res.Headers["Location"][0])
	}

	log.Info("product found")
	offer := itemOfferRe.FindSubmatch(res.Body)
	if len(offer) > 0{
		tk.Monitor.Channel <- map[string]interface{}{
			"offerid":string(offer[1]),
		}
		tk.pubcount = 0
		return nil
	}
	return nil
}

func (tk *Task) GetPX(){
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

	req, err := tk.NewRequest("POST", "https://collector-pxu6b0qd2s.px-cloud.net/api/v2/collector", []byte(fmt.Sprintf(`payload=%s&appId=%s&tag=%s&uuid=%s&ft=%s&seq=%s&en=%s&pc=%s&pxhd=%s&rsc=%s`, p2struct.Payload, "PXu6b0qd2S", p2struct.Tag, tk.pxUUID, p2struct.Ft, "0", p2struct.En, p2struct.Pc, string(tk.Client.Jar.PeekValue("pxhd")), "1")))
	if err != nil {
		log.Info("px error")
		return
	}
	req.Headers = tk.GenerateDefaultHeaders()

	log.Info("making px2 req")
	res, err := tk.Do(req)
	if err != nil {
		log.Info("px error")
		return
	}
	log.Info("made px2 req")

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

	req, err = tk.NewRequest("POST", "https://collector-pxu6b0qd2s.px-cloud.net/api/v2/collector", []byte(fmt.Sprintf(`payload=%s&appId=%s&tag=%s&uuid=%s&ft=%s&seq=%s&en=%s&cs=%s&pc=%s&sid=%s&pxhd=%s&cts=%s&rsc=%s`, p3struct.Payload, "PXu6b0qd2S", p3struct.Tag, p3struct.Uuid, p3struct.Ft, "1", p3struct.En, p3struct.Cs, p3struct.Pc, p3struct.Sid, string(tk.Client.Jar.PeekValue("pxhd")), p3struct.Cts, p3struct.Rsc)))
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

	//cookie, err := pxClient.GetCookie(tk.Ctx, &perimeterx.GetCookieRequest{PXResponse: res.Body})
	//if err != nil {
	//	log.Info("px error")
	//	return
	//}
	//tk.Client.Jar.Set("_px3", cookie.Value)

	cookie, err := pxClient.GetPXde(tk.Ctx, &perimeterx.GetCookieRequest{PXResponse: res.Body})
	if err != nil {
		log.Info("px error")
		return
	}

	log.Info("init pxde",  cookie.Value)
	tk.Client.Jar.Set("_pxde", cookie.Value)
}