package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Структура для JSON-RPC запроса
type JsonRpcRequest struct {
	Jsonrpc string        `json:"jsonrpc"`
	Id      int           `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

// Структура для ответа
type JsonRpcResponse struct {
	Jsonrpc string        `json:"jsonrpc"`
	Id      int           `json:"id"`
	Result  BalanceResult `json:"result"`
	Error   *JsonRpcError `json:"error"`
}

// Структура для результата баланса
type BalanceResult struct {
	Value int64 `json:"value"`
}

// Структура для обработки ошибок
type JsonRpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type BoostResponse struct {
	Data Data `json:"data"`
}

type Data struct {
	Boosts []Boost `json:"boosts"`
}

type Boost struct {
	Url         string   `json:"url"`
	ChainId     string   `json:"chainid"`
	TokenAdress string   `json:"tokenAddress"`
	Amount      float64  `json:"amount"`
	TotalAmount float64  `json:"totalAmount"`
	Icon        string   `json:"icon"`
	Header      string   `json:"header"`
	Description string   `json:"description"`
	Links       LinkData `json:"links"`
}

type LinkData struct {
	objs any
}

func getBalance(walletAddress string) (int64, error) {
	url := "https://api.mainnet-beta.solana.com"

	// Создаем JSON-RPC запрос
	request :=
		JsonRpcRequest{
			Jsonrpc: "2.0",
			Id:      1,
			Method:  "getBalance",
			Params:  []interface{}{walletAddress},
		}

	// Преобразуем структуру в JSON
	jsonReq, err := json.Marshal(request)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal request: %v", err)
	}

	// Отправляем запрос
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonReq))
	if err != nil {
		return 0, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("received non-200 response: %v", resp.Status)
	}

	// Декодируем ответ
	var jsonResp JsonRpcResponse
	if err := json.NewDecoder(resp.Body).Decode(&jsonResp); err != nil {
		return 0, fmt.Errorf("failed to decode response: %v", err)
	}

	// Обрабатываем ошибки, если они есть
	if jsonResp.Error != nil {
		return 0, fmt.Errorf("error from Solana: %s", jsonResp.Error.Message)
	}

	// Возвращаем баланс
	return jsonResp.Result.Value, nil
}

func main() {
	//walletAddress := "EsoQDanDVFM3bL9XxguhAtc7iTN4HBQTkJfzJ9MEMfTc" // Замените на адрес вашего кошелька
	url := "https://api.dexscreener.com/token-boosts/latest/v1"
	/*
		balance, err := getBalance(walletAddress)
		if err != nil {
			fmt.Printf("Error getting balance: %v\n", err)
			return
		}

		fmt.Printf("Balance of wallet %s: %d lamports\n", walletAddress, balance)
	*/
	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("Ошибка при выполнении запроса: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		fmt.Printf("Ошибка: Получен статус-код %d от API", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Ошибка при чтении тела ответа: %v", err)
	}

	fmt.Println("Данные из body:")
	fmt.Println(string(body))

	var boostResponse BoostResponse

	err = json.Unmarshal(body, &boostResponse)
	if err != nil {
		fmt.Printf("Ошибка при десереализации json: %v", err)
	}

	if err := json.NewDecoder(response.Body).Decode(&boostResponse); err != nil {
		fmt.Printf("Ошибка при декодировании JSON: %v", err)
	}

	for _, boost := range boostResponse.Data.Boosts {
		fmt.Printf("Токен: %s, Сумма: %.2f\n", boost.TokenAdress, boost.TotalAmount)
	}
}
