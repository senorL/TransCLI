package translate

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/senorL/TransCLI/conf"
	"golang.org/x/exp/rand"
)

type TranslationResult struct {
	From        string `json:"from"`
	To          string `json:"to"`
	TransResult []struct {
		Src string `json:"src"`
		Dst string `json:"dst"`
	} `json:"trans_result"`
}

func TranForBaidu(query string) {
	rand.Seed(uint64(time.Now().UnixNano()))
	// 生成盐值
	salt := strconv.Itoa(rand.Intn(32768) + 32768)
	// 生成签名
	sign := generateSign(query, salt)

	// 拼接 URL 编码后的数据
	params := fmt.Sprintf("q=%s&from=auto&to=zh&appid=%s&salt=%s&sign=%s",
		url.QueryEscape(query),
		conf.GetConf().TModle.APPID,
		salt,
		sign)

	// 请求的 URL
	url := "https://fanyi-api.baidu.com/api/trans/vip/translate"
	method := "POST"

	// 使用字符串构建请求体
	payload := strings.NewReader(params)

	// 创建 HTTP 客户端
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// 发送请求
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 解析 JSON 响应
	var result TranslationResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 提取并打印翻译结果
	if len(result.TransResult) > 0 {
		translatedText := result.TransResult[0].Dst
		fmt.Println(translatedText)
	} else {
		fmt.Println("未能翻译")
	}
}

func generateSign(query, salt string) string {
	signStr := conf.GetConf().TModle.APPID + query + salt + conf.GetConf().TModle.APIKey
	hash := md5.Sum([]byte(signStr))
	return hex.EncodeToString(hash[:])
}
