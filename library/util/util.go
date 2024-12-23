package util

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	craned "crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/gofrs/uuid"
	"github.com/jinzhu/now"
	"github.com/jung-kurt/gofpdf"
	"github.com/signintech/gopdf"
	"github.com/xuri/excelize/v2"
	"html"
	"io"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"unsafe"
	"weihu_server/library/config"
)

var snowNode *snowflake.Node

func init() {
	node := int64(rand.Intn(1000))
	fmt.Println("util, init", node)
	snowNode, _ = snowflake.NewNode(node)
}

// GetID 获取雪花算法ID
func GetID() string {
	// Generate a snowflake ID.
	id := snowNode.Generate()
	return fmt.Sprintf("%d", id)
}

// GetUUID 获取UUID
func GetUUID() string {
	u4, _ := uuid.NewV4()
	return u4.String()
}

// Md5 MD5加密
func Md5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

// IntToString int 转 string
func IntToString(num int) string {
	return fmt.Sprintf("%d", num)
}

// Int64ToString int64 转 string
func Int64ToString(num int64) string {
	return fmt.Sprintf("%d", num)
}

// Float64ToString float64 转 string
func Float64ToString(num float64) string {
	return fmt.Sprintf("%f", num)
}

// InStringSlice 判断字符串是否在数组内
func InStringSlice(target string, arr []string) bool {
	tmpList := make([]string, len(arr))
	// 目标的修改不会影响源
	copy(tmpList, arr)
	sort.Strings(tmpList)
	index := sort.SearchStrings(tmpList, target)
	if index < len(tmpList) && tmpList[index] == target {
		return true
	}
	return false
}

// StringToInt64 字符串转INT64
func StringToInt64(str string) int64 {
	num, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0
	}
	return num
}

// StringToFloat64 字符串转float64
func StringToFloat64(str string) float64 {
	num, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0
	}
	return num
}

// StringToBool 字符串转bool
func StringToBool(str string) bool {
	b, err := strconv.ParseBool(str)
	if err != nil {
		return false
	}
	return b
}

// CheckTimeFormat 检查日期格式是否正确
func CheckTimeFormat(date string, sType int) bool {
	timeFormatTpl := "2006-01-02 15:04:05"
	switch sType {
	case 1:
		timeFormatTpl = "2006-01-02 15:04:05.000"
		break
	case 2:
		timeFormatTpl = "2006-01-02 15:04:05"
		break
	case 3:
		timeFormatTpl = "2006-01-02 15:04"
		break
	case 4:
		timeFormatTpl = "2006-01-02 15"
		break
	case 5:
		timeFormatTpl = "2006-01-02"
		break
	case 6:
		timeFormatTpl = "2006-01"
		break
	case 7:
		timeFormatTpl = "2006"
		break
	case 8:
		timeFormatTpl = "2006/01/02 15:04:05"
		break
	}
	_, err := time.Parse(timeFormatTpl, date)
	if err != nil {
		return false
	}
	return true
}

// TimeFormat 格式化时间 1-毫秒 2-秒 3-分 4-时 5-天 6-月 7-年
func TimeFormat(t int64, sType int) string {
	if t <= 0 {
		return ""
	}
	if len(fmt.Sprintf("%d", t)) == 10 {
		t = t * 1000
	} else if len(fmt.Sprintf("%d", t)) == 13 {

	} else {
		return ""
	}
	curTime := time.UnixMilli(t)
	switch sType {
	case 1:
		return curTime.Format("2006-01-02 15:04:05.000")
	case 2:
		return curTime.Format("2006-01-02 15:04:05")
	case 3:
		return curTime.Format("2006-01-02 15:04")
	case 4:
		return curTime.Format("2006-01-02 15")
	case 5:
		return curTime.Format("2006-01-02")
	case 6:
		return curTime.Format("2006-01")
	case 7:
		return curTime.Format("2006")
	case 8:
		return curTime.Format("2006/01/02")
	case 9:
		return curTime.Format("2006/01/02 15:04:05")
	}
	return ""
}

// TimeFormatByTimezone 格式化时间 1-毫秒 2-秒 3-分 4-时 5-天 6-月 7-年
func TimeFormatByTimezone(t int64, sType int, timeZone string) string {
	if t <= 0 {
		return ""
	}
	if len(fmt.Sprintf("%d", t)) == 10 {
		t = t * 1000
	} else if len(fmt.Sprintf("%d", t)) == 13 {

	} else {
		return ""
	}

	if timeZone != "" {
		t = ConvertTimestampToTargetTimezone(t, time.Local.String(), timeZone)
	}

	curTime := time.UnixMilli(t)

	switch sType {
	case 1:
		return curTime.Format("2006-01-02 15:04:05.000")
	case 2:
		return curTime.Format("2006-01-02 15:04:05")
	case 3:
		return curTime.Format("2006-01-02 15:04")
	case 4:
		return curTime.Format("2006-01-02 15")
	case 5:
		return curTime.Format("2006-01-02")
	case 6:
		return curTime.Format("2006-01")
	case 7:
		return curTime.Format("2006")
	case 8:
		return curTime.Format("2006/01/02")
	case 9:
		return curTime.Format("2006/01/02 15:04:05")
	}
	return ""
}

// TimeFormat1 格式化时间
func TimeFormat1(t int64) string {
	if t <= 0 {
		return ""
	}
	if len(fmt.Sprintf("%d", t)) == 10 {
		t = t * 1000
	} else if len(fmt.Sprintf("%d", t)) == 13 {

	} else {
		return ""
	}
	curTime := time.UnixMilli(t).Format("2006-01-02")

	return curTime
}

// TimeFormat2 格式化时间
func TimeFormat2(t int64) string {
	if t <= 0 {
		return ""
	}
	if len(fmt.Sprintf("%d", t)) == 10 {
		t = t * 1000
	} else if len(fmt.Sprintf("%d", t)) == 13 {

	} else {
		return ""
	}
	curTime := time.UnixMilli(t).Format("2006-01-02T15:04:05")

	return curTime
}

// StringToInt string转int
func StringToInt(str string) int {
	num, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return num
}

// HmacSha256 计算HmacSha256
// key 是加密所使用的key
// data 是加密的内容
func HmacSha256(key string, data string) []byte {
	mac := hmac.New(sha256.New, []byte(key))
	_, _ = mac.Write([]byte(data))

	return mac.Sum(nil)
}

// HmacSha256ToHex 将加密后的二进制转16进制字符串
func HmacSha256ToHex(key string, data string) string {
	return hex.EncodeToString(HmacSha256(key, data))
}

// HmacSha256ToBase64 将加密后的二进制转Base64字符串
func HmacSha256ToBase64(key string, data string) string {
	return base64.URLEncoding.EncodeToString(HmacSha256(key, data))
}

// HmacSha1 计算HmacSha1
// keyStr 是加密所使用的key
// data 是加密的内容
func HmacSha1(keyStr, data string) string {
	// Crypto by HMAC-SHA1
	key := []byte(keyStr)
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(data))

	//进行base64编码
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

// HttpPost http post请求
func HttpPost(url, json string, headers map[string]string) (int, []byte, error) {
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(network, addr string) (net.Conn, error) {
				deadline := time.Now().Add(5 * time.Minute)
				c, err := net.DialTimeout(network, addr, 5*time.Minute)
				if err != nil {
					return nil, err
				}
				c.SetDeadline(deadline)
				return c, nil
			},
			ResponseHeaderTimeout: 5 * time.Minute,
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		},
	}

	request, err := http.NewRequest("POST", url, strings.NewReader(json))
	if err != nil {
		return 0, nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	if len(headers) > 0 {
		for k, v := range headers {
			request.Header.Set(k, v)
		}
	}
	//post数据并接收http响应
	resp, err := client.Do(request)
	if err != nil {
		return 0, nil, err
	}

	bys, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	log.Println(string(bys), url)

	if err != nil {
		return 0, nil, err
	}
	return resp.StatusCode, bys, nil
}

