package adapter

import (
	"context"
)

type IModel interface {
	RecognizeTextMessage(ctx context.Context, userInfo, message string) (RecognizeResult, error)
	RecognizeImageMessage(ctx context.Context, userInfo, file string) (RecognizeResult, error)
}
