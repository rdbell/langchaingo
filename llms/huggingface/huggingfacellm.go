package huggingface

import (
	"context"
	"errors"
	"os"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/huggingface/internal/huggingfaceclient"
)

const tokenEnvVarName = "HUGGINGFACEHUB_API_TOKEN"

var (
	ErrEmptyResponse = errors.New("empty response")
	ErrMissingToken  = errors.New("missing the Hugging Face API token. Set it in the HUGGINGFACEHUB_API_TOKEN environment variable") //nolint:lll
)

type LLM struct {
	client *huggingfaceclient.Client
}

var _ llms.LLM = (*LLM)(nil)

func (o *LLM) Call(ctx context.Context, prompt string, stopWords []string) (string, error) {
	r, err := o.Generate(ctx, []string{prompt}, stopWords)
	if err != nil {
		return "", err
	}
	if len(r) == 0 {
		return "", ErrEmptyResponse
	}
	return r[0].Text, nil
}

func (o *LLM) Generate(ctx context.Context, prompts []string, stopWords []string) ([]*llms.Generation, error) {
	_ = stopWords // TODO: use this
	result, err := o.client.RunInference(ctx, &huggingfaceclient.InferenceRequest{
		RepoID: "google/flan-t5-xl",
		Prompt: prompts[0],
		Task:   huggingfaceclient.InferenceTaskTextGeneration,
	})
	if err != nil {
		return nil, err
	}
	return []*llms.Generation{
		{Text: result.Text},
	}, nil
}

func New() (*LLM, error) {
	token := os.Getenv(tokenEnvVarName)
	if token == "" {
		return nil, ErrMissingToken
	}
	c, err := huggingfaceclient.New(token)
	return &LLM{
		client: c,
	}, err
}