// HttpGet http get请求
func HttpGet(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, err2 := ioutil.ReadAll(response.Body)
	if err2 != nil {
		return nil, err2
	}
	return body, err
}

func HttpRequest(url string, method string, body io.Reader, headers map[string]string) ([]byte, error) {
	client := &http.Client{}
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	if len(headers) > 0 {
		for k, v := range headers {
			request.Header.Set(k, v)
		}
	}

	//处理返回结果
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	buf, err := ioutil.ReadAll(response.Body)

	if response.StatusCode == http.StatusOK {
		return buf, nil
	} else {
		return nil, errors.New(fmt.Sprint("StatusCode=", response.StatusCode, " msg=", string(buf)))
	}
}

// ParseExcel 解析excel
func ParseExcel(file *multipart.FileHeader) ([][]string, error) {
	f, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	xlsx, err := excelize.OpenReader(f)
	if err != nil {
		return nil, err
	}
	rows, err := xlsx.GetRows("Sheet1")
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// TimeParse 时间转换
func TimeParse(date string) int64 {
	dateTime, _ := time.ParseInLocation("2006-01-02 15:04:05", date, time.Local)
	return dateTime.UnixMilli()
}

// TimeParse2 时间转换
func TimeParse2(date string) int64 {
	dateTime, _ := time.ParseInLocation("2006/01/02 15:04:05", date, time.Local)
	return dateTime.UnixMilli()
}

// TimeParseByLayout 根据类型时间转换
func TimeParseByLayout(date string, sType int64) int64 {
	if date == "" {
		return 0
	}
	var layout string
	switch sType {
	case 1:
		layout = "2006-01-02 15:04:05.000"
	case 2:
		layout = "2006-01-02 15:04:05"
	case 3:
		layout = "2006-01-02 15:04"
	case 4:
		layout = "2006-01-02 15"
	case 5:
		layout = "2006-01-02"
	case 6:
		layout = "2006-01"
	case 7:
		layout = "2006"
	}
	dateTime, _ := time.ParseInLocation(layout, date, time.Local)
	return dateTime.UnixMilli()
}

// TimeParseByLayoutByTimezone 根据类型时间转换
func TimeParseByLayoutByTimezone(date string, sType int64, timeZone string) int64 {
	if date == "" {
		return 0
	}
	var layout string
	switch sType {
	case 1:
		layout = "2006-01-02 15:04:05.000"
	case 2:
		layout = "2006-01-02 15:04:05"
	case 3:
		layout = "2006-01-02 15:04"
	case 4:
		layout = "2006-01-02 15"
	case 5:
		layout = "2006-01-02"
	case 6:
		layout = "2006-01"
	case 7:
		layout = "2006"
	}

	location := time.Local
	fromTimezone := time.Local.String()
	if timeZone != "" {
		fromLocation, err := time.LoadLocation(timeZone)
		if err != nil {
			log.Printf("Error: %v", err)
		} else {
			location = fromLocation
			fromTimezone = timeZone
		}
	}

	dateTime, _ := time.ParseInLocation(layout, date, location)

	return ConvertTimestampToTargetTimezone(dateTime.UnixMilli(), fromTimezone, time.Local.String())
}

// InsertCsv 导出csv
func InsertCsv(fileName string, titleList []string, dataList [][]string) {
	fileSavePath := fmt.Sprintf("./%s", fileName)

	file, err := os.Open(fileSavePath)
	defer file.Close()

	// 如果文件不存在，创建文件
	if err != nil && os.IsNotExist(err) {
		nfs, er := os.Create(fileSavePath)
		if er != nil {
			panic(er)
		}
		defer nfs.Close()

		nfs.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM，避免使用Microsoft Excel打开乱码
		// 写入字段标题
		w := csv.NewWriter(nfs) //创建一个新的写入文件流
		//设置属性
		w.Comma = ','
		w.UseCRLF = true

		w.Write(titleList)
		for _, data := range dataList {
			w.Write(data)
		}

		// 这里必须刷新，才能将数据写入文件。
		w.Flush()
	} else {
		// 如果文件存在，直接加在末尾
		nfs, er := os.OpenFile(fileSavePath, os.O_APPEND|os.O_RDWR, 0666)
		defer nfs.Close()
		if er != nil {
			panic(er)
		}
		nfs.Seek(0, io.SeekEnd)
		w := csv.NewWriter(nfs) //创建一个新的写入文件流
		//设置属性
		w.Comma = ','
		w.UseCRLF = true

		for _, data := range dataList {
			w.Write(data)
		}
		//这里必须刷新，才能将数据写入文件。
		w.Flush()
	}
}

// WriteExcel 写入excel文件
func WriteExcel(fileName, path string, titleList []string, dataList [][]string) {
	tmpPath := filepath.Join(config.GetString("server.tmpPath"), path)
	//判断文件夹是否存在
	_, err := os.Stat(tmpPath)
	if err != nil {
		//不存在先新建文件夹
		err = os.MkdirAll(tmpPath, 0666)
		if err != nil {
			return
		}
	}
	fileSavePath := filepath.Join(tmpPath, fileName)
	file, err := os.Open(fileSavePath)
	defer file.Close()

	// 如果文件不存在，创建文件
	if err != nil && os.IsNotExist(err) {
		nfs, er := os.Create(fileSavePath)
		defer nfs.Close()
		if er != nil {
			return
		}

		// 写入字段标题
		xlsx := excelize.NewFile()
		xlsx.SetSheetRow("Sheet1", "A1", &titleList)
		for i, data := range dataList {
			xlsx.SetSheetRow("Sheet1", fmt.Sprintf("A%d", i+2), &data)
		}
		err = xlsx.SaveAs(fileSavePath)
	} else {
		// 如果文件存在，直接加在末尾
		xlsx, err := excelize.OpenFile(fileSavePath)
		defer xlsx.Close()
		if err != nil {
			return
		}
		rows, err := xlsx.GetRows("Sheet1")
		if err != nil {
			return
		}
		for i, data := range dataList {
			xlsx.SetSheetRow("Sheet1", fmt.Sprintf("A%d", len(rows)+i+1), &data)
		}
		xlsx.SaveAs(fileSavePath)
	}
}

// JsonToString json转字符串
func JsonToString(data interface{}) string {
	bys, _ := json.Marshal(data)
	return string(bys)
}

// JsonToMap json转map
func JsonToMap(jsonStr string) map[string]interface{} {
	dataMap := make(map[string]interface{})
	// 采用 decode+UseNumber() 来实现反序列化 防止精度丢失
	decoder := json.NewDecoder(bytes.NewBufferString(jsonStr))
	decoder.UseNumber() // 指定使用 Number 类型
	err := decoder.Decode(&dataMap)
	//err := json.Unmarshal([]byte(jsonStr), &dataMap)
	if err != nil {
		fmt.Printf("Unmarshal with error: %+v\n", err)
		return nil
	}

	return dataMap
}

// ByteToMap byte转map
func ByteToMap(jsonByte []byte) map[string]interface{} {
	dataMap := make(map[string]interface{})
	// 采用 decode+UseNumber() 来实现反序列化 防止精度丢失
	decoder := json.NewDecoder(bytes.NewBuffer(jsonByte))
	decoder.UseNumber() // 指定使用 Number 类型
	err := decoder.Decode(&dataMap)

	//err := json.Unmarshal(jsonByte, &dataMap)
	if err != nil {
		fmt.Printf("Unmarshal with error: %+v\n", err)
		return nil
	}

	return dataMap
}

func DownloadFile(url string, fileName string) (filePath string, err error) {
	//上传的文件路径
	uploadDir := config.GetString("server.tmpPath")
	//判断文件夹是否存在
	_, err = os.Stat(uploadDir)
	if err != nil {
		//不存在先新建文件夹
		err = os.MkdirAll(uploadDir, 0666)
		if err != nil {
			fmt.Printf("创建文件夹失败：%s", err.Error())
			return
		}
	}
	filePath = uploadDir + "/" + fileName
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("下载文件失败：%s", err.Error())
		return
	}
	defer resp.Body.Close()
	out, err := os.Create(filePath)
	if err != nil {
		log.Printf("创建文件失败：%s", err.Error())
		return
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Printf("写入文件失败：%s", err.Error())
		return
	}
	return
}

