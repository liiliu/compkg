package loghub

import (
	sls "github.com/aliyun/aliyun-log-go-sdk"
	"log"
	"time"
	"weihu_server/library/config"
)

var logClient sls.ClientInterface

func Initial() {
	// 日志服务的服务入口
	endpoint := config.GetString("logHub.endPoint")

	// AccessKey ID和AccessKey Secret。
	accessKeyId := config.GetString("logHub.accessKeyId")
	accessKeySecret := config.GetString("logHub.accessKeySecret")

	// RAM用户角色的临时安全令牌。此处取值为空，表示不使用临时安全令牌。
	securityToken := ""
	// 创建日志服务Client。
	logClient = sls.CreateNormalInterface(endpoint, accessKeyId, accessKeySecret, securityToken)
}

// CreateProject 创建项目
func CreateProject(projectName, description string) {
	// 检查项目是否存在
	projectExist, err := logClient.CheckProjectExist(projectName)
	if err != nil {
		log.Fatalf("创建项目 : %s 失败 %v\n", projectName, err)
		return
	}

	if projectExist {
		log.Printf("项目 : %s 已经创建或在阿里云范围内具有全局名称冲突\n", projectName)
		return
	}

	project, err := logClient.CreateProject(projectName, description)
	if err != nil {
		if e, ok := err.(*sls.Error); ok && e.Code == "ProjectAlreadyExist" {
			log.Printf("项目 : %s 已经创建或在阿里云范围内具有全局名称冲突\n", projectName)
		} else {
			log.Fatalf("创建项目 : %s 失败 %v\n", projectName, err)
		}
	} else {
		log.Printf("项目 : %s 创建成功\n", project.Name)
		time.Sleep(60 * time.Second)
	}
}

// CreateStore 创建存储
func CreateStore(projectName, logStoreName string) {
	// 检查存储是否存在
	storeExist, err := logClient.CheckLogstoreExist(projectName, logStoreName)
	if err != nil {
		log.Fatalf("创建存储 : %s 失败 %v\n", projectName, err)
		return
	}

	if storeExist {
		log.Printf("存储 : %s 已存在\n", logStoreName)
		return
	}

	logStore := new(sls.LogStore)
	logStore.Name = logStoreName
	err = logClient.CreateLogStoreV2(projectName, logStore)
	if err != nil {
		if e, ok := err.(*sls.Error); ok && e.Code == "LogStoreAlreadyExist" {
			log.Printf("存储 : %s 已存在\n", logStoreName)
		} else {
			log.Fatalf("创建存储 : %s 失败 %v\n", logStoreName, err)
		}
	} else {
		log.Printf("创建存储 : %v 成功\n", logStoreName)
		time.Sleep(60 * time.Second)
	}
}

// CreateIndex 创建索引
func CreateIndex(projectName, logStoreName string, longCols []string, doubleCols []string, textCols []string) {
	indexMap := make(map[string]sls.IndexKey)
	for _, col := range longCols {
		indexMap[col] = sls.IndexKey{
			Token:         []string{" "},
			CaseSensitive: false,
			Type:          "long",
		}
	}
	for _, col := range doubleCols {
		indexMap[col] = sls.IndexKey{
			Token:         []string{" "},
			CaseSensitive: false,
			Type:          "double",
		}
	}
	for _, col := range textCols {
		indexMap[col] = sls.IndexKey{
			Token:         []string{",", ":", " "},
			CaseSensitive: false,
			Type:          "text",
		}
	}
	// 为LogStore创建索引。
	index := sls.Index{
		// 字段索引。
		Keys: indexMap,
		// 全文索引。
		Line: &sls.IndexLine{
			Token:         []string{",", ":", " "},
			CaseSensitive: false,
			IncludeKeys:   []string{},
			ExcludeKeys:   []string{},
		},
	}

	err := logClient.CreateIndex(projectName, logStoreName, index)
	if err != nil {
		if e, ok := err.(*sls.Error); ok && e.Code == "IndexAlreadyExist" {
			log.Printf("索引 : %+v 已存在\n", index)
		} else {
			log.Fatalf("创建索引失败 %v\n", err)
		}
	} else {
		log.Println("创建索引成功")
		time.Sleep(60 * time.Second)
	}
}

// PutLogs 推送日志
func PutLogs(projectName string, logStoreName string, logGroup *sls.LogGroup) {
	err := logClient.PutLogs(projectName, logStoreName, logGroup)
	if err != nil {
		log.Fatalf("推送日志失败 %v", err)
	}
	//log.Println("推送日志成功")
}
