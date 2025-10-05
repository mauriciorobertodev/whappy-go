package webhook_test

import (
	"encoding/hex"
	"testing"
	"time"

	"github.com/mauriciorobertodev/whappy-go/internal/domain/webhook"
	"github.com/mauriciorobertodev/whappy-go/internal/fake"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCaches(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Webhook Entity Suite")
}

var _ = Describe("Webhook entity", func() {
	It("should create a new webhook", func() {
		w1 := webhook.New("http://example.com/w1", []string{"event1", "event2"}, true)
		Expect(w1).ToNot(BeNil())
		Expect(w1.ID).ToNot(BeEmpty())
		Expect(w1.URL).To(Equal("http://example.com/w1"))
		Expect(w1.Events).To(Equal([]string{"event1", "event2"}))
		Expect(w1.Active).To(BeTrue())
		Expect(w1.CreatedAt).To(BeTemporally("~", time.Now().UTC(), time.Second))
		Expect(w1.UpdatedAt).To(BeTemporally("~", time.Now().UTC(), time.Second))

		w2 := webhook.New("http://example.com/w2", []string{"event3"}, false)
		Expect(w2).ToNot(BeNil())
		Expect(w2.ID).ToNot(BeEmpty())
		Expect(w2.URL).To(Equal("http://example.com/w2"))
		Expect(w2.Events).To(Equal([]string{"event3"}))
		Expect(w2.Active).To(BeFalse())
		Expect(w2.CreatedAt).To(BeTemporally("~", time.Now().UTC(), time.Second))
		Expect(w2.UpdatedAt).To(BeTemporally("~", time.Now().UTC(), time.Second))
	})

	It("should activate and deactivate a webhook", func() {
		w := webhook.New("http://example.com/w", []string{"event"}, false)
		Expect(w.Active).To(BeFalse())

		w.Activate()
		Expect(w.Active).To(BeTrue())
		Expect(w.UpdatedAt).To(BeTemporally("~", time.Now().UTC(), time.Second))

		w.Deactivate()
		Expect(w.Active).To(BeFalse())
		Expect(w.UpdatedAt).To(BeTemporally("~", time.Now().UTC(), time.Second))
	})

	It("should update a webhook", func() {
		w := webhook.New("http://example.com/w", []string{"event"}, true)
		originalUpdatedAt := w.UpdatedAt

		time.Sleep(100 * time.Millisecond) // Ensure UpdatedAt will be different

		w.Update("http://example.com/updated", []string{"event1", "event2"})
		Expect(w.URL).To(Equal("http://example.com/updated"))
		Expect(w.Events).To(Equal([]string{"event1", "event2"}))
		Expect(w.UpdatedAt).To(BeTemporally("~", time.Now().UTC(), time.Second))
		Expect(w.UpdatedAt).To(BeTemporally(">", originalUpdatedAt))
	})

	It("should renew the secret of a webhook", func() {
		w := webhook.New("http://example.com/w", []string{"event"}, true)
		originalSecret := w.GetSecret()

		time.Sleep(100 * time.Millisecond) // Ensure secret change is noticeable

		w.RenewSecret()
		newSecret := w.GetSecret()
		Expect(newSecret).ToNot(BeEmpty())
		Expect(newSecret).ToNot(Equal(originalSecret))
		Expect(w.UpdatedAt).To(BeTemporally("~", time.Now().UTC(), time.Second))
	})

	It("should set and get the secret of a webhook", func() {
		w := webhook.New("http://example.com/w", []string{"event"}, true)
		originalSecret := w.GetSecret()

		w.SetSecret("new-secret")

		newSecret := w.GetSecret()
		Expect(newSecret).To(Equal("new-secret"))
		Expect(newSecret).ToNot(Equal(originalSecret))
		Expect(w.UpdatedAt).To(BeTemporally("~", time.Now().UTC(), time.Second))
	})

	It("should attach a webhook to an instance", func() {
		w := webhook.New("http://example.com/w", []string{"event"}, true)
		Expect(w.InstanceID).To(Equal(""))

		w.AttachToInstance("instance-123")
		Expect(w.InstanceID).To(Equal("instance-123"))
		Expect(w.UpdatedAt).To(BeTemporally("~", time.Now().UTC(), time.Second))
	})

	It("should generate a valid HMAC signature for an event", func() {
		wh := webhook.New("http://example.com/webhook", []string{
			"message:*",
			"message:new/text",
			"group:participants/joined",
		}, true)

		ev := fake.NewEvent().Create()

		sig, err := wh.SignEvent(ev)
		Expect(err).To(BeNil())
		Expect(sig).ToNot(BeEmpty())

		// Signature should be hex-encoded
		bytes, err := hex.DecodeString(sig)
		Expect(err).To(BeNil())
		Expect(len(bytes)).To(Equal(32)) // SHA256 produces 32-byte HMAC
	})
})
