package zfservice

import (
	"funnel/pkg/log"
	"funnel/pkg/request"
	"go.uber.org/zap"
	"net/http"
)

func CheckCookie(cookies []*http.Cookie) bool {
	resp, err := request.NewWithoutTLS().R().SetCookies(cookies).Get(checkURL)
	if err != nil {
		return false
	}
	if resp.StatusCode() != http.StatusOK {
		log.L().Info("cookie 过期或不合法", zap.Any("cookie", cookies))
		return false
	}
	return true
}
