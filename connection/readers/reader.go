package readers

type Reader interface {
	Read() ([]byte, error)
}
