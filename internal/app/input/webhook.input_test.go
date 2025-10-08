package input_test

import (
	"github.com/mauriciorobertodev/whappy-go/internal/app/input"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/webhook"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Webhook Inputs", func() {
	Describe("CreateWebhook Input", func() {
		It("should validate successfully", func() {
			inp := &input.CreateWebhook{
				Active: true,
				URL:    "https://example.com/webhook",
				Events: []string{"event1", "event2"},
			}
			Expect(inp.Validate()).To(BeNil())
		})

		It("should return an error for invalid URL", func() {
			inp := &input.CreateWebhook{
				Active: true,
				URL:    "invalid-url",
				Events: []string{"event1", "event2"},
			}
			Expect(inp.Validate()).To(Equal(webhook.ErrInvalidURL))
		})
	})

	Describe("ToggleWebhook Input", func() {
		It("should validate successfully", func() {
			inp := &input.ToggleWebhook{
				ID:     "123e4567-e89b-12d3-a456-426614174000",
				Active: true,
			}
			Expect(inp.Validate()).To(BeNil())
		})

		It("should return an error for invalid UUID", func() {
			inp := &input.ToggleWebhook{
				ID:     "invalid-uuid",
				Active: true,
			}
			Expect(inp.Validate()).To(Equal(webhook.ErrInvalidID))
		})
	})

	Describe("UpdateWebhook Input", func() {
		It("should validate successfully", func() {
			inp := &input.UpdateWebhook{
				ID:     "123e4567-e89b-12d3-a456-426614174000",
				Active: true,
				URL:    "https://example.com/webhook",
				Events: []string{"event1", "event2"},
			}
			Expect(inp.Validate()).To(BeNil())
		})

		It("should return an error for invalid URL", func() {
			inp := &input.UpdateWebhook{
				ID:     "123e4567-e89b-12d3-a456-426614174000",
				Active: true,
				URL:    "invalid-url",
				Events: []string{"event1", "event2"},
			}
			Expect(inp.Validate()).To(Equal(webhook.ErrInvalidURL))
		})

		It("should return an error for invalid UUID", func() {
			inp := &input.UpdateWebhook{
				ID:     "invalid-uuid",
				Active: true,
				URL:    "https://example.com/webhook",
				Events: []string{"event1", "event2"},
			}
			Expect(inp.Validate()).To(Equal(webhook.ErrInvalidID))
		})
	})

	Describe("GetWebhook Input", func() {
		It("should validate successfully", func() {
			inp := &input.GetWebhook{
				ID: "123e4567-e89b-12d3-a456-426614174000",
			}
			Expect(inp.Validate()).To(BeNil())
		})

		It("should return an error for invalid UUID", func() {
			inp := &input.GetWebhook{
				ID: "invalid-uuid",
			}
			Expect(inp.Validate()).To(Equal(webhook.ErrInvalidID))
		})
	})

	Describe("DeleteWebhook Input", func() {
		It("should validate successfully", func() {
			inp := &input.DeleteWebhook{
				ID: "123e4567-e89b-12d3-a456-426614174000",
			}
			Expect(inp.Validate()).To(BeNil())
		})

		It("should return an error for invalid UUID", func() {
			inp := &input.DeleteWebhook{
				ID: "invalid-uuid",
			}
			Expect(inp.Validate()).To(Equal(webhook.ErrInvalidID))
		})
	})
})
