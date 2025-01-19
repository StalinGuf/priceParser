package rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Структура для JSON-RPC запроса
type JsonRpcRequest struct {
	Jsonrpc string   `json:"jsonrpc"`
	Id      int      `json:"id"`
	Method  string   `json:"method"`
	Params  []string `json:"params"`
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

func GetBalance(walletAddress string) (int64, error) {
	url := "https://api.mainnet-beta.solana.com"

	// Создаем JSON-RPC запрос
	request := JsonRpcRequest{
		Jsonrpc: "2.0",
		Id:      1,
		Method:  "getBalance",
		Params:  []string{walletAddress},
	}

	jsonReq, err := json.Marshal(request)
	if err != nil {
		return 0, fmt.Errorf("ошибка при маршалинге: %v", err)
	}

	// Отправляем запрос
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonReq))
	if err != nil {
		return 0, fmt.Errorf("ошибка при отправке запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("не получили успешный 200 код ответа: %v", resp.Status)
	}

	var jsonResp JsonRpcResponse
	if err := json.NewDecoder(resp.Body).Decode(&jsonResp); err != nil {
		return 0, fmt.Errorf("не смогли декодировать ответ: %v", err)
	}
	if jsonResp.Error != nil {
		return 0, fmt.Errorf("ошибка в сети: %s", jsonResp.Error.Message)
	}

	return jsonResp.Result.Value, nil
}
