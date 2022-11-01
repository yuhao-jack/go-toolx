package fun

import (
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"
	"strconv"
)

// IsBlank
//
//	@Description: 判断入参是否是零值
//	@Author yuhao <yuhao@mini1.cn>
//	@Data 2022-11-01 09:56:43
//	@param i
//	@return bool
func IsBlank(i interface{}) bool {
	value := reflect.ValueOf(i)
	switch value.Kind() {
	case reflect.String:
		return value.Len() == 0
	case reflect.Bool:
		return !value.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return value.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return value.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return value.IsNil()
	}
	return reflect.DeepEqual(value.Interface(), reflect.Zero(value.Type()).Interface())
}

// IfOr
//
//	@Description: 模拟三目运算
//	@Author yuhao <yuhao@mini1.cn>
//	@Data 2022-11-01 09:58:24
//	@param condition
//	@param a
//	@param b
//	@return T
func IfOr[T any](condition bool, a, b T) T {
	if condition {
		return a
	} else {
		return b
	}
}

// StrVal
//
//	 @Description: 	获取变量的字符串值
//						浮点型 3.0将会转换成字符串3, "3"
//						非数值或字符类型的变量将会被转换成JSON格式字符串
//	 @Author yuhao <yuhao@mini1.cn>
//	 @Data 2022-11-01 10:00:38
//	 @param value
//	 @return string
func StrVal(value interface{}) string {
	var key string
	if value == nil {
		return key
	}
	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}
	return key
}

// CheckRecover
//
//	@Description:  检测 recover并打印错误 stack
//	@Author yuhao <yuhao@mini1.cn>
//	@Data 2022-11-01 10:03:55
//	@return string
func CheckRecover() string {
	err := recover()
	if err != nil {
		buf := make([]byte, 4096)
		buf = buf[:runtime.Stack(buf, false)]
		return fmt.Sprintf("err:%v,\n%s", err, buf)
	}
	return ""
}
