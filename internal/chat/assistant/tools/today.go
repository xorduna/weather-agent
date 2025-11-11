package tools

import (
	"context"
	"github.com/openai/openai-go/v2"
	"time"
)

type TodayTool struct{}

func (t *TodayTool) Name() string {
	return "get_today_date"
}

func (t *TodayTool) Description() string {
	return "Get the current date and time in RFC3339 format."
}

func (t *TodayTool) Parameters() openai.FunctionParameters {
	return openai.FunctionParameters{}
}

func (t *TodayTool) Execute(ctx context.Context, args ...string) (string, error) {
	return time.Now().Format(time.RFC3339), nil
}
