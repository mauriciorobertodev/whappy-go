package input_test

import (
	"github.com/mauriciorobertodev/whappy-go/internal/app/input"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/picture"
	"github.com/mauriciorobertodev/whappy-go/internal/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Picture Inputs", func() {
	Describe("GetPictureInput Input", func() {
		It("should validate successfully", func() {
			inp := &input.GetPictureInput{
				PhoneOrJID:  "valid_jid@example.com",
				Preview:     utils.BoolPtr(true),
				IsCommunity: utils.BoolPtr(false),
			}
			Expect(inp.Validate()).To(BeNil())
		})

		It("should fail validation for empty PhoneOrJID field", func() {
			inp := &input.GetPictureInput{
				PhoneOrJID:  "",
				Preview:     utils.BoolPtr(true),
				IsCommunity: utils.BoolPtr(false),
			}
			Expect(inp.Validate()).To(Equal(picture.ErrInvalidJID))
		})
	})
})
