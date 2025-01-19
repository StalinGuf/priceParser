package main

import (
	"fmt"
	restapi "solMod/restApi"
	"solMod/rpc"
	"solMod/utils"
)

func main() {
	//Случайный адресс
	wallet := "7LFYBv2FQDY6aw5hL5SBCW2PiHG8TmENxBnfUnCksRbs"

	v, err := rpc.GetBalance(wallet)
	if err != nil {
		fmt.Printf("Ошибка при выполнении rpc запроса получения баланса: %v", err)
		return
	}
	fmt.Printf("Баланс равен %d\n", v)

	utils.MyDummy()

	fmt.Println("***Идём за инфой по по рекламируемым токенам***")

	err = restapi.GetBoostedTokens()
	if err != nil {
		fmt.Printf("Ошибка при выполнении запроса получения токенов: %v", err)
		return
	}

}