// ParseExcelFile 解析excel文件
func ParseExcelFile(filePath string) ([][]string, error) {
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	xlsx, err := excelize.OpenReader(file)
	if err != nil {
		return nil, err
	}
	rows, err := xlsx.GetRows("Sheet1")
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// GenerateRandomString 生成随机字符串
func GenerateRandomString(length int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// GenerateRandomPwd 生成随机密码
func GenerateRandomPwd(length int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// GenerateRandomInt 生成指定长度随机整数
func GenerateRandomInt(length int) int {
	const letters = "0123456789"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	num, _ := strconv.Atoi(string(b))
	return num
}

// UniqueSlice 去除切片中的重复值
func UniqueSlice(list []string) (result []string) {
	result = make([]string, 0)
	// 存储每个元素的数量
	countEle := make(map[string]int)
	for _, value := range list {
		countEle[value] += 1
		if countEle[value] > 1 {
			continue
		}
		result = append(result, value)
	}
	return
}

// UniqueIntSlice 去除int切片中的重复值
func UniqueIntSlice(list []int64) (result []int64) {
	// 存储每个元素的数量
	countEle := make(map[int64]int)
	for _, value := range list {
		countEle[value] += 1
		if countEle[value] > 1 {
			continue
		}
		result = append(result, value)
	}
	return
}

// DiffSlice 计算2个切片的差集
func DiffSlice(slice1, slice2 []string) (result []string) {
	for _, v := range slice1 {
		if !InStrSlice(slice2, v) {
			result = append(result, v)
		}
	}
	return
}

// MergeSlice 合并多个切片
func MergeSlice(list ...[]string) (result []string) {
	for _, value := range list {
		result = append(result, value...)
	}
	return
}

// MergeInt64Slice 合并多个Int64切片
func MergeInt64Slice(list ...[]int64) (result []int64) {
	for _, value := range list {
		result = append(result, value...)
	}
	return
}

// IsPhoneNew 判断是否电话号，包含国际主流格式号码
func IsPhoneNew(phone string) bool {
	reg := regexp.MustCompile(`^(?:\+?(\d{1,4}))?[-. (]*(\d{1,3})[-. )]*(\d{1,4}[-. ]*){1,2}(\d{1,9})$`)
	return reg.MatchString(phone)
}

// IsPhone 判断字符串是否手机号
func IsPhone(phone string) bool {
	reg := regexp.MustCompile(`^1[3456789]\d{9}$`)
	return reg.MatchString(phone)
}

// IsPhoneUS 判断字符串是否美国手机号
func IsPhoneUS(phone string) bool {
	str := []string{`^\(?(\d{3})\)?[-.\s]?\d{3}[-.\s]?\d{4}$`, `^(\+?1[-.\s]?)?\d{3}-\d{3}-\d{4}$`}
	for _, v := range str {
		reg := regexp.MustCompile(v)
		if reg.MatchString(phone) {
			return true
		}
	}
	return false
}

// IsName 判断字符串是否姓名
func IsName(name string) bool {
	reg := regexp.MustCompile("^[\u4e00-\u9fa5a-zA-Z0-9]{2,4}$")
	return reg.MatchString(name)
}

// IsWechat 判断字符串是否微信号
func IsWechat(wechat string) bool {
	reg := regexp.MustCompile(`^[a-zA-Z\d_]{5,}$`)
	return reg.MatchString(wechat)
}

// Int64SliceToString int64切片转为字符串
func Int64SliceToString(list []int64) string {
	var result []string
	for _, value := range list {
		result = append(result, strconv.FormatInt(value, 10))
	}
	return strings.Join(result, ",")
}

// Int64SliceToStrSlice ing64切片转string切片
func Int64SliceToStrSlice(list []int64) []string {
	var result []string
	for _, value := range list {
		result = append(result, strconv.FormatInt(value, 10))
	}
	return result
}

// CheckMobile 检验手机号
func CheckMobile(phone string) bool {
	// 匹配规则
	// ^1第一位为一
	// [3456789]{1} 后接一位3456789 的数字
	// \\d \d的转义 表示数字 {9} 接9位
	// $ 结束符
	regRuler := "^1[3456789]{1}\\d{9}$"

	// 正则调用规则
	reg := regexp.MustCompile(regRuler)

	// 返回 MatchString 是否匹配
	return reg.MatchString(phone)
}

// CheckQQ 检验QQ号
func CheckQQ(qq string) bool {
	// $ 结束符
	regRuler := "^[1-9][0-9]{4,10}$"

	// 正则调用规则
	reg := regexp.MustCompile(regRuler)

	// 返回 MatchString 是否匹配
	return reg.MatchString(qq)
}

// CheckStandingCard 检验身份证
func CheckStandingCard(no string) bool {
	// $ 结束符
	//regRuler := "(\\d{15}$|\\d{18}$|\\d{17}(\\d|X|x))"
	regRuler := "^[1-9]\\d{5}(18|19|20)\\d{2}(0[1-9]|1[0-2])(0[1-9]|[1-2][0-9]|3[0-1])\\d{3}(\\d|X|x)$"

	// 正则调用规则
	reg := regexp.MustCompile(regRuler)

	// 返回 MatchString 是否匹配
	return reg.MatchString(no)
}

// CheckPassport 检验护照
func CheckPassport(no string) bool {
	// $ 结束符
	regRuler := "^[a-zA-Z0-9]{3,21}$"

	// 正则调用规则
	reg := regexp.MustCompile(regRuler)

	// 返回 MatchString 是否匹配
	return reg.MatchString(no)
}

// CheckOfficerCard 检验军官证
func CheckOfficerCard(no string) bool {
	// $ 结束符
	regRuler := "^[a-zA-Z0-9]{7,21}$"

	// 正则调用规则
	reg := regexp.MustCompile(regRuler)

	// 返回 MatchString 是否匹配
	return reg.MatchString(no)
}

// CheckOrganizationCode 检验组织机构代码
func CheckOrganizationCode(no string) bool {
	// $ 结束符
	regRuler := "^[A-Z0-9]{8}-[A-Z0-9]$"

	// 正则调用规则
	reg := regexp.MustCompile(regRuler)

	// 返回 MatchString 是否匹配
	return reg.MatchString(no)
}

// CheckBusinessLicense 检验营业执照
func CheckBusinessLicense(no string) bool {
	// $ 结束符
	regRuler := "^[a-zA-Z0-9]{10,20}$"

	// 正则调用规则
	reg := regexp.MustCompile(regRuler)

	// 返回 MatchString 是否匹配
	return reg.MatchString(no)
}

// CheckCreditCode 检验统一社会信用代码
func CheckCreditCode(no string) bool {
	// $ 结束符
	regRuler := "[A-Z0-9]{18}$"

	// 正则调用规则
	reg := regexp.MustCompile(regRuler)

	// 返回 MatchString 是否匹配
	return reg.MatchString(no)
}

// InSlice 存在于切片
func InSlice(s []int64, e int64) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

// InStrSlice 存在于切片
func InStrSlice(s []string, e string) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

// GetWeekStartAndEnd 获取本周起止时间
func GetWeekStartAndEnd() (start, end string) {
	now := time.Now()
	offset := int(time.Monday - now.Weekday())
	if offset > 0 {
		offset = -6
	}
	start = now.AddDate(0, 0, offset).Format("2006-01-02")
	end = now.AddDate(0, 0, offset+6).Format("2006-01-02")
	return
}

// GetTodayStartAndEnd 今日起止时间
func GetTodayStartAndEnd() (start, end string) {
	now := time.Now()
	start = now.Format("2006-01-02")
	end = now.AddDate(0, 0, 1).Format("2006-01-02")
	return
}

// GetMonthStartAndEnd 本月起止时间
func GetMonthStartAndEnd() (start, end string) {
	now := time.Now()
	start = now.AddDate(0, 0, -now.Day()+1).Format("2006-01-02")
	end = now.AddDate(0, 1, -now.Day()).Format("2006-01-02")
	return
}

func Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.1f", value), 64)
	return value
}

func GetRate(a int64, b int64) float64 {
	if b == 0 {
		return 0
	}
	return Decimal(float64(a) / float64(b) * 100)
}

// AddZero 数字前面补0
func AddZero(num int64, length int) string {
	return fmt.Sprintf("%0"+strconv.Itoa(length)+"d", num)
}

func Substr(str string, start int, length int) string {
	if str == "" {
		return ""
	}
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}

	return string(rs[start:end])
}

// AddDate Time.AddDate方法对月年做加减时，若当前天数大于目标月的天数，则会自动再加上多余的天数,重写此方法
func AddDate(t time.Time, year, month, day int) time.Time {
	//先跳到目标月的1号
	targetDate := t.AddDate(year, month, -t.Day()+1)
	//获取目标月的临界值
	targetDay := targetDate.AddDate(0, 1, -1).Day()
	//对比临界值与源日期值，取最小的值
	if targetDay > t.Day() {
		targetDay = t.Day()
	}
	//最后用目标月的1号加上目标值和入参的天数
	targetDate = targetDate.AddDate(0, 0, targetDay-1+day)
	return targetDate
}

func GetLocDateString(timestamp int64) string {
	if timestamp == 0 {
		return ""
	}
	timeLen := len(fmt.Sprintf("%d", timestamp))
	if timeLen == 10 {
		return time.Unix(timestamp, 0).Format("2006年01月02日")
	} else if timeLen == 13 {
		return time.UnixMilli(timestamp).Format("2006年01月02日")
	}
	return ""
}

// RandSlice 随机获取切片中的一个值
func RandSlice(s []int64) int64 {
	rand.Seed(time.Now().UnixNano())
	return s[rand.Intn(len(s))]
}

// GetDayList 获取最近一个月时间中文字符串切片 12月1日
func GetDayList() []string {
	var list []string
	now := time.Now()
	for i := 29; i >= 0; i-- {
		list = append(list, now.AddDate(0, 0, -i).Format("01月02日"))
	}
	return list
}

// GetMonthList 获取最近1年月份字符串切片
func GetMonthList() []string {
	var list []string
	now := time.Now()
	for i := 11; i >= 0; i-- {
		list = append(list, AddDate(now, 0, -i, 0).Format("2006年01月"))
	}
	return list
}

func customAddDate(t time.Time, year, month, day int) time.Time {
	//先跳到目标月的1号
	targetDate := t.AddDate(year, month, -t.Day()+1)
	//获取目标月的临界值
	targetDay := targetDate.AddDate(0, 1, -1).Day()
	//对比临界值与源日期值，取最小的值
	if targetDay > t.Day() {
		targetDay = t.Day()
	}
	//最后用目标月的1号加上目标值和入参的天数
	targetDate = targetDate.AddDate(0, 0, targetDay-1+day)
	return targetDate
}

// GetRandomInt 获取随机整数
func GetRandomInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

// TimeFormatParseTimestamp 获取当前时间的当天开始时间或结束时间
// date 时间
// sType 1-开始时间 2-结束时间
func TimeFormatParseTimestamp(date string, sType int64) int64 {
	dueTimeCst, _ := time.ParseInLocation("2006-01-02", date, time.Local)
	var timestamp time.Time
	switch sType {
	case 1:
		timestamp = now.New(dueTimeCst).BeginningOfDay()
	case 2:
		timestamp = now.New(dueTimeCst).EndOfDay()
	}
	return timestamp.UnixMilli()
}

// TimeFormatParseTimestampByTimezone 获取当前时间的当天开始时间或结束时间
// date 时间
// sType 1-开始时间 2-结束时间
func TimeFormatParseTimestampByTimezone(date string, sType int64, timeZone string) int64 {
	location := time.Local
	fromTimezone := time.Local.String()
	if timeZone != "" {
		fromLocation, err := time.LoadLocation(timeZone)
		if err != nil {
			log.Printf("Error: %v", err)
		} else {
			location = fromLocation
			fromTimezone = timeZone
		}
	}

	dueTimeCst, _ := time.ParseInLocation("2006-01-02", date, location)
	var timestamp time.Time
	switch sType {
	case 1:
		timestamp = now.New(dueTimeCst).BeginningOfDay()
	case 2:
		timestamp = now.New(dueTimeCst).EndOfDay()
	}

	return ConvertTimestampToTargetTimezone(timestamp.UnixMilli(), fromTimezone, time.Local.String())
}

// TrimSpace 去除空格
func TrimSpace(val string) string {
	reg := regexp.MustCompile(`( )+|(\n)+`)
	str := reg.ReplaceAllString(strings.TrimSpace(val), "$1$2")

	return str
}

// TimeStampToMillisecond 判断时间戳秒转毫秒
func TimeStampToMillisecond(timestamp int64) int64 {
	if timestamp < 10000000000 {
		timestamp = timestamp * 1000
	}
	return timestamp
}

// SortInt64Slice int64切片指定排序
func SortInt64Slice(slice []int64, sortType int) []int64 {
	switch sortType {
	case 1:
		sort.Slice(slice, func(i, j int) bool {
			return slice[i] < slice[j]
		})
	case 2:
		sort.Slice(slice, func(i, j int) bool {
			return slice[i] > slice[j]
		})
	}
	return slice
}

// RemoveSliceElements 删除切片多个元素
func RemoveSliceElements(slice []int64, elems []int64) []int64 {
	var result []int64
	for _, v := range slice {
		if !InSlice(elems, v) {
			result = append(result, v)
		}
	}
	return result
}

// RemoveStrSliceElements 删除字符串切片多个元素
func RemoveStrSliceElements(slice []string, elems []string) []string {
	var result []string
	for _, v := range slice {
		if !InStrSlice(elems, v) {
			result = append(result, v)
		}
	}
	return result
}

// StrReplace 字符串拼接替换
func StrReplace(slice string) string {
	return strings.ReplaceAll(slice, ",", "','")
}

// EqualStrSlice 判断两个切片元素是否相同
func EqualStrSlice(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for _, v := range a {
		if !InStrSlice(b, v) {
			return false
		}
	}
	return true
}

// EqualInt64Slice 判断两个int64切片是否相同
func EqualInt64Slice(a, b []int64) bool {
	if len(a) != len(b) {
		return false
	}
	for _, v := range a {
		if !InSlice(b, v) {
			return false
		}
	}
	return true
}

// RemovePhonePrefix 移除手机号前的区号
func RemovePhonePrefix(phone string) string {
	if len(phone) > 11 {
		return phone[len(phone)-11:]
	}
	return phone
}

// ResolveDomain 通过指定dns地址解析域名
func ResolveDomain(domain string, dns string) (string, error) {
	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Second * 5,
			}
			return d.DialContext(ctx, "udp", dns)
		},
	}
	ips, err := resolver.LookupHost(context.Background(), domain)
	if err != nil {
		return "", err
	}
	fmt.Printf("ips: %v\n", ips)
	return ips[0], nil
}

