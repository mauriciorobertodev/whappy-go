package cache_test

import (
	"github.com/joho/godotenv"
	intf "github.com/mauriciorobertodev/whappy-go/internal/app/cache"
	"github.com/mauriciorobertodev/whappy-go/internal/app/logger"
	"github.com/mauriciorobertodev/whappy-go/internal/infra/cache"
	"github.com/mauriciorobertodev/whappy-go/internal/infra/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = DescribeTableSubtree("Cache helpers", func(driver string) {
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
					Host: config.GetEnvString("REDIS_HOST", "localhost"),
					Port: config.GetEnvInt("REDIS_PORT", 6379),
				},
			})
		}
	})

	AfterEach(func() {
		err := c.Flush()
		Expect(err).To(BeNil())
	})

	It("should set and get a value using helpers", func() {
		key := "test:key"
		value := "value"

		err := cache.Set(c, key, value, intf.DefaultTTL)
		Expect(err).To(BeNil())

		got, err := cache.Get[string](c, key)
		Expect(err).To(BeNil())
		Expect(got).To(Equal(value))

		keyNoTTL := "test:key:no-ttl"
		err = cache.Forever(c, keyNoTTL, value)
		Expect(err).To(BeNil())

		got, err = cache.Get[string](c, keyNoTTL)
		Expect(err).To(BeNil())
		Expect(got).To(Equal(value))
	})

	It("should return ErrNotFound for non-existing key using helpers", func() {
		_, err := cache.Get[string](c, "non-existing-key")
		Expect(err).ToNot(BeNil())
		Expect(err).To(Equal(intf.ErrNotFound))
	})

	It("should delete a key using helpers", func() {
		key := "test:key"
		value := "value"

		err := cache.Set(c, key, value, intf.DefaultTTL)
		Expect(err).To(BeNil())

		got, err := cache.Get[string](c, key)
		Expect(err).To(BeNil())
		Expect(got).To(Equal(value))

		err = c.Delete(key)
		Expect(err).To(BeNil())

		_, err = cache.Get[string](c, key)
		Expect(err).ToNot(BeNil())
		Expect(err).To(Equal(intf.ErrNotFound))
	})

	It("should set a value forever using helpers", func() {
		key := "test:forever"
		value := "forever value"

		err := cache.Forever(c, key, value)
		Expect(err).To(BeNil())

		got, err := cache.Get[string](c, key)
		Expect(err).To(BeNil())
		Expect(got).To(Equal(value))
	})

	It("should serialize and unserialize complex types using helpers", func() {
		type Object struct {
			ID   int
			Flag bool
		}
		type NestedObject struct {
			Title string
			Value float64
			Inner Object
		}
		type Complex struct {
			Name   string
			Age    int
			Email  string
			Object Object
			Map    map[string]int
			Slice  []string
			Nested NestedObject
		}

		key := "test:complex"
		value := Complex{
			Name:  "John Doe",
			Age:   30,
			Email: "john.doe@example.com",
			Object: Object{
				ID:   1,
				Flag: true,
			},
			Map: map[string]int{
				"one": 1,
				"two": 2,
			},
			Slice: []string{"a", "b", "c"},
			Nested: NestedObject{
				Title: "Nested Title",
				Value: 3.14,
				Inner: Object{
					ID:   2,
					Flag: false,
				},
			},
		}
		err := cache.Set(c, key, value, intf.DefaultTTL)
		Expect(err).To(BeNil())

		got, err := cache.Get[Complex](c, key)
		Expect(err).To(BeNil())
		Expect(got).To(Equal(value))
	})
}, Entry("with Memory", "memory"), Entry("with Redis", "redis"))
