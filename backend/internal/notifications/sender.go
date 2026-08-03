package notifications

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	webpush "github.com/SherClockHolmes/webpush-go"
)

type Sender interface {
	Send(ctx context.Context, delivery Delivery) (SendResult, error)
}

type WebPushSender struct {
	options webpush.Options
}

func NewWebPushSender(publicKey, privateKey, subject string) *WebPushSender {
	return &WebPushSender{
		options: webpush.Options{
			Subscriber:      subject,
			VAPIDPublicKey:  publicKey,
			VAPIDPrivateKey: privateKey,
			TTL:             24 * 60 * 60,
			Urgency:         webpush.UrgencyNormal,
			HTTPClient: &http.Client{
				Timeout: 15 * time.Second,
			},
		},
	}
}

func (s *WebPushSender) Send(ctx context.Context, delivery Delivery) (SendResult, error) {
	subscription := &webpush.Subscription{
		Endpoint: delivery.Endpoint,
		Keys: webpush.Keys{
			P256dh: delivery.P256DH,
			Auth:   delivery.Auth,
		},
	}

	response, err := webpush.SendNotificationWithContext(
		ctx,
		delivery.Payload,
		subscription,
		&s.options,
	)
	if err != nil {
		return SendResult{}, fmt.Errorf("send web push: %w", err)
	}
	defer response.Body.Close()
	_, _ = io.Copy(io.Discard, response.Body)

	return SendResult{StatusCode: response.StatusCode}, nil
}