// GetWeekStartTimeAndEndTime 获取周开始时间和结束时间
func GetWeekStartTimeAndEndTime(d time.Time) (time.Time, time.Time) {
	offset := int(time.Monday - d.Weekday())
	if offset > 0 {
		offset = -6
	}
	startTime := AddDate(d, 0, 0, offset)
	endTime := AddDate(startTime, 0, 0, 6)
	return startTime, endTime
}

// GetMonthStartTimeAndEndTime 获取月开始时间和结束时间
func GetMonthStartTimeAndEndTime(d time.Time) (time.Time, time.Time) {
	startTime := AddDate(d, 0, 0, -d.Day()+1)
	endTime := AddDate(startTime, 0, 1, -1)

	return startTime, endTime
}

// GetYearStartTimeAndEndTime 获取年开始时间和结束时间
func GetYearStartTimeAndEndTime(d time.Time) (time.Time, time.Time) {
	startTime := AddDate(d, 0, 0, -d.YearDay()+1)
	endTime := AddDate(startTime, 1, 0, -1)

	return startTime, endTime
}

// GetBetweenDates 获取天时间列表
func GetBetweenDates(startTime, endTime time.Time, timeFormatTpl string) []string {
	d := make([]string, 0)
	if endTime.Before(startTime) {
		// 如果结束时间小于开始时间，异常
		return d
	}

	// 输出日期格式固定
	//timeFormatTpl := "2006-01-02"
	if timeFormatTpl == "" {
		timeFormatTpl = "01/02"
	}
	date2Str := endTime.Format(timeFormatTpl)
	d = append(d, startTime.Format(timeFormatTpl))
	if startTime.Format(timeFormatTpl) == date2Str {
		return d
	}
	for {
		startTime = AddDate(startTime, 0, 0, 1)
		dateStr := startTime.Format(timeFormatTpl)
		d = append(d, dateStr)
		if dateStr == date2Str {
			break
		}
	}
	return d
}

