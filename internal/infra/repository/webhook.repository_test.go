package repository_test

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/mauriciorobertodev/whappy-go/internal/app/logger"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/instance"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/webhook"
	"github.com/mauriciorobertodev/whappy-go/internal/fake"
	"github.com/mauriciorobertodev/whappy-go/internal/infra/config"
	"github.com/mauriciorobertodev/whappy-go/internal/infra/database"
	"github.com/mauriciorobertodev/whappy-go/internal/infra/repository"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = DescribeTableSubtree("TokenRepository", func(driver string) {
	Expect(godotenv.Load("./../../../.env")).ToNot(HaveOccurred())
	config.LoadLoggers(logger.LevelNone)

	var (
		repo     webhook.WebhookRepository
		instRepo instance.InstanceRepository
		db       *sqlx.DB
		migrator *database.Migrator
	)

	BeforeEach(func() {
		var conf config.DatabaseConfig

		if driver == "sqlite" {
			conf = config.DatabaseConfig{
				Driver: config.DatabaseDriverSQLite,
				DbName: ":memory:",
			}
		}

		if driver == "postgres" {
			conf = config.DatabaseConfig{
				Driver: config.DatabaseDriverPostgres,
				DbName: config.GetEnvString("DB_NAME", ""),
				DbUser: config.GetEnvString("DB_USER", ""),
				DbPass: config.GetEnvString("DB_PASS", ""),
				DbHost: config.GetEnvString("DB_HOST", ""),
				DbPort: config.GetEnvString("DB_PORT", ""),
			}
		}

		db = database.New(&conf)

		migrator = database.NewMigrator(db, conf.CodeDriver())

		migrator.Reset()

		repo = repository.NewWebhookRepository(db)
		instRepo = repository.NewInstanceRepository(db)
	})

	AfterEach(func() {
		db.Close()
	})

	It("should insert and find a webhook by ID", func() {
		inst := fake.InstanceFactory().WithID("instance-1").Create()
		Expect(instRepo.Insert(inst)).To(Succeed())

		w1 := fake.WebhookFactory().WithInstanceID(inst.ID).Create()

		Expect(repo.Insert(w1)).To(Succeed())

		got, err := repo.Get(webhook.WhereID(w1.ID))
		Expect(err).ToNot(HaveOccurred())
		Expect(got.ID).To(Equal(w1.ID))
		Expect(got.InstanceID).To(Equal(w1.InstanceID))
		Expect(got.URL).To(Equal(w1.URL))
		Expect(got.Active).To(Equal(w1.Active))
		Expect(got.Events).To(Equal(w1.Events))
		Expect(got.GetSecret()).To(Equal(w1.GetSecret()))
		Expect(got.CreatedAt).To(BeTemporally("~", w1.CreatedAt, time.Second))
		Expect(got.UpdatedAt).To(BeTemporally("~", w1.UpdatedAt, time.Second))
	})

	It("should get webhooks by filter", func() {
		instRepo.InsertMany([]*instance.Instance{
			fake.InstanceFactory().WithID("instance-1").Create(),
			fake.InstanceFactory().WithID("instance-2").Create(),
			fake.InstanceFactory().WithID("instance-3").Create(),
		})

		wt1 := fake.WebhookFactory().WithID("id-1").Inactive().WithInstanceID("instance-1").Create()
		wt2 := fake.WebhookFactory().WithID("id-2").Inactive().WithInstanceID("instance-2").Create()
		wt3 := fake.WebhookFactory().WithID("id-3").Active().WithInstanceID("instance-3").Create()
		err := repo.InsertMany([]*webhook.Webhook{
			wt1, wt2, wt3,
		})

		Expect(err).To(BeNil())

		fetched, err := repo.Get(webhook.WhereID("id-1"))
		Expect(err).To(BeNil())
		Expect(fetched).ToNot(BeNil())
		Expect(fetched.ID).To(Equal("id-1"))

		fetched, err = repo.Get(webhook.WhereInstanceID("instance-2"))
		Expect(err).To(BeNil())
		Expect(fetched).ToNot(BeNil())
		Expect(fetched.ID).To(Equal("id-2"))

		fetched, err = repo.Get(webhook.WhereActive(true))
		fmt.Println(fetched)

		Expect(err).To(BeNil())
		Expect(fetched).ToNot(BeNil())
		Expect(fetched.ID).To(Equal("id-3"))
	})

	It("should list webhooks by filter", func() {
		instRepo.InsertMany([]*instance.Instance{
			fake.InstanceFactory().WithID("instance-1").Create(),
			fake.InstanceFactory().WithID("instance-2").Create(),
			fake.InstanceFactory().WithID("instance-3").Create(),
			fake.InstanceFactory().WithID("instance-4").Create(),
		})

		wt1 := fake.WebhookFactory().WithID("id-1").Inactive().WithInstanceID("instance-1").Create()
		wt2 := fake.WebhookFactory().WithID("id-2").Inactive().WithInstanceID("instance-2").Create()
		wt3 := fake.WebhookFactory().WithID("id-3").Active().WithInstanceID("instance-3").Create()
		wt4 := fake.WebhookFactory().WithID("id-4").Active().WithInstanceID("instance-4").Create()
		wt5 := fake.WebhookFactory().WithID("id-5").Active().WithInstanceID("instance-4").Create()
		err := repo.InsertMany([]*webhook.Webhook{
			wt1, wt2, wt3, wt4, wt5,
		})

		Expect(err).To(BeNil())

		got, err := repo.List(webhook.WhereID("id-1"))
		Expect(err).To(BeNil())
		Expect(got).ToNot(BeNil())
		Expect(len(got)).To(Equal(1))

		got, err = repo.List(webhook.WhereInstanceID("instance-2"))
		Expect(err).To(BeNil())
		Expect(got).ToNot(BeNil())
		Expect(len(got)).To(Equal(1))

		got, err = repo.List(webhook.WhereInstanceID("instance-4"))
		Expect(err).To(BeNil())
		Expect(got).ToNot(BeNil())
		Expect(len(got)).To(Equal(2))

		got, err = repo.List()
		Expect(err).To(BeNil())
		Expect(got).ToNot(BeNil())
		Expect(len(got)).To(Equal(5))
	})

	It("should return nil when webhook is not found by ID", func() {
		got, err := repo.Get(webhook.WhereID("not-exist"))
		Expect(err).NotTo(HaveOccurred())
		Expect(got).To(BeNil())
	})

	It("should insert and find webhooks by InstanceID", func() {
		i1 := fake.InstanceFactory().WithID("instance-1").Create()
		i2 := fake.InstanceFactory().WithID("instance-2").Create()

		err := instRepo.InsertMany([]*instance.Instance{i1, i2})
		Expect(err).To(BeNil())

		w1 := fake.WebhookFactory().WithID("w1").WithInstanceID("instance-1").WithCreatedAt(time.Now().Add(time.Minute * 1)).WithSecret("secret-1").Create()
		w2 := fake.WebhookFactory().WithID("w2").WithInstanceID("instance-1").WithCreatedAt(time.Now().Add(time.Minute * 2)).WithSecret("secret-2").Create()
		w3 := fake.WebhookFactory().WithID("w3").WithInstanceID("instance-2").WithCreatedAt(time.Now().Add(time.Minute * 3)).WithSecret("secret-3").Create()

		err = repo.InsertMany([]*webhook.Webhook{w1, w2, w3})
		Expect(err).To(BeNil())

		got, err := repo.Get(webhook.WhereInstanceID("instance-1"))

		Expect(err).ToNot(HaveOccurred())
		Expect(got.ID).To(Equal(w2.ID)) // lasted
		Expect(got.InstanceID).To(Equal(w2.InstanceID))
		Expect(got.URL).To(Equal(w2.URL))
		Expect(got.Active).To(Equal(w2.Active))
		Expect(got.Events).To(Equal(w2.Events))
		Expect(got.GetSecret()).To(Equal(w2.GetSecret()))
		Expect(got.CreatedAt).To(BeTemporally("~", w2.CreatedAt, time.Second))
		Expect(got.UpdatedAt).To(BeTemporally("~", w2.UpdatedAt, time.Second))
	})

	It("should not allow duplicate IDs", func() {
		inst := fake.InstanceFactory().WithID("instance-1").Create()
		Expect(instRepo.Insert(inst)).To(Succeed())

		w1 := fake.WebhookFactory().WithInstanceID(inst.ID).Create()
		Expect(repo.Insert(w1)).To(Succeed())

		Expect(repo.Insert(w1)).To(HaveOccurred())
	})

	It("should delete a webhook", func() {
		inst := fake.InstanceFactory().WithID("instance-1").Create()
		Expect(instRepo.Insert(inst)).To(Succeed())

		w1 := fake.WebhookFactory().WithInstanceID(inst.ID).Create()
		Expect(repo.Insert(w1)).To(Succeed())

		err := repo.Delete(webhook.WhereID(w1.ID))
		Expect(err).ToNot(HaveOccurred())

		w1, err = repo.Get(webhook.WhereID(w1.ID))
		Expect(err).NotTo(HaveOccurred())
		Expect(w1).To(BeNil())
	})

	It("should return nil when deleting non-existent webhook", func() {
		err := repo.Delete(webhook.WhereID("not-exist"))
		Expect(err).To(BeNil())
	})
}, Entry("with SQLite", "sqlite"), Entry("with Postgres", "postgres"))
