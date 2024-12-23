// Package jwt 处理 JWT 认证
package jwt

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strings"
	"time"
	"weihu_server/library/cache"
	"weihu_server/library/common"
	"weihu_server/library/config"

	jwtPkg "github.com/golang-jwt/jwt/v4"
)

var (
	ErrTokenExpired           = errors.New("令牌已过期")
	ErrTokenExpiredMaxRefresh = errors.New("令牌已过最大刷新时间")
	ErrTokenMalformed         = errors.New("请求令牌格式有误")
	ErrTokenInvalid           = errors.New("请求令牌无效")
	ErrHeaderEmpty            = errors.New("需要认证才能访问！")
	ErrHeaderMalformed        = errors.New("请求头中 Authorization 格式有误")
)

// JWT 定义一个jwt对象
type JWT struct {

	// 秘钥，用以加密 JWT，读取配置信息 app.key
	SignKey []byte

	// 刷新 Token 的最大过期时间
	MaxRefresh time.Duration
}

// CustomClaims 自定义载荷
type CustomClaims struct {
	UserID       string `json:"user_id"`
	TeamID       string `json:"team_id"`
	SysType      string `json:"sys_type"`
	IsSuper      bool   `json:"is_super"`
	ExpireAtTime int64  `json:"expire_time"`
	jwtPkg.RegisteredClaims
}

func NewJWT() *JWT {
	return &JWT{
		SignKey:    []byte(config.GetString("jwt.secret")),
		MaxRefresh: config.GetDuration("jwt.maxRefreshTime") * time.Minute,
	}
}

// createToken 创建 Token，内部使用，外部请调用 IssueToken
func (jwt *JWT) createToken(claims CustomClaims) (string, error) {
	// 使用HS256算法进行token生成
	token := jwtPkg.NewWithClaims(jwtPkg.SigningMethodHS256, claims)
	return token.SignedString(jwt.SignKey)
}

// expireAtTime 过期时间
func (jwt *JWT) expireAtTime() time.Time {

	timeNow := time.Now()

	expireTime := config.GetDuration("jwt.expireTime")

	expire := expireTime * time.Minute
	return timeNow.Add(expire)
}

// ParserToken 解析 Token，中间件中调用
func (jwt *JWT) ParserToken(c *fiber.Ctx) (*CustomClaims, error) {
	tokenString, parseErr := jwt.getTokenFromHeader(c)
	if parseErr != nil {
		return nil, parseErr
	}

	// 1. 调用 jwt 库解析用户传参的 Token
	token, err := jwt.parseTokenString(tokenString)

	// 2. 解析出错
	if err != nil {
		validationErr, ok := err.(*jwtPkg.ValidationError)
		if ok {
			if validationErr.Errors == jwtPkg.ValidationErrorMalformed {
				return nil, ErrTokenMalformed
			} else if validationErr.Errors == jwtPkg.ValidationErrorExpired {
				return nil, ErrTokenExpired
			}
		}
		return nil, ErrTokenInvalid
	}

	// 3. 将 token 中的 claims 信息解析出来和 CustomClaims 数据结构进行校验
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		if !claims.IsSuper {
			if !jwt.tokenInWhitelist(claims.UserID, claims.SysType, tokenString) {
				return nil, ErrTokenInvalid
			}

			// 修改Token有效期
			jwt.whitelistAdd(claims.UserID, claims.SysType, tokenString, config.GetInt64("jwt.appExpireTime"))

		}

		return claims, nil
	}

	return nil, ErrTokenInvalid
}

// RefreshToken 更新 Token，用以提供 refresh token 接口
func (jwt *JWT) RefreshToken(c *fiber.Ctx) (string, error) {
	// 1. 从 Header 里获取 token
	tokenString, parseErr := jwt.getTokenFromHeader(c)
	if parseErr != nil {
		return "", parseErr
	}

	// 2. 调用 jwt 库解析用户传参的 Token
	token, err := jwt.parseTokenString(tokenString)

	// 3. 解析出错，未报错证明是合法的 Token（甚至未到过期时间）
	if err != nil {
		validationErr, ok := err.(*jwtPkg.ValidationError)
		// 满足 refresh 的条件：只是单一的报错 ValidationErrorExpired
		if !ok || validationErr.Errors != jwtPkg.ValidationErrorExpired {
			return "", err
		}
	}

	// 4. 解析 CustomClaims 的数据
	claims := token.Claims.(*CustomClaims)

	if !jwt.tokenInWhitelist(claims.UserID, claims.SysType, tokenString) {
		return "", ErrTokenInvalid
	}

	// 5. 检查是否过了『最大允许刷新的时间』
	x := time.Now().Add(-jwt.MaxRefresh).Unix()
	if claims.IssuedAt.Unix() > x {
		// 删除旧Token
		jwt.whitelistDel(claims.UserID, claims.SysType)

		return jwt.createToken(*claims)
	}

	return "", ErrTokenExpiredMaxRefresh
}

// getTokenFromHeader 使用 jwtPkg.ParseWithClaims 解析 Token
// Authorization:Bearer xxxxx
func (jwt *JWT) getTokenFromHeader(c *fiber.Ctx) (string, error) {

	authHeader := string(c.Request().Header.Peek("Authorization"))

	if authHeader == "" {
		return "", ErrHeaderEmpty
	}
	// 按空格分割
	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return "", ErrHeaderMalformed
	}
	return parts[1], nil
}

// parseTokenString 使用 jwtPkg.ParseWithClaims 解析 Token
func (jwt *JWT) parseTokenString(tokenString string) (*jwtPkg.Token, error) {
	return jwtPkg.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwtPkg.Token) (interface{}, error) {
		return jwt.SignKey, nil
	})
}

