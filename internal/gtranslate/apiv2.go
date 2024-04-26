package gtranslate

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/text/language"
)

var ttk TranslationToken

func init() {
	ttk = NewTranslationToken()
}

const (
	defaultNumberOfRetries = 2
)

func translate(text, from, to string, withVerification bool, tries int, delay time.Duration) ([]string, error) {
	if tries == 0 {
		tries = defaultNumberOfRetries
	}

	if withVerification {
		if _, err := language.Parse(from); err != nil && from != "auto" {
			log.Println("[WARNING], '" + from + "' is a invalid language, switching to 'auto'")
			from = "auto"
		}
		if _, err := language.Parse(to); err != nil {
			log.Println("[WARNING], '" + to + "' is a invalid language, switching to 'en'")
			to = "en"
		}
	}

	urll := fmt.Sprintf("https://translate.%s/translate_a/single", GoogleHost)

	//token, err := ttk.Get(text)
	//if err != nil {
	//	return nil, err
	//}

	data := map[string]string{
		"client": "gtx",
		"sl":     from,
		"tl":     to,
		"hl":     to,
		"ie":     "UTF-8",
		"oe":     "UTF-8",
		"otf":    "1",
		"ssel":   "0",
		"tsel":   "0",
		"kc":     "7",
		"q":      text,
	}

	u, err := url.Parse(urll)
	if err != nil {
		return nil, nil
	}

	parameters := url.Values{}

	for k, v := range data {
		parameters.Add(k, v)
	}
	for _, v := range []string{"at", "bd", "ex", "ld", "md", "qca", "rw", "rm", "ss", "t"} {
		parameters.Add("dt", v)
	}

	//parameters.Add("tk", token)
	u.RawQuery = parameters.Encode()

	var r *http.Response

	for tries > 0 {
		r, err = http.Get(u.String())
		log.Println(u.String())
		if err != nil {
			if err == http.ErrHandlerTimeout {
				return nil, errBadNetwork
			}
			return nil, err
		}

		if r.StatusCode == http.StatusOK {
			break
		}

		if r.StatusCode == http.StatusForbidden {
			tries--
			time.Sleep(delay)
		}
	}

	raw, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var resp []interface{}

	err = json.Unmarshal([]byte(raw), &resp)
	if err != nil {
		return nil, err
	}

	responseTextArr := make([]string, 1)
	if resp[0] != nil {

	}
	a, ok := resp[0].([]interface{})
	if ok {
		for _, obj := range a {
			if len(obj.([]interface{})) == 0 {
				break
			}

			t, ok := obj.([]interface{})[0].(string)
			if ok {
				responseTextArr[0] += t
			}
		}
	}
	obj1, ok := resp[5].([]interface{})
	if ok {
		for _, obj2 := range obj1 {
			if len(obj2.([]interface{})) == 0 {
				break
			}
			obj3, ok := obj2.([]interface{})[2].([]interface{})
			if ok {
				for i, obj4 := range obj3 {
					obj5, ok := obj4.([]interface{})[0].(string)
					if ok {
						if len(responseTextArr) < i+2 {
							responseTextArr = append(responseTextArr, "")
						}
						responseTextArr[i+1] += obj5
					}
				}
			}
		}
	}

	return responseTextArr, nil
}
