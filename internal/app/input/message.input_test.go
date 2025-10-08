package input_test

import (
	"github.com/mauriciorobertodev/whappy-go/internal/app/input"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/message"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Message Inputs", func() {
	Describe("ReadMessagesInput Input", func() {
		It("should validate successfully", func() {
			inp := &input.ReadMessagesInput{
				Chat:   "valid_chat_id",
				Sender: "valid_sender_id",
				IDs:    []string{"valid_message_id_1", "valid_message_id_2"},
			}
			Expect(inp.Validate()).To(BeNil())
		})
		It("should fail validation for empty Chat field", func() {

			inp := &input.ReadMessagesInput{
				Chat:   "",
				Sender: "valid_sender_id",
				IDs:    []string{"valid_message_id_1", "valid_message_id_2"},
			}
			Expect(inp.Validate()).To(Equal(message.ErrInvalidJID))
		})

		It("should fail validation for empty Sender field", func() {
			inp := &input.ReadMessagesInput{
				Chat:   "valid_chat_id",
				Sender: "",
				IDs:    []string{"valid_message_id_1", "valid_message_id_2"},
			}
			Expect(inp.Validate()).To(Equal(message.ErrInvalidJID))
		})

		It("should fail validation for empty IDs field", func() {
			inp := &input.ReadMessagesInput{
				Chat:   "valid_chat_id",
				Sender: "valid_sender_id",
				IDs:    []string{},
			}
			Expect(inp.Validate()).To(Equal(message.ErrEmptyMessageIDs))
		})
	})

	Describe("GenerateMessageIDs Input", func() {
		It("should validate successfully", func() {
			inp := &input.GenerateMessageIDs{
				Quantity: 5,
			}
			Expect(inp.Validate()).To(BeNil())
		})

		It("should fail validation for zero Quantity", func() {
			inp := &input.GenerateMessageIDs{
				Quantity: 0,
			}
			Expect(inp.Validate()).To(Equal(message.ErrInvalidQuantity))
		})

		It("should fail validation for negative Quantity", func() {
			inp := &input.GenerateMessageIDs{
				Quantity: -3,
			}
			Expect(inp.Validate()).To(Equal(message.ErrInvalidQuantity))
		})

		It("should fail validation for excessive Quantity", func() {
			inp := &input.GenerateMessageIDs{
				Quantity: message.MaxGenerateMessageIDs + 1,
			}
			Expect(inp.Validate()).To(Equal(message.ErrInvalidQuantity))
		})
	})
})
