# README.md

**Install:**

```bash
go install github.com/242617/ycgpt@latest
```

**Описание:**

Этот проект представляет собой программу на языке программирования Go, которая использует API Yandex Cloud для генерации текста с помощью модели Yandex GPT-Lite. Программа позволяет генерировать текст на основе заданного запроса и текста.

**Установка:**

Для запуска программы необходимо установить Go и убедиться, что переменные окружения YC_FOLDER_ID и YC_IAM_TOKEN установлены корректно. Переменные окружения можно установить в настройках вашего аккаунта Yandex Cloud.

**Использование:**

Программа принимает два аргумента: текст для генерации и (необязательный) текст для добавления к запросу. Если вы хотите использовать только один аргумент, укажите пустую строку в качестве второго аргумента.

Пример использования:
```
$ go run main.go "Hello, world!"
"Good morning, world! It's a beautiful day outside."
```

В этом примере программа сгенерирует текст "Good morning, world! It's a beautiful day outside.".

**Параметры:**

- **YC_FOLDER_ID**: идентификатор папки в Yandex Cloud, в которой находится модель Yandex GPT-Lite.
- **YC_IAM_TOKEN**: токен IAM, который предоставляет доступ к API Yandex Cloud.

Эти параметры можно получить из переменных окружения.

**Пример использования API:**

1. Создайте объект ycCompletitionRequest с параметрами запроса (например, modelURI, completionOptions и messages).
2. Отправьте запрос на API с помощью http.NewRequest и http.DefaultClient.
3. Обработайте ответ от API, используя json.NewDecoder.
4. Получите результат от API и используйте его для генерации текста.
5. Выведите результат.

Вот пример реализации этого алгоритма:
```go
func retrieve(text, data string) string {
    // Создание объекта ycCompletionRequest
    request := ycCompletitionRequest{
        ModelURI: fmt.Sprintf("gpt://%s/yandexgpt-lite", YCFolderID),
        CompletionOptions: struct {
            Temperature float64 `json:"temperature"`
            MaxTokens string  `json:"maxTokens"`
        }{
            Temperature: .6,
            MaxTokens: "2000",
        },
        Messages: append(request.Messages,
            Message{
                Role: "system",
                Text: text,
            },
        ),
    }
    if data != "" {
        request.Messages = append(request.Messages,
            Message{
                Role: "user",
                Text: data,
            },
        )
    }

    // Отправка запроса на API
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
```
Этот код отправляет запрос на API и получает результат. Затем он использует результат для генерации текста. Результат выводится на экран.

Это лишь пример использования программы. Вы можете изменить код, чтобы удовлетворить свои потребности.

**Заключение:**

Эта программа представляет собой пример использования API Yandex Cloud для генерации текста. Она может быть использована для создания простых приложений, которые используют модели NLP.

Обратите внимание, что некоторые параметры берутся из переменных окружения, например, YC_FOLDER_ID и YC_IAM_TOKEN. Убедитесь, что эти переменные установлены корректно перед запуском программы.
