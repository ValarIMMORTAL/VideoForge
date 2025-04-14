package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

// ConvertStruct 将源结构体 src 转换为目标结构体 dst 的指针
// 如果字段名称或类型不匹配（且无法自动转换），则返回错误
func ConvertByJSONTag(src, dst interface{}) error {
	dstValue := reflect.ValueOf(dst)
	// 检查 dst 是否是指向结构体的指针
	if dstValue.Kind() != reflect.Ptr || dstValue.Elem().Kind() != reflect.Struct {
		return errors.New("dst must be a pointer to struct")
	}

	// 检查 src 是否是结构体或指针
	srcValue := reflect.ValueOf(src)
	if srcValue.Kind() == reflect.Ptr {
		srcValue = srcValue.Elem()
	}
	if srcValue.Kind() != reflect.Struct {
		return errors.New("src must be a struct or pointer to struct")
	}

	// 将 src 转为 JSON
	jsonData, err := json.Marshal(src)
	if err != nil {
		return fmt.Errorf("failed to marshal src to JSON: %v", err)
	}

	// 将 JSON 解析到 dst
	if err := json.Unmarshal(jsonData, dst); err != nil {
		return fmt.Errorf("failed to unmarshal JSON to dst: %v", err)
	}

	return nil
}
