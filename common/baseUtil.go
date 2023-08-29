package common

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"time"
	"unicode/utf8"
)

// 结构体转一维数组
func StructToArray(data interface{}) []interface{} {
	value := reflect.ValueOf(data)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	if value.Kind() != reflect.Struct {
		return nil
	}

	var result []interface{}
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		result = append(result, field.Interface())
	}

	return result
}

// 去重
func RemoveDuplicate(old []string) []string {
	result := []string{}
	temp := map[string]struct{}{}
	for _, item := range old {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

// 获取当前目录
func GetCurrentPath() (dir string, err error) {
	return os.Getwd()
}

// 创建结果目录
func CreateResultDir() string {
	CurrentRunPath, _ = GetCurrentPath()
	ResultPath = fmt.Sprintf("%s/result/", CurrentRunPath)
	exists := FolderExists(ResultPath)
	if !exists {
		DirCreate(ResultPath)
	}
	return ResultPath
}

// 获取当前时间字符串
func GetCurrentTimeString() string {
	return time.Now().Format("2006_01_02_15_04_05")
}

// 获取当前时间字符串 时间格式
func GetCurrentTimeStringTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// GetTempPathFileName 获取临时文件名字
func GetTempPathFileName() (pathFileName string) {
	return filepath.Join(os.TempDir(), fmt.Sprintf("%s.tmp", GetRandomString2(16)))
}

// GetRandomString2 随机字符
func GetRandomString2(n int) string {
	randBytes := make([]byte, n/2)
	rand.Read(randBytes)
	return fmt.Sprintf("%x", randBytes)
}

// 扫描Target格式处理
func GetScanTarget(file string) (targets []string) {
	f, err := ReadFile(file)
	if err == nil {
		for _, target := range f {
			target = TargetStrip(target)
			targets = append(targets, target)
		}
	}
	return targets
}

// 自定义类型  StringSlice
type StringSlice []string

// 实现 Value 方法，用于将自定义类型转换为数据库支持的原始类型
func (s StringSlice) Value() (driver.Value, error) {
	jsonData, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	return string(jsonData), nil
}

// 实现 Scan 方法，将数据库中的原始类型转换为自定义类型
func (s *StringSlice) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal StringSlice value")
	}
	return json.Unmarshal(b, &s)
}

// 是否还有无效字符，用于对中文检测
func HasInvalidChars(s string) bool {
	for _, r := range s {
		if r == utf8.RuneError {
			return true
		}
	}
	return false
}

// 通用函数，将结构体切片转换为二维数组切片
func ConvertStructSliceTo2D(slice interface{}) ([][]interface{}, error) {
	sliceValue := reflect.ValueOf(slice)
	// 检查输入参数是否为切片类型
	if sliceValue.Kind() != reflect.Slice {
		return nil, fmt.Errorf("Input is not a slice")
	}
	// 获取结构体类型和字段数量
	elemType := sliceValue.Type().Elem()
	numFields := elemType.NumField()
	// 定义二维数组切片
	var result [][]interface{}
	// 遍历结构体切片并将元素添加到二维数组切片中
	for i := 0; i < sliceValue.Len(); i++ {
		// 定义一维数组切片
		row := make([]interface{}, numFields)

		// 遍历结构体字段并将值添加到一维数组切片中
		for j := 0; j < numFields; j++ {
			row[j] = sliceValue.Index(i).Field(j).Interface()
		}
		result = append(result, row)
	}
	return result, nil
}

// 通用函数，将字符切片转换为二维数组切片
func ConvertStringSliceTo2D(slice []string) [][]string {
	// 获取切片长度
	numRows := len(slice)
	// 定义二维数组切片
	result := make([][]string, numRows)
	// 遍历切片并将元素转换为一维数组切片
	for i, value := range slice {
		// 将每个元素转换为一维数组切片
		result[i] = []string{value}
	}
	return result
}

// 通用函数，将[][]string转换为[][]interface{}
func ConvertStringSlice2DToInterface(slice [][]string) [][]interface{} {
	// 获取二维切片的行数和列数
	numRows := len(slice)
	numCols := len(slice[0])
	// 定义目标二维切片
	result := make([][]interface{}, numRows)
	// 遍历二维切片并将元素转换为interface{}
	for i, row := range slice {
		// 创建一维切片
		result[i] = make([]interface{}, numCols)
		for j, value := range row {
			// 将string转换为interface{}
			result[i][j] = value
		}
	}
	return result
}

// 通用函数，将[][]interface{}转换为[][]string
func ConvertInterfaceSlice2DToStringSlice(slice [][]interface{}) [][]string {
	// 获取二维切片的行数和列数
	numRows := len(slice)
	numCols := len(slice[0])
	// 定义目标二维切片
	result := make([][]string, numRows)
	// 遍历二维切片并将元素转换为string
	for i, row := range slice {
		// 创建一维切片
		result[i] = make([]string, numCols)
		for j, value := range row {
			// 将interface{}转换为string
			//result[i][j] = value.(string)
			result[i][j] = ToString(value)
		}
	}
	return result
}

// 将任何类型的值转换为字符串
func ToString(value interface{}) string {
	return fmt.Sprintf("%v", value)
}

// 输出n个回车
func PrintNewLines(count int) {
	for i := 0; i < count; i++ {
		fmt.Println()
	}
}

func BannerInit() {
	fmt.Println(banner)
	fmt.Println(ffpAbout)
}
