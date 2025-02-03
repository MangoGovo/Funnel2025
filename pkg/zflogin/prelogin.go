package zflogin

import (
	"encoding/json"
	"fmt"

	"funnel/internal/common/zf"
	"funnel/internal/exceptions"
	"funnel/pkg/config"
	_ "funnel/pkg/log"
	rdb "funnel/pkg/redis"
	"funnel/pkg/request"
	"funnel/pkg/schedule"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

const queueName = "PreLoginData"

func init() {
	RunProducer()
}

var maxQueueLength = config.Config.GetInt("preLogin.cacheCapacity")
var frequency = config.Config.GetInt("preLogin.produceFrequency")
var spec = fmt.Sprintf("@every %ds", frequency)

func RunProducer() {
	schedule.Register(spec, func() {
		err := producePreLoginData()
		if err != nil {
			zap.L().Error("生产预登录数据失败", zap.Error(err))
		}
	})
}

// 下面生产和消费的函数是核心, 需考虑代码的鲁棒性
func producePreLoginData() error {
	if err := cleanExpiredData(); err != nil {
		return err
	}
	if getQueueLength() >= int64(maxQueueLength) {
		zap.L().Debug("预登录队列数据已经满载")
		return nil
	}
	data, err := getPreLoginData()
	if err != nil {
		return err
	}

	bytes, err := json.Marshal(&data)
	if err != nil {
		return err
	}

	// TODO[blog] 写到博客里
	// 使用 SortedSet 解决 redis 中 List 无法给每个元素设置单独过期时间的痛点
	return rdb.Client.ZAdd(rdb.Ctx, queueName, &redis.Z{
		Score:  float64(data.ExpiredAt),
		Member: string(bytes),
	}).Err()
}

func consumePreLoginData() (*zf.PreLoginData, error) {
	// TODO[test] 测试缓存过期
	// 处理缓存过期
	if err := cleanExpiredData(); err != nil {
		return getPreLoginData()
	}
	// 从缓存里读取数据
	data := &zf.PreLoginData{}
	val, err := rdb.Client.ZPopMin(rdb.Ctx, queueName).Result()
	if err != nil {
		zap.L().Error("消费PreLoginData缓存失败")
		return getPreLoginData()
	}
	if len(val) == 0 {
		// 缓存为空
		return getPreLoginData()
	}
	err = json.Unmarshal([]byte(fmt.Sprint(val[0].Member)), data)
	if err != nil {
		zap.L().Error("反序列化PreLoginData缓存失败")
		return getPreLoginData()
	}
	zap.L().Debug("消费预登录缓存成功")
	return data, err
}

func cleanExpiredData() error {
	cnt, err := rdb.Client.ZRemRangeByScore(rdb.Ctx, queueName, "0", getTimestampStr()).Result()
	if err != nil {
		zap.L().Error("清理PreLoginData缓存失败")
		// 消费失败则重新获取, 下同
		return err
	}
	if cnt > 0 {
		zap.L().Info(fmt.Sprintf("PreLoginData缓存过期了%d个", cnt))
	}
	return nil
}

func getQueueLength() int64 {
	return rdb.Client.ZCard(rdb.Ctx, queueName).Val()
}

// getPreLoginData 获取在登录前可以被提前获取的数据
// 过完验证码的cookie
// 获取的对应的公钥

var expiredHours = config.Config.GetInt("preLogin.expireHours")

func getPreLoginData() (*zf.PreLoginData, error) {
	data := &zf.PreLoginData{}
	// 1. 初始化登录
	resp, err := request.New().R().Get(zf.LoginURL)
	if err != nil {
		return data, err
	}

	cookies := resp.Cookies()
	if len(cookies) == 0 {
		zap.L().Error("登录初始化失败")
		return data, exceptions.ZFLoginError
	}
	// 提取CSRFToken
	data.CSRFToken = ExtractCSRFToken(resp.String())

	// 2. 破解登录
	if !CrackCaptcha(cookies) {
		zap.L().Error("验证码破解失败")
		return data, exceptions.ZFLoginError
	}
	//zap.L().Debug("验证码破解成功", zap.Any("cookies", cookies))

	for _, cookie := range cookies {
		switch cookie.Name {
		case "JSESSIONID":
			data.JSESSIONID = cookie.Value
		case "route":
			data.Route = cookie.Value
		}
	}

	// 3. 获取公钥
	_, err = request.NewReqWithCookies(cookies).
		SetResult(data).
		SetQueryParams(map[string]string{
			"t": getTimestampStr(),
			"_": getTimestampStr(),
		}).
		Get(zf.PubKeyURL)
	data.ExpiredAt = getTimestamp() + int64(expiredHours*60*60*1000)
	return data, err
}
