package input_test

import (
	"github.com/mauriciorobertodev/whappy-go/internal/app/input"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/contact"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Contact Inputs", func() {
	Describe("CheckPhones Input", func() {
		It("should validate successfully", func() {
			inp := &input.CheckPhones{
				Phones: []string{"user@example.com"},
			}
			Expect(inp.Validate()).To(BeNil())
		})

		It("should return an error for empty Phones", func() {
			inp := &input.CheckPhones{
				Phones: []string{},
			}
			Expect(inp.Validate()).To(Equal(contact.EmptyPhones))
		})
	})

	Describe("GetContact Input", func() {
		It("should validate successfully", func() {
			inp := &input.GetContact{
				PhoneOrJID: "user@example.com",
			}
			Expect(inp.Validate()).To(BeNil())
		})

		It("should return an error for empty PhoneOrJID", func() {
			inp := &input.GetContact{
				PhoneOrJID: "",
			}
			Expect(inp.Validate()).To(Equal(contact.ErrInvalidJID))
		})
	})
})
