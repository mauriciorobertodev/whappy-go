package repository_test

import (
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/mauriciorobertodev/whappy-go/internal/app/logger"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/file"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/instance"
	"github.com/mauriciorobertodev/whappy-go/internal/fake"
	"github.com/mauriciorobertodev/whappy-go/internal/infra/config"
	"github.com/mauriciorobertodev/whappy-go/internal/infra/database"
	"github.com/mauriciorobertodev/whappy-go/internal/infra/repository"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestRepositories(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Repositories Suite")
}

var _ = DescribeTableSubtree("FileRepository", func(driver string) {
	Expect(godotenv.Load("./../../../.env.test")).ToNot(HaveOccurred())
	config.LoadLoggers(logger.LevelNone)

	var (
		repo     file.FileRepository
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

		repo = repository.NewFileRepository(db)
		instRepo = repository.NewInstanceRepository(db)
	})

	AfterEach(func() {
		db.Close()
	})

	It("should insert a file", func() {
		id := "file-id"
		name := "test.txt"
		mime := "text/plain"
		extension := "txt"
		size := uint64(1234)
		sha256 := "sha256hash"
		path := "/files/test.txt"
		url := "http://localhost/files/test.txt"
		width := uint32(800)
		height := uint32(600)
		duration := uint32(0)
		pages := uint32(0)
		instanceID := "instance-id"

		inst := fake.InstanceFactory().WithID(instanceID).Create()
		Expect(instRepo.Insert(inst)).To(Succeed())
		// defer instRepo.Delete(instance.WithID(instanceID))

		testFile := fake.FileFactory().
			WithID(id).
			WithName(name).
			WithMime(mime).
			WithExtension(extension).
			WithSize(size).
			WithSha256(sha256).
			WithPath(path).
			WithURL(url).
			WithWidth(&width).
			WithHeight(&height).
			WithDuration(&duration).
			WithPages(&pages).
			WithInstanceID(&instanceID).
			Create()
		err := repo.Insert(testFile)

		Expect(err).ToNot(HaveOccurred())

		f, err := repo.Get(file.WhereID(testFile.ID))
		Expect(err).ToNot(HaveOccurred())
		Expect(f).ToNot(BeNil())
		Expect(f.ID).To(Equal(id))
		Expect(f.Name).To(Equal(name))
		Expect(f.Mime).To(Equal(mime))
		Expect(f.Extension).To(Equal(extension))
		Expect(f.Size).To(Equal(size))
		Expect(f.Sha256).To(Equal(sha256))
		Expect(f.Path).To(Equal(path))
		Expect(f.URL).To(Equal(url))
		Expect(f.Width).To(Equal(&width))
		Expect(f.Height).To(Equal(&height))
		Expect(f.Duration).To(Equal(&duration))
		Expect(f.Pages).To(Equal(&pages))
		Expect(f.InstanceID).To(Equal(&instanceID))
	})

	It("should update a file", func() {
		testFile := fake.FileFactory().Create()
		Expect(repo.Insert(testFile)).To(Succeed())

		testFile.Name = "updated.txt"
		Expect(repo.Update(testFile)).To(Succeed())

		f, _ := repo.Get(file.WhereID(testFile.ID))
		Expect(f.Name).To(Equal("updated.txt"))
	})

	It("should get a file by ID", func() {
		testFile := fake.FileFactory().Create()
		Expect(repo.Insert(testFile)).To(Succeed())

		f, err := repo.Get(file.WhereID(testFile.ID))
		Expect(err).ToNot(HaveOccurred())
		Expect(f).ToNot(BeNil())
		Expect(f.ID).To(Equal(testFile.ID))
	})

	It("should get a file by SHA256", func() {
		testFile := fake.FileFactory().Create()
		Expect(repo.Insert(testFile)).To(Succeed())

		f, err := repo.Get(file.WhereSha256(testFile.Sha256))
		Expect(err).ToNot(HaveOccurred())
		Expect(f).ToNot(BeNil())
		Expect(f.Sha256).To(Equal(testFile.Sha256))
	})

	It("should return error when file not found", func() {
		f, err := repo.Get(file.WhereID("non-existent-id"))
		Expect(err).To(MatchError(file.ErrFileNotFound))
		Expect(f).To(BeNil())
	})

	It("should list files with cursor", func() {
		Expect(repo.InsertMany(fake.FileFactory().CreateMany(100))).To(Succeed())

		var cursor *time.Time = nil

		for i := 1; i <= 10; i++ {
			files, err := repo.List(file.WithCursor(cursor, 11))

			Expect(err).ToNot(HaveOccurred())
			Expect(files).ToNot(BeEmpty())

			cursor = files[len(files)-1].CreatedAt

			if i == 10 {
				Expect(len(files)).To(Equal(10))
			} else {
				Expect(len(files)).To(Equal(11))
			}
		}
	})

	It("should delete a file", func() {
		testFile := fake.FileFactory().Create()
		Expect(repo.Insert(testFile)).To(Succeed())

		Expect(repo.Get(file.WhereID(testFile.ID))).ToNot(BeNil())

		Expect(repo.Delete(file.WhereID(testFile.ID))).To(Succeed())

		f, err := repo.Get(file.WhereID(testFile.ID))
		Expect(f).To(BeNil())
		Expect(err).To(MatchError(file.ErrFileNotFound))
	})

	// Thumbnail relationship
	It("should insert file with thumbnail", func() {
		thumbFile := fake.FileFactory().Image().Create()
		Expect(repo.Insert(thumbFile)).To(Succeed())
		thumbnail, err := thumbFile.ToImageFile()
		Expect(err).ToNot(HaveOccurred())

		fileWithThumbnail := fake.FileFactory().WithThumbnail(thumbnail).Create()
		Expect(repo.Insert(fileWithThumbnail)).To(Succeed())

		fetchedFile, err := repo.Get(file.WhereID(fileWithThumbnail.ID), file.WithThumbnail())
		Expect(err).ToNot(HaveOccurred())
		Expect(fetchedFile).ToNot(BeNil())
		Expect(fetchedFile.Thumbnail).ToNot(BeNil())
		Expect(fetchedFile.Thumbnail.ID).To(Equal(thumbnail.ID))
	})

	It("should list files with thumbnails and loaded thumbnails", func() {
		Expect(repo.InsertMany(fake.FileFactory().CreateMany(5))).To(Succeed())

		Expect(repo.Count()).To(Equal(uint64(5)))
		Expect(repo.Count(file.WhereHasThumbnail())).To(Equal(uint64(0)))
		Expect(repo.Count(file.WhereDoesNotHaveThumbnail())).To(Equal(uint64(5)))

		thumbnails := fake.FileFactory().Image().CreateMany(5)
		Expect(repo.InsertMany(thumbnails)).To(Succeed())

		Expect(repo.Count()).To(Equal(uint64(10)))
		Expect(repo.Count(file.WhereHasThumbnail())).To(Equal(uint64(0)))
		Expect(repo.Count(file.WhereDoesNotHaveThumbnail())).To(Equal(uint64(10)))

		for i := 0; i < 5; i++ {
			thumbnail, err := thumbnails[i].ToImageFile()
			Expect(err).ToNot(HaveOccurred())

			f := fake.FileFactory().WithThumbnail(thumbnail).Create()
			Expect(repo.Insert(f)).To(Succeed())
		}

		Expect(repo.Count()).To(Equal(uint64(15)))
		Expect(repo.Count(file.WhereHasThumbnail())).To(Equal(uint64(5)))
		Expect(repo.Count(file.WhereDoesNotHaveThumbnail())).To(Equal(uint64(10)))

		files, err := repo.List(file.WhereHasThumbnail(), file.WithThumbnail())
		Expect(err).ToNot(HaveOccurred())
		Expect(len(files)).To(Equal(5))

		for i, f := range files {
			Expect(f.Thumbnail).ToNot(BeNil())
			Expect(f.Thumbnail.ID).ToNot(BeEmpty())
			Expect(f.Thumbnail.ID).To(Equal(thumbnails[(5-i)-1].ID))
		}
	})

	It("should list files without thumbnails", func() {
		// Without thumbnail
		Expect(repo.InsertMany(fake.FileFactory().CreateMany(10))).To(Succeed())

		Expect(repo.Count()).To(Equal(uint64(10)))
		Expect(repo.Count(file.WhereHasThumbnail())).To(Equal(uint64(0)))
		Expect(repo.Count(file.WhereDoesNotHaveThumbnail())).To(Equal(uint64(10)))

		thumbnails := fake.FileFactory().Image().CreateMany(5)
		Expect(repo.InsertMany(thumbnails)).To(Succeed())

		Expect(repo.Count()).To(Equal(uint64(15)))
		Expect(repo.Count(file.WhereHasThumbnail())).To(Equal(uint64(0)))
		Expect(repo.Count(file.WhereDoesNotHaveThumbnail())).To(Equal(uint64(15)))

		for i := 0; i < 5; i++ {
			thumbnail, err := thumbnails[i].ToImageFile()
			Expect(err).ToNot(HaveOccurred())

			f := fake.FileFactory().WithThumbnail(thumbnail).Create()
			Expect(repo.Insert(f)).To(Succeed())
		}

		Expect(repo.Count()).To(Equal(uint64(20)))
		Expect(repo.Count(file.WhereHasThumbnail())).To(Equal(uint64(5)))
		Expect(repo.Count(file.WhereDoesNotHaveThumbnail())).To(Equal(uint64(15)))

		files, err := repo.List(file.WhereDoesNotHaveThumbnail())
		Expect(err).ToNot(HaveOccurred())
		Expect(len(files)).To(Equal(15))

		for _, f := range files {
			Expect(f.Thumbnail).To(BeNil())
		}
	})

	It("should list files loading thumbnails with cursor", func() {
		thumbnails := fake.FileFactory().Image().CreateMany(50)
		Expect(repo.InsertMany(thumbnails)).To(Succeed())
		for i := 0; i < 100; i++ {
			createdAt := time.Now().Add(time.Duration(i+51) * time.Minute).UTC()

			if i < 50 {
				thumbnail, err := thumbnails[i].ToImageFile()
				Expect(err).ToNot(HaveOccurred())
				Expect(repo.Insert(fake.FileFactory().WithThumbnail(thumbnail).WithCreatedAt(&createdAt).Create())).To(Succeed())
				continue
			}

			Expect(repo.Insert(fake.FileFactory().WithCreatedAt(&createdAt).Create())).To(Succeed())
		}

		Expect(repo.Count()).To(Equal(uint64(150)))
		Expect(repo.Count(file.WhereHasThumbnail())).To(Equal(uint64(50)))
		Expect(repo.Count(file.WhereDoesNotHaveThumbnail())).To(Equal(uint64(100)))

		var cursor *time.Time = nil

		for i := 1; i <= 10; i++ {
			files, err := repo.List(file.WithThumbnail(), file.WithCursor(cursor, 11))

			Expect(err).ToNot(HaveOccurred())
			Expect(files).ToNot(BeEmpty())

			cursor = files[len(files)-1].CreatedAt

			if i < 5 {
				// all here does have thumbnails
				for _, f := range files {
					Expect(f.Thumbnail).To(BeNil())
				}
			} else if i > 5 && i < 10 {
				// half of them have thumbnails
				for _, f := range files {
					Expect(f.Thumbnail).ToNot(BeNil())
				}
			} else if i > 10 {
				// all here does have thumbnails
				for _, f := range files {
					Expect(f.Thumbnail).To(BeNil())
				}
			}

			if i == 15 {
				Expect(len(files)).To(Equal(10))
			} else {
				Expect(len(files)).To(Equal(11))
			}
		}
	})
}, Entry("with SQLite", "sqlite"), Entry("with Postgres", "postgres"))
