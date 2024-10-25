package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
)

const YCCompletitionURL = "https://llm.api.cloud.yandex.net/foundationModels/v1/completion"

var (
	YCFolderID = os.Getenv("YC_FOLDER_ID")
	YCIAMToken = os.Getenv("YC_IAM_TOKEN")
)

func init() { log.SetFlags(log.Lshortfile) }
func main() {
	if len(os.Args) <= 2 {
		log.Fatal("unexpeted arguments number")
	}
	text := strings.Join(os.Args[1:], " ")
	if len(text) < 10 {
		log.Fatal("text too short")
	}

	info, err := os.Stdin.Stat()
	die(err)

	var data string
	if info.Size() > 0 {
		b, err := io.ReadAll(os.Stdin)
		die(err)
		data = strings.TrimSpace(string(b))
	}

	fmt.Println(retrieve(text, data))
}

func retrieve(text, data string) string {
	request := ycCompletitionRequest{ModelURI: fmt.Sprintf("gpt://%s/yandexgpt-lite", YCFolderID)}
	request.CompletionOptions.Temperature = .6
	request.CompletionOptions.MaxTokens = "2000"
	request.Messages = append(request.Messages,
		Message{
			Role: "system",
			Text: text,
		},
	)
	if data != "" {
		request.Messages = append(request.Messages,
			Message{
				Role: "user",
				Text: data,
			},
		)
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(request)
	die(err)

	req, err := http.NewRequest(http.MethodPost, YCCompletitionURL, &buf)
	die(err)
	req.Header.Add("Authorization", "Bearer "+YCIAMToken)
	req.Header.Add("x-folder-id", YCFolderID)

	res, err := http.DefaultClient.Do(req)
	die(err)
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		die(fmt.Errorf("unexpected status %d", res.StatusCode))
	}

	var response ycCompletitionResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	die(err)

	return func() string {
		for _, alternative := range response.Result.Alternatives {
			if alternative.Message.Role == "assistant" {
				return alternative.Message.Text
			}
		}
		die(errors.New("no assistant text found"))
		return ""
	}()
}

type ycCompletitionRequest struct {
	ModelURI          string `json:"modelUri"`
	CompletionOptions struct {
		Temperature float64 `json:"temperature"`
		MaxTokens   string  `json:"maxTokens"`
	} `json:"completionOptions"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role string `json:"role"`
	Text string `json:"text"`
}

type ycCompletitionResponse struct {
	Result struct {
		Alternatives []struct {
			Message struct {
				Role string `json:"role"`
				Text string `json:"text"`
			} `json:"message"`
		} `json:"alternatives"`
	} `json:"result"`
}

func die(err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("%s:%d: %s", file, line, err.Error())
		os.Exit(1)
	}
}
