package main

import (
	"log"

	"monelog/hello/internal/api"
	"monelog/hello/internal/handlers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// 1. Inisialisasi handler yang sudah kita buat
	txHandler := handlers.NewTransactionHandler()

	// 2. Konversi ke strict handler (sesuai dengan opsi generate)
	strictHandler := api.NewStrictHandler(txHandler, nil) // Nilai kedua untuk error handler (opsional)

	// 3. Daftarkan handler ke Echo router

	// Fungsi RegisterHandlers ini di-generate oleh oapi-codegen

	api.RegisterHandlers(e, strictHandler)

	// 4. Jalankan server
	log.Fatal(e.Start(":8080"))

}
