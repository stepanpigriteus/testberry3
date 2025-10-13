package domain

import (
	"context"
)

type Service interface {
	CreateShorten(ctx context.Context, link string) error
	GetShorten(ctx context.Context, link string) (ShortLink, error)
	GetAnalytics(ctx context.Context, link string, group string) (VisitStats, error)
	SaveVisit(ctx context.Context, visit Visit) error
}
