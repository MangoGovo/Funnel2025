package zfservice

import "funnel/pkg/config"

var (
	baseURL    = config.Config.GetString("zf.baseURL")
	captchaURL = baseURL + "/jwglxt/zfcaptchaLogin"
	loginURL   = baseURL + "/jwglxt/xtgl/login_slogin.html"
	pubKeyURL  = baseURL + "/jwglxt/xtgl/login_getPublicKey.html"
	checkURL   = baseURL + "/jwglxt/xtgl/index_cxYhxxIndex.html"
)

// captchaData 验证码信息
type captchaData struct {
	Msg    string `json:"msg"`
	T      int64  `json:"t"`
	Si     string `json:"si"`
	Imtk   string `json:"imtk"`
	Mi     string `json:"mi"`
	Vs     string `json:"vs"`
	Status string `json:"status"`
}

// move 滑块在某个时刻的位置
type move struct {
	X int   `json:"x"`
	Y int   `json:"y"`
	T int64 `json:"t"`
}

// captchaResult 滑块验证结果
type captchaResult struct {
	Msg    string `json:"msg"`
	Vs     string `json:"vs"`
	Status string `json:"status"`
}

// slideResult 滑块匹配结果
type slideResult struct {
	TargetX int
	TargetY int
	Target  []int
}

// pubKeyData RSA公钥
type pubKeyData struct {
	Modulus  string `json:"modulus"`
	Exponent string `json:"exponent"`
}

// preLoginData 用于维护登录前一些可以被缓存的数据
type preLoginData struct {
	// 记录数据过期的时间
	ExpiredAt int64 `json:"-"`

	// 公钥
	Modulus  string `json:"modulus"`
	Exponent string `json:"exponent"`

	// 过完验证码的Cookie
	JSESSIONID string `json:"JSESSIONID"`
	Route      string `json:"route"`

	// 其他登录需要的参数
	CSRFToken string `json:"csrftoken"`
}
