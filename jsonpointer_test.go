// Package jsonpointer
// Created by zc on 2021/8/27.
package jsonpointer

import (
	"reflect"
	"testing"
)

type Test struct {
	Name     string  `json:"name"`
	Children []Child `json:"children"`
}

type Child struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	Data *Data  `json:"data"`
}

type Data struct {
	Message string            `json:"message"`
	Labels  map[string]string `json:"labels"`
	Line    *int              `json:"line"`
}

var (
	line = 5
	data = Test{
		Name: "张三",
		Children: []Child{
			{
				Name: "李四",
				Age:  20,
				Data: &Data{
					Message: "a12312",
					Labels: map[string]string{
						"a/b": "data",
					},
					Line: &line,
				},
			},
		},
	}
)

func TestGetPointersData(t *testing.T) {
	type args struct {
		data interface{}
		refs []string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				data: data,
				refs: []string{
					"/name",
					"/children/0",
					"/children/0/age",
					"/children/0/data",
					"/children/0/data/line",
					"/children/0/data/labels/a~1b",
				},
			},
			want: map[string]interface{}{
				"/name":                        data.Name,
				"/children/0":                  data.Children[0],
				"/children/0/age":              data.Children[0].Age,
				"/children/0/data":             data.Children[0].Data,
				"/children/0/data/line":        data.Children[0].Data.Line,
				"/children/0/data/labels/a~1b": data.Children[0].Data.Labels["a/b"],
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetPointersData(tt.args.data, tt.args.refs)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPointersData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPointersData() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getPointersData(t *testing.T) {
	type args struct {
		rv   reflect.Value
		refs []string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				rv: reflect.ValueOf(data),
				refs: []string{
					"/name",
					"/children/0",
					"/children/0/age",
					"/children/0/data",
					"/children/0/data/line",
					"/children/0/data/labels/a~1b",
				},
			},
			want: map[string]interface{}{
				"/name":                        data.Name,
				"/children/0":                  data.Children[0],
				"/children/0/age":              data.Children[0].Age,
				"/children/0/data":             data.Children[0].Data,
				"/children/0/data/line":        data.Children[0].Data.Line,
				"/children/0/data/labels/a~1b": data.Children[0].Data.Labels["a/b"],
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getPointersData(tt.args.rv, tt.args.refs)
			if (err != nil) != tt.wantErr {
				t.Errorf("getPointersData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getPointersData() got = %v, want %v", got, tt.want)
			}
		})
	}
}
