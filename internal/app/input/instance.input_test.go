package input_test

import (
	"github.com/mauriciorobertodev/whappy-go/internal/app/input"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/instance"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Instance Inputs", func() {
	Describe("CreateInstance Input", func() {
		It("should validate successfully", func() {
			inp := &input.CreateInstance{
				Name: "Test Instance",
			}
			Expect(inp.Validate()).To(BeNil())
		})

		It("should fail validation when Name is empty", func() {
			inp := &input.CreateInstance{
				Name: "",
			}
			Expect(inp.Validate()).To(Equal(instance.ErrNameTooShort))
		})

		It("should fail validation when Name is too long", func() {
			longName := ""
			for i := 0; i < instance.MaxNameLength+1; i++ {
				longName += "a"
			}
			inp := &input.CreateInstance{
				Name: longName,
			}
			Expect(inp.Validate()).To(Equal(instance.ErrNameTooLong))
		})
	})

	Describe("GetInstance Input", func() {
		It("should validate successfully", func() {
			inp := &input.GetInstance{
				ID: "valid-id",
			}
			Expect(inp.Validate()).To(BeNil())
		})

		It("should fail validation when ID is empty", func() {
			inp := &input.GetInstance{
				ID: "",
			}
			Expect(inp.Validate()).To(Equal(instance.ErrInvalidID))
		})
	})

	Describe("RenewInstanceToken Input", func() {
		It("should validate successfully", func() {
			inp := &input.RenewInstanceToken{
				ID: "valid-id",
			}
			Expect(inp.Validate()).To(BeNil())
		})

		It("should fail validation when ID is empty", func() {
			inp := &input.RenewInstanceToken{
				ID: "",
			}
			Expect(inp.Validate()).To(Equal(instance.ErrInvalidID))
		})
	})
})
