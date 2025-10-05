package service

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/mauriciorobertodev/whappy-go/internal/app"
	"github.com/mauriciorobertodev/whappy-go/internal/app/storage"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/file"
	"github.com/mauriciorobertodev/whappy-go/internal/utils"
)

type FileService struct {
	storage  storage.Storage
	fileRepo file.FileRepository
}

func NewFileService(storage storage.Storage, fileRepo file.FileRepository) *FileService {
	return &FileService{
		storage:  storage,
		fileRepo: fileRepo,
	}
}

func (p *FileService) SaveStream(ctx context.Context, r io.Reader, mime *string) (*file.File, error) {
	l := app.GetFileServiceLogger()

	// Lê os primeiros 512 bytes para MIME
	header := make([]byte, 512)
	n, err := io.ReadFull(r, header)
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		return nil, err
	}

	var mimeType = new(string)
	if mime == nil || *mime == "" {
		*mimeType = http.DetectContentType(header[:n])
	} else {
		*mimeType = *mime
	}

	// Cria pipe para streaming direto para Storage
	pr, pw := io.Pipe()

	hash := sha256.New()
	counter := &countWriter{}

	// Cria um MultiWriter que envia para hash e contador enquanto escreve no pipe
	writer := io.MultiWriter(pw, hash, counter)

	// Goroutine que escreve o header + resto do stream no MultiWriter
	go func() {
		defer pw.Close()
		_, _ = writer.Write(header[:n])
		_, _ = io.Copy(writer, r)
	}()

	f := file.NewFromMime(*mimeType)
	err = p.storage.Save(ctx, f.Path, pr)
	if err != nil {
		return nil, err
	}

	f.Sha256 = hex.EncodeToString(hash.Sum(nil))
	f.Size = uint64(counter.Count)
	url, err := p.storage.URL(ctx, f.Path)
	if err != nil {
		l.Error("Error getting file URL from storage", "error", err.Error())
	} else {
		f.URL = url
	}

	return f, nil
}

func (p *FileService) SaveFromBase64(ctx context.Context, b64 string, mimeType *string) (*file.File, error) {
	b64 = clearBase64(b64)
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(b64))
	return p.SaveStream(ctx, reader, mimeType)
}

func (p *FileService) SaveFromURL(ctx context.Context, url string, mimeType *string) (*file.File, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return p.SaveStream(ctx, resp.Body, mimeType)
}

func (p *FileService) SaveFromBytes(ctx context.Context, data []byte, mimeType *string) (*file.File, error) {
	return p.SaveStream(ctx, bytes.NewReader(data), mimeType)
}

func (s *FileService) LoadFrom(ctx context.Context, source string) (*file.File, io.ReadCloser, error) {
	l := app.GetFileServiceLogger()

	if source == "" {
		return nil, nil, file.ErrFileSourceEmpty
	}

	var f *file.File
	var stream io.ReadCloser
	var err error

	if utils.IsValidURL(source) {
		l.Debug("Loading file from URL", "source", source)
		f, stream, err = s.LoadFromURL(ctx, source)
		if err != nil {
			l.Error("Error loading file", "error", err)
			return nil, nil, err
		}
		return f, stream, nil
	}

	if utils.IsUUID(source) {
		l.Debug("Loading file from cache by UUID", "source", source)
		f, stream, err = s.LoadFromUploads(ctx, source)
		if err != nil {
			l.Error("Error loading file", "error", err)
			return nil, nil, err
		}
		return f, stream, nil
	}

	l.Debug("Loading file from Base64")
	f, stream, err = s.LoadFromBase64(ctx, source)
	if err != nil {
		l.Error("Error loading file", "error", err)
		return nil, nil, err
	}
	return f, stream, nil
}

func (p *FileService) LoadFromBase64(ctx context.Context, b64 string) (*file.File, io.ReadCloser, error) {
	b64 = clearBase64(b64)
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(b64))
	mimeType, stream, err := detectMimeTypeFromStream(io.NopCloser(reader))
	if err != nil {
		return nil, nil, file.ErrFileUnreachable
	}

	f := file.NewFromMime(mimeType)

	return f, stream, nil
}

func (p *FileService) LoadFromURL(ctx context.Context, url string) (*file.File, io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, nil, file.ErrFileUnreachable
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, nil, file.ErrFileUnreachable
	}

	mimeType, stream, err := detectMimeTypeFromStream(resp.Body)
	if err != nil {
		return nil, nil, file.ErrFileUnreachable
	}

	f := file.NewFromMime(mimeType)

	return f, stream, nil
}

func (p *FileService) LoadFromBytes(ctx context.Context, data []byte) (*file.File, io.ReadCloser, error) {

	mimeType, stream, err := detectMimeTypeFromStream(io.NopCloser(bytes.NewReader(data)))
	if err != nil {
		return nil, nil, file.ErrFileUnreachable
	}

	f := file.NewFromMime(mimeType)
	return f, stream, nil
}

