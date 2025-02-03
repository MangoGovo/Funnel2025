package zf

// PreLoginData 用于维护登录前一些可以被缓存的数据
type PreLoginData struct {
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

const (
	BaseURL   = "http://www.gdjw.zjut.edu.cn"
	LoginURL  = BaseURL + "/jwglxt/xtgl/login_slogin.html"
	PubKeyURL = BaseURL + "/jwglxt/xtgl/login_getPublicKey.html"
)