// EmojiEncode Emoji表情转码
func EmojiEncode(s string) string {
	ret := ""
	rs := []rune(s)
	for i := 0; i < len(rs); i++ {
		if len(string(rs[i])) == 4 {
			u := `[\u` + strconv.FormatInt(int64(rs[i]), 16) + `]`
			ret += u
		} else {
			ret += string(rs[i])
		}
	}

	return ret
}

// EmojiDecode Emoji表情解码
func EmojiDecode(s string) string {
	//emoji表情的数据表达式
	re := regexp.MustCompile("\\[[\\\\u0-9a-zA-Z]+\\]") //[u1f602]
	//提取emoji数据表达式
	reg := regexp.MustCompile("\\[\\\\u|]")
	src := re.FindAllString(s, -1)
	for i := 0; i < len(src); i++ {
		e := reg.ReplaceAllString(src[i], "")
		p, err := strconv.ParseInt(e, 16, 32)
		if err == nil {
			s = strings.Replace(s, src[i], string(rune(p)), -1)
		}
	}

	return s
}

// IsEmoji 判断是否包含表情
func IsEmoji(s string) bool {
	rs := []rune(s)
	for i := 0; i < len(rs); i++ {
		if len(string(rs[i])) == 4 {
			return true
		}
	}
	return false
}

// TimeLayoutParse 时间格式转换
func TimeLayoutParse(date string) string {
	dateTime, _ := time.ParseInLocation("2006-01-02", date, time.Local)
	return dateTime.Format("01/02")
}

// TimestampToStrDate 获取当前时间字符串 格式yyyyMMddHHmmssSSS 24小时制
func TimestampToStrDate() string {
	// 当前时间
	currentTime := time.Now().In(time.FixedZone("CST", 8*3600)) // 东八
	// 输出日期格式固定
	timeFormatTpl := "20060102150405.000"
	return strings.Replace(currentTime.Format(timeFormatTpl), ".", "", -1)
}

// Sha1 sha1加密
func Sha1(str string) string {
	sha1Ctx := sha1.New()
	sha1Ctx.Write([]byte(str))
	return hex.EncodeToString(sha1Ctx.Sum(nil))
}

// MaxInt64Slice 获取int64切片最大值
func MaxInt64Slice(slice []int64) int64 {
	if len(slice) == 0 {
		return 0
	}
	max := slice[0]
	for _, v := range slice {
		if v > max {
			max = v
		}
	}
	return max
}

// GetMonthDayList 根据时间获取本月初到月底的日期
func GetMonthDayList(date string) []string {
	var outs []string
	nowT, _ := time.Parse("2006-01-02", date)
	year, month, _ := nowT.Date()
	firstDay := time.Date(year, month, 1, 0, 0, 0, 0, nowT.Location())
	lastDay := firstDay.AddDate(0, 1, -1)
	for d := firstDay; !d.After(lastDay); d = d.AddDate(0, 0, 1) {
		outs = append(outs, d.Format("2006-01-02"))
	}
	return outs
}

// GetMonthStartAndEndTimestampByDate 根据时间获取本月初与月底的时间戳
func GetMonthStartAndEndTimestampByDate(date string) (int64, int64) {
	t, _ := time.Parse("2006-01-02", date)
	year, month, _ := t.Date()
	//获取本月开始时间
	monthStart := time.Date(year, month, 1, 0, 0, 0, 0, t.Location())
	//获取本月结束时间
	monthEnd := monthStart.AddDate(0, 1, -1)
	return monthStart.UnixMilli(), monthEnd.UnixMilli()
}

// GetFileNameByUrl 根据url获取文件名称(不含后缀)
func GetFileNameByUrl(url string) string {
	if url == "" {
		return ""
	}
	split := strings.Split(url, "/")
	if len(split) == 0 {
		return GetUUID()
	}
	return strings.Split(split[len(split)-1], ".")[0]
}

// GetFileExtensionFromURL 根据url获取文件后缀
func GetFileExtensionFromURL(urlStr string) string {
	// 解析 URL
	u, err := url.Parse(urlStr)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return ""
	}

	// 获取路径部分
	path := u.Path

	// 提取文件名
	fileName := filepath.Base(path)

	// 获取文件扩展名
	ext := filepath.Ext(fileName)

	// 去掉点（`.`）符号
	return strings.TrimPrefix(ext, ".")
}

// GetFileName 根据url获取文件名称
func GetFileName(url string) string {
	if url == "" {
		return ""
	}
	split := strings.Split(url, "/")
	if len(split) == 0 {
		return GetUUID()
	}
	return split[len(split)-1]
}

// GetAlertDateStr 毫秒时间戳转换为天小时分钟
func GetAlertDateStr(timestamp int64) string {
	alertDateStr := ""
	if timestamp == 0 {
		return alertDateStr
	}

	day := timestamp / 86400000
	if day > 0 {
		alertDateStr += fmt.Sprintf("%d天", day)
	}
	hour := (timestamp - day*86400000) / 3600000
	if hour > 0 {
		alertDateStr += fmt.Sprintf("%d小时", hour)
	}
	minute := (timestamp - day*86400000 - hour*3600000) / 60000
	if minute > 0 {
		alertDateStr += fmt.Sprintf("%d分钟", minute)
	}
	return alertDateStr
}

// GetAlertDateDay 毫秒时间戳转天
func GetAlertDateDay(timestamp int64) int64 {
	return timestamp / 86400000
}

// DivUp 除法运算并向上取整
func DivUp(a, b int64) int64 {
	if a%b == 0 {
		return a / b
	}
	return a/b + 1
}

