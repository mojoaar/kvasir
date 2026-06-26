package storage

import "fmt"

func (s *Store) ListTags() ([]Tag, error) {
	var tags []Tag
	if err := s.DB.Select(&tags, `SELECT * FROM tags ORDER BY name`); err != nil {
		return nil, fmt.Errorf("storage: list tags: %w", err)
	}
	if tags == nil {
		tags = []Tag{}
	}
	return tags, nil
}

func (s *Store) CreateTag(tag *Tag) error {
	const query = `INSERT INTO tags (name, color) VALUES (:name, :color)`
	result, err := s.DB.NamedExec(query, tag)
	if err != nil {
		return fmt.Errorf("storage: create tag: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("storage: create tag last id: %w", err)
	}
	tag.ID = id
	return s.DB.Get(tag, `SELECT * FROM tags WHERE id = ?`, id)
}

func (s *Store) GetTag(id int64) (*Tag, error) {
	var tag Tag
	if err := s.DB.Get(&tag, `SELECT * FROM tags WHERE id = ?`, id); err != nil {
		return nil, fmt.Errorf("storage: get tag: %w", err)
	}
	return &tag, nil
}

func (s *Store) UpdateTag(tag *Tag) error {
	const query = `UPDATE tags SET name = :name, color = :color WHERE id = :id`
	_, err := s.DB.NamedExec(query, tag)
	if err != nil {
		return fmt.Errorf("storage: update tag: %w", err)
	}
	return s.DB.Get(tag, `SELECT * FROM tags WHERE id = ?`, tag.ID)
}

func (s *Store) DeleteTag(id int64) error {
	tx, err := s.DB.Beginx()
	if err != nil {
		return fmt.Errorf("storage: delete tag tx: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`DELETE FROM note_tags WHERE tag_id = ?`, id); err != nil {
		return fmt.Errorf("storage: delete tag note_tags: %w", err)
	}
	result, err := tx.Exec(`DELETE FROM tags WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("storage: delete tag: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("storage: delete tag rows: %w", err)
	}
	if affected == 0 {
		return fmt.Errorf("storage: tag %d not found", id)
	}
	return tx.Commit()
}

func (s *Store) AddTagToNote(noteID, tagID int64) error {
	const query = `INSERT OR IGNORE INTO note_tags (note_id, tag_id) VALUES (?, ?)`
	_, err := s.DB.Exec(query, noteID, tagID)
	if err != nil {
		return fmt.Errorf("storage: add tag to note: %w", err)
	}
	return nil
}

func (s *Store) RemoveTagFromNote(noteID, tagID int64) error {
	_, err := s.DB.Exec(`DELETE FROM note_tags WHERE note_id = ? AND tag_id = ?`, noteID, tagID)
	if err != nil {
		return fmt.Errorf("storage: remove tag from note: %w", err)
	}
	return nil
}

func (s *Store) GetNoteTags(noteID int64) ([]Tag, error) {
	const query = `
		SELECT t.* FROM tags t
		JOIN note_tags nt ON t.id = nt.tag_id
		WHERE nt.note_id = ?
		ORDER BY t.name
	`
	var tags []Tag
	if err := s.DB.Select(&tags, query, noteID); err != nil {
		return nil, fmt.Errorf("storage: get note tags: %w", err)
	}
	if tags == nil {
		tags = []Tag{}
	}
	return tags, nil
}
