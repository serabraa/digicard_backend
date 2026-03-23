package main

import (
	"businesscard-wallet-backend/delivery/http/handlers"
	"businesscard-wallet-backend/infrastructure/assets"
	"businesscard-wallet-backend/infrastructure/config"
	"businesscard-wallet-backend/infrastructure/packager"
	"businesscard-wallet-backend/infrastructure/signer"
	"businesscard-wallet-backend/infrastructure/template"
	"businesscard-wallet-backend/usecase/generatepass"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("godotenv load warning: %v", err)
	}
	cfg := config.Load()

	certPaths, err := prepareSignerFiles()
	if err != nil {
		log.Fatalf("failed to prepare signer files: %v", err)
	}

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
		certPaths.passCertPath,
		certPaths.passKeyPath,
		certPaths.wwdrPath,
		filepath.Join(os.TempDir(), "digicard-work"),
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
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
	}

	log.Printf("server started on :%s", cfg.Port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

type signerFilePaths struct {
	passCertPath string
	passKeyPath  string
	wwdrPath     string
}

func prepareSignerFiles() (signerFilePaths, error) {
	passCertPEM := strings.ReplaceAll(os.Getenv("PASS_CERT_PEM"), `\n`, "\n")
	passKeyPEM := strings.ReplaceAll(os.Getenv("PASS_KEY_PEM"), `\n`, "\n")
	wwdrPEM := strings.ReplaceAll(os.Getenv("WWDR_PEM"), `\n`, "\n")

	if passCertPEM == "" {
		return signerFilePaths{}, fmt.Errorf("PASS_CERT_PEM is empty")
	}
	if passKeyPEM == "" {
		return signerFilePaths{}, fmt.Errorf("PASS_KEY_PEM is empty")
	}
	if wwdrPEM == "" {
		return signerFilePaths{}, fmt.Errorf("WWDR_PEM is empty")
	}

	baseDir := filepath.Join(os.TempDir(), "digicard-certs")
	if err := os.MkdirAll(baseDir, 0o700); err != nil {
		return signerFilePaths{}, fmt.Errorf("create temp cert dir: %w", err)
	}

	passCertPath := filepath.Join(baseDir, "pass-cert.pem")
	passKeyPath := filepath.Join(baseDir, "pass-key.pem")
	wwdrPath := filepath.Join(baseDir, "wwdr.pem")

	if err := os.WriteFile(passCertPath, []byte(passCertPEM), 0o600); err != nil {
		return signerFilePaths{}, fmt.Errorf("write pass cert pem: %w", err)
	}
	if err := os.WriteFile(passKeyPath, []byte(passKeyPEM), 0o600); err != nil {
		return signerFilePaths{}, fmt.Errorf("write pass key pem: %w", err)
	}
	if err := os.WriteFile(wwdrPath, []byte(wwdrPEM), 0o600); err != nil {
		return signerFilePaths{}, fmt.Errorf("write wwdr pem: %w", err)
	}

	return signerFilePaths{
		passCertPath: passCertPath,
		passKeyPath:  passKeyPath,
		wwdrPath:     wwdrPath,
	}, nil
}
