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

type Error struct {
	Ref     string      `json:"ref"`
	Message string      `json:"message"`
	Default interface{} `json:"default"`
}

func (e *Error) Error() string {
	var str string
	if e.Ref != "" {
		str += "path: " + e.Ref + ", "
	}
	return str + "error: " + e.Message
}

func (e *Error) WithRef(ref string) *Error {
	e.Ref = ref
	return e
}

func (e *Error) WithDefault(value interface{}) *Error {
	e.Default = value
	return e
}

func NewError(message string) *Error {
	return &Error{
		Message: message,
	}
}
