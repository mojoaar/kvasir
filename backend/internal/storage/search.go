package storage

import "fmt"

func (s *Store) SearchByTag(query string) ([]Note, error) {
	const sql = `
		SELECT DISTINCT n.*
		FROM notes n
		JOIN note_tags nt ON n.id = nt.note_id
		JOIN tags t ON t.id = nt.tag_id
		WHERE t.name LIKE ? AND n.deleted_at IS NULL
		ORDER BY n.updated_at DESC
	`

	var notes []Note
	if err := s.DB.Select(&notes, sql, "%"+query+"%"); err != nil {
		return nil, fmt.Errorf("storage: search by tag: %w", err)
	}

	if notes == nil {
		notes = []Note{}
	}

	return notes, nil
}
