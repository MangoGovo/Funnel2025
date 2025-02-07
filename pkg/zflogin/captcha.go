//nolint:all
package zflogin

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"image"
	"net/http"

	"funnel/pkg/request"
	"github.com/avast/retry-go"
	"go.uber.org/zap"
	"gocv.io/x/gocv"
)

// CrackCaptcha 识别正方验证码
func CrackCaptcha(cookies []*http.Cookie) bool {
	// 1. 获取验证码token 无论重试与否 只需获得一次即可
	rtk := getRTK(cookies)
	if rtk == "" {
		return false
	}

	err := retry.Do(func() error {
		return retryCaptcha(cookies, rtk)
	}, retry.Attempts(5))

	return err == nil
}

// retryCaptcha 错误重试
func retryCaptcha(cookies []*http.Cookie, rtk string) error {
	// 2. 获取验证码参数
	captchaData, err := getCaptchaData(cookies, rtk)
	if err != nil {
		return err
	}

	// 3. 下载验证码图片
	imtk := captchaData.Imtk
	bgURL := captchaData.Si
	targetURL := captchaData.Mi
	bg, target, err := fetchCaptcha(cookies, bgURL, targetURL, imtk)
	if err != nil {
		return err
	}

	// 4. 识别并生成参数
	movement, err := generateMovement(bg, target)
	if err != nil {
		return err
	}

	// 5. 打请求过验证码
	extend := getExtend()
	return fetchVerify(cookies, rtk, movement, extend)
}

// getRTK 获取验证码token
func getRTK(cookies []*http.Cookie) string {
	params := map[string]string{
		"type":       "resource",
		"instanceId": "zfcaptchaLogin",
		"name":       "zfdun_captcha.js",
	}

	resp, err := request.NewReqWithCookies(cookies).
		SetQueryParams(params).
		Get(captchaURL)
	if err != nil {
		return "err"
	}

	rtk := extractRTK(resp.String())
	return rtk
}

// getCaptchaData 获取验证码参数
func getCaptchaData(cookies []*http.Cookie, rtk string) (data *captchaData, err error) {
	params := map[string]string{
		"type":       "refresh",
		"rtk":        rtk,
		"time":       getTimestampStr(),
		"instanceId": "zfcaptchaLogin",
	}

	_, err = request.NewReqWithCookies(cookies).
		SetQueryParams(params).
		SetResult(&data).
		Get(captchaURL)
	if err != nil {
		return nil, err
	}
	return data, err
}

// getExtend 获取extend参数
func getExtend() string {
	userAgent := request.New().Header.Get("User-Agent")
	data := map[string]string{
		"appName":    "Netscape",
		"appVersion": userAgent,
		"userAgent":  userAgent,
	}
	bytes, _ := json.Marshal(data)
	//	base64
	return base64.StdEncoding.EncodeToString(bytes)
}

// fetchImage 下载图片
func fetchImage(cookies []*http.Cookie, url, imtk string) ([]byte, error) {
	params := map[string]string{
		"type":       "image",
		"id":         url,
		"imtk":       imtk,
		"t":          getTimestampStr(),
		"instanceId": "zfcaptchaLogin",
	}

	resp, err := request.NewReqWithCookies(cookies).
		SetQueryParams(params).
		Get(captchaURL)

	if err != nil {
		return nil, err
	}
	return resp.Body(), nil
}

// fetchCaptcha 下载验证码图片
func fetchCaptcha(cookies []*http.Cookie, bgURL string, targetURL string, imtk string) (bg []byte, target []byte, err error) {
	bg, err = fetchImage(cookies, bgURL, imtk)
	if err != nil {
		return nil, nil, err
	}

	target, err = fetchImage(cookies, targetURL, imtk)
	if err != nil {
		return nil, nil, err
	}
	return bg, target, nil
}

// generateMovement 生成滑块轨迹参数
func generateMovement(bg []byte, target []byte) (string, error) {
	movement := make([]move, 0)
	start := getTimestamp()

	// 滑块识别
	result, err := slideMatch(bg, target)
	if err != nil {
		return "", err
	}
	xEnd := result.Target[0]

	// 生成轨迹
	xMove := 0
	for xMove < xEnd {
		step := randomInt(1, 10)
		// 确保不会超出目标
		if xMove+step > xEnd {
			step = xEnd - xMove
		}
		point := move{
			X: 50 + xMove,
			Y: 50,
			T: start,
		}
		movement = append(movement, point)
		xMove += step
		start += int64(randomInt(100, 250)) // 时间间隔
	}

	// 将轨迹序列化为 JSON
	bytes, err := json.Marshal(movement)
	if err != nil {
		return "", err
	}

	// 转换为 base64 编码
	return base64.StdEncoding.EncodeToString(bytes), nil
}

// fetchVerify 打请求过验证
func fetchVerify(cookies []*http.Cookie, rtk, movement, extend string) error {
	var result captchaResult
	resp, err := request.NewReqWithCookies(cookies).
		SetFormData(map[string]string{
			"type":       "verify",
			"rtk":        rtk,
			"time":       getTimestampStr(),
			"mt":         movement,
			"instanceId": "zfcaptchaLogin",
			"extend":     extend,
		}).
		SetResult(&result).
		Post(captchaURL)
	zap.L().Debug("fetchVerify", zap.Any("resp", resp))
	if err != nil {
		return err
	}
	if result.Status != "success" {
		return errors.New(result.Msg)
	}
	return nil
}

// slideMatch 滑块匹配
func slideMatch(backgroundBytes []byte, targetBytes []byte) (*slideResult, error) {
	// 解码目标图像
	targetMat, err := gocv.IMDecode(targetBytes, gocv.IMReadAnyColor)
	if err != nil {
		return nil, err
	}
	defer targetMat.Close()

	// 解码背景图像
	backgroundMat, err := gocv.IMDecode(backgroundBytes, gocv.IMReadAnyColor)
	if err != nil {
		return nil, err
	}
	defer backgroundMat.Close()

	// 对图像应用 Canny 边缘检测
	targetCanny := gocv.NewMat()
	defer targetCanny.Close()
	gocv.Canny(targetMat, &targetCanny, 100, 200)

	backgroundCanny := gocv.NewMat()
	defer backgroundCanny.Close()
	gocv.Canny(backgroundMat, &backgroundCanny, 100, 200)

	// 将图像转换为 RGB
	targetRGB := gocv.NewMat()
	defer targetRGB.Close()
	gocv.CvtColor(targetCanny, &targetRGB, gocv.ColorGrayToBGR)

	backgroundRGB := gocv.NewMat()
	defer backgroundRGB.Close()
	gocv.CvtColor(backgroundCanny, &backgroundRGB, gocv.ColorGrayToBGR)

	// 模板匹配
	mask := gocv.NewMat()
	defer mask.Close()
	result := gocv.NewMat()
	defer result.Close()
	gocv.MatchTemplate(backgroundRGB, targetRGB, &result, gocv.TmCcoeffNormed, mask)

	// 获取匹配结果
	_, maxVal, _, maxLoc := gocv.MinMaxLoc(result)
	if maxVal == 0 {
		return nil, err
	}

	// 计算匹配框的右下角坐标
	h, w := targetRGB.Rows(), targetRGB.Cols()
	bottomRight := image.Point{X: maxLoc.X + w, Y: maxLoc.Y + h}

	return &slideResult{
		TargetX: 0,
		TargetY: 0,
		Target:  []int{maxLoc.X, maxLoc.Y, bottomRight.X, bottomRight.Y},
	}, nil
}
