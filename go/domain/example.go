package domain

import (
	"github.com/AlekSi/pointer"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type ExampleMsg struct {
	Data Example
	Msg  *types.Message
}

type Example struct {
	Filed1 string
	Filed2 int
	Filed3 bool
	Field4 float64
	Field5 []string
	Field6 map[string]string
	Field7 *string
}

func NewExample() *Example {
	return &Example{
		Filed1: "filed1",
		Filed2: 1,
		Filed3: true,
		Field4: 1.1,
		Field5: []string{"field5"},
		Field6: map[string]string{"field6": "field6"},
		Field7: pointer.ToString("field7"),
	}
}
