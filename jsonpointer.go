// Package jsonpointer

// Copyright © 2021 zc2638 <zc2638@qq.com>.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package jsonpointer

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// GetPointersData 根据json-pointer组获取结构体下的结构指针
func GetPointersData(data interface{}, refs []string) (map[string]interface{}, error) {
	rv := reflect.ValueOf(data)
	return getPointersData(rv, refs)
}

func getPointersData(rv reflect.Value, refs []string) (map[string]interface{}, error) {
	switch rv.Type().Kind() {
	case reflect.Ptr:
		return getPointersData(rv.Elem(), refs)
	case reflect.Slice:
	case reflect.Map:
	case reflect.Struct:
	default:
		return nil, errors.New("data type not support")
	}
	result := make(map[string]interface{})
	for _, ref := range refs {
		key := strings.TrimPrefix(ref, "/")
		refPaths := strings.Split(key, "/")
		refData, err := getPointerData(rv, refPaths)
		if err != nil {
			return nil, fmt.Errorf("path: %s, error: %s", ref, err)
		}
		result[ref] = refData
	}
	return result, nil
}

func getPointerData(rv reflect.Value, refPaths []string) (interface{}, error) {
	if len(refPaths) == 0 {
		if rv.Kind() == reflect.Invalid {
			return nil, nil
		}
		return rv.Interface(), nil
	}
	key := refPaths[0]
	// 处理 ~0 和 ~1 的转换
	key = strings.ReplaceAll(key, "~1", "/")
	key = strings.ReplaceAll(key, "~0", "~")
	switch rv.Type().Kind() {
	case reflect.Ptr:
		if rv.IsNil() {
			return nil, errors.New("ptr value is nil")
		}
		return getPointerData(rv.Elem(), refPaths)
	case reflect.Slice, reflect.Array:
		if rv.IsNil() {
			return nil, errors.New("array/slice value is nil")
		}
		i, err := strconv.Atoi(key)
		if err != nil {
			return nil, errors.New("not a valid array index")
		}
		if i >= rv.Len() {
			return nil, errors.New("index out of range")
		}
		return getPointerData(rv.Index(i), refPaths[1:])
	case reflect.Map:
		if rv.IsNil() {
			return nil, errors.New("map value is nil")
		}
		value := rv.MapIndex(reflect.ValueOf(key))
		if value.Kind() == reflect.Invalid {
			valT := rv.Type().Elem()
			value = reflect.New(valT).Elem()
		}
		return getPointerData(value, refPaths[1:])
	case reflect.Struct:
		for i := 0; i < rv.NumField(); i++ {
			field := rv.Type().Field(i)
			current, ok := field.Tag.Lookup("json")
			if !ok {
				continue
			}
			if current != key {
				continue
			}
			return getPointerData(rv.Field(i), refPaths[1:])
		}
		return nil, errors.New("field mismatch")
	default:
		return nil, errors.New("type mismatch")
	}
}
