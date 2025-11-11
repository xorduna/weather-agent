package tools

import (
    "time"
)

type TimeTool struct{}

func (t *TimeTool) Name() string {
    return "time"
}

func (t *TimeTool) Execute(args ...string) (string, error) {
    return time.Now().Format(time.RFC3339), nil
}
package tools

type EchoTool struct{}

func (e *EchoTool) Name() string {
    return "echo"
}

func (e *EchoTool) Execute(args ...string) (string, error) {
    return "Echo: " + joinArgs(args), nil
}

func joinArgs(args []string) string {
    result := ""
    for i, arg := range args {
        if i > 0 {
            result += " "
        }
        result += arg
    }
    return result
}

