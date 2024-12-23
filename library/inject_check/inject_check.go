package inject_check

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"html"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"weihu_server/library/util"
)

// CheckInject 防注入检测
func CheckInject(c *fiber.Ctx) error {
	errMsg, ret := CheckSqlInject(c.Body())
	if ret != 0 {
		return c.JSON(errMsg)
	}
	return c.Next()
}

func CheckSqlInject(params ...interface{}) (errMsg string, ret int) {
	for _, param := range params {
		// 使用反射获取参数的类型
		//paramType := reflect.TypeOf(param)
		// 使用反射获取参数的名称
		paramName := reflect.TypeOf(param).Name()
		// 使用反射获取参数值
		paramValue := reflect.ValueOf(param)

		switch paramValue.Kind() {
		case reflect.Invalid:
			fmt.Printf("%s = invalid\n", paramName)
		case reflect.String:
			// 校验字符串,过滤特殊字符

		case reflect.Slice:
			// 校验slice元素
			for i := 0; i < paramValue.Len(); i++ {

			}
		case reflect.Struct:
			// 递归校验struct字段
			for i := 0; i < paramValue.NumField(); i++ {
				field := paramValue.Field(i)
				if field.Kind() == reflect.String {
					field.SetString(html.EscapeString(field.String()))
				}
			}
		case reflect.Array:
			// 递归校验Array字段
			for i := 0; i < paramValue.Len(); i++ {

			}
		case reflect.Map:
			// 递归校验Map字段
			//for _, key := range paramValue.MapKeys() {
			//
			//}
		case reflect.Interface:
			// 递归校验Interface字段
			if paramValue.IsNil() {
				fmt.Printf("%s = nil\n", paramName)
			} else {

			}
		default:
			// 其他类型,可按需校验
		}
	}

	return
}

// cleanXSS 清理XSS
func cleanXSS(param string) string {
	temp := param

	// 替换小于号(英文换成中文)
	param = strings.ReplaceAll(param, "<", "＜")
	// 替换大于号(英文换成中文)
	param = strings.ReplaceAll(param, ">", "＞")
	// 替换左括号(英文换成中文)
	param = strings.ReplaceAll(param, "\\(", "（")
	// 替换右括号(英文换成中文)
	param = strings.ReplaceAll(param, "\\)", "）")
	// 替换单引号(英文换成中文)
	param = strings.ReplaceAll(param, "'", "＇")
	// 替换分号(英文换成中文)
	param = strings.ReplaceAll(param, ";", "；")

	/**-------------------------javascript替换------------------------*/
	// 替换小于号(转义字符)
	param = strings.ReplaceAll(param, "<", "&lt;")
	// 替换大于号(转义字符)
	param = strings.ReplaceAll(param, ">", "&gt;")
	// 替换左括号(转义字符)
	param = strings.ReplaceAll(param, "\\(", "&#40;")
	// 替换右括号(转义字符)
	param = strings.ReplaceAll(param, "\\)", "&#41")
	param = strings.ReplaceAll(param, "eval\\((.*)\\)", "")
	param = strings.ReplaceAll(param, "[\\\"'][\\s]*javascript:(.*)[\\\"']", "\"\"")
	param = strings.ReplaceAll(param, "script", "")
	param = strings.ReplaceAll(param, "link", "")
	param = strings.ReplaceAll(param, "frame", "")
	param = strings.ReplaceAll(param, ";", "")
	param = strings.ReplaceAll(param, "0x0d", "")
	param = strings.ReplaceAll(param, "0x0a", "")
	/**-----------------------javascript替换--------------------------*/
	// 正则表达式替换
	reg := regexp.MustCompile("(eval\\((.*)\\)|script)")
	param = reg.ReplaceAllString(param, "")

	reg = regexp.MustCompile("[\"']\\s*javascript:(.*)[\"']")
	param = reg.ReplaceAllString(param, "")

	if temp != param {
		// System.out.println("输入信息存在xss攻击！");
		// System.out.println("原始输入信息-->" + temp);
		// System.out.println("处理后信息-->" + src);

		fmt.Printf("xss攻击检查：参数含有非法攻击字符，已禁止继续访问！！\n")
		fmt.Printf("原始输入信息-->%s\n", temp)
		return fmt.Sprintf("xss攻击检查：参数含有非法攻击字符，已禁止继续访问！！")
	}
	return param
}

// sqlInject SQL校验
func sqlInject(param string) bool {
	param = util.TrimSpace(param)
	if param == "" {
		return false
	}

	//去掉'|"|;|\字符
	param = strings.ReplaceAll(param, "'", "")
	param = strings.ReplaceAll(param, "\"", "")
	param = strings.ReplaceAll(param, ";", "")
	param = strings.ReplaceAll(param, "\\", "")

	//转换成小写
	param = strings.ToLower(param)

	//非法字符
	keywords := []string{"master", "truncate", "insert", "select", "delete", "update", "declare", "alert", "alter", "drop"}

	//判断是否包含非法字符
	for _, keyword := range keywords {
		if strings.Contains(param, keyword) {
			return true
		}
	}

	return false
}

// xss攻击拦截
func xssClean(param string) (string, error) {
	param = util.TrimSpace(param)
	if param == "" {
		return param, nil
	}

	//非法字符
	keywords := []string{"<", ">", "<>", "()", ")", "(", "javascript:", "script", "alter", "''", "'"}

	//判断是否包含非法字符
	for _, keyword := range keywords {
		if strings.Contains(param, keyword) {
			return "", fmt.Errorf("参数含有非法攻击字符，已禁止继续访问！")
		}
	}

	return param, nil
}

// FormatAtom 格式化数据
func FormatAtom(FieldValue reflect.Value) string {
	switch FieldValue.Kind() {
	case reflect.Invalid:
		return "invalid"
	case reflect.String:
		return FieldValue.String()
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		return strconv.FormatInt(FieldValue.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(FieldValue.Uint(), 10)
	// ...floating-point and complex cases omitted for brevity...
	case reflect.Bool:
		return strconv.FormatBool(FieldValue.Bool())
	case reflect.Chan, reflect.Func, reflect.Ptr, reflect.Slice, reflect.Map:
		return FieldValue.Type().String() + " 0x" +
			strconv.FormatUint(uint64(FieldValue.Pointer()), 16)
	default: // reflect.Array, reflect.Struct, reflect.Interface
		return FieldValue.Type().String() + " value"
	}
}

// SanitizeInterface 过滤特殊字符
func SanitizeInterface(input interface{}) interface{} {
	value := reflect.ValueOf(input)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	if value.Kind() != reflect.Struct {
		return input
	}

	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		if field.Kind() == reflect.String {
			field.SetString(html.EscapeString(field.String()))
		}
	}

	return input
}
