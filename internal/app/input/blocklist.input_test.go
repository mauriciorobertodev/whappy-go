package input_test

import (
	"github.com/mauriciorobertodev/whappy-go/internal/app/input"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/blocklist"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Blocklist Inputs", func() {
	Describe("Block Input", func() {
		It("should validate successfully", func() {
			inp := &input.Block{
				PhoneOrJID: "user@example.com",
			}
			Expect(inp.Validate()).To(BeNil())
		})

		It("should return an error for empty PhoneOrJID", func() {
			inp := &input.Block{
				PhoneOrJID: "",
			}
			Expect(inp.Validate()).To(Equal(blocklist.ErrInvalidJID))
		})
	})

	Describe("Unblock Input", func() {
		It("should validate successfully", func() {
			inp := &input.Unblock{
				PhoneOrJID: "user@example.com",
			}
			Expect(inp.Validate()).To(BeNil())
		})

		It("should return an error for empty PhoneOrJID", func() {
			inp := &input.Unblock{
				PhoneOrJID: "",
			}
			Expect(inp.Validate()).To(Equal(blocklist.ErrInvalidJID))
		})
	})
})
