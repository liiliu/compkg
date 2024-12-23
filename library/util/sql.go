package util

import (
	"database/sql"
	"fmt"
	"gorm.io/gorm"
)

// ScanRows2map 取出不需要结构体的数据
func ScanRows2map(rows *sql.Rows) []map[string]string {
	res := make([]map[string]string, 0)            //  定义结果 map
	colTypes, _ := rows.ColumnTypes()              // 获取列信息
	rowParam := make([]interface{}, len(colTypes)) // 参数数组
	rowValue := make([]interface{}, len(colTypes)) // 接收数据的一行值

	// 初始化 rowParam 和 rowValue
	for i := range colTypes {
		rowValue[i] = new(sql.RawBytes) // 使用通用的 RawBytes 类型处理
		rowParam[i] = &rowValue[i]
	}

	// 遍历行
	for rows.Next() {
		rows.Scan(rowParam...) // 扫描行数据
		record := make(map[string]string)
		for i, colType := range colTypes {
			if rowValue[i] == nil {
				record[colType.Name()] = ""
			} else {
				record[colType.Name()] = fmt.Sprintf("%v", rowValue[i])
			}
		}
		res = append(res, record)
	}
	return res
}

func QueryResultAsMaps(db *gorm.DB, sql string) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	rows, err := db.Raw(sql).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	values := make([]interface{}, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err := rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}

		record := make(map[string]interface{})
		for k, value := range values {
			record[columns[k]] = value
		}
		results = append(results, record)
	}

	return results, rows.Err()
}
