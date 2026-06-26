package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

type Note struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	VaultID   *int64 `json:"vaultId,omitempty"`
	ParentID  *int64 `json:"parentId,omitempty"`
	IsFolder  bool   `json:"isFolder"`
	SortOrder int    `json:"sortOrder"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	DeletedAt *string `json:"deletedAt,omitempty"`
}

type SearchResult struct {
	Note
	Rank    float64 `json:"rank"`
	Snippet string  `json:"snippet"`
}

type createNoteRequest struct {
	Title    string `json:"title"`
	Content  string `json:"content,omitempty"`
	VaultID  *int64 `json:"vaultId,omitempty"`
	ParentID *int64 `json:"parentId,omitempty"`
	IsFolder bool   `json:"isFolder"`
}

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{},
	}
}

func (c *Client) ListNotes() ([]Note, error) {
	url := c.BaseURL + "/api/v1/notes"
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("list notes: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("list notes: read body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("list notes: %s (status %d)", string(body), resp.StatusCode)
	}

	var notes []Note
	if err := json.Unmarshal(body, &notes); err != nil {
		return nil, fmt.Errorf("list notes: parse response: %w", err)
	}

	return notes, nil
}

func (c *Client) GetNote(id int64) (*Note, error) {
	url := fmt.Sprintf("%s/api/v1/notes/%d", c.BaseURL, id)
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("get note: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("get note: read body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get note: %s (status %d)", string(body), resp.StatusCode)
	}

	var note Note
	if err := json.Unmarshal(body, &note); err != nil {
		return nil, fmt.Errorf("get note: parse response: %w", err)
	}

	return &note, nil
}

func (c *Client) CreateNote(title, content string) (*Note, error) {
	req := createNoteRequest{
		Title:   title,
		Content: content,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("create note: marshal: %w", err)
	}

	url := c.BaseURL + "/api/v1/notes"
	resp, err := c.HTTPClient.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create note: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("create note: read body: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("create note: %s (status %d)", string(respBody), resp.StatusCode)
	}

	var note Note
	if err := json.Unmarshal(respBody, &note); err != nil {
		return nil, fmt.Errorf("create note: parse response: %w", err)
	}

	return &note, nil
}

func (c *Client) Search(query string, limit int) ([]SearchResult, error) {
	if limit <= 0 {
		limit = 20
	}

	url := fmt.Sprintf("%s/api/v1/search?q=%s&limit=%d", c.BaseURL, query, limit)
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("search: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("search: read body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search: %s (status %d)", string(respBody), resp.StatusCode)
	}

	var results []SearchResult
	if err := json.Unmarshal(respBody, &results); err != nil {
		return nil, fmt.Errorf("search: parse response: %w", err)
	}

	return results, nil
}
