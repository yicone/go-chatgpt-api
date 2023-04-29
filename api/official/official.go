package official

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yicone/go-chatgpt-api/api"
)

const (
	apiUrl             = "https://api.openai.com"
	apiChatCompletions = apiUrl + "/v1/chat/completions"
	apiCheckUsage      = apiUrl + "/dashboard/billing/credit_grants"
)

var client *http.Client

func init() {
	client = &http.Client{
		Timeout: 0,
	}
}

type ChatCompletionsRequest struct {
	Model    string                   `json:"model"`
	Messages []ChatCompletionsMessage `json:"messages"`
	Stream   bool                     `json:"stream"`
}

type ChatCompletionsMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

//goland:noinspection GoUnhandledErrorResult
func ChatCompletions(c *gin.Context) {
	var chatCompletionsRequest ChatCompletionsRequest
	c.ShouldBindJSON(&chatCompletionsRequest)
	data, _ := json.Marshal(chatCompletionsRequest)
	req, _ := http.NewRequest("POST", apiChatCompletions, bytes.NewBuffer(data))
	req.Header.Set("Authorization", api.GetAccessToken(c.GetHeader(api.AuthorizationHeader)))
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Content-Type", "application/json")
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		} else {
			c.Writer.Write([]byte(line))
			c.Writer.Flush()
		}
	}
}

//goland:noinspection GoUnhandledErrorResult
func CheckUsage(c *gin.Context) {
	req, _ := http.NewRequest("GET", apiCheckUsage, nil)
	req.Header.Set("Authorization", api.GetAccessToken(c.GetHeader(api.AuthorizationHeader)))
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	c.Writer.Write(body)
}
