package consumer

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/mauriciorobertodev/whappy-go/internal/app"
	"github.com/mauriciorobertodev/whappy-go/internal/app/cache"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/events"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/webhook"
	c "github.com/mauriciorobertodev/whappy-go/internal/infra/cache"
)

// TODO: implement retry logic with exponential backoff ?
// TODO: implement dead letter queue for failed webhooks ?
// TODO: implement circuit breaker pattern ?
// TODO: implement metrics for webhook delivery ?
// TODO: implement alerting for webhook delivery failures ?

type WebhookConsumer struct {
	webRepo webhook.WebhookRepository
	cache   cache.Cache
}

func NewWebhookConsumer(
	webRepo webhook.WebhookRepository,
	cache cache.Cache,
) *WebhookConsumer {
	return &WebhookConsumer{
		webRepo: webRepo,
		cache:   cache,
	}
}

func (w *WebhookConsumer) Handle(event events.Event) {
	l := app.GetWebhookLogger()

	if event.InstanceID == nil {
		l.Debug("event has no instance_id, skipping webhook delivery", "event", event.Name)
		return
	}

	cacheKey := cache.CacheKeyWebhooksPrefix + *event.InstanceID

	webhooks, err := c.Get[[]c.CachedWebhook](w.cache, cacheKey)

	if err != nil {
		databaseWebhooks, err := w.webRepo.List(webhook.WhereInstanceID(*event.InstanceID), webhook.WhereActive(true))
		if err != nil {
			l.Error("failed to fetch webhooks from database", "instance_id", *event.InstanceID, "error", err)
			return
		}

		for _, wh := range databaseWebhooks {
			webhooks = append(webhooks, c.ToCachedWebhook(wh))
		}

		_ = c.Set(w.cache, cacheKey, webhooks, cache.DefaultTTL*6) // 30m
	}

	for _, wh := range webhooks {
		go func(wh *webhook.Webhook) {
			if event.Matches(wh.Events) {
				l.Debug("sending webhook", "webhook_id", wh.ID, "url", wh.URL, "event", event.Name)
				if err := w.Send(wh, event); err != nil {
					l.Error("failed to send webhook", "webhook_id", wh.ID, "error", err)
				}
			} else {
				l.Debug("event does not match webhook events, skipping", "webhook_id", wh.ID, "event", event.Name)
			}
		}(c.FromCachedWebhook(&wh))
	}
}

func (w *WebhookConsumer) Send(wh *webhook.Webhook, event events.Event) error {
	l := app.GetWebhookLogger()

	timestamp := event.OccurredAt.Unix()

	signature, err := wh.SignEvent(event, timestamp)
	if err != nil {
		l.Error("failed to sign webhook event", "error", err)
		return err
	}

	body, err := event.ToJSON()
	if err != nil {
		l.Error("failed to marshal webhook event", "error", err)
		return err
	}

	req, err := http.NewRequest(http.MethodPost, wh.URL, bytes.NewReader(body))
	if err != nil {
		l.Error("failed to create webhook request", "error", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Whappy GO Webhook/1.0")
	req.Header.Set("X-Whappy-Event", string(event.Name))
	req.Header.Set("X-Whappy-Signature", signature)
	req.Header.Set("X-Whappy-Timestamp", fmt.Sprintf("%d", timestamp))

	// Limit the timeout to 5 seconds, to avoid hanging requests
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		l.Error("failed to send webhook request", "url", wh.URL, "error", err)
		return err
	}
	defer resp.Body.Close()

	bodyResp, err := io.ReadAll(resp.Body)
	if err != nil {
		l.Error("failed to read webhook response", "url", wh.URL, "error", err)
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		l.Warn("webhook returned non-2xx status",
			"url", wh.URL,
			"status", resp.StatusCode,
			"response", string(bodyResp),
		)
		return fmt.Errorf("webhook returned %d: %s", resp.StatusCode, string(bodyResp))
	}

	l.Info("webhook delivered successfully",
		"url", wh.URL,
		"status", resp.StatusCode,
	)

	return nil
}
