package storage

import "fmt"

func (s *Store) ListNotes(vaultID *int64, parentID *int64, offset, limit int) ([]Note, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}

	query := `SELECT * FROM notes WHERE deleted_at IS NULL`
	var args []interface{}

	if vaultID != nil {
		query += ` AND vault_id = ?`
		args = append(args, *vaultID)
	}
	if parentID != nil {
		query += ` AND parent_id = ?`
		args = append(args, *parentID)
	}
	query += ` ORDER BY is_folder DESC, sort_order, title LIMIT ? OFFSET ?`
	args = append(args, limit, offset)

	var notes []Note
	if err := s.DB.Select(&notes, query, args...); err != nil {
		return nil, fmt.Errorf("storage: list notes: %w", err)
	}
	if notes == nil {
		notes = []Note{}
	}
	return notes, nil
}

func (s *Store) CreateNote(note *Note) error {
	const query = `
		INSERT INTO notes (title, content, vault_id, parent_id, is_folder, sort_order)
		VALUES (:title, :content, :vault_id, :parent_id, :is_folder, :sort_order)
	`
	result, err := s.DB.NamedExec(query, note)
	if err != nil {
		return fmt.Errorf("storage: create note: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("storage: create note last id: %w", err)
	}
	note.ID = id
	return s.DB.Get(note, `SELECT * FROM notes WHERE id = ?`, id)
}

func (s *Store) GetNote(id int64) (*Note, error) {
	var note Note
	if err := s.DB.Get(&note, `SELECT * FROM notes WHERE id = ? AND deleted_at IS NULL`, id); err != nil {
		return nil, fmt.Errorf("storage: get note: %w", err)
	}
	return &note, nil
}

func (s *Store) UpdateNote(note *Note) error {
	const query = `
		UPDATE notes SET
			title      = :title,
			content    = :content,
			vault_id   = :vault_id,
			parent_id  = :parent_id,
			is_folder  = :is_folder,
			sort_order = :sort_order,
			updated_at = datetime('now')
		WHERE id = :id AND deleted_at IS NULL
	`
	_, err := s.DB.NamedExec(query, note)
	if err != nil {
		return fmt.Errorf("storage: update note: %w", err)
	}
	return s.DB.Get(note, `SELECT * FROM notes WHERE id = ?`, note.ID)
}

func (s *Store) SoftDeleteNote(id int64) error {
	const query = `UPDATE notes SET deleted_at = datetime('now'), updated_at = datetime('now') WHERE id = ? AND deleted_at IS NULL`
	result, err := s.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("storage: soft delete note: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("storage: soft delete rows: %w", err)
	}
	if affected == 0 {
		return fmt.Errorf("storage: note %d not found or already deleted", id)
	}
	return nil
}
