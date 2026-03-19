package signer

type NoopSigner struct{}

func NewNoopSigner() *NoopSigner {
	return &NoopSigner{}
}

func (s *NoopSigner) Sign(manifest []byte) ([]byte, error) {
	return []byte("SIGNATURE_PLACEHOLDER"), nil
}
