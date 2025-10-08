package repository_test

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/mauriciorobertodev/whappy-go/internal/app/logger"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/instance"
	"github.com/mauriciorobertodev/whappy-go/internal/fake"
	"github.com/mauriciorobertodev/whappy-go/internal/infra/config"
	"github.com/mauriciorobertodev/whappy-go/internal/infra/database"
	"github.com/mauriciorobertodev/whappy-go/internal/infra/repository"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = DescribeTableSubtree("InstanceRepository", func(driver string) {
	Expect(godotenv.Load("./../../../.env")).ToNot(HaveOccurred())
	config.LoadLoggers(logger.LevelNone)

	var (
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

		instRepo = repository.NewInstanceRepository(db)
	})

	AfterEach(func() {
		db.Close()
	})

	It("should insert an instance", func() {
		id := "instance-id"
		name := "Test Instance"

		inst := fake.InstanceFactory().WithID(id).WithName(name).Create()

		err := instRepo.Insert(inst)
		Expect(err).To(BeNil())

		fetched, err := instRepo.Get(instance.WhereID(id))
		Expect(err).To(BeNil())
		Expect(fetched).ToNot(BeNil())
		Expect(fetched.ID).To(Equal(id))
		Expect(fetched.Name).To(Equal(name))
	})

	It("should insert many instances", func() {
		instances := []*instance.Instance{
			fake.InstanceFactory().WithID("id-1").Connected().Create(),
			fake.InstanceFactory().WithID("id-2").Connected().Create(),
			fake.InstanceFactory().WithID("id-3").Connected().Create(),
		}

		err := instRepo.InsertMany(instances)
		Expect(err).To(BeNil())

		for _, inst := range instances {
			fetched, err := instRepo.Get(instance.WhereID(inst.ID))
			Expect(err).To(BeNil())
			Expect(fetched).ToNot(BeNil())
			Expect(fetched.ID).To(Equal(inst.ID))
			Expect(fetched.Name).To(Equal(inst.Name))
			Expect(fetched.Phone).To(Equal(inst.Phone))
			Expect(fetched.JID).To(Equal(inst.JID))
			Expect(fetched.LID).To(Equal(inst.LID))
			Expect(fetched.Device).To(Equal(inst.Device))
			Expect(fetched.Status).To(Equal(inst.Status))
			Expect(*fetched.LastLoginAt).To(Equal(*inst.LastLoginAt))
			Expect(*fetched.LastConnectedAt).To(Equal(*inst.LastConnectedAt))
			Expect(fetched.BannedAt).To(BeNil())
			Expect(fetched.BanExpiresAt).To(BeNil())
			Expect(fetched.CreatedAt).To(Equal(inst.CreatedAt))
			Expect(fetched.UpdatedAt).To(Equal(inst.UpdatedAt))
		}
	})

	It("should update an instance", func() {
		id := "instance-id"
		name := "Test Instance"

		inst := fake.InstanceFactory().WithID(id).WithName(name).Create()

		err := instRepo.Insert(inst)
		Expect(err).To(BeNil())

		fetched, err := instRepo.Get(instance.WhereID(id))
		Expect(err).To(BeNil())
		Expect(fetched).ToNot(BeNil())
		Expect(fetched.ID).To(Equal(id))
		Expect(fetched.Name).To(Equal(name))

		// Update
		newName := "Updated Instance"
		fetched.Name = newName
		err = instRepo.Update(fetched)
		Expect(err).To(BeNil())

		updated, err := instRepo.Get(instance.WhereID(id))
		Expect(err).To(BeNil())
		Expect(updated).ToNot(BeNil())
		Expect(updated.ID).To(Equal(id))
		Expect(updated.Name).To(Equal(newName))
	})

	It("should delete an instance", func() {
		inst := fake.InstanceFactory().Create()

		err := instRepo.Insert(inst)
		Expect(err).To(BeNil())

		fetched, err := instRepo.Get(instance.WhereID(inst.ID))
		Expect(err).To(BeNil())
		Expect(fetched).ToNot(BeNil())
		Expect(fetched.ID).To(Equal(inst.ID))

		// Delete
		err = instRepo.Delete(instance.WhereID(fetched.ID))
		Expect(err).To(BeNil())

		deleted, err := instRepo.Get(instance.WhereID(fetched.ID))
		Expect(err).To(Equal(instance.ErrInstanceNotFound))
		Expect(deleted).To(BeNil())
	})

	It("should get instance by filter", func() {
		err := instRepo.InsertMany([]*instance.Instance{
			fake.InstanceFactory().WithID("id-1").Create(),
			fake.InstanceFactory().WithID("id-2").WithName("Instance Two").Create(),
			fake.InstanceFactory().WithID("id-3").WithJID("jid-3").Create(),
			fake.InstanceFactory().WithID("id-4").WithLID("lid-4").Create(),
			fake.InstanceFactory().WithID("id-5").WithPhone("phone-5").Create(),
			fake.InstanceFactory().WithID("id-6").WithDevice("device-jid-6").Create(),
			fake.InstanceFactory().WithID("id-7").WithStatus("status-7").Create(),
		})

		Expect(err).To(BeNil())

		fetched, err := instRepo.Get(instance.WhereID("id-1"))
		Expect(err).To(BeNil())
		Expect(fetched).ToNot(BeNil())
		Expect(fetched.ID).To(Equal("id-1"))

		fetched, err = instRepo.Get(instance.WhereName("Instance Two"))
		Expect(err).To(BeNil())
		Expect(fetched).ToNot(BeNil())
		Expect(fetched.ID).To(Equal("id-2"))
		Expect(fetched.Name).To(Equal("Instance Two"))

		fetched, err = instRepo.Get(instance.WhereJID("jid-3"))
		Expect(err).To(BeNil())
		Expect(fetched).ToNot(BeNil())
		Expect(fetched.ID).To(Equal("id-3"))
		Expect(fetched.JID).To(Equal("jid-3"))

		fetched, err = instRepo.Get(instance.WhereLID("lid-4"))
		Expect(err).To(BeNil())
		Expect(fetched).ToNot(BeNil())
		Expect(fetched.ID).To(Equal("id-4"))
		Expect(fetched.LID).To(Equal("lid-4"))

		fetched, err = instRepo.Get(instance.WherePhone("phone-5"))
		Expect(err).To(BeNil())
		Expect(fetched).ToNot(BeNil())
		Expect(fetched.ID).To(Equal("id-5"))
		Expect(fetched.Phone).To(Equal("phone-5"))

		fetched, err = instRepo.Get(instance.WhereDevice("device-jid-6"))
		Expect(err).To(BeNil())
		Expect(fetched).ToNot(BeNil())
		Expect(fetched.ID).To(Equal("id-6"))
		Expect(fetched.Device).To(Equal("device-jid-6"))

		// TODO: put new filters here
	})

	It("should return error when instance not found", func() {
		fetched, err := instRepo.Get(instance.WhereID("non-existing-id"))
		Expect(err).To(Equal(instance.ErrInstanceNotFound))
		Expect(fetched).To(BeNil())
	})

	It("should list instances", func() {
		err := instRepo.InsertMany([]*instance.Instance{
			fake.InstanceFactory().WithID("id-1").Create(),
			fake.InstanceFactory().WithID("id-2").Create(),
			fake.InstanceFactory().WithID("id-3").Create(),
		})

		Expect(err).To(BeNil())

		list, err := instRepo.List()
		Expect(err).To(BeNil())
		Expect(list).To(HaveLen(3))
	})

	It("should return empty list when no instances", func() {
		list, err := instRepo.List()
		Expect(err).To(BeNil())
		Expect(list).To(HaveLen(0))
	})

	It("should list instances with filters", func() {
		err := instRepo.InsertMany([]*instance.Instance{
			fake.InstanceFactory().WithID("id-1").WithName("Instance X").WithStatus("CONNECTED").Create(),
			fake.InstanceFactory().WithID("id-2").WithName("Instance X").WithStatus("CONNECTING").Create(),
			fake.InstanceFactory().WithID("id-3").WithName("Instance X").WithStatus("CONNECTED").Create(),
		})

		Expect(err).To(BeNil())

		list, err := instRepo.List(instance.WhereID("id-1"))
		Expect(err).To(BeNil())
		Expect(list).To(HaveLen(1))
		Expect(list[0].ID).To(Equal("id-1"))

		list, err = instRepo.List(instance.WhereName("Instance X"))
		Expect(err).To(BeNil())
		Expect(list).To(HaveLen(3))

		list, err = instRepo.List(instance.WhereStatus("CONNECTED"))
		Expect(err).To(BeNil())
		Expect(list).To(HaveLen(2))

		list, err = instRepo.List(instance.WhereStatus("CONNECTING"))
		Expect(err).To(BeNil())
		Expect(list).To(HaveLen(1))
		Expect(list[0].ID).To(Equal("id-2"))

		list, err = instRepo.List(instance.Limit(2))
		Expect(err).To(BeNil())
		Expect(list).To(HaveLen(2))
	})

	It("should list ordered by created_at desc", func() {
		err := instRepo.InsertMany([]*instance.Instance{
			fake.InstanceFactory().WithID("id-1").WithCreatedAt(time.Now().Add(-1 * time.Second)).Create(),
			fake.InstanceFactory().WithID("id-2").WithCreatedAt(time.Now().Add(1 * time.Second)).Create(),
			fake.InstanceFactory().WithID("id-3").WithCreatedAt(time.Now()).Create(),
		})

		Expect(err).To(BeNil())

		list, err := instRepo.List()
		Expect(err).To(BeNil())
		Expect(list).To(HaveLen(3))
		Expect(list[0].ID).To(Equal("id-2"))
		Expect(list[1].ID).To(Equal("id-3"))
		Expect(list[2].ID).To(Equal("id-1"))
	})

	It("should count instances", func() {
		err := instRepo.InsertMany([]*instance.Instance{
			fake.InstanceFactory().WithID("id-1").WithName("Instance X").WithStatus("CONNECTED").Create(),
			fake.InstanceFactory().WithID("id-2").WithName("Instance X").WithStatus("CONNECTING").Create(),
			fake.InstanceFactory().WithID("id-3").WithName("Instance X").WithStatus("CONNECTED").Create(),
		})

		Expect(err).To(BeNil())

		count := instRepo.Count()
		Expect(count).To(Equal(3))

		count = instRepo.Count(instance.WhereID("id-1"))
		Expect(count).To(Equal(1))

		count = instRepo.Count(instance.WhereName("Instance X"))
		Expect(count).To(Equal(3))

		count = instRepo.Count(instance.WhereStatus("CONNECTED"))
		Expect(count).To(Equal(2))

		count = instRepo.Count(instance.WhereStatus("CONNECTING"))
		Expect(count).To(Equal(1))

		count = instRepo.Count(instance.WhereStatus("NON-EXISTING"))
		Expect(count).To(Equal(0))
	})
}, Entry("with SQLite", "sqlite"), Entry("with Postgres", "postgres"))
