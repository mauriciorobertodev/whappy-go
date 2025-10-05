package webhook

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
)

type Payload interface {
	ToJSON() ([]byte, error)
}

type Webhook struct {
	ID         string   `json:"id"`
	Active     bool     `json:"active"`
	URL        string   `json:"url"`
	Events     []string `json:"events"`
	secret     string
	InstanceID string    `json:"instance_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func New(url string, events []string, active bool) *Webhook {
	id, _ := uuid.NewV7()
	return &Webhook{
		ID:         id.String(),
		URL:        url,
		Events:     events,
		Active:     active,
		secret:     generateSecret(),
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
		InstanceID: "",
	}
}

func (w *Webhook) Activate() {
	w.Active = true
	w.UpdatedAt = time.Now().UTC()
}

func (w *Webhook) Deactivate() {
	w.Active = false
	w.UpdatedAt = time.Now().UTC()
}

func (w *Webhook) Update(url string, events []string) {
	w.URL = url
	w.Events = events
	w.UpdatedAt = time.Now().UTC()
}

func (w *Webhook) GetSecret() string {
	return w.secret
}

func (w *Webhook) SetSecret(secret string) {
	w.secret = secret
	w.UpdatedAt = time.Now().UTC()
}

func (w *Webhook) RenewSecret() {
	w.secret = generateSecret()
	w.UpdatedAt = time.Now().UTC()
}

func (w *Webhook) AttachToInstance(instanceID string) {
	w.InstanceID = instanceID
	w.UpdatedAt = time.Now().UTC()
}

func (w *Webhook) SignEvent(p Payload) (string, error) {
	payload, err := p.ToJSON()
	if err != nil {
		return "", err
	}

	h := hmac.New(sha256.New, []byte(w.secret))
	h.Write(payload)
	return hex.EncodeToString(h.Sum(nil)), nil
}

func generateSecret() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		panic("failed to generate secret")
	}
	return hex.EncodeToString(b)
}
