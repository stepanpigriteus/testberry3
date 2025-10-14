package storage

import (
	"context"
	"treethree/domain"

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

func (s *StorageImpl) CreateComments(ctx context.Context, comment domain.Comment) error {
	const query = `
		INSERT INTO comments (parent_id, text, author)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at;
	`

	row := s.db.QueryRowContext(ctx, query,
		comment.ParentID,
		comment.Text,
		comment.Author,
	)

	err := row.Scan(&comment.ID, &comment.CreatedAt, &comment.UpdatedAt)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to insert comment")
		return err
	}

	s.logger.Info().Msgf("comment created with id=%d", comment.ID)
	return nil
}

func (s *StorageImpl) GetComments(ctx context.Context, id int) (domain.Comment, error) {
	var root domain.Comment

	query := `
		SELECT id, parent_id, author, text AS content, created_at, updated_at, deleted_at
		FROM comments
		WHERE id = $1
	`
	row := s.db.QueryRowContext(ctx, query, id)
	err := row.Scan(&root.ID, &root.ParentID, &root.Author, &root.Text, &root.CreatedAt, &root.UpdatedAt, &root.DeletedAt)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to get root comment")
		return domain.Comment{}, err
	}

	children, err := s.getChildren(ctx, int(root.ID))
	if err != nil {
		return domain.Comment{}, err
	}
	root.Children = children

	return root, nil
}

func (s *StorageImpl) getChildren(ctx context.Context, parentID int) ([]*domain.Comment, error) {
	query := `
		SELECT id, parent_id, author, text AS content, created_at, updated_at, deleted_at
		FROM comments
		WHERE parent_id = $1
	`
	rows, err := s.db.QueryContext(ctx, query, parentID)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to query child comments")
		return nil, err
	}
	defer rows.Close()

	var children []*domain.Comment

	for rows.Next() {
		var c domain.Comment
		if err := rows.Scan(&c.ID, &c.ParentID, &c.Author, &c.Text, &c.CreatedAt, &c.UpdatedAt, &c.DeletedAt); err != nil {
			s.logger.Error().Err(err).Msg("failed to scan child comment")
			return nil, err
		}

		cChildren, err := s.getChildren(ctx, int(c.ID))
		if err != nil {
			return nil, err
		}
		c.Children = cChildren

		children = append(children, &c)
	}

	return children, nil
}
