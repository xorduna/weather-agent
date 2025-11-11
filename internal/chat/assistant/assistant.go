package assistant

import (
	"context"
	"errors"
	"github.com/acai-travel/tech-challenge/internal/chat/assistant/tools"
	"github.com/acai-travel/tech-challenge/internal/chat/model"
	"github.com/openai/openai-go/v2"
	"go.opentelemetry.io/otel"
	"log/slog"
	"strings"
)

var tracer = otel.Tracer("assistant")

type Assistant struct {
	cli             openai.Client
	registeredTools map[string]Tool
	openaiTools     []openai.ChatCompletionToolUnionParam
}

type Tool interface {
	Name() string
	Description() string
	Parameters() openai.FunctionParameters
	Execute(ctx context.Context, args ...string) (string, error)
}

func New() *Assistant {

	usedTools := []Tool{
		&tools.WeatherTool{},
		&tools.TodayTool{},
	}

	openaiTools := []openai.ChatCompletionToolUnionParam{}
	registeredTools := map[string]Tool{}

	for _, t := range usedTools {
		registeredTools[t.Name()] = t
		openaiTools = append(openaiTools, openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
			Name:        t.Name(),
			Description: openai.String(t.Description()),
			Parameters:  t.Parameters(),
		}))
	}

	return &Assistant{
		cli:             openai.NewClient(),
		registeredTools: registeredTools,
		openaiTools:     openaiTools,
	}
}

func (a *Assistant) Title(ctx context.Context, conv *model.Conversation) (string, error) {
	ctx, span := tracer.Start(ctx, "Assistant.Reply")
	defer span.End()

	if len(conv.Messages) == 0 {
		return "An empty conversation", nil
	}

	slog.InfoContext(ctx, "Generating title for conversation", "conversation_id", conv.ID)

	msgs := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage("Generate a concise, descriptive title for the conversation. It should reflect users intention based on the user message. The title should be a single line, no more than 80 characters, and should not include any special characters or emojis."),
	}

	for _, m := range conv.Messages {
		msgs = append(msgs, openai.UserMessage(m.Content))
	}

	resp, err := a.cli.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:    openai.ChatModelO1,
		Messages: msgs,
	})

	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 || strings.TrimSpace(resp.Choices[0].Message.Content) == "" {
		return "", errors.New("empty response from OpenAI for title generation")
	}

	title := resp.Choices[0].Message.Content
	title = strings.ReplaceAll(title, "\n", " ")
	title = strings.Trim(title, " \t\r\n-\"'")

	if len(title) > 80 {
		title = title[:80]
	}

	slog.InfoContext(ctx, "Generated title for conversation", "conversation_id", conv.ID, "title", title)

	return title, nil
}

func (a *Assistant) Reply(ctx context.Context, conv *model.Conversation) (string, error) {
	ctx, span := tracer.Start(ctx, "Assistant.Reply")
	defer span.End()

	if len(conv.Messages) == 0 {
		return "", errors.New("conversation has no messages")
	}

	slog.InfoContext(ctx, "Generating reply for conversation", "conversation_id", conv.ID)

	msgs := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage("You are a helpful, concise AI assistant. Provide accurate, safe, and clear responses."),
	}

	for _, m := range conv.Messages {
		switch m.Role {
		case model.RoleUser:
			msgs = append(msgs, openai.UserMessage(m.Content))
		case model.RoleAssistant:
			msgs = append(msgs, openai.AssistantMessage(m.Content))
		}
	}

	for i := 0; i < 15; i++ {
		resp, err := a.cli.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
			Model:    openai.ChatModelGPT4_1,
			Messages: msgs,
			Tools:    a.openaiTools,
		})

		if err != nil {
			return "", err
		}

		if len(resp.Choices) == 0 {
			return "", errors.New("no choices returned by OpenAI")
		}

		if message := resp.Choices[0].Message; len(message.ToolCalls) > 0 {
			msgs = append(msgs, message.ToParam())

			for _, call := range message.ToolCalls {
				slog.InfoContext(ctx, "Tool call received", "name", call.Function.Name, "args", call.Function.Arguments)

				if _, ok := a.registeredTools[call.Function.Name]; !ok {
					return "", errors.New("unknown tool call: " + call.Function.Name)
				} else {
					slog.InfoContext(ctx, "Executing tool", "name", call.Function.Name)

					// For the sake of simplicity, we ignore the error from tool execution here.
					answer, _ := a.registeredTools[call.Function.Name].Execute(ctx, call.Function.Arguments)
					msgs = append(msgs, openai.ToolMessage(answer, call.ID))
				}
			}

			continue
		}

		return resp.Choices[0].Message.Content, nil
	}

	return "", errors.New("too many tool calls, unable to generate reply")
}
