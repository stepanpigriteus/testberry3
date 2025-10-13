package domain

import "context"

type Storage interface {
	CreateShorten(ctx context.Context, link ShortLink) error
	GetShorten(ctx context.Context, link string) (ShortLink, error)
	GetAnalytics(ctx context.Context, olink string, group string) (VisitStats, error)
	UpdateClickCount(ctx context.Context, shortLinkID int64, shortLinkClick int) error
	SaveVisit(ctx context.Context, visit Visit) error
}
