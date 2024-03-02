package server

import (
	"fmt"

	"github.com/Aadithya-V/qp-gen/routers.go"
)

func Start() error {
	fmt.Println("Starting server...")
	r := routers.Router()
	err := r.Run(":8001")
	if err != nil {
		fmt.Printf("Failed to start server: %v", err)
		return err
	}
	return nil
}
