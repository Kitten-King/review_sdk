package review_sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient создает экземпляр SDK-клиента для сервиса отзывов
func NewClient(baseURL string) Client {
	return &client{
		baseURL: baseURL, // Например, "http://localhost:8082/query"
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// Универсальная структура для GraphQL-запроса
type gqlRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

// Универсальный контейнер для разбора GraphQL-ответа
type gqlResponse struct {
	Data   json.RawMessage `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

func (c *client) CreateReview(ctx context.Context, input NewReview) (*Review, error) {
	// Пишем мутацию точно так же, как в песочнице (Playground)
	mutation := `
		mutation CreateReview($input: NewReview!) {
			createReview(input: $input) {
				id
				userId
				cityId
				rating
				comment
				createdAt
			}
		}
	`

	reqBody := gqlRequest{
		Query: mutation,
		Variables: map[string]interface{}{
			"input": input,
		},
	}

	var respData struct {
		CreateReview Review `json:"createReview"`
	}

	err := c.doRequest(ctx, reqBody, &respData)
	if err != nil {
		return nil, err
	}

	return &respData.CreateReview, nil
}

func (c *client) GetCityReviews(ctx context.Context, cityID int) ([]Review, error) {
	query := `
		query GetCityReviews($cityId: Int!) {
			getCityReviews(cityId: $cityId) {
				id
				userId
				cityId
				rating
				comment
				createdAt
			}
		}
	`

	reqBody := gqlRequest{
		Query: query,
		Variables: map[string]interface{}{
			"cityId": cityID,
		},
	}

	var respData struct {
		GetCityReviews []Review `json:"getCityReviews"`
	}

	err := c.doRequest(ctx, reqBody, &respData)
	if err != nil {
		return nil, err
	}

	return respData.GetCityReviews, nil
}

// Вспомогательный метод для отправки HTTP POST запроса
func (c *client) doRequest(ctx context.Context, gqlReq gqlRequest, target interface{}) error {
	marshaled, err := json.Marshal(gqlReq)
	if err != nil {
		return fmt.Errorf("failed to marshal gql request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL, bytes.NewBuffer(marshaled))
	if err != nil {
		return fmt.Errorf("failed to create http request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var gqlResp gqlResponse
	if err := json.NewDecoder(resp.Body).Decode(&gqlResp); err != nil {
		return fmt.Errorf("failed to decode gql response: %w", err)
	}

	// Если GraphQL вернул ошибки внутри схемы
	if len(gqlResp.Errors) > 0 {
		return fmt.Errorf("graphql error: %s", gqlResp.Errors[0].Message)
	}

	// Анмаршалим конкретную data-структуру в наш таргет
	if err := json.Unmarshal(gqlResp.Data, target); err != nil {
		return fmt.Errorf("failed to unmarshal data block: %w", err)
	}

	return nil
}
