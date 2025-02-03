package zflogin

import (
	"fmt"
	"net/http"
	"time"

	"funnel/pkg/request"
	"go.uber.org/zap"
)

// Login 正方登陆
func Login(username, password string) ([]*http.Cookie, error) {
	preLoginData, err := consumePreLoginData()
	if err != nil {
		return []*http.Cookie{}, err
	}

	pubKey := &pubKeyData{
		Modulus:  preLoginData.Modulus,
		Exponent: preLoginData.Exponent,
	}
	encryptedPwd, err := encryptPwd(pubKey, password)
	if err != nil {
		zap.L().Error("密码加密失败", zap.Error(err))
		return []*http.Cookie{}, err
	}

	// 发送登陆请求
	loginData := map[string]string{
		"csrftoken": preLoginData.CSRFToken,
		"language":  "zh_CN",
		"yhm":       username,
		"mm":        encryptedPwd,
	}
	startTime := time.Now()
	resp, err := request.NewReqWithCookies([]*http.Cookie{
		{Name: "JSESSIONID", Value: preLoginData.JSESSIONID},
		{Name: "route", Value: preLoginData.Route},
	}).
		SetFormData(loginData).
		Post(loginURL)
	if err != nil {
		return nil, err
	}
	fmt.Println(time.Since(startTime))
	zap.L().Info(fmt.Sprint(resp.RawResponse.Request.URL))
	/* TODO 解析错误提示
	<p id="tips" class="bg_danger sl_danger">
		<span class="glyphicon glyphicon-minus-sign"></span>用户名或密码不正确，请重新输入！
	</p>
	*/

	return nil, nil
}
