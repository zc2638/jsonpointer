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

type Parser struct {
	rv *reflect.Value
}

func NewParser(data interface{}) (*Parser, error) {
	rv := reflect.ValueOf(data)
	switch rv.Type().Kind() {
	case reflect.Ptr:
	case reflect.Slice:
	case reflect.Map:
	case reflect.Struct:
	default:
		return nil, errors.New("data type not support")
	}
	return &Parser{rv: &rv}, nil
}

func (p *Parser) Check(ref string) bool {
	if _, err := p.Get(ref); err != nil {
		return false
	}
	return true
}

func (p *Parser) Get(ref string) (interface{}, error) {
	return getPointerData(*p.rv, refPaths(ref))
}

func (p *Parser) Batch(refs []string) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	for _, ref := range refs {
		refData, err := getPointerData(*p.rv, refPaths(ref))
		if err != nil {
			// TODO 增加自定义error结构便于获取ref值
			return nil, fmt.Errorf("path: %s, error: %s", ref, err)
		}
		result[ref] = refData
	}
	return result, nil
}

func refPaths(ref string) []string {
	key := strings.TrimPrefix(ref, "/")
	return strings.Split(key, "/")
}

func transferPointer(key string) string {
	key = strings.ReplaceAll(key, "~1", "/")
	return strings.ReplaceAll(key, "~0", "~")
}

func getPointerData(rv reflect.Value, refPaths []string) (interface{}, error) {
	if len(refPaths) == 0 {
		if !rv.CanInterface() {
			return nil, errors.New("invalid value")
		}
		return rv.Interface(), nil
	}
	key := transferPointer(refPaths[0])
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
