package config

type StorageDriver string

const (
	StorageDriverLocal StorageDriver = "local"
	StorageDriverS3    StorageDriver = "s3"
	StorageDriverNone  StorageDriver = "none"
)

func (d StorageDriver) IsValid() bool {
	switch d {
	case StorageDriverLocal, StorageDriverS3, StorageDriverNone:
		return true
	default:
		return false
	}
}

type StorageConfig struct {
	Driver StorageDriver
	// Just for local storage
	Path string
	// Just for S3 storage
	Key       string
	Secret    string
	Region    string
	Bucket    string
	Endpoint  string
	PathStyle bool
	// Shared
	URL string
}

func LoadStorageConfig() *StorageConfig {
	driver := GetEnvString("STORAGE_DRIVER", "local")

	if driver == "" || driver == "none" {
		return nil
	}

	url := GetEnvString("S3_URL", "")

	if driver == "local" {
		url = GetEnvURL("APP_URL", "")
	}

	conf := &StorageConfig{
		Driver:    StorageDriver(GetEnvString("STORAGE_DRIVER", "local")), // local, s3
		Path:      GetEnvString("STORAGE_PATH", "/storage"),
		Key:       GetEnvString("S3_KEY", ""),
		Secret:    GetEnvString("S3_SECRET", ""),
		Region:    GetEnvString("S3_REGION", ""),
		Bucket:    GetEnvString("S3_BUCKET", ""),
		Endpoint:  GetEnvURL("S3_ENDPOINT", ""),
		PathStyle: GetEnvBool("S3_PATH_STYLE", false),
		URL:       url,
	}

	if !conf.IsConfigured() {
		return nil
	}

	return conf
}

func (c *StorageConfig) IsConfigured() bool {
	if c == nil {
		return false
	}

	if !c.Driver.IsValid() {
		panic("Invalid Storage Driver: " + string(c.Driver))
	}

	if c.Driver == StorageDriverS3 {
		return c.Key != "" && c.Secret != "" && c.Region != "" && c.Bucket != ""
	}

	return c.Path != "" || c.URL != ""
}

func (c *StorageConfig) IsS3() bool {
	return c != nil && c.Driver == StorageDriverS3
}

func (c *StorageConfig) IsLocal() bool {
	return c != nil && c.Driver == StorageDriverLocal
}
