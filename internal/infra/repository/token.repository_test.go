package repository_test

import (
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/mauriciorobertodev/whappy-go/internal/app/logger"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/instance"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/token"
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
		repo     token.TokenRepository
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

		repo = repository.NewTokenRepository(db)
		instRepo = repository.NewInstanceRepository(db)
	})

	AfterEach(func() {
		db.Close()
	})

	It("should insert and find a token by ID", func() {
		inst := fake.InstanceFactory().WithID("instance-1").Create()
		Expect(instRepo.Insert(inst)).To(Succeed())

		tok := fake.TokenFactory().WithInstanceID(inst.ID).Create()

		Expect(repo.Insert(tok)).To(Succeed())

		got, err := repo.FindByID(tok.ID)
		Expect(err).ToNot(HaveOccurred())
		Expect(got).To(Equal(tok))
	})

	It("should return nil when token is not found by ID", func() {
		got, err := repo.FindByID("not-exist")
		Expect(err).NotTo(HaveOccurred())
		Expect(got).To(BeNil())
	})

	It("should insert and find tokens by InstanceID", func() {
		i1 := fake.InstanceFactory().WithID("instance-1").Create()
		i2 := fake.InstanceFactory().WithID("instance-2").Create()

		Expect(instRepo.Insert(i1)).To(Succeed())
		Expect(instRepo.Insert(i2)).To(Succeed())

		t1 := fake.TokenFactory().WithID("t1").WithInstanceID("instance-1").WithHash("hash-1").Create()
		t2 := fake.TokenFactory().WithID("t2").WithInstanceID("instance-1").WithHash("hash-2").Create()
		t3 := fake.TokenFactory().WithID("t3").WithInstanceID("instance-2").WithHash("hash-3").Create()

		Expect(repo.Insert(t1)).To(Succeed())
		Expect(repo.Insert(t2)).To(Succeed())
		Expect(repo.Insert(t3)).To(Succeed())

		tokens, err := repo.FindByInstanceID("instance-1")

		Expect(err).ToNot(HaveOccurred())
		Expect(tokens).To(HaveLen(2))
		Expect(tokens).To(ContainElements(t1, t2))
	})

	It("should not allow duplicate IDs", func() {
		inst := fake.InstanceFactory().WithID("instance-1").Create()
		Expect(instRepo.Insert(inst)).To(Succeed())

		tok := fake.TokenFactory().WithInstanceID(inst.ID).Create()
		Expect(repo.Insert(tok)).To(Succeed())

		Expect(repo.Insert(tok)).To(HaveOccurred())
	})

	It("should delete a token", func() {
		inst := fake.InstanceFactory().WithID("instance-1").Create()
		Expect(instRepo.Insert(inst)).To(Succeed())

		tok := fake.TokenFactory().WithInstanceID(inst.ID).Create()
		Expect(repo.Insert(tok)).To(Succeed())

		err := repo.Delete(tok.ID)
		Expect(err).ToNot(HaveOccurred())

		tok, err = repo.FindByID(tok.ID)
		Expect(err).NotTo(HaveOccurred())
		Expect(tok).To(BeNil())
	})

	It("should return nil when deleting non-existent token", func() {
		err := repo.Delete("not-exist")
		Expect(err).To(BeNil())
	})
}, Entry("with SQLite", "sqlite"), Entry("with Postgres", "postgres"))
