package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type CategoryHandler struct {
	db *sqlx.DB
}

type Category struct {
	ID           string    `db:"id" json:"id"`
	Name         string    `db:"name" json:"name"`
	Slug         string    `db:"slug" json:"slug"`
	Description  *string   `db:"description" json:"description"`
	Icon         *string   `db:"icon" json:"icon"`
	Color        *string   `db:"color" json:"color"`
	DisplayOrder int       `db:"display_order" json:"display_order"`
	IsActive     bool      `db:"is_active" json:"is_active"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

func NewCategoryHandler(db *sqlx.DB) *CategoryHandler {
	return &CategoryHandler{db: db}
}

func (h *CategoryHandler) GetCategories(c *gin.Context) {
	query := `
		SELECT id, name, slug, description, icon, color, display_order, is_active, created_at, updated_at
		FROM categories
		WHERE is_active = true
		ORDER BY display_order ASC, name ASC
	`

	rows, err := h.db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to fetch categories",
			"message": err.Error(),
		})
		return
	}
	defer rows.Close()

	categories := []Category{}
	for rows.Next() {
		var cat Category
		err := rows.Scan(
			&cat.ID,
			&cat.Name,
			&cat.Slug,
			&cat.Description,
			&cat.Icon,
			&cat.Color,
			&cat.DisplayOrder,
			&cat.IsActive,
			&cat.CreatedAt,
			&cat.UpdatedAt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Failed to scan category",
				"message": err.Error(),
			})
			return
		}
		categories = append(categories, cat)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    categories,
	})
}