// HttpPostMultiFile httpPost多文件上传
func HttpPostMultiFile(url string, headers map[string]string, params map[string]string, files map[string][]string, cookies []*http.Cookie) ([]byte, error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	for k, v := range params {
		_ = bodyWriter.WriteField(k, v)
	}

	for k, vs := range files {
		for _, v := range vs {
			fileWriter, err := bodyWriter.CreateFormFile(k, v)
			if err != nil {
				return nil, err
			}
			fh, err := os.Open(v)
			if err != nil {
				return nil, err
			}
			_, err = io.Copy(fileWriter, fh)
			if err != nil {
				return nil, err
			}
			fh.Close()
		}
	}

	contentType := bodyWriter.FormDataContentType()
	_ = bodyWriter.Close()

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bodyBuf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	for _, v := range cookies {
		req.AddCookie(v)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

// AesEncrypt AES加密
func AesEncrypt(content string, aesKey string) (res string, err error) {
	if content == "" {
		err = errors.New("content is nil")
		return
	}

	//base64解码
	keyByte, err := base64.StdEncoding.DecodeString(aesKey)
	if err != nil {
		return
	}
	aesBlockEncryptor, err := aes.NewCipher(keyByte)
	if err != nil {
		return
	}

	//128位全0向量
	iv := make([]byte, 16)
	//cbc模式
	aesEncryptor := cipher.NewCBCEncrypter(aesBlockEncryptor, iv)

	//aes加密
	contentByte := PKCS5Padding([]byte(content), aesBlockEncryptor.BlockSize())
	encrypted := make([]byte, len(contentByte))
	aesEncryptor.CryptBlocks(encrypted, contentByte)
	//加密结果base64编码
	res = base64.StdEncoding.EncodeToString(encrypted)
	return
}

// AesDecrypt Aes解密
func AesDecrypt(content string, aesKey string) (res string, err error) {
	if content == "" {
		err = errors.New("content is nil")
		return
	}

	//base64解码
	keyByte, err := base64.StdEncoding.DecodeString(aesKey)
	if err != nil {
		return
	}
	aesBlockDecryptor, err := aes.NewCipher(keyByte)
	if err != nil {
		return
	}

	//128位全0向量
	iv := make([]byte, 16)
	//cbc模式
	aesDecrypter := cipher.NewCBCDecrypter(aesBlockDecryptor, iv)

	//aes加密
	contentByte, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return
	}
	decrypted := make([]byte, len(contentByte))
	aesDecrypter.CryptBlocks(decrypted, contentByte)
	res = string(PKCS5Trimming(decrypted))
	return
}

// PKCS5Padding PKCS5包装
func PKCS5Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

// PKCS5Trimming 解包装
func PKCS5Trimming(encrypt []byte) []byte {
	padding := encrypt[len(encrypt)-1]
	return encrypt[:len(encrypt)-int(padding)]
}

// GetFileSize 获取文件大小
func GetFileSize(filePath string) int64 {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return 0
	}
	return fileInfo.Size()
}

// GetFileMd5 获取文件md5
func GetFileMd5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	md5Hash := md5.New()
	if _, err := io.Copy(md5Hash, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(md5Hash.Sum(nil)), nil
}

// ToString 将任意对象转字符串
func ToString(v interface{}) string {
	switch v.(type) {
	case string:
		return v.(string)
	case int:
		return strconv.Itoa(v.(int))
	case int64:
		return strconv.FormatInt(v.(int64), 10)
	case float64:
		return strconv.FormatFloat(v.(float64), 'f', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(v.(float32)), 'f', -1, 32)
	case bool:
		return strconv.FormatBool(v.(bool))
	default:
		return ""
	}
}

// InterfaceToStruct json转struct
func InterfaceToStruct(obj interface{}, out interface{}) error {
	index, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(bytes.NewBufferString(string(index)))
	decoder.UseNumber() // 指定使用 Number 类型
	err = decoder.Decode(out)
	if err != nil {
		fmt.Printf("Unmarshal with error: %+v\n", err)
		return nil
	}
	return nil
}

// JsonToStruct json转struct
func JsonToStruct(obj string, out interface{}) error {
	decoder := json.NewDecoder(bytes.NewBufferString(obj))
	decoder.UseNumber() // 指定使用 Number 类型
	err := decoder.Decode(out)
	if err != nil {
		fmt.Printf("Unmarshal with error: %+v\n", err)
		return nil
	}
	return nil
}

// GetRandomPassword 生成大小写字母+数字+特殊字符（@或下划线)三种组合，长度10~18位的随机密码
func GetRandomPassword() string {
	var password string
	var length int
	var randNum int
	var temp string
	var randStr = []string{"abcdefghijklmnopqrstuvwxyz", "ABCDEFGHIJKLMNOPQRSTUVWXYZ", "0123456789", "@_"}
	for i := 0; i < 3; i++ {
		randNum = rand.Intn(4)
		temp = randStr[randNum]
		randNum = rand.Intn(len(temp))
		password += string(temp[randNum])
	}
	length = rand.Intn(9) + 10
	for i := 0; i < length-3; i++ {
		randNum = rand.Intn(3)
		temp = randStr[randNum]
		randNum = rand.Intn(len(temp))
		password += string(temp[randNum])
	}
	return password
}

// GeneratePassword 随机生成一个包含至少一个数字和一个字母的密码
func GeneratePassword(length int) string {
	if length < 2 {
		return ""
	}

	// 定义字符集
	const (
		letters  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
		digits   = "0123456789"
		allChars = letters + digits
	)

	rand.Seed(time.Now().UnixNano())

	// 确保密码中至少包含一个字母和一个数字
	password := []rune{
		rune(letters[rand.Intn(len(letters))]),
		rune(digits[rand.Intn(len(digits))]),
	}

	// 生成剩余的字符
	for i := 2; i < length; i++ {
		password = append(password, rune(allChars[rand.Intn(len(allChars))]))
	}

	// 打乱密码字符顺序
	rand.Shuffle(len(password), func(i, j int) {
		password[i], password[j] = password[j], password[i]
	})

	return string(password)
}

// EmptyToNone 将空字符串转换为无
func EmptyToNone(str string) string {
	if str == "" {
		return "无"
	}
	return str
}

// ConvertImageToPdf 将多个图片转为pdf
func ConvertImageToPdf(imagePaths []string, pdfPath string) error {
	defer func() {
		fmt.Println("耗时：", time.Since(time.Now()).Nanoseconds())
	}()
	pdf := gofpdf.New("P", "mm", "A4", "")
	for _, v := range imagePaths {
		pdf.AddPage()
		pdf.ImageOptions(v, 0, 0, 210, 297, false, gofpdf.ImageOptions{ReadDpi: true}, 0, "")
	}
	return pdf.OutputFileAndClose(pdfPath)
}

// ConvertImageToPdf2 将多个图片转为pdf
func ConvertImageToPdf2(imagePaths []string, pdfPath string) error {
	defer func() {
		fmt.Println("耗时：", time.Since(time.Now()).Nanoseconds())
	}()
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	defer pdf.Close()
	for _, path := range imagePaths {
		pdf.AddPage()
		err := pdf.Image(path, 0, 0, nil)
		if err != nil {
			fmt.Println("Error adding image:", err)
			return err
		}
	}

	//pdf.SetCompressLevel(-10)
	err := pdf.WritePdf(pdfPath)
	if err != nil {
		fmt.Println("Error writing PDF:", err)
		return err
	}

	fmt.Println("PDF created successfully")
	return nil
}

// EscapeString 将字符串转义
func EscapeString(str string) string {
	if str == "" {
		return ""
	}
	return html.EscapeString(str)
}

// IsRepeat 判断数组里有效值是否有重复元素
func IsRepeat(arr []string) bool {
	var m = make(map[string]bool)
	for _, v := range arr {
		if _, ok := m[v]; ok {
			return true
		}
		m[v] = true
	}
	return false
}

func isXSS(input string) bool {
	// 检测是否包含<、>、&、;这些XSS攻击常见字符
	input = strings.ToLower(input)
	dangerousChars := []string{"script>", "</script", "javascript:", "alert(", "><", "truncate ", "insert ", "select ", "delete ", "update ", "declare ", "alert ", "alter ", "drop "}

	for _, char := range dangerousChars {
		if strings.Contains(input, char) {
			return true
		}
	}
	return false
}

func HaveXssAttack(in interface{}) bool {
	if in == nil {
		return false
	}
	jsonPost, _ := json.Marshal(in)
	body := string(jsonPost)
	if body != "" && isXSS(body) {
		return true
	}
	return false
}

// IsSame 判断两个切片是否相同
func IsSame(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for _, v := range a {
		if InStrSlice(b, v) == false {
			return false
		}
	}
	return true
}

