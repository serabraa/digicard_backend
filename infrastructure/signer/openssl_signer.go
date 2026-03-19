package signer

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type OpenSSLSigner struct {
	passCertPEM string
	passKeyPEM  string
	wwdrPEM     string
	workDir     string
}

func NewOpenSSLSigner(passCertPEM, passKeyPEM, wwdrPEM, workDir string) *OpenSSLSigner {
	return &OpenSSLSigner{
		passCertPEM: passCertPEM,
		passKeyPEM:  passKeyPEM,
		wwdrPEM:     wwdrPEM,
		workDir:     workDir,
	}
}

func (s *OpenSSLSigner) Sign(manifest []byte) ([]byte, error) {
	runDir := filepath.Join(s.workDir, fmt.Sprintf("sign-%d", time.Now().UnixNano()))
	if err := os.MkdirAll(runDir, 0o755); err != nil {
		return nil, fmt.Errorf("create signer workdir: %w", err)
	}
	defer os.RemoveAll(runDir)

	manifestPath := filepath.Join(runDir, "manifest.json")
	signaturePath := filepath.Join(runDir, "signature")

	if err := os.WriteFile(manifestPath, manifest, 0o644); err != nil {
		return nil, fmt.Errorf("write manifest: %w", err)
	}

	cmd := exec.Command(
		"openssl", "smime", "-binary", "-sign",
		"-signer", s.passCertPEM,
		"-inkey", s.passKeyPEM,
		"-certfile", s.wwdrPEM,
		"-in", manifestPath,
		"-out", signaturePath,
		"-outform", "DER",
		"-nodetach",
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("openssl sign failed: %w: %s", err, string(output))
	}

	signature, err := os.ReadFile(signaturePath)
	if err != nil {
		return nil, fmt.Errorf("read signature: %w", err)
	}

	return signature, nil
}
