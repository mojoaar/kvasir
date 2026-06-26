package storage

type Note struct {
	ID        int64    `db:"id"          json:"id"`
	Title     string   `db:"title"       json:"title"`
	Content   string   `db:"content"     json:"content"`
	VaultID   *int64   `db:"vault_id"    json:"vaultId,omitempty"`
	ParentID  *int64   `db:"parent_id"   json:"parentId,omitempty"`
	IsFolder  bool     `db:"is_folder"   json:"isFolder"`
	SortOrder int      `db:"sort_order"  json:"sortOrder"`
	CreatedAt SQLTime  `db:"created_at"  json:"createdAt"`
	UpdatedAt SQLTime  `db:"updated_at"  json:"updatedAt"`
	DeletedAt *SQLTime `db:"deleted_at"  json:"deletedAt,omitempty"`
}

type Tag struct {
	ID        int64   `db:"id"         json:"id"`
	Name      string  `db:"name"       json:"name"`
	Color     string  `db:"color"      json:"color"`
	CreatedAt SQLTime `db:"created_at" json:"createdAt"`
}

type NoteTag struct {
	NoteID int64 `db:"note_id" json:"noteId"`
	TagID  int64 `db:"tag_id"  json:"tagId"`
}

type Vault struct {
	ID           int64   `db:"id"             json:"id"`
	Name         string  `db:"name"           json:"name"`
	Description  string  `db:"description"    json:"description"`
	PasswordHash string  `db:"password_hash"  json:"-"`
	CreatedAt    SQLTime `db:"created_at"     json:"createdAt"`
	UpdatedAt    SQLTime `db:"updated_at"     json:"updatedAt"`
}

type Attachment struct {
	ID        int64   `db:"id"        json:"id"`
	NoteID    int64   `db:"note_id"   json:"noteId"`
	Filename  string  `db:"filename"  json:"filename"`
	MimeType  string  `db:"mime_type" json:"mimeType"`
	Size      int64   `db:"size"      json:"size"`
	Path      string  `db:"path"      json:"path"`
	CreatedAt SQLTime `db:"created_at" json:"createdAt"`
}

type Version struct {
	ID         int64   `db:"id"          json:"id"`
	NoteID     int64   `db:"note_id"     json:"noteId"`
	Content    string  `db:"content"     json:"content"`
	VersionNum int     `db:"version_num" json:"versionNum"`
	CreatedAt  SQLTime `db:"created_at"  json:"createdAt"`
}

type Theme struct {
	ID         int64   `db:"id"          json:"id"`
	Name       string  `db:"name"        json:"name"`
	ConfigJSON string  `db:"config_json" json:"configJson"`
	IsBuiltin  bool    `db:"is_builtin"  json:"isBuiltin"`
	CreatedAt  SQLTime `db:"created_at"  json:"createdAt"`
}

type Plugin struct {
	ID           int64   `db:"id"            json:"id"`
	Name         string  `db:"name"          json:"name"`
	ManifestJSON string  `db:"manifest_json" json:"manifestJson"`
	Enabled      bool    `db:"enabled"       json:"enabled"`
	CreatedAt    SQLTime `db:"created_at"    json:"createdAt"`
}

type PluginPermission struct {
	ID         int64  `db:"id"          json:"id"`
	PluginID   int64  `db:"plugin_id"   json:"pluginId"`
	Permission string `db:"permission"  json:"permission"`
}
