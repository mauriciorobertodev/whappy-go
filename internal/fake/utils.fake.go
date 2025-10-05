package fake

import (
	"math/rand"
)

func FakePhone() string {
	return "55" + FakeDigits(11)
}

func FakeJID() string {
	return FakePhone() + "@s.whatsapp.net"
}

func FakeLID() string {
	return FakeDigits(20) + "@lid"
}

func FakeDigits(n int) string {
	if n <= 0 {
		return ""
	}

	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = "0123456789"[rand.Intn(10)]
	}
	return string(b)
}

func FakeSha256() string {
	return FakeHex(64)
}

func FakeHex(n int) string {
	if n <= 0 {
		return ""
	}

	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = "0123456789abcdef"[rand.Intn(16)]
	}
	return string(b)
}
