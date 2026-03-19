package ports

type PackageFile struct {
	Name string
	Data []byte
}

type Packager interface {
	Package(files []PackageFile) ([]byte, error)
}
