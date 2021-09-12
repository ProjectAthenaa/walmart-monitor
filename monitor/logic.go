package monitor

import (
	"fmt"
	"github.com/json-iterator/go"
	"regexp"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
	itemOfferRe = regexp.MustCompile(`offers":\["(\w+)"]`)
)

func (tk *Task) iteration() error{
	//get px. set px in jar

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

	offer := itemOfferRe.FindSubmatch(res.Body)
	if len(offer) > 0{
		tk.Monitor.Channel <- map[string]interface{}{
			"offerid":string(offer[1]),
		}
		return nil
	}
	return nil
}
