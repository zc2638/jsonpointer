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
	dataReflect = reflect.ValueOf(data)
)

func TestNewParser(t *testing.T) {
	type args struct {
		data interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    *Parser
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				data: data,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewParser(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewParser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestParser_Batch(t *testing.T) {
	type fields struct {
		rv *reflect.Value
	}
	type args struct {
		refs []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				rv: &dataReflect,
			},
			args: args{
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
			p := &Parser{
				rv: tt.fields.rv,
			}
			got, err := p.Batch(tt.args.refs)
			if (err != nil) != tt.wantErr {
				t.Errorf("Batch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Batch() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_Check(t *testing.T) {
	type fields struct {
		rv *reflect.Value
	}
	type args struct {
		ref string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "ok",
			fields: fields{
				rv: &dataReflect,
			},
			args: args{
				ref: "/children/0/data",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				rv: tt.fields.rv,
			}
			if got := p.Check(tt.args.ref); got != tt.want {
				t.Errorf("Check() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_Get(t *testing.T) {
	type fields struct {
		rv *reflect.Value
	}
	type args struct {
		ref string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				rv: &dataReflect,
			},
			args: args{
				ref: "/children/0/data",
			},
			want:    data.Children[0].Data,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				rv: tt.fields.rv,
			}
			got, err := p.Get(tt.args.ref)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getPointerData(t *testing.T) {
	type args struct {
		rv       reflect.Value
		refPaths []string
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				rv:       dataReflect,
				refPaths: refPaths("/children/0/data"),
			},
			want:    data.Children[0].Data,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getPointerData(tt.args.rv, tt.args.refPaths)
			if (err != nil) != tt.wantErr {
				t.Errorf("getPointerData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getPointerData() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_refPaths(t *testing.T) {
	type args struct {
		ref string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "ok",
			args: args{
				ref: "/children/0/data",
			},
			want: []string{
				"children",
				"0",
				"data",
			},
		},
		{
			name: "pointer",
			args: args{
				ref: "/children/0/data/labels/a~1b",
			},
			want: []string{
				"children",
				"0",
				"data",
				"labels",
				"a~1b",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := refPaths(tt.args.ref); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("refPaths() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_transferPointer(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "ok",
			args: args{
				key: "children",
			},
			want: "children",
		},
		{
			name: "pointer",
			args: args{
				key: "a~0b",
			},
			want: "a~b",
		},
		{
			name: "pointer_slash",
			args: args{
				key: "a~1b",
			},
			want: "a/b",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := transferPointer(tt.args.key); got != tt.want {
				t.Errorf("transferPointer() = %v, want %v", got, tt.want)
			}
		})
	}
}
