package gateway

import (
	"os"
	"os/exec"
	"strings"
)

const fixedAlertMessage = "わーにんぐ！\n食べログ検索botからのエラーメッセージです。\n"

type AlertBot struct{}

func NewAlertBot() *AlertBot {
	return &AlertBot{}
}

func (a *AlertBot) AlertMessage(message string) error {
	return exec.Command("alert-linebot", "-userid", os.Getenv("ALERT_LINEBOT_USERID"), "-message", strings.Join([]string{fixedAlertMessage, message}, "")).Run()
}
