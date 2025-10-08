package input_test

import (
	"github.com/mauriciorobertodev/whappy-go/internal/app/input"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/token"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Token Inputs", func() {
	Describe("GetToken Input", func() {
		It("should validate successfully", func() {
			inp := &input.GetToken{
				ID: "123e4567-e89b-12d3-a456-426614174000",
			}
			Expect(inp.Validate()).To(BeNil())
		})

		It("should fail validation for invalid UUID", func() {
			inp := &input.GetToken{
				ID: "invalid-uuid",
			}
			Expect(inp.Validate()).To(Equal(token.ErrInvalidID))
		})
	})
})
