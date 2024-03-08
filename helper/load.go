package helper

import (
	"fmt"

	"github.com/joho/godotenv"
)

func Envload() {
	err := godotenv.Load()
	if err != nil {
		panic("Failed to connect env")
	}
	fmt.Println("Connected to env")
}
