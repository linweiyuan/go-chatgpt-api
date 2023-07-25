package imitate

import (
	"strings"
)

//goland:noinspection SpellCheckingInspection
func ConvertToString(chatgptResponse *ChatGPTResponse, previousText *StringStruct, role bool, id string, model string) string {
	text := strings.ReplaceAll(chatgptResponse.Message.Content.Parts[0], *&previousText.Text, "")
	translatedResponse := NewChatCompletionChunk(text, id, model)
	if role {
		translatedResponse.Choices[0].Delta.Role = chatgptResponse.Message.Author.Role
	}
	previousText.Text = chatgptResponse.Message.Content.Parts[0]
	return "data: " + translatedResponse.String() + "\n\n"
}
