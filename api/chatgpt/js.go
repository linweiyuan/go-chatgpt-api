package chatgpt

import (
	"fmt"
	"strings"
)

func getGetScript(url string, accessToken string, errorMessage string) string {
	return fmt.Sprintf(`
		fetch('%s', {
			headers: {
				'Authorization': '%s'
			}
		})
		.then(response => {
			if (!response.ok) {
				throw new Error('%s');
			}
			return response.text();
		})
		.then(text => {
			arguments[0](text);
		})
		.catch(err => {
			arguments[0](err.message);
		});
	`, url, accessToken, errorMessage)
}

func getPostScriptForStartConversation(url string, accessToken string, jsonString string, messageID string) string {
	return strings.ReplaceAll(fmt.Sprintf(`
		// get the whole data again to make sure get the endTurn message back
		const getEndTurnMessage = (dataArray) => {
			dataArray.pop(); // empty
			dataArray.pop(); // data: [DONE]
			return '!' + dataArray.pop().substring(6); // endTurn message
		};

		const xhr = new XMLHttpRequest();
		xhr.open('POST', '%s');
		xhr.setRequestHeader('Accept', 'text/event-stream');
		xhr.setRequestHeader('Authorization', '%s');
		xhr.setRequestHeader('Content-Type', 'application/json');
		xhr.onreadystatechange = function() {
			switch (xhr.readyState) {
				case xhr.LOADING: {
					switch (xhr.status) {
						case 200: {
							const dataArray = xhr.responseText.substr(xhr.seenBytes).split("\n\n");
							dataArray.pop(); // empty string
							if (dataArray.length) {
								let data = dataArray.pop(); // target data
								if (data === 'data: [DONE]') { // this DONE will break the ending handling
									data = getEndTurnMessage(xhr.responseText.split("\n\n"));
								} else if (data.startsWith('event')) {
									data = data.substring(49);
								}
								if (data) {
									if (data.startsWith('!')) {
										conversationMap.set('[MESSAGE_!@#_ID]', data);
									} else {
										conversationMap.set('[MESSAGE_!@#_ID]', data.substring(6));
									}
								}
							}
							break;
						}
						case 401: {
							conversationMap.set('[MESSAGE_!@#_ID]', xhr.status + 'Access token has expired.');
							break;
						}
						case 403: {
							conversationMap.set('[MESSAGE_!@#_ID]', xhr.status + 'Something went wrong. If this issue persists please contact us through our help center at help.openai.com.');
							break;
						}
						case 404: {
							conversationMap.set('[MESSAGE_!@#_ID]', xhr.status + JSON.parse(xhr.responseText).detail);
							break;
						}
						case 413: {
							conversationMap.set('[MESSAGE_!@#_ID]', xhr.status + JSON.parse(xhr.responseText).detail.message);
							break;
						}
						case 422: {
							const detail = JSON.parse(xhr.responseText).detail[0];
							conversationMap.set('[MESSAGE_!@#_ID]', xhr.status + detail.loc + ' -> ' + detail.msg);
							break;
						}
						case 429: {
							conversationMap.set('[MESSAGE_!@#_ID]', xhr.status + JSON.parse(xhr.responseText).detail);
							break;
						}
						case 500: {
							conversationMap.set('[MESSAGE_!@#_ID]', xhr.status + 'Unknown error.');
							break;
						}
					}
					xhr.seenBytes = xhr.responseText.length;
					break;
				}
				case xhr.DONE:
					// keep exception handling
					const conversationData = conversationMap.get('[MESSAGE_!@#_ID]');
					if (conversationData) {
						if (!conversationData.startsWith('4') && !conversationData.startsWith('5')) {
							conversationMap.set('[MESSAGE_!@#_ID]', getEndTurnMessage(xhr.responseText.split("\n\n")));
						}
					}
					break;
			}
		};
		xhr.send(JSON.stringify(%s));
	`, url, accessToken, jsonString), "[MESSAGE_!@#_ID]", messageID)
}

func getPostScript(url string, accessToken string, jsonString string, errorMessage string) string {
	return fmt.Sprintf(`
		fetch('%s', {
			method: 'POST',
			headers: {
				'Authorization': '%s',
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(%s)
		})
		.then(response => {
			if (!response.ok) {
				throw new Error('%s');
			}
			return response.text();
		})
		.then(text => {
			arguments[0](text);
		})
		.catch(err => {
			arguments[0](err.message);
		});
	`, url, accessToken, jsonString, errorMessage)
}

func getPatchScript(url string, accessToken string, jsonString string, errorMessage string) string {
	return fmt.Sprintf(`
		fetch('%s', {
			method: 'PATCH',
			headers: {
				'Authorization': '%s',
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(%s)
		})
		.then(response => {
			if (!response.ok) {
				throw new Error('%s');
			}
			return response.text();
		})
		.then(text => {
			arguments[0](text);
		})
		.catch(err => {
			arguments[0](err.message);
		});
	`, url, accessToken, jsonString, errorMessage)
}
