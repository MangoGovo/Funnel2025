package zflogin

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/hex"
	"math/big"
	random "math/rand"
	"regexp"
	"strconv"
	"time"

	_ "funnel/pkg/log"
)

func reExtract(pattern string, str string) string {
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(str)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func randomInt(min, max int) int {
	return random.Intn(max-min) + min
}

func ExtractCSRFToken(html string) string {
	return reExtract(
		`<input[^>]+id=\"csrftoken\"[^>]+value=\"([^\"]+)\"`,
		html)
}

func extractRTK(js string) string {
	return reExtract(`rtk:'([a-f0-9-]+)'`, js)
}

func getTimestamp() int64 {
	return time.Now().UnixMilli()
}

func getTimestampStr() string {
	return strconv.FormatInt(getTimestamp(), 10)
}

// RSA
func encryptPwd(key *pubKeyData, password string) (string, error) {
	nString, _ := base64.StdEncoding.DecodeString(key.Modulus)
	n, _ := new(big.Int).SetString(hex.EncodeToString(nString), 16)
	eString, _ := base64.StdEncoding.DecodeString(key.Exponent)
	e, _ := strconv.ParseInt(hex.EncodeToString(eString), 16, 32)
	pub := rsa.PublicKey{E: int(e), N: n}
	cc, err := rsa.EncryptPKCS1v15(rand.Reader, &pub, []byte(password))
	return base64.StdEncoding.EncodeToString(cc), err
}
