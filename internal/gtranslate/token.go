package gtranslate

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// TranslationToken encapsulates operations related to caching the translation token.
type TranslationToken struct {
	token string
}

// encodedTransform function applies a series of transformations to the input number
// based on the provided transformation string, and returns the result.
func (t *TranslationToken) encodedTransform(inputNum uint, transformStr string) uint {
	// Loop over the transformation string in steps of 3
	for i := 0; i < len(transformStr)-2; i += 3 {
		// Get the current character
		curChar := transformStr[i+2]

		// Determine the offset: if the character is a letter, subtract 87,
		// otherwise convert it to a number
		var offset int
		if 'a' <= curChar {
			offset = int(curChar - 87)
		} else {
			offset, _ = strconv.Atoi(string(curChar))
		}

		// Apply the transformation: if the next character is '+', shift right,
		// otherwise shift left
		var transformed uint
		if transformStr[i+1] == '+' {
			transformed = inputNum >> offset
		} else {
			transformed = inputNum << offset
		}

		// Apply the operation specified by the current character:
		// if it's '+', add the transformed value, otherwise XOR it
		if transformStr[i] == '+' {
			inputNum = inputNum + transformed
		} else {
			inputNum = inputNum ^ transformed
		}
	}

	return inputNum
}

// encodeString function transforms the input string based on the current config and
// returns the encoded string result.
func (t *TranslationToken) encodeString(inputStr string) string {
	var configValue string
	if t.token != "" {
		configValue = t.token
	} else {
		configValue = t.token
		if configValue == "" {
			configValue = "0"
		}
		t.token = configValue
	}
	tkkKey := "ttk"
	paramKey := "&" + tkkKey + "="
	configParts := strings.Split(configValue, ".")
	configValueNum, _ := strconv.Atoi(configParts[0])
	utf8Vals := []rune{}
	for idx := 0; idx < len(inputStr); idx++ {
		charCode := int(inputStr[idx])
		if charCode < 128 {
			utf8Vals = append(utf8Vals, rune(charCode))
		} else {
			if charCode < 2048 {
				utf8Vals = append(utf8Vals, rune(charCode>>6|192))
			} else {
				if (charCode&64512) == 55296 && idx+1 < len(inputStr) && (int(inputStr[idx+1])&64512) == 56320 {
					charCode = 65536 + ((charCode & 1023) << 10) + (int(inputStr[idx+1]) & 1023)
					utf8Vals = append(utf8Vals, rune(charCode>>18|240))
					utf8Vals = append(utf8Vals, rune(charCode>>12&63|128))
					idx++
				} else {
					utf8Vals = append(utf8Vals, rune(charCode>>12|224))
				}
				utf8Vals = append(utf8Vals, rune(charCode>>6&63|128))
			}
			utf8Vals = append(utf8Vals, rune(charCode&63|128))
		}
	}
	sum := configValueNum
	for i := 0; i < len(utf8Vals); i++ {
		sum += int(utf8Vals[i])
		sum = int(t.encodedTransform(uint(sum), "+-a^+6"))
	}
	sum = int(t.encodedTransform(uint(sum), "+-3^+b+-f"))
	sum ^= configValueNum
	if sum < 0 {
		sum = (sum & 2147483647) + 2147483648
	}
	sum %= 1e6
	sumStr := strconv.Itoa(sum)
	return paramKey + (sumStr + "." + strconv.Itoa(sum^configValueNum))
}

func (t *TranslationToken) update() error {
	tm := time.Now().UnixNano() / 3600000
	now := math.Floor(float64(tm))
	ttk, err := strconv.ParseFloat(t.token, 64)
	if err != nil {
		return err
	}

	if ttk == now {
		return nil
	}

	resp, err := http.Get(fmt.Sprintf("https://translate.%s", GoogleHost))
	if err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	matches := regexp.MustCompile(`tkk:\s?'(.+?)'`).FindStringSubmatch(string(body))
	if len(matches) > 0 {
		t.token = matches[0]
	}
	return nil
}

func (t *TranslationToken) Get(text string) (string, error) {
	t.update()
	tk := t.encodeString(text)
	tk = strings.Replace(tk, "&ttk=", "", -1)
	return tk, nil
}

func NewTranslationToken() TranslationToken {
	return TranslationToken{token: "0"}
}