// GetTxImSign 获取腾讯IM签名
func GetTxImSign(appId string, identifier string) string {
	out, _ := exec.Command("./signature", "config/txim/txim_private_key", appId, identifier).Output()
	return strings.Replace(string(out), "\n", "", -1)
}

// InterfaceToMd5 将interface转为md5
func InterfaceToMd5(in interface{}) string {
	jsonPost, _ := json.Marshal(in)
	return fmt.Sprintf("%x", md5.Sum(jsonPost))
}

// CompareVersion 比较两个版本号大小，v1>v2返回1，v1<v2返回-1，v1=v2返回0
func CompareVersion(v1, v2 string) int {
	v1Arr := strings.Split(v1, ".")
	v2Arr := strings.Split(v2, ".")
	for i := 0; i < len(v1Arr); i++ {
		if i >= len(v2Arr) {
			return 1
		}
		if v1Arr[i] > v2Arr[i] {
			return 1
		} else if v1Arr[i] < v2Arr[i] {
			return -1
		}
	}
	if len(v1Arr) < len(v2Arr) {
		return -1
	}
	return 0
}

// GetLastSevenDate 获取最近七天的日期，形如01/02
func GetLastSevenDate() []string {
	var dateSlice []string
	nowTime := time.Now()
	for i := 0; i < 7; i++ {
		dateSlice = append(dateSlice, nowTime.AddDate(0, 0, -i).Format("01/02"))
	}
	sort.Strings(dateSlice)
	return dateSlice
}

// GetThisWeekDate 获取本周截止当天的日期,形如01/02
func GetThisWeekDate() []string {
	var dateSlice []string
	nowTime := time.Now()
	weekday := nowTime.Weekday()
	for i := 0; i < int(weekday); i++ {
		dateSlice = append(dateSlice, nowTime.AddDate(0, 0, -i).Format("01/02"))
	}
	sort.Strings(dateSlice)
	return dateSlice
}

// GetThisWeekMondayToSunday 获取本周一到周天的日期，形如01/02
func GetThisWeekMondayToSunday() []string {
	var dateSlice []string
	nowTime := time.Now()
	weekday := nowTime.Weekday()
	for i := 0; i < 7; i++ {
		dateSlice = append(dateSlice, nowTime.AddDate(0, 0, i-int(weekday)+1).Format("01/02"))
	}
	sort.Strings(dateSlice)
	return dateSlice
}

// GetThisMonthStartToEnd 获取本月初到月末的日期，形如01/02
func GetThisMonthStartToEnd() []string {
	var dateSlice []string
	// 获取当前月份
	year, month, day := time.Now().Date()
	// 获取当前月份的总天数
	daysInMonth := time.Date(year, month+1, 0, 0, 0, 0, 0, time.Now().Location()).Day()
	for i := 0; i < daysInMonth; i++ {
		dateSlice = append(dateSlice, time.Now().AddDate(0, 0, i-day+1).Format("01/02"))
	}
	sort.Strings(dateSlice)
	return dateSlice
}

// GetThisMonthDate 获取本月截止当天的日期,形如01/02
func GetThisMonthDate() []string {
	var dateSlice []string
	nowTime := time.Now()
	for i := 0; i < nowTime.Day(); i++ {
		dateSlice = append(dateSlice, nowTime.AddDate(0, 0, -i).Format("01/02"))
	}
	sort.Strings(dateSlice)
	return dateSlice
}

// GetThisMonthStartToToday 获取本月初到月末的日期，形如01/02
func GetThisMonthStartToToday() []string {
	var dateSlice []string
	nowTime := time.Now()
	for i := 0; i < nowTime.Day(); i++ {
		dateSlice = append(dateSlice, nowTime.AddDate(0, 0, i-nowTime.Day()+1).Format("01/02"))
	}
	sort.Strings(dateSlice)
	return dateSlice
}

// GetDateSlice 根据起止毫秒时间戳,获取时间日期切片,形如01/02
func GetDateSlice(startTime, endTime int64) []string {
	var dateSlice []string
	startTime = startTime / 1000
	endTime = endTime / 1000
	startTime = startTime - startTime%86400
	endTime = endTime - endTime%86400
	for i := startTime; i <= endTime; i += 86400 {
		dateSlice = append(dateSlice, time.Unix(i, 0).Format("01/02"))
	}
	sort.Strings(dateSlice)
	return dateSlice
}

// Decimal2 保留两位小数
func Decimal2(value float64) float64 {
	return math.Trunc(value*1e2+0.5) * 1e-2
}

// Desensitize 脱敏仅展示前几位和后几位，其余用*代替
func Desensitize(str string, start, end int) string {
	if len(str) <= start+end {
		return str
	}
	return str[:start] + strings.Repeat("*", len(str)-start-end) + str[len(str)-end:]
}

// GetBirthdayAndSex 根据身份证号提取生日的年-月-日字符串及性别1男2女
func GetBirthdayAndSex(idCardNo string) (string, int) {
	if len(idCardNo) != 18 {
		return "", 0
	}
	sex := idCardNo[16:17]
	sexInt, _ := strconv.Atoi(sex)
	sexStr := 1
	if sexInt%2 == 0 {
		sexStr = 2
	}
	birthday := idCardNo[6:14]
	birthdayStr := fmt.Sprintf("%s-%s-%s", birthday[0:4], birthday[4:6], birthday[6:8])
	return birthdayStr, sexStr
}

// GetCdpId 获取cdpId
func GetCdpId(id int64) string {
	return fmt.Sprintf("8%0"+strconv.Itoa(10)+"d", id)
}

func IsTimeFormat(str string) bool {
	_, err := time.Parse("2006年01月", str)
	return err == nil
}

// GetValueByIndex 依次获取数组中的值
func GetValueByIndex(index int, arr []int64) int64 {
	return arr[index%len(arr)]
}

// GetMinKey 对map的值排序,取出最小的key
func GetMinKey(m map[int64]int64) int64 {
	var keys []int64
	for k := range m {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return m[keys[i]] < m[keys[j]]
	})
	return keys[0]
}

// GetSortKeys 对map的值排序,支持指定增序或倒序，返回排序后的keys
func GetSortKeys(m map[int64]int64, desc bool) []int64 {
	var keys []int64
	for k := range m {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		if desc {
			return m[keys[i]] > m[keys[j]]
		} else {
			return m[keys[i]] < m[keys[j]]
		}
	})

	sort.Slice(keys, func(i, j int) bool {
		if m[keys[i]] == m[keys[j]] {
			return keys[i] < keys[j]
		}
		return false
	})

	return keys
}

