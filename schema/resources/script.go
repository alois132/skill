package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/alois132/skill/util"
)

// 脚本

type ScriptFunc[I, O any] func(ctx context.Context, input I) (output O, err error)

type Script interface {
	Run(ctx context.Context, args string) (result string, err error)
	GetName() string
	GetUsage() string
}

type EasyScript[I, O any] struct {
	Name  string `json:"name"`
	Usage string `json:"usage"`
	Fn    ScriptFunc[I, O]
}

func (s *EasyScript[I, O]) Run(ctx context.Context, args string) (result string, err error) {
	input := util.NewInstance[I]()
	err = json.Unmarshal([]byte(args), &input)
	if err != nil {
		return "", err
	}

	output, err := s.Fn(ctx, input)
	if err != nil {
		return "", err
	}

	resultByte, err := json.Marshal(output)
	if err != nil {
		return "", err
	}

	return string(resultByte), nil
}

func (s *EasyScript[I, O]) GetName() string {
	return s.Name
}

func (s *EasyScript[I, O]) GetUsage() string {
	// 如果已经设置了使用说明，直接返回
	if s.Usage != "" {
		return s.Usage
	}

	// 从泛型类型中提取输入输出类型信息
	inType := util.TypeOf[I]()
	outType := util.TypeOf[O]()

	// 生成使用说明
	s.Usage = fmt.Sprintf("Input: %s, Output: %s", inType.String(), outType.String())
	return s.Usage
}

// NewEasyScript creates a new EasyScript with the given name and function
func NewEasyScript[I, O any](name string, fn ScriptFunc[I, O]) *EasyScript[I, O] {
	return &EasyScript[I, O]{
		Name: name,
		Fn:   fn,
	}
}

// WithUsage sets the usage description for the script
func (s *EasyScript[I, O]) WithUsage(usage string) *EasyScript[I, O] {
	s.Usage = usage
	return s
}

// TypeInfo returns information about the input and output types
func TypeInfo[I, O any]() (string, string) {
	inType := util.TypeOf[I]()
	outType := util.TypeOf[O]()
	return inType.String(), outType.String()
}

// IsEmptyType checks if a type is effectively empty (like interface{} for output)
func IsEmptyType(t reflect.Type) bool {
	if t.Kind() == reflect.Interface && t.NumMethod() == 0 {
		return true
	}
	return false
}
