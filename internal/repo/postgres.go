package repo

import (
	"database/sql"
	"myblog/internal/models"

	_ "github.com/lib/pq"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreatePost(p *models.Post) error {
	return r.db.QueryRow(
		"INSERT INTO posts(title, content) VALUES($1, $2) RETURNING id, created_at",
		p.Title, p.Content,
	).Scan(&p.ID, &p.CreatedAt)
}

func (r *Repository) GetAllPosts() ([]models.Post, error) {
	rows, err := r.db.Query("SELECT id, title, created_at FROM posts ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var p models.Post
		if err := rows.Scan(&p.ID, &p.Title, &p.CreatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}

func (r *Repository) GetPostByID(id int) (models.Post, error) {
	var p models.Post
	err := r.db.QueryRow("SELECT id, title, content, created_at FROM posts WHERE id = $1", id).
		Scan(&p.ID, &p.Title, &p.Content, &p.CreatedAt)
	return p, err
}

func (r *Repository) UpdatePost(p *models.Post) error {
	_, err := r.db.Exec("UPDATE posts SET title = $1, content = $2 WHERE id = $3",
		p.Title, p.Content, p.ID)
	return err
}

func (r *Repository) DeletePost(id int) error {
	_, err := r.db.Exec("DELETE FROM posts WHERE id = $1", id)
	return err
}