// IsEmail 校验邮箱格式
func IsEmail(email string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`).MatchString(email)
}

// Byte2Str []byte转string
func Byte2Str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// GenerateResetToken 生成邮件重置令牌
func GenerateResetToken() (string, error) {
	b := make([]byte, 32) // 生成32字节的随机数据
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// RemoveMultiChar 移除字符串内包含的多个字符
func RemoveMultiChar(str string, chars []string) string {
	for _, char := range chars {
		str = strings.ReplaceAll(str, char, "")
	}
	return str
}

// JsonStringToSlice json字符串转[]string
func JsonStringToSlice(jsonStr string) ([]string, error) {
	var strSlice []string
	err := json.Unmarshal([]byte(jsonStr), &strSlice)
	if err != nil {
		return nil, err
	}
	return strSlice, nil
}

// AnyToString 把任意类型转字符串
func AnyToString(v interface{}) string {
	return fmt.Sprintf("%v", v)
}

// ExtractFileNameFromUrl 从URL中提取文件名
func ExtractFileNameFromUrl(fileUrl string) (string, error) {
	parsedURL, err := url.Parse(fileUrl)
	if err != nil {
		return "", err
	}
	path := parsedURL.Path
	fileName := path[strings.LastIndex(path, "/")+1:]
	return fileName, nil
}

// ReadRemoteFile 根据远程 URL 直接读取文件内容并返回为字符串
func ReadRemoteFile(url string) (string, error) {
	// 发送 GET 请求
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("无法访问 URL: %w", err)
	}
	defer resp.Body.Close()

	// 检查 HTTP 响应状态码
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("请求失败，状态码: %d", resp.StatusCode)
	}

	// 读取响应体
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("无法读取响应内容: %w", err)
	}

	// 返回内容作为字符串
	return string(content), nil
}

// Substring 截取字符串,如果英文就大写
func Substring(str string, start int, end int) string {
	rs := []rune(str)
	if start < 0 || start >= len(rs) {
		return ""
	}
	if end < 0 || end > len(rs) {
		return ""
	}

	return strings.ToUpper(string(rs[start:end]))
}

func LastCharToUpper(str string) string {
	if str == "" {
		return ""
	}
	rs := []rune(str)
	if len(rs) == 1 {
		return strings.ToUpper(string(rs))
	}
	return strings.ToUpper(string(rs[len(rs)-1:]))
}

// IsPhoneNumber 判断字符串是不是手机号,首位可能是+号
func IsPhoneNumber(str string) bool {
	return regexp.MustCompile(`^\+?\d{1,11}$`).MatchString(str)
}

// ConvertTimestampToTargetTimezone 将指定时区的时间戳转为目标时区时间戳
func ConvertTimestampToTargetTimezone(timestamp int64, fromTimezone, toTimezone string) int64 {
	if timestamp <= 0 {
		return 0
	}
	if len(fmt.Sprintf("%d", timestamp)) == 10 {
		timestamp = timestamp * 1000
	} else if len(fmt.Sprintf("%d", timestamp)) != 13 {
		return 0
	}
	if fromTimezone == toTimezone {
		return timestamp
	}
	fromLocation, err := time.LoadLocation(fromTimezone)
	if err != nil {
		return 0
	}
	toLocation, err := time.LoadLocation(toTimezone)
	if err != nil {
		log.Println(err.Error())
		return 0
	}
	fromTime := time.UnixMilli(timestamp).In(fromLocation)
	toTime := fromTime.In(toLocation)
	return toTime.UnixMilli()
}

// SecondToHourMinute 秒转小时及分钟
func SecondToHourMinute(second int64) string {
	var str string
	hour := second / 3600
	if hour > 0 {
		str += fmt.Sprintf("%dh", hour)
	}
	minute := (second - hour*3600) / 60
	if minute > 0 {
		str += fmt.Sprintf("%dmin", minute)
	}
	return str
}

// SecondToHourMinuteSecond 秒转x时x分x秒
func SecondToHourMinuteSecond(seconds int64) string {
	var hours, minutes, secs int64

	hours = seconds / 3600
	seconds %= 3600
	minutes = seconds / 60
	secs = seconds % 60

	if hours > 0 {
		return fmt.Sprintf("%d小时%d分%d秒", hours, minutes, secs)
	} else if minutes > 0 {
		return fmt.Sprintf("%d分%d秒", minutes, secs)
	} else if secs > 0 {
		return fmt.Sprintf("%d秒", secs)
	} else {
		return ""
	}
}

// StrToZero 字符串为空时转为0
func StrToZero(str string) string {
	if str == "" {
		return "0"
	}
	return str
}

// SecondToHourMinuteSecond2 秒转xHxMinxS
func SecondToHourMinuteSecond2(seconds int64) string {
	var hours, minutes, secs int64

	hours = seconds / 3600
	seconds %= 3600
	minutes = seconds / 60
	secs = seconds % 60

	if hours > 0 {
		return fmt.Sprintf("%dh%dmin%ds", hours, minutes, secs)
	} else if minutes > 0 {
		return fmt.Sprintf("%dmin%ds", minutes, secs)
	} else if secs > 0 {
		return fmt.Sprintf("%ds", secs)
	} else {
		return ""
	}
}

// ConvertBytesToReadable 将字节大小转换为可读的字符串格式
func ConvertBytesToReadable(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 1
	units := []string{"B", "KB", "MB", "GB"}
	for _, u := range units {
		if float64(bytes) < math.Pow(float64(unit), float64(exp)) {
			return fmt.Sprintf("%.2f %s", float64(bytes)/math.Pow(float64(div), float64(exp-1)), u)
		}
		exp++
	}
	return fmt.Sprintf("%.2f %s", float64(bytes)/math.Pow(float64(div), float64(exp-1)), units[len(units)-1])
}

// FormatSecondsToTime 将秒数格式化为 00:00:00 格式
func FormatSecondsToTime(seconds int) string {
	// 将秒数转为 time.Duration 类型
	duration := time.Duration(seconds) * time.Second
	// 格式化为 00:00:00 格式
	return fmt.Sprintf("%02d:%02d:%02d", int(duration.Hours()), int(duration.Minutes())%60, int(duration.Seconds())%60)
}

// Encrypt 加密函数，接收字符串、密钥并返回加密后的密文
func Encrypt(plaintext string, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(craned.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plaintext))

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

// Decrypt 解密函数，将密文解密为原始字符串
func Decrypt(ciphertext string, key []byte) (string, error) {
	ciphertextBytes, err := base64.URLEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(ciphertextBytes) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	iv := ciphertextBytes[:aes.BlockSize]
	ciphertextBytes = ciphertextBytes[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertextBytes, ciphertextBytes)

	return string(ciphertextBytes), nil
}

// 将 map 序列化为字符串后加密
func encryptMap(params map[string]interface{}, key []byte) (string, error) {
	jsonData, err := json.Marshal(params) // 处理 interface{} 类型
	if err != nil {
		return "", err
	}
	return Encrypt(string(jsonData), key)
}

// 解密后将字符串反序列化为 map
func decryptMap(ciphertext string, key []byte) (map[string]interface{}, error) {
	decrypted, err := Decrypt(ciphertext, key)
	if err != nil {
		return nil, err
	}

	var params map[string]interface{}
	err = json.Unmarshal([]byte(decrypted), &params) // 反序列化为 map[string]interface{}
	return params, err
}

// FormatMessageTime
// 格式化时间戳，
// 当天的消息，只显示具体时/分
// 昨天的消息，显示 昨天 时/分
// 前天的消息，显示 星期几 时/分
// 再之前的消息，显示 X月X日 时/分
func FormatMessageTime(timestamp int64, timeZone string) string {
	curTime := time.Now()

	timestamp = ConvertTimestampToTargetTimezone(timestamp, time.Local.String(), timeZone)
	msgTime := time.UnixMilli(timestamp)

	if msgTime.Format("2006-01-02") == curTime.Format("2006-01-02") {
		return fmt.Sprintf("%02d:%02d", msgTime.Hour(), msgTime.Minute())
	} else if msgTime.AddDate(0, 0, 1).Format("2006-01-02") == curTime.Format("2006-01-02") {
		return "昨天 " + fmt.Sprintf("%02d:%02d", msgTime.Hour(), msgTime.Minute())
	} else if msgTime.AddDate(0, 0, 2).Format("2006-01-02") == curTime.Format("2006-01-02") {
		daysOfWeek := []string{"星期日", "星期一", "星期二", "星期三", "星期四", "星期五", "星期六"}
		return daysOfWeek[msgTime.Weekday()] + " " + fmt.Sprintf("%02d:%02d", msgTime.Hour(), msgTime.Minute())
	} else {
		return fmt.Sprintf("%02d月%02d日 %02d:%02d", msgTime.Month(), msgTime.Day(), msgTime.Hour(), msgTime.Minute())
	}
}

// CalculatePercentage 计算两个数的百分比，并在分子小于 1 时返回 1
func CalculatePercentage(numerator float64, denominator float64) int64 {
	if denominator == 0 {
		fmt.Println("Warning: Division by zero is undefined.")
		return 1 // 避免除零错误，返回1
	}

	// 计算百分比
	percentage := (numerator / denominator) * 100

	if percentage < 1 {
		return 1
	}

	return int64(percentage)
}
