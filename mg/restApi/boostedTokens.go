package restapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type tokenInfo struct {
	Url         string              `json:"url"`
	ChainId     string              `json:"chainid"`
	TokenAdress string              `json:"tokenAddress"`
	Amount      float64             `json:"amount"`
	TotalAmount float64             `json:"totalAmount"`
	Icon        string              `json:"icon"`
	Header      string              `json:"header"`
	Description string              `json:"description"`
	Links       []map[string]string `json:"links"`
}

func GetBoostedTokens() error {

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	url := "https://api.dexscreener.com/token-boosts/latest/v1"

	response, err := client.Get(url)
	if err != nil {
		fmt.Printf("Ошибка при выполнении запроса: %v", err)
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		fmt.Printf("Ошибка: Получен статус-код %d от API", response.StatusCode)
		return err
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Ошибка при чтении тела ответа: %v", err)
		return err
	}
	var tokens []tokenInfo
	err = json.Unmarshal(body, &tokens)
	if err != nil {
		fmt.Printf("Ошибка при десереализации json: %v", err)
		return err
	}

	for _, token := range tokens {
		//Смотрим только Solana сеть
		if token.ChainId == "solana" {
			fmt.Printf("Токен: %s, Сумма %.3f, Общая сумма: %.2f\n", token.TokenAdress, token.Amount, token.TotalAmount)

			fmt.Println("линки:")
			for i, link := range token.Links {
				fmt.Printf("  линк #%d:\n", i+1)
				for k, v := range link {
					fmt.Printf("    %s: %s\n", k, v)
				}
			}
		}

	}
	return nil
}
