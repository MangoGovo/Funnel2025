package zflogin

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"funnel/pkg/request"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

// Login 正方登陆
func Login(username, password string) ([]*http.Cookie, error) {
	cachedLoginData, err := consumePreLoginData()
	if err != nil {
		return []*http.Cookie{}, err
	}

	pubKey := &pubKeyData{
		Modulus:  cachedLoginData.Modulus,
		Exponent: cachedLoginData.Exponent,
	}
	encryptedPwd, err := encryptPwd(pubKey, password)
	if err != nil {
		zap.L().Error("密码加密失败", zap.Error(err))
		return []*http.Cookie{}, err
	}

	// 发送登陆请求
	loginData := map[string]string{
		"csrftoken": cachedLoginData.CSRFToken,
		"language":  "zh_CN",
		"yhm":       username,
		"mm":        encryptedPwd,
	}
	cookies := []*http.Cookie{
		{Name: "JSESSIONID", Value: cachedLoginData.JSESSIONID},
		{Name: "route", Value: cachedLoginData.Route},
	}
	startTime := time.Now()
	resp, err := request.NewWithoutRedirect().R().
		SetCookies(cookies).
		SetFormData(loginData).
		SetQueryParam("time", getTimestampStr()).
		Post(loginURL)
	if err != nil && !errors.Is(err, resty.ErrAutoRedirectDisabled) {
		return nil, err
	}

	zap.L().Debug("登陆请求耗时" + fmt.Sprint(time.Since(startTime)))
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "JSESSIONID" {
			cookies[0].Value = cookie.Value
			return cookies, nil
		}
	}
	// TODO 解析错误提示
	// TODO 弹性补充机制
	/*
		<p id="tips" class="bg_danger sl_danger">
			<span class="glyphicon glyphicon-minus-sign"></span>用户名或密码不正确，请重新输入！
		</p>
	*/
	zap.L().Info(username + "登陆失败")
	return nil, nil
}
