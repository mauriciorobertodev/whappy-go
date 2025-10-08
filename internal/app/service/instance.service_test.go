package service_test

import (
	"github.com/google/uuid"
	"github.com/mauriciorobertodev/whappy-go/internal/app"
	"github.com/mauriciorobertodev/whappy-go/internal/app/input"
	"github.com/mauriciorobertodev/whappy-go/internal/app/logger"
	"github.com/mauriciorobertodev/whappy-go/internal/app/service"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/instance"
	"github.com/mauriciorobertodev/whappy-go/internal/fake"
	"github.com/mauriciorobertodev/whappy-go/internal/infra/config"
	"github.com/mauriciorobertodev/whappy-go/internal/infra/database"
	"github.com/mauriciorobertodev/whappy-go/internal/infra/registry"
	"github.com/mauriciorobertodev/whappy-go/internal/infra/repository"
	"github.com/mauriciorobertodev/whappy-go/internal/infra/token"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Instance Service", func() {
	config.LoadLoggers(logger.LevelNone)

	db := database.New(&config.DatabaseConfig{
		Driver: config.DatabaseDriverSQLite,
		DbName: "test",
	})

	tokenRepo := repository.NewTokenRepository(db)
	instRepo := repository.NewInstanceRepository(db)

	instRegistry := registry.NewInMemoryInstanceRegistry()

	generator := token.NewGenerator()
	hasher := token.NewHasher(&config.TokenConfig{Hasher: config.HasherSimple})
	bus := fake.NewFakeEventBus()
	ca := fake.NewFakeCache()

	tokenService := service.NewTokenService(tokenRepo, hasher, generator, bus, ca)
	instanceService := service.NewInstanceService(tokenService, instRepo, instRegistry, bus)

	migrator := database.NewMigrator(db, db.DriverName())

	BeforeEach(func() {
		migrator.Reset()
		ca.Flush()
		instRegistry.Clear()
		bus.Clear()
	})

	Describe("CreateInstance", func() {
		It("should create an instance", func() {
			Expect(instRepo.Count()).To(Equal(0))
			Expect(tokenRepo.Count()).To(Equal(0))

			inst, tok, appErr := instanceService.Create(GinkgoT().Context(), input.CreateInstance{
				Name: "My Instance",
			})

			Expect(appErr).To(BeNil())
			Expect(inst).NotTo(BeNil())
			Expect(tok).NotTo(BeNil())

			Expect(instRepo.Count()).To(Equal(1))
			Expect(tokenRepo.Count()).To(Equal(1))

			savedInst, err := instRepo.Get(instance.WhereID(inst.ID))
			Expect(err).To(BeNil())
			Expect(savedInst).NotTo(BeNil())
			Expect(savedInst.ID).To(Equal(inst.ID))

			savedTok, err := tokenRepo.FindByInstanceID(inst.ID)
			Expect(err).To(BeNil())
			Expect(savedTok).NotTo(BeNil())

			// Deve ter disparado o evento de criação
			Expect(bus.Published()).To(HaveLen(1))
			Expect(bus.HasPublished(instance.EventCreated)).To(BeTrue())
		})

		It("should return error when try to create instance with empty name", func() {
			Expect(instRepo.Count()).To(Equal(0))

			inst, tok, appErr := instanceService.Create(GinkgoT().Context(), input.CreateInstance{
				Name: "",
			})

			Expect(appErr).NotTo(BeNil())
			Expect(appErr.Err).To(Equal(instance.ErrNameTooShort))
			Expect(appErr.Code).To(Equal(app.CodeInstanceNameTooShort))
			Expect(inst).To(BeNil())
			Expect(tok).To(BeNil())

			Expect(instRepo.Count()).To(Equal(0))

			// Should not have published any event
			Expect(bus.Published()).To(HaveLen(0))
		})

		It("should return error when try to create instance with too long name", func() {
			Expect(instRepo.Count()).To(Equal(0))

			longName := ""
			for i := 0; i <= instance.MaxNameLength; i++ {
				longName += "a"
			}

			inst, tok, appErr := instanceService.Create(GinkgoT().Context(), input.CreateInstance{
				Name: longName,
			})

			Expect(appErr).NotTo(BeNil())
			Expect(appErr.Err).To(Equal(instance.ErrNameTooLong))
			Expect(appErr.Code).To(Equal(app.CodeInstanceNameTooLong))
			Expect(inst).To(BeNil())
			Expect(tok).To(BeNil())

			Expect(instRepo.Count()).To(Equal(0))

			// Should not have published any event
			Expect(bus.Published()).To(HaveLen(0))
		})
	})

	Describe("ListInstances", func() {
		It("should list instances", func() {
			Expect(instRepo.Count()).To(Equal(0))

			// Create some instances
			for i := 1; i <= 3; i++ {
				_, _, appErr := instanceService.Create(GinkgoT().Context(), input.CreateInstance{
					Name: "Instance " + uuid.New().String(),
				})
				Expect(appErr).To(BeNil())
			}

			Expect(instRepo.Count()).To(Equal(3))

			instances, appErr := instanceService.List(GinkgoT().Context())
			Expect(appErr).To(BeNil())
			Expect(instances).NotTo(BeNil())
			Expect(len(instances)).To(Equal(3))
		})
	})

	Describe("GetInstance", func() {
		It("should get instance by id", func() {
			Expect(instRepo.Count()).To(Equal(0))

			inst, _, appErr := instanceService.Create(GinkgoT().Context(), input.CreateInstance{
				Name: "My Instance",
			})
			Expect(appErr).To(BeNil())
			Expect(inst).NotTo(BeNil())

			Expect(instRepo.Count()).To(Equal(1))

			gotInst, appErr := instanceService.Get(GinkgoT().Context(), input.GetInstance{
				ID: inst.ID,
			})
			Expect(appErr).To(BeNil())
			Expect(gotInst).NotTo(BeNil())
			Expect(gotInst.ID).To(Equal(inst.ID))
		})

		It("should return error when try to get instance with invalid id", func() {
			gotInst, appErr := instanceService.Get(GinkgoT().Context(), input.GetInstance{
				ID: "invalid-uuid",
			})
			Expect(appErr).NotTo(BeNil())
			Expect(appErr.Err).To(Equal(instance.ErrInstanceNotFound))
			Expect(appErr.Code).To(Equal(app.CodeInstanceNotFound))
			Expect(gotInst).To(BeNil())
		})

		It("should return error when try to get instance that does not exist", func() {
			gotInst, appErr := instanceService.Get(GinkgoT().Context(), input.GetInstance{
				ID: uuid.New().String(),
			})
			Expect(appErr).NotTo(BeNil())
			Expect(appErr.Err).To(Equal(instance.ErrInstanceNotFound))
			Expect(appErr.Code).To(Equal(app.CodeInstanceNotFound))
			Expect(gotInst).To(BeNil())
		})
	})
})
