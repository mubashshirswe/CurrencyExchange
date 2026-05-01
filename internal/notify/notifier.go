package notify

import "context"

// DeliveredUser sends FCM notifications to transaction delivered users (delivered_user_id).
type DeliveredUser interface {
	NotifyPendingDelivery(ctx context.Context, deliveredUserID *int64, txnID int64, phone, details string)
	NotifyDeliveryCompleted(ctx context.Context, deliveredUserID int64, txnID int64, details string)
}

type NoopDeliveredUser struct{}

func (NoopDeliveredUser) NotifyPendingDelivery(context.Context, *int64, int64, string, string) {}

func (NoopDeliveredUser) NotifyDeliveryCompleted(context.Context, int64, int64, string) {}
