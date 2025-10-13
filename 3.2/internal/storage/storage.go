package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
	"treeTwo/domain"

	"github.com/rs/zerolog"
	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/zlog"
)

type StorageImpl struct {
	db     *dbpg.DB
	logger zerolog.Logger
}

func NewStorage(ctx context.Context, masterDSN string, slaveDSNs []string, logger zerolog.Logger) *StorageImpl {
	opts := &dbpg.Options{MaxOpenConns: 10, MaxIdleConns: 5}
	db, err := dbpg.New(masterDSN, slaveDSNs, opts)
	if err != nil {
		zlog.Logger.Error().Msgf("init database error %s", err)
	}
	return &StorageImpl{
		db:     db,
		logger: logger,
	}
}

func (s *StorageImpl) CreateShorten(ctx context.Context, link domain.ShortLink) error {
	if link == (domain.ShortLink{}) {
		return fmt.Errorf("empty shortlink")
	}
	query := `
        INSERT INTO short_links (short_code, original_url)
        VALUES ($1, $2)
        RETURNING id, short_code, original_url, created_at, click_count
    `

	err := s.db.QueryRowContext(ctx, query, link.ShortCode, link.OriginalURL).
		Scan(&link.ID, &link.ShortCode, &link.OriginalURL, &link.CreatedAt, &link.ClickCount)
	if err != nil {
		return fmt.Errorf("failed to insert short link: %w", err)
	}

	return nil
}

func (s *StorageImpl) GetShorten(ctx context.Context, link string) (domain.ShortLink, error) {
	var shortLink domain.ShortLink
	query := `
        SELECT * FROM short_links WHERE short_code = $1
    `
	err := s.db.QueryRowContext(ctx, query, link).Scan(
		&shortLink.ID,
		&shortLink.ShortCode,
		&shortLink.OriginalURL,
		&shortLink.CreatedAt,
		&shortLink.ClickCount,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ShortLink{}, fmt.Errorf("not found in db")
		}
		return domain.ShortLink{}, fmt.Errorf("query failed: %w", err)
	}
	return shortLink, nil

}

