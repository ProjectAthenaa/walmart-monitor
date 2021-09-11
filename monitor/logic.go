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
	for {
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

		offer := itemOfferRe.FindSubmatch(res.Body)
		if len(offer) > 0{
			tk.Monitor.Channel <- map[string]interface{}{
				"offerid":string(offer[1]),
			}
			return nil
		}
	}
}