func (p *FileService) LoadFromUploads(ctx context.Context, fileID string) (*file.File, io.ReadCloser, error) {
	f, err := p.fileRepo.Get(file.WhereID(fileID))
	if err != nil {
		return nil, nil, file.ErrFileUnreachable
	}

	if f == nil {
		return nil, nil, file.ErrFileNotFound
	}

	stream, err := p.storage.Load(ctx, f.Path)
	if err != nil {
		return nil, nil, file.ErrFileUnreachable
	}

	return f, stream, nil
}

func (s *FileService) GetFrom(ctx context.Context, source string) (*file.File, *[]byte, error) {
	l := app.GetFileServiceLogger()

	if source == "" {
		return nil, nil, file.ErrFileSourceEmpty
	}

	if utils.IsValidURL(source) {
		l.Debug("Getting file from URL", "source", source)
		f, data, err := s.GetFromURL(ctx, source)
		if err != nil {
			l.Error("Error getting file", "error", err)
			return nil, nil, err
		}
		return f, data, nil
	}

	if utils.IsUUID(source) {
		l.Debug("Getting file from cache by UUID", "source", source)
		f, data, err := s.GetFromUploads(ctx, source)
		if err != nil {
			l.Error("Error getting file", "error", err)
			return nil, nil, err
		}
		return f, data, nil
	}

	l.Debug("Getting file from Base64")
	f, data, err := s.GetFromBase64(ctx, source)
	if err != nil {
		l.Error("Error getting file", "error", err)
		return nil, nil, err
	}
	return f, data, nil
}

func (p *FileService) GetFromBase64(ctx context.Context, b64 string) (*file.File, *[]byte, error) {
	b64 = clearBase64(b64)
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(b64))
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, nil, file.ErrFileUnreachable
	}

	f := file.NewFromMime(http.DetectContentType(data))
	f.Sha256 = fmt.Sprintf("%x", sha256.Sum256(data))
	f.Size = uint64(len(data))

	return f, &data, nil
}

func (p *FileService) GetFromURL(ctx context.Context, url string) (*file.File, *[]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, nil, file.ErrFileUnreachable
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, nil, file.ErrFileUnreachable
	}

	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, file.ErrFileUnreachable
	}

	f := file.NewFromMime(http.DetectContentType(data))
	f.Sha256 = fmt.Sprintf("%x", sha256.Sum256(data))
	f.Size = uint64(len(data))

	return f, &data, nil
}

func (p *FileService) GetFromBytes(ctx context.Context, data []byte) (*file.File, *[]byte, error) {
	out := make([]byte, len(data))
	copy(out, data)

	f := file.NewFromMime(http.DetectContentType(data))
	f.Sha256 = fmt.Sprintf("%x", sha256.Sum256(data))
	f.Size = uint64(len(data))

	return f, &out, nil
}

func (p *FileService) GetFromUploads(ctx context.Context, fileID string) (*file.File, *[]byte, error) {
	f, err := p.fileRepo.Get(file.WhereID(fileID))
	if err != nil {
		return nil, nil, file.ErrFileUnreachable
	}

	if f == nil {
		return nil, nil, file.ErrFileNotFound
	}

	data, err := p.storage.Get(ctx, f.Path)
	if err != nil {
		return nil, nil, file.ErrFileUnreachable
	}

	return f, &data, nil
}

// --------------------- HELPERS ---------------------

type countWriter struct {
	Count int64
}

func (c *countWriter) Write(p []byte) (int, error) {
	n := len(p)
	c.Count += int64(n)
	return n, nil
}

func detectMimeTypeFromStream(stream io.ReadCloser) (string, io.ReadCloser, error) {
	// to lendo os 512 primeiros bytes para detectar o MIME
	buf := make([]byte, 512)
	n, err := stream.Read(buf)
	if err != nil && err != io.EOF {
		stream.Close()
		return "", nil, err
	}

	// aqui eu to detectando o MIME
	mime := http.DetectContentType(buf[:n])

	// preciso criar um novo ReadCloser que primeiro lê o que já foi lido (buf[:n]) e depois continua lendo do stream original
	// io.MultiReader cria um Reader que lê de múltiplos Readers em sequência
	// io.NopCloser transforma um Reader em um ReadCloser (com um Close que não faz nada)

	// então tipo como eu já li 512 agora eu concateno as streams grando um stream que vai ler a primeira stream com os bytes já lidos e depois vai continuar lendo a outra stream com o resto dos bytes
	newStream := io.NopCloser(io.MultiReader(bytes.NewReader(buf[:n]), stream))

	// Obs.: não posos fechar o stream original aqui, porque o multiReader ainda vai precisar dela

	return mime, newStream, nil
}

// Não sei porque mas o decode gera um erro se eu não remover isso, o que pra mim é estranho
func clearBase64(b64 string) string {
	if idx := strings.Index(b64, ","); idx != -1 {
		b64 = b64[idx+1:]
	}
	return b64
}
