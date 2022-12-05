// Copyright 2021 dudaodong@gmail.com. All rights reserved.
// Use of this source code is governed by MIT license

// Package convertor implements some functions to convert data.
package convertor

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
)

// ToBool convert string to a boolean
func ToBool(s string) (bool, error) {
	return strconv.ParseBool(s)
}

// ToBytes convert interface to bytes
func ToBytes(value any) ([]byte, error) {
	v := reflect.ValueOf(value)

	switch value.(type) {
	case int, int8, int16, int32, int64:
		number := v.Int()
		buf := bytes.NewBuffer([]byte{})
		buf.Reset()
		err := binary.Write(buf, binary.BigEndian, number)
		return buf.Bytes(), err
	case uint, uint8, uint16, uint32, uint64:
		number := v.Uint()
		buf := bytes.NewBuffer([]byte{})
		buf.Reset()
		err := binary.Write(buf, binary.BigEndian, number)
		return buf.Bytes(), err
	case float32:
		number := float32(v.Float())
		bits := math.Float32bits(number)
		bytes := make([]byte, 4)
		binary.BigEndian.PutUint32(bytes, bits)
		return bytes, nil
	case float64:
		number := v.Float()
		bits := math.Float64bits(number)
		bytes := make([]byte, 8)
		binary.BigEndian.PutUint64(bytes, bits)
		return bytes, nil
	case bool:
		return strconv.AppendBool([]byte{}, v.Bool()), nil
	case string:
		return []byte(v.String()), nil
	case []byte:
		return v.Bytes(), nil
	case decimal.Decimal:
		return value.(decimal.Decimal).MarshalJSON()
	default:
		newValue, err := json.Marshal(value)
		return newValue, err
	}
}

// ToChar convert string to char slice
func ToChar(s string) []string {
	c := make([]string, 0)
	if len(s) == 0 {
		c = append(c, "")
	}
	for _, v := range s {
		c = append(c, string(v))
	}
	return c
}

// ToChannel convert a array of elements to a read-only channels
func ToChannel[T any](array []T) <-chan T {
	ch := make(chan T)

	go func() {
		for _, item := range array {
			ch <- item
		}
		close(ch)
	}()

	return ch
}

