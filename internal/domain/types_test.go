package domain

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConversation(t *testing.T) {
	conv := Conversation{
		ConversationID: "c1",
		Turns: []Turn{
			{
				Query:  "Hi",
				Answer: "Hello",
			},
		},
		Annotation: "Pass",
		Extra: map[string]json.RawMessage{
			"first_test_field":  json.RawMessage(`"first"`),
			"second_test_field": json.RawMessage(`"second"`),
		},
	}

	data, err := json.Marshal(&conv)
	assert.NoError(t, err, "No error expected when deserialize conversation")

	var m map[string]json.RawMessage
	err = json.Unmarshal(data, &m)
	require.Contains(t, m, "conversation_id")
	assert.Equal(t, json.RawMessage(`"c1"`), m["conversation_id"])

	require.Contains(t, m, "first_test_field")
	assert.Equal(t, json.RawMessage(`"first"`), m["first_test_field"])

	require.Contains(t, m, "second_test_field")
	assert.Equal(t, json.RawMessage(`"second"`), m["second_test_field"])
}

func TestConversationMarshalUnmarshal(t *testing.T) {

	conv := Conversation{
		ConversationID: "c1",
		Turns: []Turn{
			{
				Query:  "Hi",
				Answer: "Hello",
			},
		},
		Annotation: "Pass",
		Extra: map[string]json.RawMessage{
			"first_test_field":  json.RawMessage(`"first"`),
			"second_test_field": json.RawMessage(`"second"`),
		},
	}

	data, err := json.Marshal(&conv)
	assert.NoError(t, err, "No error expected when deserialize conversation")

	var resultConv Conversation
	err = json.Unmarshal(data, &resultConv)
	assert.NoError(t, err, "No error expected when serialize conversation")

	assert.Equal(t, conv.ConversationID, resultConv.ConversationID)
	assert.Equal(t, conv.Annotation, resultConv.Annotation)
	for i := 0; i < len(conv.Turns); i++ {
		assert.Equal(t, conv.Turns[i].Query, resultConv.Turns[i].Query)
		assert.Equal(t, conv.Turns[i].Answer, resultConv.Turns[i].Answer)
	}
}

func TestConversationMarshalJSON_OmitsEmptyAnnotation(t *testing.T) {
	conv := Conversation{
		ConversationID: "c1",
		Turns: []Turn{
			{Query: "Hi", Answer: "Hello"},
		},
		Annotation: "",
	}

	data, err := json.Marshal(&conv)
	require.NoError(t, err)

	var m map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(data, &m))

	assert.NotContains(t, m, "human_annotation")
	require.Contains(t, m, "conversation_id")
	assert.Equal(t, json.RawMessage(`"c1"`), m["conversation_id"])
}
