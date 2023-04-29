package webdriver

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	_ "github.com/yicone/go-chatgpt-api/env"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"github.com/yicone/go-chatgpt-api/api"
	"github.com/yicone/go-chatgpt-api/util/logger"
)

var WebDriver selenium.WebDriver

//goland:noinspection GoUnhandledErrorResult,SpellCheckingInspection
func init() {
	var wg sync.WaitGroup
	wg.Add(1)

	// Create a channel for receiving signals
	sigs := make(chan os.Signal, 1)
	// Use the signal package to notify the operating system which signals we want to receive
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		logger.Warn("Received signal: " + sig.String())
		wg.Wait()
		if WebDriver != nil {
			WebDriver.Quit()
		}
		os.Exit(0)
	}()

	chatgptProxyServer := os.Getenv("CHATGPT_PROXY_SERVER")
	if chatgptProxyServer == "" {
		logger.Error("CHATGPT_PROXY_SERVER is empty")
		return
	}
	logger.Info("CHATGPT_PROXY_SERVER is: " + chatgptProxyServer)

	chromeArgs := []string{
		"--no-sandbox",
		"--disable-gpu",
		"--disable-dev-shm-usage",
		"--disable-blink-features=AutomationControlled",
		"--incognito",
	}

	mode := os.Getenv("CHATGPT_API_MODE")
	if mode == "development" {
		// Processing Debug Version Logic
		logger.Info("CHATGPT_API_MODE is: " + mode)
	} else {
		chromeArgs = append(chromeArgs, "--headless=new")
	}

	networkProxyServer := os.Getenv("NETWORK_PROXY_SERVER")
	if networkProxyServer != "" {
		logger.Info("NETWORK_PROXY_SERVER is: " + networkProxyServer)
		chromeArgs = append(chromeArgs, "--proxy-server="+networkProxyServer)
	}

	WebDriver, _ = selenium.NewRemote(selenium.Capabilities{
		"chromeOptions": chrome.Capabilities{
			Args:            chromeArgs,
			ExcludeSwitches: []string{"enable-automation"},
		},
	}, chatgptProxyServer)
	wg.Done() // Ensure that WebDriver is started successfully

	if WebDriver == nil {
		logger.Error("Please make sure chatgpt proxy service is running")
		os.Exit(1)
		return
	}

	WebDriver.Get(api.ChatGPTUrl)

	if isReady(WebDriver) {
		logger.Info(api.ChatGPTWelcomeText)
		OpenNewTabAndChangeBackToOldTab()
	} else {
		if !isAccessDenied(WebDriver) {
			if HandleCaptcha(WebDriver) {
				logger.Info(api.ChatGPTWelcomeText)
				OpenNewTabAndChangeBackToOldTab()
			}
		}
	}
}

//goland:noinspection GoUnhandledErrorResult
func OpenNewTabAndChangeBackToOldTab() {
	WebDriver.ExecuteScript(fmt.Sprintf("open('%s');", api.ChatGPTUrl), nil)
	handles, _ := WebDriver.WindowHandles()
	WebDriver.SwitchWindow(handles[0])

	InitXhrMap()
	InitConversationMap()
}

func InitXhrMap() {
	_, err := WebDriver.ExecuteScript("window.xhrMap = new Map();", nil)
	if err != nil {
		logger.Error("Failed to init xhrMap, please restart go-chatgpt-api.")
	}
}

func InitConversationMap() {
	// to save conversations, (k,v): {"request message id": "response message data"}
	_, err := WebDriver.ExecuteScript("window.conversationMap = new Map();", nil)
	if err != nil {
		logger.Error("Failed to init conversationMap, please restart go-chatgpt-api.")
	}
}