// ToString convert value to string
// for number, string, []byte, will convert to string
// for other type (slice, map, array, struct) will call json.Marshal
func ToString(value any) string {
	if value == nil {
		return ""
	}

	switch value.(type) {
	case float32:
		return strconv.FormatFloat(float64(value.(float32)), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(value.(float64), 'f', -1, 64)
	case int:
		return strconv.FormatInt(int64(value.(int)), 10)
	case int8:
		return strconv.FormatInt(int64(value.(int8)), 10)
	case int16:
		return strconv.FormatInt(int64(value.(int16)), 10)
	case int32:
		return strconv.FormatInt(int64(value.(int32)), 10)
	case int64:
		return strconv.FormatInt(value.(int64), 10)
	case uint:
		return strconv.FormatUint(uint64(value.(uint)), 10)
	case uint8:
		return strconv.FormatUint(uint64(value.(uint8)), 10)
	case uint16:
		return strconv.FormatUint(uint64(value.(uint16)), 10)
	case uint32:
		return strconv.FormatUint(uint64(value.(uint32)), 10)
	case uint64:
		return strconv.FormatUint(value.(uint64), 10)
	case string:
		return value.(string)
	case []byte:
		return string(value.([]byte))
	case decimal.Decimal:
		return value.(decimal.Decimal).String()
	default:
		newValue, _ := json.Marshal(value)
		return string(newValue)

		// todo: maybe we should't supprt other type conversion
		// v := reflect.ValueOf(value)
		// log.Panicf("Unsupported data type: %s ", v.String())
		// return ""
	}
}

// ToJson convert value to a valid json string
func ToJson(value any) (string, error) {
	result, err := json.Marshal(value)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

// ToFloat convert value to a float64, if input is not a float return 0.0 and error
func ToFloat(value any) (float64, error) {
	v := reflect.ValueOf(value)

	result := 0.0
	err := fmt.Errorf("ToInt: unvalid interface type %T", value)
	switch value.(type) {
	case int, int8, int16, int32, int64:
		result = float64(v.Int())
		return result, nil
	case uint, uint8, uint16, uint32, uint64:
		result = float64(v.Uint())
		return result, nil
	case float32, float64:
		result = v.Float()
		return result, nil
	case string:
		result, err = strconv.ParseFloat(v.String(), 64)
		if err != nil {
			result = 0.0
		}
		return result, err
	case decimal.Decimal:
		return value.(decimal.Decimal).InexactFloat64(), nil
	default:
		return result, err
	}
}

// ToInt convert value to a int64, if input is not a numeric format return 0 and error
func ToInt(value any) (int64, error) {
	v := reflect.ValueOf(value)

	var result int64
	err := fmt.Errorf("ToInt: invalid interface type %T", value)
	switch value.(type) {
	case int, int8, int16, int32, int64:
		result = v.Int()
		return result, nil
	case uint, uint8, uint16, uint32, uint64:
		result = int64(v.Uint())
		return result, nil
	case float32, float64:
		result = int64(v.Float())
		return result, nil
	case string:
		result, err = strconv.ParseInt(v.String(), 0, 64)
		if err != nil {
			result = 0
		}
		return result, err
	case decimal.Decimal:
		return value.(decimal.Decimal).IntPart(), nil
	default:
		return result, err
	}
}

// ToPointer returns a pointer to this value
func ToPointer[T any](value T) *T {
	return &value
}

// ToMap convert a slice or an array of structs to a map based on iteratee function
func ToMap[T any, K comparable, V any](array []T, iteratee func(T) (K, V)) map[K]V {
	result := make(map[K]V, len(array))
	for _, item := range array {
		k, v := iteratee(item)
		result[k] = v
	}

	return result
}

// StructToMap convert struct to map, only convert exported struct field
// map key is specified same as struct field tag `json` value
func StructToMap(value any) (map[string]any, error) {
	v := reflect.ValueOf(value)
	t := reflect.TypeOf(value)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("data type %T not support, shuld be struct or pointer to struct", value)
	}

	result := make(map[string]any)

	fieldNum := t.NumField()
	pattern := `^[A-Z]`
	regex := regexp.MustCompile(pattern)
	for i := 0; i < fieldNum; i++ {
		name := t.Field(i).Name
		tag := t.Field(i).Tag.Get("json")
		if regex.MatchString(name) && tag != "" {
			// result[name] = v.Field(i).Interface()
			result[tag] = v.Field(i).Interface()
		}
	}

	return result, nil
}

// MapToSlice convert a map to a slice based on iteratee function
func MapToSlice[T any, K comparable, V any](aMap map[K]V, iteratee func(K, V) T) []T {
	result := make([]T, 0, len(aMap))

	for k, v := range aMap {
		result = append(result, iteratee(k, v))
	}

	return result
}

// ColorHexToRGB convert hex color to rgb color
func ColorHexToRGB(colorHex string) (red, green, blue int) {
	colorHex = strings.TrimPrefix(colorHex, "#")
	color64, err := strconv.ParseInt(colorHex, 16, 32)
	if err != nil {
		return
	}
	color := int(color64)
	return color >> 16, (color & 0x00FF00) >> 8, color & 0x0000FF
}

// ColorRGBToHex convert rgb color to hex color
func ColorRGBToHex(red, green, blue int) string {
	r := strconv.FormatInt(int64(red), 16)
	g := strconv.FormatInt(int64(green), 16)
	b := strconv.FormatInt(int64(blue), 16)

	if len(r) == 1 {
		r = "0" + r
	}
	if len(g) == 1 {
		g = "0" + g
	}
	if len(b) == 1 {
		b = "0" + b
	}

	return "#" + r + g + b
}

// EncodeByte encode data to byte
func EncodeByte(data any) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(data)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// DecodeByte decode byte data to target object
func DecodeByte(data []byte, target any) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(target)
}

// ToDecimal Converts any value to the decimal type and returns 0 if the conversion fails
func ToDecimal(value any) decimal.Decimal {
	v := reflect.ValueOf(value)
	var result decimal.Decimal
	switch value.(type) {
	case int, int8, int16, int32, int64:
		result = decimal.NewFromInt(v.Int())
		return result
	case uint, uint8, uint16, uint32, uint64:
		result = decimal.NewFromInt(int64(v.Uint()))
		return result
	case float32:
		result = decimal.NewFromFloat32(value.(float32))
		return result
	case float64:
		result = decimal.NewFromFloat(v.Float())
		return result
	case string:
		result, err := decimal.NewFromString(v.String())
		if err != nil {
			return decimal.Decimal{}
		}
		return result
	case decimal.Decimal:
		return value.(decimal.Decimal)
	default:
		return result
	}
}
