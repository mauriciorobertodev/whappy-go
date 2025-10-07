package consumer_test

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/mauriciorobertodev/whappy-go/internal/app/logger"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/events"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/instance"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/webhook"
	"github.com/mauriciorobertodev/whappy-go/internal/fake"
	"github.com/mauriciorobertodev/whappy-go/internal/infra/cache"
	"github.com/mauriciorobertodev/whappy-go/internal/infra/config"
	"github.com/mauriciorobertodev/whappy-go/internal/infra/consumer"
	"github.com/mauriciorobertodev/whappy-go/internal/infra/database"
	"github.com/mauriciorobertodev/whappy-go/internal/infra/eventbus"
	"github.com/mauriciorobertodev/whappy-go/internal/infra/repository"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestConsumes(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Consumers Suite")
}

var _ = AfterSuite(func() {
	_ = os.Remove("test.db")
})

var _ = Describe("Webhook consumer", func() {
	config.LoadLoggers(logger.LevelNone)

	bus := eventbus.New(&config.EventBusConfig{
		Driver: config.EventBusDriverInMemory,
	})

	db := database.New(&config.DatabaseConfig{
		Driver: config.DatabaseDriverSQLite,
		DbName: "test",
	})

	cache := cache.New(&config.CacheConfig{
		Driver: config.CacheDriverInMemory,
	})

	instRepo := repository.NewInstanceRepository(db)
	webRepo := repository.NewWebhookRepository(db)

	webhookConsumer := consumer.NewWebhookConsumer(webRepo, cache)

	bus.SubscribeAll(webhookConsumer.Handle)

	migrator := database.NewMigrator(db, db.DriverName())

	BeforeEach(func() {
		migrator.Reset()
		cache.Flush()
	})

	It("should consume an event and send webhooks", func() {
		evt1 := fake.NewEvent().WithInstanceID("instance-1").WithName("fake:event/batata").Create()
		evt2 := fake.NewEvent().WithInstanceID("instance-1").WithName("fake:event/tomate").Create()

		received := false
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer GinkgoRecover()
			received = true
			// fmt.Println("✅ Recebi essa bagaçaaaaaaaaaaaaa!")

			Expect(r.Header.Get("Content-Type")).To(Equal("application/json"))
			Expect(r.Header.Get("X-Whappy-Signature")).ToNot(BeEmpty())
			Expect(r.Header.Get("X-Whappy-Timestamp")).ToNot(BeEmpty())
			Expect(r.Header.Get("X-Whappy-Event")).To(Equal("fake:event/batata"))
			Expect(r.Method).To(Equal(http.MethodPost))

			wh, err := webRepo.Get(webhook.WhereInstanceID("instance-1"), webhook.WhereActive(true))
			Expect(err).To(BeNil())
			Expect(wh).ToNot(BeNil())

			signature, err := wh.SignEvent(evt1, evt1.OccurredAt.Unix())
			Expect(err).To(BeNil())
			Expect(signature).ToNot(BeEmpty())
			Expect(r.Header.Get("X-Whappy-Signature")).To(Equal(signature))
			Expect(r.Header.Get("X-Whappy-Timestamp")).To(Equal(fmt.Sprintf("%d", evt1.OccurredAt.Unix())))

			body, err := io.ReadAll(r.Body)
			Expect(err).To(BeNil())

			// Json to event
			var evt events.Event
			err = json.Unmarshal(body, &evt)
			Expect(err).To(BeNil())
			Expect(evt.Name).To(Equal(evt1.Name))
			Expect(evt.OccurredAt).To(Equal(evt1.OccurredAt))
			Expect(evt.InstanceID).ToNot(BeNil())
			Expect(*evt.InstanceID).To(Equal("instance-1"))
			Expect(evt.Payload).To(Equal(evt1.Payload))

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		}))
		defer ts.Close()

		instRepo.InsertMany([]*instance.Instance{
			fake.InstanceFactory().WithID("instance-1").Create(),
		})

		webRepo.InsertMany([]*webhook.Webhook{
			fake.WebhookFactory().WithURL(ts.URL).WithEvents([]string{"fake:event/batata"}).Active().WithInstanceID("instance-1").Create(),
		})

		bus.Publish(evt1)
		bus.Publish(evt2)

		Eventually(func() bool { return received }, "2s", "100ms").Should(BeTrue())
	})

	It("should validate webhook signature correctly", func() {
		secret := "super-hiper-mega-secret"
		evt := fake.NewEvent().WithInstanceID("instance-1").WithName("fake:event/batata").Create()

		received := false
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer GinkgoRecover()
			received = true

			body, err := io.ReadAll(r.Body)
			Expect(err).To(BeNil())

			sentSignature := r.Header.Get("X-Whappy-Signature")
			Expect(sentSignature).ToNot(BeEmpty())

			sentTimestamp := r.Header.Get("X-Whappy-Timestamp")
			Expect(sentTimestamp).ToNot(BeEmpty())

			wh, err := webRepo.Get(webhook.WhereInstanceID("instance-1"), webhook.WhereActive(true))
			Expect(err).To(BeNil())
			Expect(wh).ToNot(BeNil())

			message := append(body, []byte(r.Header.Get("X-Whappy-Timestamp"))...)

			h := hmac.New(sha256.New, []byte(secret))
			h.Write(message)

			Expect(hex.EncodeToString(h.Sum(nil))).To(Equal(sentSignature))

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		}))
		defer ts.Close()

		instRepo.InsertMany([]*instance.Instance{
			fake.InstanceFactory().WithID("instance-1").Create(),
		})

		webRepo.InsertMany([]*webhook.Webhook{
			fake.WebhookFactory().WithURL(ts.URL).WithEvents([]string{"fake:event/batata"}).Active().WithInstanceID("instance-1").WithSecret(secret).Create(),
		})

		bus.Publish(evt)

		Eventually(func() bool { return received }, "2s", "100ms").Should(BeTrue())
	})

	It("should handle a burst of 100 webhook events", func() {
		numEvents := 100
		receivedCount := 0
		secret := "super-hiper-mega-secret"

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer GinkgoRecover()
			receivedCount++

			// fmt.Println("✅ Recebi essa bagaçaaaaaaaaaaaaa!", receivedCount)

			body, err := io.ReadAll(r.Body)
			Expect(err).To(BeNil())

			sentSignature := r.Header.Get("X-Whappy-Signature")
			Expect(sentSignature).ToNot(BeEmpty())

			sentTimestamp := r.Header.Get("X-Whappy-Timestamp")
			Expect(sentTimestamp).ToNot(BeEmpty())

			wh, err := webRepo.Get(webhook.WhereInstanceID("instance-1"), webhook.WhereActive(true))
			Expect(err).To(BeNil())
			Expect(wh).ToNot(BeNil())

			var event events.Event
			err = json.Unmarshal(body, &event)
			Expect(err).To(BeNil())
			Expect(event.InstanceID).ToNot(BeNil())
			Expect(*event.InstanceID).To(Equal("instance-1"))
			Expect(event.OccurredAt).ToNot(BeZero())
			Expect(fmt.Sprintf("%d", event.OccurredAt.Unix())).To(Equal(sentTimestamp))

			Expect(wh.GetSecret()).To(Equal(secret))

			message := append(body, []byte(r.Header.Get("X-Whappy-Timestamp"))...)

			h := hmac.New(sha256.New, []byte(secret))
			h.Write(message)

			Expect(hex.EncodeToString(h.Sum(nil))).To(Equal(sentSignature))
			Expect(wh.SignEvent(event, event.OccurredAt.Unix())).To(Equal(sentSignature))

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		}))
		defer ts.Close()

		instRepo.InsertMany([]*instance.Instance{
			fake.InstanceFactory().WithID("instance-1").Create(),
		})

		webRepo.InsertMany([]*webhook.Webhook{
			fake.WebhookFactory().
				WithURL(ts.URL).
				WithEvents([]string{"fake:event/*"}).
				Active().
				WithInstanceID("instance-1").
				WithSecret(secret).
				Create(),
		})

		wh, err := webRepo.Get(webhook.WhereInstanceID("instance-1"), webhook.WhereActive(true))
		Expect(err).To(BeNil())
		Expect(wh).ToNot(BeNil())
		Expect(wh.GetSecret()).To(Equal(secret))

		for i := 0; i < numEvents; i++ {
			ev := fake.NewEvent().WithInstanceID("instance-1").WithRandomName("fake:event/").Create()
			bus.Publish(ev)
		}

		Eventually(func() bool { return receivedCount == numEvents }, "10s", "100ms").Should(BeTrue())
	})

})
