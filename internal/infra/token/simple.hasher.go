package token

type SimpleHasher struct{}

func NewSimpleHasher() *SimpleHasher {
	return &SimpleHasher{}
}

func (h *SimpleHasher) Hash(input string) (string, error) {
	return input, nil
}

func (h *SimpleHasher) Compare(hash, input string) bool {
	return hash == input
}
