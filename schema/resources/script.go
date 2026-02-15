package resources

import (
	"context"
	"encoding/json"
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

}
