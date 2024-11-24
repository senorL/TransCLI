package transfer

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/senorL/TransCLI/conf"
	"golang.org/x/exp/rand"
)

func TranForBaidu(query []string) {
	var builder strings.Builder
	for _, q := range query {
		builder.WriteString(q)
	}
	result, _ := translate(builder.String())
	fmt.Println(result)
}

func generateSign(query, salt string) string {
	signStr := conf.GetConf().TModle.APPID + query + salt + conf.GetConf().TModle.APIKey
	hash := md5.Sum([]byte(signStr))
	return hex.EncodeToString(hash[:])
}
func translate(query string) (string, error) {
	rand.Seed(uint64(time.Now().UnixNano()))
	salt := strconv.Itoa(rand.Intn(32768) + 32768)
	sign := generateSign(query, salt)

	params := map[string]string{
		"q":     query,
		"from":  "en",
		"to":    "zh",
		"appid": conf.GetConf().TModle.APPID,
		"salt":  salt,
		"sign":  sign,
	}

	url := "https://fanyi-api.baidu.com/api/trans/vip/translate"
	resp, err := http.Get(url + "?" + encodeParams(params))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if transResult, ok := result["trans_result"].([]interface{}); ok {
		if firstResult, ok := transResult[0].(map[string]interface{}); ok {
			if dst, ok := firstResult["dst"].(string); ok {
				return dst, nil
			}
		}
	}

	if errorMsg, ok := result["error_msg"].(string); ok {
		return "", fmt.Errorf("error: %s", errorMsg)
	}

	return "", fmt.Errorf("unknown error")
}

func encodeParams(params map[string]string) string {
	var encoded string
	for key, value := range params {
		encoded += key + "=" + value + "&"
	}
	return encoded[:len(encoded)-1]
}
