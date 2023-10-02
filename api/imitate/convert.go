package imitate

import (
	"fmt"
	"strings"
)

func ConvertToString(chatgptResponse *ChatGPTResponse, previousText *StringStruct, role bool, id string, model string) string {
	var text string

	if len(chatgptResponse.Message.Content.Parts) == 1 {
		if part, ok := chatgptResponse.Message.Content.Parts[0].(string); ok {
			text = strings.ReplaceAll(part, previousText.Text, "")
			previousText.Text = part
		} else {
			text = fmt.Sprintf("%v", chatgptResponse.Message.Content.Parts[0])
		}
	} else {
		// When using GPT-4 messages with images (multimodal_text), the length of 'parts' might be 2.
		// Since the chatgpt API currently does not support multimodal content
		// and there is no official format for multimodal content,
		// the content is temporarily returned as is.
		var parts []string
		for _, part := range chatgptResponse.Message.Content.Parts {
			parts = append(parts, fmt.Sprintf("%v", part))
		}
		text = strings.Join(parts, ", ")
	}

	translatedResponse := NewChatCompletionChunk(text, id, model)
	if role {
		translatedResponse.Choices[0].Delta.Role = chatgptResponse.Message.Author.Role
	}

	return "data: " + translatedResponse.String() + "\n\n"
}
