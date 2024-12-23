package captcha

import (
	config2 "github.com/TestsLing/aj-captcha-go/config"
	constant "github.com/TestsLing/aj-captcha-go/const"
	"github.com/TestsLing/aj-captcha-go/service"
	conf "weihu_server/library/config"
)

var factory *service.CaptchaServiceFactory

func Init() {
	// 滑动模块配置
	var blockPuzzleConfig = &config2.BlockPuzzleConfig{Offset: 10}
	var config = config2.BuildConfig(constant.RedisCacheKey, constant.DefaultResourceRoot, nil,
		nil, blockPuzzleConfig, 2*60)
	//var config = config2.NewConfig()
	factory = service.NewCaptchaServiceFactory(config)
	//注册自定义配置redis数据库
	factory.RegisterCache(constant.RedisCacheKey, service.NewConfigRedisCacheService([]string{conf.GetString("redis.host")},
		"", conf.GetString("redis.password"), false, conf.GetInt("redis.db")))

	// 这里默认是注册了 内存缓存，但是不足以应对生产环境，希望自行注册缓存驱动 实现缓存接口即可替换（CacheType就是注册进去的 key）
	//factory.RegisterCache(constant.MemCacheKey, service.NewMemCacheService(20)) // 这里20指的是缓存阈值

	// 注册了验证码服务
	factory.RegisterService(constant.BlockPuzzleCaptcha, service.NewBlockPuzzleCaptchaService(factory))

}

// CreateBlockPuzzleCaptcha 创建验证码
func CreateBlockPuzzleCaptcha() (map[string]interface{}, error) {
	return factory.GetService(constant.BlockPuzzleCaptcha).Get()
}

// CheckBlockPuzzleCaptcha 核对验证码
func CheckBlockPuzzleCaptcha(token, pointJson string) error {
	return factory.GetService(constant.BlockPuzzleCaptcha).Check(token, pointJson)
}

// VerifyBlockPuzzleCaptcha 二次校验验证码
func VerifyBlockPuzzleCaptcha(token, pointJson string) error {
	return factory.GetService(constant.BlockPuzzleCaptcha).Verification(token, pointJson)
}
