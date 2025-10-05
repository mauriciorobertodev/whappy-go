package cache_test

import (
	"testing"
	"time"

	"github.com/joho/godotenv"
	intf "github.com/mauriciorobertodev/whappy-go/internal/app/cache"
	"github.com/mauriciorobertodev/whappy-go/internal/app/logger"
	"github.com/mauriciorobertodev/whappy-go/internal/infra/cache"
	"github.com/mauriciorobertodev/whappy-go/internal/infra/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCaches(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Caches Suite")
}

var _ = DescribeTableSubtree("Cache", func(driver string) {
	Expect(godotenv.Load("./../../../.env")).ToNot(HaveOccurred())
	config.LoadLoggers(logger.LevelNone)

	var (
		c intf.Cache
	)

	BeforeEach(func() {
		if driver == "memory" {
			c = cache.New(&config.CacheConfig{
				Driver: config.CacheDriverInMemory,
			})
		}

		if driver == "redis" {
			c = cache.New(&config.CacheConfig{
				Driver: config.CacheDriverRedis,
				RedisConfig: &config.RedisConfig{
					Host: config.GetEnvString("REDIS_HOST", ""),
					Port: config.GetEnvInt("REDIS_PORT", 0),
				},
			})
		}
	})

	AfterEach(func() {
		err := c.Flush()
		Expect(err).To(BeNil())
	})

	It("should set and get a value", func() {
		key := "test:key"
		value := []byte("value")

		err := c.Set(key, value, intf.DefaultTTL)
		Expect(err).To(BeNil())

		got, err := c.Get(key)
		Expect(err).To(BeNil())
		Expect(got).To(Equal(value))
	})

	It("should return ErrNotFound for non-existing key", func() {
		_, err := c.Get("non-existing-key")
		Expect(err).ToNot(BeNil())
		Expect(err).To(Equal(intf.ErrNotFound))
	})

	It("should delete a key", func() {
		key := "test:key"
		value := []byte("value")

		err := c.Set(key, value, intf.DefaultTTL)
		Expect(err).To(BeNil())

		got, err := c.Get(key)
		Expect(err).To(BeNil())
		Expect(got).To(Equal(value))

		err = c.Delete(key)
		Expect(err).To(BeNil())

		_, err = c.Get(key)
		Expect(err).ToNot(BeNil())
		Expect(err).To(Equal(intf.ErrNotFound))
	})

	It("should set a value forever", func() {
		key := "test:forever"
		value := []byte("forever value")

		err := c.Forever(key, value)
		Expect(err).To(BeNil())

		got, err := c.Get(key)
		Expect(err).To(BeNil())
		Expect(got).To(Equal(value))
	})

	It("should expire a key", func() {
		key := "test:expire"
		value := []byte("temp value")

		err := c.Set(key, value, 1*time.Second)
		Expect(err).To(BeNil())

		got, err := c.Get(key)
		Expect(err).To(BeNil())
		Expect(got).To(Equal(value))

		time.Sleep(1 * time.Second)

		_, err = c.Get(key)
		Expect(err).ToNot(BeNil())
		Expect(err).To(Equal(intf.ErrNotFound))
	})

	It("should handle zero TTL as forever", func() {
		key := "test:zero-ttl"
		value := []byte("zero ttl value")

		err := c.Set(key, value, 0)
		Expect(err).To(BeNil())

		got, err := c.Get(key)
		Expect(err).To(BeNil())
		Expect(got).To(Equal(value))

		time.Sleep(1 * time.Second)

		got, err = c.Get(key)
		Expect(err).To(BeNil())
		Expect(got).To(Equal(value))
	})

	It("should overwrite existing key", func() {
		key := "test:overwrite"
		value1 := []byte("value1")
		value2 := []byte("value2")

		err := c.Set(key, value1, intf.DefaultTTL)
		Expect(err).To(BeNil())

		got, err := c.Get(key)
		Expect(err).To(BeNil())
		Expect(got).To(Equal(value1))

		err = c.Set(key, value2, intf.DefaultTTL)
		Expect(err).To(BeNil())

		got, err = c.Get(key)
		Expect(err).To(BeNil())
		Expect(got).To(Equal(value2))
	})

	It("should handle multiple keys", func() {
		keys := []string{"key1", "key2", "key3"}
		values := [][]byte{[]byte("value1"), []byte("value2"), []byte("value3")}

		for i, key := range keys {
			err := c.Set(key, values[i], intf.DefaultTTL)
			Expect(err).To(BeNil())
		}

		for i, key := range keys {
			got, err := c.Get(key)
			Expect(err).To(BeNil())
			Expect(got).To(Equal(values[i]))
		}
	})
}, Entry("with Memory", "memory"), Entry("with Redis", "redis"))