func (s *StorageImpl) GetAnalytics(ctx context.Context, shortCode string, group string) (domain.VisitStats, error) {
	var stats domain.VisitStats

	query := `
        SELECT
            COUNT(*) AS total,
            COUNT(DISTINCT ip_address) AS unique_ips
        FROM link_visits v
        JOIN short_links l ON v.short_link_id = l.id
        WHERE l.short_code = $1
    `
	err := s.db.QueryRowContext(ctx, query, shortCode).Scan(
		&stats.TotalVisits,
		&stats.UniqueIPs,
	)
	if err != nil {
		return stats, err
	}

	deviceRows, err := s.db.QueryContext(ctx, `
        SELECT device_type, COUNT(*) 
        FROM link_visits v
        JOIN short_links l ON v.short_link_id = l.id
        WHERE l.short_code = $1
        GROUP BY device_type
    `, shortCode)
	if err != nil {
		return stats, err
	}
	defer deviceRows.Close()

	stats.DeviceStats = make(map[string]int)
	for deviceRows.Next() {
		var device string
		var count int
		if err := deviceRows.Scan(&device, &count); err != nil {
			continue
		}
		stats.DeviceStats[device] = count
	}

	switch group {
	case "day":
		dailyRows, err := s.db.QueryContext(ctx, `
            SELECT 
                DATE(visited_at) as date,
                COUNT(*) as count
            FROM link_visits v
            JOIN short_links l ON v.short_link_id = l.id
            WHERE l.short_code = $1
            GROUP BY DATE(visited_at)
            ORDER BY date DESC
        `, shortCode)
		if err != nil {
			return stats, err
		}
		defer dailyRows.Close()

		stats.DailyActivity = make(map[string]int)
		for dailyRows.Next() {
			var date time.Time
			var count int
			if err := dailyRows.Scan(&date, &count); err != nil {
				continue
			}
			stats.DailyActivity[date.Format("2006-01-02")] = count
		}

	case "month":
		monthlyRows, err := s.db.QueryContext(ctx, `
            SELECT 
                TO_CHAR(visited_at, 'YYYY-MM') as month,
                COUNT(*) as count
            FROM link_visits v
            JOIN short_links l ON v.short_link_id = l.id
            WHERE l.short_code = $1
            GROUP BY TO_CHAR(visited_at, 'YYYY-MM')
            ORDER BY month DESC
        `, shortCode)
		if err != nil {
			return stats, err
		}
		defer monthlyRows.Close()

		stats.MonthlyActivity = make(map[string]int)
		for monthlyRows.Next() {
			var month string
			var count int
			if err := monthlyRows.Scan(&month, &count); err != nil {
				continue
			}
			stats.MonthlyActivity[month] = count
		}

	case "detailed":
		dailyRows, err := s.db.QueryContext(ctx, `
            SELECT 
                DATE(visited_at) as date,
                COUNT(*) as count
            FROM link_visits v
            JOIN short_links l ON v.short_link_id = l.id
            WHERE l.short_code = $1
            GROUP BY DATE(visited_at)
            ORDER BY date DESC
        `, shortCode)
		if err != nil {
			return stats, err
		}
		defer dailyRows.Close()

		stats.DailyActivity = make(map[string]int)
		for dailyRows.Next() {
			var date time.Time
			var count int
			if err := dailyRows.Scan(&date, &count); err != nil {
				continue
			}
			stats.DailyActivity[date.Format("2006-01-02")] = count
		}

		monthlyRows, err := s.db.QueryContext(ctx, `
            SELECT 
                TO_CHAR(visited_at, 'YYYY-MM') as month,
                COUNT(*) as count
            FROM link_visits v
            JOIN short_links l ON v.short_link_id = l.id
            WHERE l.short_code = $1
            GROUP BY TO_CHAR(visited_at, 'YYYY-MM')
            ORDER BY month DESC
        `, shortCode)
		if err != nil {
			return stats, err
		}
		defer monthlyRows.Close()

		stats.MonthlyActivity = make(map[string]int)
		for monthlyRows.Next() {
			var month string
			var count int
			if err := monthlyRows.Scan(&month, &count); err != nil {
				continue
			}
			stats.MonthlyActivity[month] = count
		}

		visitRows, err := s.db.QueryContext(ctx, `
            SELECT 
                v.id,
                v.short_link_id,
                v.visited_at,
                v.user_agent,
                v.ip_address,
                COALESCE(v.device_type, '') as device_type
            FROM link_visits v
            JOIN short_links l ON v.short_link_id = l.id
            WHERE l.short_code = $1
            ORDER BY v.visited_at DESC
            LIMIT 100
        `, shortCode)
		if err != nil {
			return stats, err
		}
		defer visitRows.Close()

		stats.Visits = []domain.Visit{}
		for visitRows.Next() {
			var visit domain.Visit
			if err := visitRows.Scan(
				&visit.ID,
				&visit.ShortLinkID,
				&visit.VisitedAt,
				&visit.UserAgent,
				&visit.IPAddress,
				&visit.DeviceType,
			); err != nil {
				continue
			}
			stats.Visits = append(stats.Visits, visit)
		}

	default:
		var todayCount int
		err := s.db.QueryRowContext(ctx, `
            SELECT COUNT(*) 
            FROM link_visits v
            JOIN short_links l ON v.short_link_id = l.id
            WHERE l.short_code = $1 
            AND DATE(visited_at) = CURRENT_DATE
        `, shortCode).Scan(&todayCount)

		if err != nil && err != sql.ErrNoRows {
			return stats, err
		}

		stats.DailyActivity = make(map[string]int)
		if todayCount > 0 {
			today := time.Now().Format("2006-01-02")
			stats.DailyActivity[today] = todayCount
		}
	}

	return stats, nil
}

func (s *StorageImpl) UpdateClickCount(ctx context.Context, id int64, count int) error {
	query := `UPDATE short_links SET click_count = $1 WHERE id = $2`
	_, err := s.db.ExecContext(ctx, query, count, id)
	return err
}

func (s *StorageImpl) SaveVisit(ctx context.Context, visit domain.Visit) error {
	fmt.Println(visit)
	query := `
		INSERT INTO link_visits (short_link_id, visited_at, user_agent, ip_address, device_type)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := s.db.ExecContext(ctx, query,
		visit.ShortLinkID,
		visit.VisitedAt,
		visit.UserAgent,
		visit.IPAddress,
		visit.DeviceType,
	)
	return err
}
