package zflogin

const (
	baseURL    = "http://www.gdjw.zjut.edu.cn"
	captchaURL = baseURL + "/jwglxt/zfcaptchaLogin"
	loginURL   = baseURL + "/jwglxt/xtgl/login_slogin.html"
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
