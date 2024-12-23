package language

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"sort"
	"strconv"
	"strings"
)

var bundle *i18n.Bundle

func Initial() {
	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	for _, lang := range []string{"en", "zh"} {
		bundle.MustLoadMessageFile(fmt.Sprintf("./config/locales/%s.json", lang))
	}
}

// Translate 执行翻译
func Translate(key string, data interface{}, locale string) string {
	localize := i18n.NewLocalizer(bundle, locale)
	// 执行翻译
	translation := localize.MustLocalize(&i18n.LocalizeConfig{
		MessageID:    key,
		TemplateData: data,
	})
	return translation
}

// Middleware 中间件提取 Accept-Language 头部的语言信息并存储在上下文中
func Middleware(c *fiber.Ctx) error {
	acceptLang := c.Get("Accept-Language")
	locale, _, _ := ParseAcceptLanguage(acceptLang)
	c.Locals("locale", locale)
	return c.Next()
}

// Info 结构体存储语言和其优先级
type Info struct {
	Language string
	Quality  float64
}

// ParseAcceptLanguage 解析 Accept-Language 头部
func ParseAcceptLanguage(header string) (bestLanguage string, bestQuality float64, err error) {
	languages := make([]Info, 0)

	// 分割成多个语言标签
	tokens := strings.Split(header, ",")
	for _, token := range tokens {
		parts := strings.SplitN(token, ";", 2)
		lang := strings.TrimSpace(parts[0])

		quality := 1.0
		if len(parts) > 1 && strings.HasPrefix(strings.TrimSpace(parts[1]), "q=") {
			q, err := strconv.ParseFloat(strings.TrimPrefix(strings.TrimSpace(parts[1]), "q="), 64)
			if err != nil {
				return "", 0, err
			}
			quality = q
		}

		languages = append(languages, Info{Language: lang, Quality: quality})
	}

	// 对语言列表按质量降序排序
	sort.Slice(languages, func(i, j int) bool {
		return languages[i].Quality > languages[j].Quality
	})

	if len(languages) > 0 {
		return languages[0].Language, languages[0].Quality, nil
	}

	return "", 0, nil
}
