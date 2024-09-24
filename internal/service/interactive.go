package service

import "context"

type InteractiveService interface {
	IncrReadCnt(ctx context.Context, id int64) error
}
