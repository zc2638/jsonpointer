// Package jsonpointer
// Created by zc on 2021/8/29.
package jsonpointer

import (
	"reflect"
	"testing"
)

func TestError_Error(t *testing.T) {
	type fields struct {
		Ref     string
		Message string
		Default interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "ok",
			fields: fields{
				Ref:     "/data",
				Message: "data type not support",
				Default: nil,
			},
			want: "path: /data, error: data type not support",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Error{
				Ref:     tt.fields.Ref,
				Message: tt.fields.Message,
				Default: tt.fields.Default,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_WithDefault(t *testing.T) {
	type fields struct {
		Ref     string
		Message string
		Default interface{}
	}
	type args struct {
		value interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Error
	}{
		{
			name: "ok",
			fields: fields{
				Message: "data type not support",
				Default: "test data",
			},
			args: args{
				value: "test data",
			},
			want: &Error{
				Message: "data type not support",
				Default: "test data",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Error{
				Ref:     tt.fields.Ref,
				Message: tt.fields.Message,
				Default: tt.fields.Default,
			}
			if got := e.WithDefault(tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_WithRef(t *testing.T) {
	type fields struct {
		Ref     string
		Message string
		Default interface{}
	}
	type args struct {
		ref string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Error
	}{
		{
			name: "ok",
			fields: fields{
				Message: "data type not support",
			},
			args: args{
				ref: "/data",
			},
			want: &Error{
				Ref:     "/data",
				Message: "data type not support",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Error{
				Ref:     tt.fields.Ref,
				Message: tt.fields.Message,
				Default: tt.fields.Default,
			}
			if got := e.WithRef(tt.args.ref); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithRef() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewError(t *testing.T) {
	type args struct {
		message string
	}
	tests := []struct {
		name string
		args args
		want *Error
	}{
		{
			name: "ok",
			args: args{
				message: "data type not support",
			},
			want: &Error{
				Message: "data type not support",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewError(tt.args.message); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewError() = %v, want %v", got, tt.want)
			}
		})
	}
}
