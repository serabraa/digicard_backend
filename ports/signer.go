package ports

type Signer interface {
	Sign(manifest []byte) ([]byte, error)
}
