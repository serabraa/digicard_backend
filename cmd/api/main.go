package main

import (
	"businesscard-wallet-backend/delivery/http/handlers"
	"businesscard-wallet-backend/infrastructure/assets"
	"businesscard-wallet-backend/infrastructure/config"
	"businesscard-wallet-backend/infrastructure/packager"
	"businesscard-wallet-backend/infrastructure/signer"
	"businesscard-wallet-backend/infrastructure/template"
	"businesscard-wallet-backend/usecase/generatepass"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	cfg := config.Load()

	templateProvider := template.NewStaticTemplateProvider(
		cfg.Organization,
		cfg.PassTypeID,
		cfg.TeamID,
		cfg.BackgroundRGB,
		cfg.LabelColorRGB,
		cfg.ForegroundRGB,
	)

	assetProvider := assets.NewLocalAssetProvider(cfg.AssetsDir)
	signerService := signer.NewOpenSSLSigner(
		"certs/tmp/pass-cert.pem",
		"certs/tmp/pass-key.pem",
		"certs/tmp/wwdr.pem",
		"certs/tmp/work",
	)
	packagerService := packager.NewZipPackager()

	generatePassService := generatepass.NewService(
		templateProvider,
		assetProvider,
		signerService,
		packagerService,
	)

	generatePassHandler := handlers.NewGeneratePassHandler(generatePassService)

	mux := http.NewServeMux()
	mux.Handle("/generate-pass", generatePassHandler)

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
	}

	log.Printf("server started on :%s", cfg.Port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
