package eventbus_test

import (
	"testing"
	"time"

	"github.com/joho/godotenv"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/mauriciorobertodev/whappy-go/internal/app/logger"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/events"
	"github.com/mauriciorobertodev/whappy-go/internal/fake"
	"github.com/mauriciorobertodev/whappy-go/internal/infra/config"
	"github.com/mauriciorobertodev/whappy-go/internal/infra/eventbus"
)

func TestEventBus(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Event Bus Suite")
}

var _ = DescribeTableSubtree("Event Bus", func(driver string) {
	Expect(godotenv.Load("./../../../.env")).ToNot(HaveOccurred())
	config.LoadLoggers(logger.LevelNone)

	var (
		bus events.EventBus
	)

	BeforeEach(func() {
		if driver == "memory" {
			bus = eventbus.New(&config.EventBusConfig{
				Driver: config.EventBusDriverInMemory,
			})
		}

		if driver == "redis" {
			bus = eventbus.New(&config.EventBusConfig{
				Driver: config.EventBusDriverRedis,
				RedisConfig: &config.RedisConfig{
					Host: config.GetEnvString("REDIS_HOST", "localhost"),
					Port: config.GetEnvInt("REDIS_PORT", 6379),
				},
			})
		}
	})

	It("should deliver events to a specific handler", func() {
		received1 := make(chan events.Event, 1)
		received2 := make(chan events.Event, 1)

		bus.Subscribe("test.event.1", func(e events.Event) { received1 <- e })
		bus.Subscribe("test.event.2", func(e events.Event) { received2 <- e })
		bus.Subscribe("test.event.3", func(e events.Event) { received2 <- e })

		bus.Publish(fake.NewEvent().WithName("test.event.1").WithPayload("test-payload.1").Create())
		time.Sleep(10 * time.Millisecond) // small sleep to ensure order of events
		bus.Publish(fake.NewEvent().WithName("test.event.2").WithPayload("test-payload.2").Create())
		time.Sleep(10 * time.Millisecond) // small sleep to ensure order of events
		bus.Publish(fake.NewEvent().WithName("test.event.3").WithPayload("test-payload.3").Create())

		var got1, got2, got3 events.Event
		Eventually(received1, 100*time.Millisecond).Should(Receive(&got1))
		Eventually(received2, 100*time.Millisecond).Should(Receive(&got2))
		Eventually(received2, 100*time.Millisecond).Should(Receive(&got3))

		Expect(string(got1.Name)).To(Equal("test.event.1"))
		Expect(got1.Payload).To(Equal("test-payload.1"))

		Expect(string(got2.Name)).To(Equal("test.event.2"))
		Expect(got2.Payload).To(Equal("test-payload.2"))

		Expect(string(got3.Name)).To(Equal("test.event.3"))
		Expect(got3.Payload).To(Equal("test-payload.3"))
	})

	It("should deliver events to all handlers", func() {
		received1 := make(chan events.Event, 1)
		received2 := make(chan events.Event, 1)

		bus.SubscribeAll(func(e events.Event) { received1 <- e })
		bus.SubscribeAll(func(e events.Event) { received2 <- e })

		bus.Publish(fake.NewEvent().WithName("test.event.1").WithPayload("test-payload.1").Create())
		bus.Publish(fake.NewEvent().WithName("test.event.2").WithPayload("test-payload.2").Create())
		bus.Publish(fake.NewEvent().WithName("test.event.3").WithPayload("test-payload.3").Create())

		var r1got1, r1got2, r1got3 events.Event
		Eventually(received1, 100*time.Millisecond).Should(Receive(&r1got1))
		Eventually(received1, 100*time.Millisecond).Should(Receive(&r1got2))
		Eventually(received1, 100*time.Millisecond).Should(Receive(&r1got3))

		Expect(string(r1got1.Name)).To(Equal("test.event.1"))
		Expect(r1got1.Payload).To(Equal("test-payload.1"))
		Expect(string(r1got2.Name)).To(Equal("test.event.2"))
		Expect(r1got2.Payload).To(Equal("test-payload.2"))
		Expect(string(r1got3.Name)).To(Equal("test.event.3"))
		Expect(r1got3.Payload).To(Equal("test-payload.3"))

		var r2got1, r2got2, r2got3 events.Event
		Eventually(received2, 100*time.Millisecond).Should(Receive(&r2got1))
		Eventually(received2, 100*time.Millisecond).Should(Receive(&r2got2))
		Eventually(received2, 100*time.Millisecond).Should(Receive(&r2got3))

		Expect(string(r2got1.Name)).To(Equal("test.event.1"))
		Expect(r2got1.Payload).To(Equal("test-payload.1"))
		Expect(string(r2got2.Name)).To(Equal("test.event.2"))
		Expect(r2got2.Payload).To(Equal("test-payload.2"))
		Expect(string(r2got3.Name)).To(Equal("test.event.3"))
		Expect(r2got3.Payload).To(Equal("test-payload.3"))
	})
}, Entry("with Memory", "memory"), Entry("with Redis", "redis"))
