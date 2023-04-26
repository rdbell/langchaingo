package memory

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tmc/langchaingo/schema"
)

func TestBufferMemory(t *testing.T) {
	t.Parallel()

	m := NewBuffer()
	result1, err := m.LoadMemoryVariables(map[string]any{})
	require.NoError(t, err)
	expected1 := map[string]any{"history": ""}
	assert.Equal(t, expected1, result1)

	err = m.SaveContext(map[string]any{"foo": "bar"}, map[string]any{"bar": "foo"})
	require.NoError(t, err)

	result2, err := m.LoadMemoryVariables(map[string]any{})
	require.NoError(t, err)

	expected2 := map[string]any{"history": "Human: bar\nAI: foo"}
	assert.Equal(t, expected2, result2)
}

func TestBufferMemoryReturnMessage(t *testing.T) {
	t.Parallel()

	m := NewBuffer()
	m.ReturnMessages = true
	expected1 := map[string]any{"history": []schema.ChatMessage{}}
	result1, err := m.LoadMemoryVariables(map[string]any{})
	require.NoError(t, err)
	assert.Equal(t, expected1, result1)

	err = m.SaveContext(map[string]any{"foo": "bar"}, map[string]any{"bar": "foo"})
	require.NoError(t, err)

	result2, err := m.LoadMemoryVariables(map[string]any{})
	require.NoError(t, err)

	expectedChatHistory := NewChatMessageHistory(
		WithPreviousMessages([]schema.ChatMessage{
			schema.HumanChatMessage{Text: "bar"},
			schema.AIChatMessage{Text: "foo"},
		}),
	)

	expected2 := map[string]any{"history": expectedChatHistory.Messages()}
	assert.Equal(t, expected2, result2)
}

func TestBufferMemoryWithPreLoadedHistory(t *testing.T) {
	t.Parallel()

	m := NewBuffer()
	m.ChatHistory = NewChatMessageHistory(
		WithPreviousMessages([]schema.ChatMessage{
			schema.HumanChatMessage{Text: "bar"},
			schema.AIChatMessage{Text: "foo"},
		}),
	)

	result, err := m.LoadMemoryVariables(map[string]any{})
	require.NoError(t, err)
	expected := map[string]any{"history": "Human: bar\nAI: foo"}
	assert.Equal(t, expected, result)
}
