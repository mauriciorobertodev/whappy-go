package input_test

import (
	"github.com/mauriciorobertodev/whappy-go/internal/app/input"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/chat"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Chat Inputs", func() {
	Describe("SendChatPresence Input", func() {
		It("should validate successfully", func() {
			inp := &input.SendChatPresenceInput{
				To:   "valid_jid@example.com",
				Type: chat.ChatPresenceTyping,
			}
			Expect(inp.Validate()).To(BeNil())
		})

		It("should fail validation for empty To field", func() {
			inp := &input.SendChatPresenceInput{
				To:   "",
				Type: chat.ChatPresenceTyping,
			}
			Expect(inp.Validate()).To(Equal(chat.ErrInvalidJID))
		})

		It("should fail validation for invalid Type field", func() {
			inp := &input.SendChatPresenceInput{
				To:   "valid_jid@example.com",
				Type: "invalid_type",
			}
			Expect(inp.Validate()).To(Equal(chat.ErrInvalidChatPresenceType))
		})
	})
})
