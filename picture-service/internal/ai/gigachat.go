package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"picture-service/internal/config"
	"strings"

	gigago "github.com/AlexandrVIvanov/gigago"
)

type AiRepository struct {
	client     *gigago.Client
	httpClient *http.Client
	authKey    string
}

func InitAiRepository(cfg *config.Config) (*AiRepository, error) {
	ctx := context.Background()

	client, err := gigago.NewClient(
		ctx,
		cfg.GigaChat.AuthKey,
		gigago.WithCustomInsecureSkipVerify(true),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create GigaChat client: %w", err)
	}

	return &AiRepository{
		client:     client,
		httpClient: &http.Client{},
		authKey:    cfg.GigaChat.AuthKey,
	}, nil
}

// SplitTextToScenes разбивает текст на логические сцены
func (r *AiRepository) SplitTextToScenes(ctx context.Context, text string) ([]string, error) {
	model := r.client.GenerativeModel("GigaChat")

	prompt := fmt.Sprintf(`Ты — ИИ-ассистент для анализа художественных текстов.
Разбей следующий текст на логические сцены. Каждая сцена должна быть отдельным блоком текста.
Верни ответ в формате JSON: {"scenes": ["сцена 1", "сцена 2", ...]}

Текст: %s

Только JSON, без лишних слов.`, text)

	messages := []gigago.Message{
		{Role: gigago.RoleUser, Content: prompt},
	}

	resp, err := model.Generate(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("failed to split text: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("empty response from GigaChat")
	}

	var result struct {
		Scenes []string `json:"scenes"`
	}

	if err := json.Unmarshal([]byte(resp.Choices[0].Message.Content), &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return result.Scenes, nil
}

// GeneratePrompt создаёт промпт для генерации изображения на основе текста сцены
func (r *AiRepository) GeneratePrompt(ctx context.Context, sceneText string) (string, error) {
	model := r.client.GenerativeModel("GigaChat")

	prompt := fmt.Sprintf(`На основе описания сцены создай короткий промпт для генерации изображения.
Промпт должен быть на русском языке, 30-50 слов.

Сцена: %s

Только промпт, без лишних слов.`, sceneText)

	messages := []gigago.Message{
		{Role: gigago.RoleUser, Content: prompt},
	}

	resp, err := model.Generate(ctx, messages)
	if err != nil {
		return "", fmt.Errorf("failed to generate prompt: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("empty response from GigaChat")
	}

	return resp.Choices[0].Message.Content, nil
}

// GenerateImage генерирует изображение через Kandinsky
func (r *AiRepository) GenerateImage(ctx context.Context, prompt string) ([]byte, error) {
	finalPrompt := prompt
	if !strings.Contains(strings.ToLower(prompt), "нарисуй") {
		finalPrompt = "Нарисуй " + prompt
	}

	// Используем прямой HTTP запрос к API GigaChat с параметром function_call: auto
	fileID, err := r.requestImageGeneration(ctx, finalPrompt)
	if err != nil {
		return nil, fmt.Errorf("failed to request image generation: %w", err)
	}

	// Скачиваем изображение по полученному file_id
	imageData, err := r.downloadFile(ctx, fileID)
	if err != nil {
		return nil, fmt.Errorf("failed to download image: %w", err)
	}

	return imageData, nil
}

// requestImageGeneration отправляет запрос к GigaChat API на генерацию изображения
func (r *AiRepository) requestImageGeneration(ctx context.Context, prompt string) (string, error) {
	url := "https://gigachat.devices.sberbank.ru/api/v1/chat/completions"

	body := map[string]interface{}{
		"model": "Kandinsky",
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"function_call": "auto",
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+r.authKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("empty response from API")
	}

	content := result.Choices[0].Message.Content
	fileID := extractFileID(content)
	if fileID == "" {
		return "", fmt.Errorf("failed to extract file ID from response: %s", content)
	}

	return fileID, nil
}

// downloadFile скачивает файл из GigaChat по file_id
func (r *AiRepository) downloadFile(ctx context.Context, fileID string) ([]byte, error) {
	url := fmt.Sprintf("https://gigachat.devices.sberbank.ru/api/v1/files/%s/content", fileID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+r.authKey)
	req.Header.Set("Accept", "application/octet-stream")

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download file: status %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

// extractFileID извлекает UUID изображения из HTML-строки
func extractFileID(content string) string {
	start := strings.Index(content, "src=\"")
	if start == -1 {
		return ""
	}
	start += 5
	end := strings.Index(content[start:], "\"")
	if end == -1 {
		return ""
	}
	fileID := content[start : start+end]
	return strings.TrimPrefix(fileID, "/")
}