// InvalidateToken 无效Token
func (jwt *JWT) InvalidateToken(c *fiber.Ctx) error {
	// 1. 从 Header 里获取 token
	tokenString, parseErr := jwt.getTokenFromHeader(c)
	if parseErr != nil {
		return parseErr
	}

	// 2. 调用 jwt 库解析用户传参的 Token
	token, err := jwt.parseTokenString(tokenString)

	// 3. 解析出错，未报错证明是合法的 Token（甚至未到过期时间）
	if err != nil {
		validationErr, ok := err.(*jwtPkg.ValidationError)
		// 满足 refresh 的条件：只是单一的报错 ValidationErrorExpired
		if !ok || validationErr.Errors != jwtPkg.ValidationErrorExpired {
			return err
		}
	}

	// 4. 解析 CustomClaims 的数据
	claims := token.Claims.(*CustomClaims)

	// 删除Token
	jwt.whitelistDel(claims.UserID, claims.SysType)

	return nil
}

// 添加白名单
func (jwt *JWT) whitelistAdd(userId string, systemType string, token string, expiresIn int64) error {
	key := fmt.Sprintf("%s_%s_%s", common.CacheTokenPrefix, systemType, userId)
	cache.SaveString(key, token, time.Duration(expiresIn)*time.Minute)
	return nil
}

// 删除白名单
func (jwt *JWT) whitelistDel(userId string, systemType string) {
	key := fmt.Sprintf("%s_%s_%s", common.CacheTokenPrefix, systemType, userId)
	cache.Remove(key)
}

// 查询白名单
func (jwt *JWT) tokenInWhitelist(userId string, systemType string, token string) bool {
	key := fmt.Sprintf("%s_%s_%s", common.CacheTokenPrefix, systemType, userId)

	return cache.GetString(key) == token
}

// IssueClientToken 生成ClientToken，在登录成功时调用
func (jwt *JWT) IssueClientToken(userID, teamID string, isSuper bool, expiresTimeByMinute int64) string {
	claims := CustomClaims{
		userID,
		teamID,
		common.SysClient,
		isSuper,
		0,
		jwtPkg.RegisteredClaims{
			IssuedAt:  jwtPkg.NewNumericDate(time.Now()),
			NotBefore: jwtPkg.NewNumericDate(time.Now()),
			Issuer:    config.GetString("server.name"),
			// Subject:   "somebody",
			ID: userID,
			// Audience:  []string{"somebody_else"},
		},
	}

	// 2. 根据 claims 生成token对象
	token, err := jwt.createToken(claims)
	if err != nil {
		return ""
	}

	if !isSuper {
		// 加入白名单
		if err := jwt.whitelistAdd(userID, common.SysClient, token, expiresTimeByMinute); err != nil {
			return ""
		}
	}

	return token
}

// IssueBackendToken 生成BackendToken，在登录成功时调用
func (jwt *JWT) IssueBackendToken(userID string, isSuper bool) string {
	claims := CustomClaims{
		userID,
		"",
		common.SysBackend,
		isSuper,
		0,
		jwtPkg.RegisteredClaims{
			IssuedAt:  jwtPkg.NewNumericDate(time.Now()),
			NotBefore: jwtPkg.NewNumericDate(time.Now()),
			Issuer:    config.GetString("server.name"),
			// Subject:   "somebody",
			ID: userID,
			// Audience:  []string{"somebody_else"},
		},
	}

	// 2. 根据 claims 生成token对象
	token, err := jwt.createToken(claims)
	if err != nil {
		return ""
	}

	if !isSuper {
		// 加入白名单
		if err := jwt.whitelistAdd(userID, common.SysBackend, token, config.GetInt64("jwt.appExpireTime")); err != nil {
			return ""
		}
	}

	return token
}

// IssueWebToken 生成WebToken，在登录成功时调用
func (jwt *JWT) IssueWebToken(userID, teamID string, isSuper bool) string {
	claims := CustomClaims{
		userID,
		teamID,
		common.SysWeb,
		isSuper,
		0,
		jwtPkg.RegisteredClaims{
			IssuedAt:  jwtPkg.NewNumericDate(time.Now()),
			NotBefore: jwtPkg.NewNumericDate(time.Now()),
			Issuer:    config.GetString("server.name"),
			// Subject:   "somebody",
			ID: userID,
			// Audience:  []string{"somebody_else"},
		},
	}

	// 2. 根据 claims 生成token对象
	token, err := jwt.createToken(claims)
	if err != nil {
		return ""
	}

	if !isSuper {
		// 加入白名单
		if err := jwt.whitelistAdd(userID, common.SysWeb, token, config.GetInt64("jwt.appExpireTime")); err != nil {
			return ""
		}
	}

	return token
}

func (jwt *JWT) IssueLiveKitToken(userID string, expiresTimeByMinute int64) string {
	claims := CustomClaims{
		userID,
		"",
		common.SysLiveKit,
		false,
		0,
		jwtPkg.RegisteredClaims{
			IssuedAt:  jwtPkg.NewNumericDate(time.Now()),
			NotBefore: jwtPkg.NewNumericDate(time.Now()),
			Issuer:    config.GetString("server.name"),
			ID:        userID,
		},
	}

	// 2. 根据 claims 生成token对象
	token, err := jwt.createToken(claims)
	if err != nil {
		return ""
	}

	// 加入白名单
	if err = jwt.whitelistAdd(userID, common.SysLiveKit, token, expiresTimeByMinute); err != nil {
		return ""
	}

	return token
}
