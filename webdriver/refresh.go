package webdriver

import (
	"strings"

	"github.com/yicone/go-chatgpt-api/api"
	"github.com/yicone/go-chatgpt-api/util/logger"
)

//goland:noinspection GoUnhandledErrorResult
func Refresh() {
	if err := WebDriver.Refresh(); err != nil {
		errorMessage := err.Error()
		if strings.HasSuffix(errorMessage, "connect: connection refused") {
			logger.Error("Please make sure chatgpt-proxy-server is running, if running, restart it")
		} else if strings.HasSuffix(errorMessage, "invalid session id") {
			logger.Warn("Service chatgpt-proxy-server is detected, go-chatgpt-api is trying to resume")
			newRefresh()
		}
	} else {
		HandleCaptcha(WebDriver)
	}
}

//goland:noinspection GoUnhandledErrorResult
func NewSessionAndRefresh() {
	if _, err := WebDriver.PageSource(); err != nil {
		if err.Error() == "invalid session id: invalid session id" {
			newRefresh()
		}
	}
}

//goland:noinspection GoUnhandledErrorResult
func newRefresh() {
	WebDriver.NewSession()
	WebDriver.Get(api.ChatGPTUrl)
	HandleCaptcha(WebDriver)
}
