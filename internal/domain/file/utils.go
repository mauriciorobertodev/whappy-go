package file

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"mime"
	"strings"
)

// I created a map of preferred extensions for certain MIME types,
// because mime.ExtensionsByType can return multiple options,
// and sometimes the defaults are uncommon or not ideal.
// For example, on my Mac Mini M4, "image/jpeg" returns [".jpe", ".jpeg", ".jpg"],
// where ".jpe" comes first, but I prefer ".jpg".
// This map ensures we always select the preferred extension.
//
// Note: I didn’t include MIME types that already return the expected extension,
// such as "image/png" which correctly gives ".png" first.
// If on your machine a common type doesn’t return the expected extension,
// feel free to add it here.
// — Mauricio Roberto
var PREFERRED_EXTENSIONS = map[string]string{
	// ---- Images ----
	"image/jpeg": "jpg", // jpe

	// ---- Videos ----
	"video/mp4": "mp4", // m4v

	// ---- Audios ----
	"audio/mpeg":            "mp3", // m2a
	"audio/wav":             "wav", // bin
	"application/ogg":       "ogg", // ogx
	"audio/ogg":             "ogg", // oga
	"audio/ogg codecs=opus": "ogg", // bin

	// ---- Documents ----
	"text/plain": "txt",  // conf
	"text/html":  "html", // htm

	// ---- Compacteds ----
	"application/gzip": "gz", // bin

}

func GetPreferredExtension(mimeType string) string {
	if ext, exists := PREFERRED_EXTENSIONS[mimeType]; exists {
		return ext
	}
	return ""
}

func DetectExtension(mimeType string) string {
	if ext := GetPreferredExtension(mimeType); ext != "" {
		return ext
	}

	extensions, err := mime.ExtensionsByType(mimeType)

	if err != nil || len(extensions) == 0 {
		return "bin"
	}

	return strings.TrimLeft(extensions[0], ".")
}

func DetectDimensions(data *[]byte) (*uint32, *uint32, error) {
	width, height, err := DecodeImageConfig(data)
	if err != nil {
		return nil, nil, err
	}

	return width, height, nil
}

func DecodeImageConfig(data *[]byte) (*uint32, *uint32, error) {
	cfg, _, err := image.DecodeConfig(strings.NewReader(string(*data)))
	if err != nil {
		return nil, nil, ErrCorruptedFile
	}

	width := uint32(cfg.Width)
	height := uint32(cfg.Height)

	return &width, &height, nil
}
