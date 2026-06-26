package storage

import "fmt"

type SearchResult struct {
	Note
	Rank    float64 `db:"rank"    json:"rank"`
	Snippet string  `db:"snippet" json:"snippet"`
}

func (s *Store) Search(query string, limit int) ([]SearchResult, error) {
	if limit <= 0 {
		limit = 20
	}

	const sql = `
		SELECT n.*, f.rank, snippet(notes_fts, 1, '<mark>', '</mark>', '...', 32) AS snippet
		FROM notes_fts f
		JOIN notes n ON n.id = f.rowid
		WHERE notes_fts MATCH ? AND n.deleted_at IS NULL
		ORDER BY rank
		LIMIT ?
	`

	var results []SearchResult
	if err := s.DB.Select(&results, sql, query, limit); err != nil {
		return nil, fmt.Errorf("storage: search: %w", err)
	}

	if results == nil {
		results = []SearchResult{}
	}

	return results, nil
}

func (s *Store) RebuildFTS() error {
	const sql = `INSERT INTO notes_fts(notes_fts) VALUES('rebuild')`
	_, err := s.DB.Exec(sql)
	if err != nil {
		return fmt.Errorf("storage: rebuild fts: %w", err)
	}
	return nil
}
