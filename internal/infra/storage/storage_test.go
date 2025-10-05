package storage_test

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/joho/godotenv"
	"github.com/mauriciorobertodev/whappy-go/internal/app/logger"
	intf "github.com/mauriciorobertodev/whappy-go/internal/app/storage"
	"github.com/mauriciorobertodev/whappy-go/internal/infra/config"
	"github.com/mauriciorobertodev/whappy-go/internal/infra/storage"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestStorages(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Storages Suite")
}

var _ = DescribeTableSubtree("Storage", func(driver string) {
	Expect(godotenv.Load("./../../../.env")).ToNot(HaveOccurred())
	config.LoadLoggers(logger.LevelNone)

	var (
		ctx      context.Context
		testKey  string
		testData []byte
		store    intf.Storage
	)

	BeforeEach(func() {
		ctx = context.Background()

		testKey = "file.txt"
		testData = []byte("Hello LocalStorage!")

		if driver == "s3" {
			store = storage.New(&config.StorageConfig{
				Driver:    config.StorageDriverS3,
				Key:       config.GetEnvString("S3_KEY", ""),
				Secret:    config.GetEnvString("S3_SECRET", ""),
				Region:    config.GetEnvString("S3_REGION", ""),
				Bucket:    config.GetEnvString("S3_BUCKET", ""),
				Endpoint:  config.GetEnvString("S3_ENDPOINT", ""),
				URL:       config.GetEnvString("S3_URL", ""),
				PathStyle: true,
			})
		}

		if driver == "local" {
			store = storage.New(&config.StorageConfig{
				Driver: config.StorageDriverLocal,
				Path:   "./../../../" + config.GetEnvString("STORAGE_PATH", "") + "/tests",
				URL:    config.GetEnvURL("APP_URL", ""),
			})
		}
	})

	AfterEach(func() {
		_ = store.Delete(ctx, testKey)
	})

	It("should report healthy status", func() {
		Expect(store.Healthy(ctx)).ToNot(HaveOccurred())
	})

	It("should save and load a file", func() {
		exists, err := store.Exists(ctx, testKey)
		Expect(err).ToNot(HaveOccurred())
		Expect(exists).To(BeFalse())

		Expect(store.Save(ctx, testKey, bytes.NewReader(testData))).To(Succeed())

		reader, err := store.Load(ctx, testKey)
		Expect(err).ToNot(HaveOccurred())
		defer reader.Close()

		data, err := io.ReadAll(reader)
		Expect(err).ToNot(HaveOccurred())
		Expect(data).To(Equal(testData))
	})

	It("should put and get data as bytes", func() {
		exists, err := store.Exists(ctx, testKey)
		Expect(err).ToNot(HaveOccurred())
		Expect(exists).To(BeFalse())

		Expect(store.Put(ctx, testKey, testData)).To(Succeed())

		data, err := store.Get(ctx, testKey)
		Expect(err).ToNot(HaveOccurred())
		Expect(data).To(Equal(testData))
	})

	It("should check existence of a file", func() {
		exists, err := store.Exists(ctx, testKey)
		Expect(err).ToNot(HaveOccurred())
		Expect(exists).To(BeFalse())

		Expect(store.Put(ctx, testKey, testData)).To(Succeed())

		exists, err = store.Exists(ctx, testKey)
		Expect(err).ToNot(HaveOccurred())
		Expect(exists).To(BeTrue())

		Expect(store.Delete(ctx, testKey)).To(Succeed())

		exists, err = store.Exists(ctx, testKey)
		Expect(err).ToNot(HaveOccurred())
		Expect(exists).To(BeFalse())
	})

	It("should delete a file", func() {
		Expect(store.Put(ctx, testKey, testData)).To(Succeed())

		exists, err := store.Exists(ctx, testKey)
		Expect(err).ToNot(HaveOccurred())
		Expect(exists).To(BeTrue())

		Expect(store.Delete(ctx, testKey)).To(Succeed())

		exists, err = store.Exists(ctx, testKey)
		Expect(err).ToNot(HaveOccurred())
		Expect(exists).To(BeFalse())
	})

	It("should generate a correct URL", func() {
		url, err := store.URL(ctx, testKey)
		Expect(err).ToNot(HaveOccurred())
		Expect(url).To(ContainSubstring(testKey))
	})
}, Entry("with Local", "local"), Entry("with S3", "s3"))
