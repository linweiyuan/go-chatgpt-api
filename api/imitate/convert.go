package imitate

import (
	"strings"
)

//goland:noinspection SpellCheckingInspection
func ConvertToString(chatgptResponse *ChatGPTResponse, previousText *StringStruct, role bool) string {
	translatedResponse := NewChatCompletionChunk(strings.ReplaceAll(chatgptResponse.Message.Content.Parts[0], *&previousText.Text, ""))
	if role {
		translatedResponse.Choices[0].Delta.Role = chatgptResponse.Message.Author.Role
	}
	previousText.Text = chatgptResponse.Message.Content.Parts[0]
	return "data: " + translatedResponse.String() + "\n\n"
}
