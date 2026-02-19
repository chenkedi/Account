package repository

import (
	"account/internal/business/models"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var (
	ErrCategoryNotFound = errors.New("category not found")
)

type CategoryRepository struct {
	db *sqlx.DB
}

func NewCategoryRepository(db *sqlx.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) Create(userID uuid.UUID, name string, categoryType models.CategoryType, parentID *uuid.UUID, icon string) (*models.Category, error) {
	now := time.Now().UTC()
	category := &models.Category{
		ID:             uuid.New(),
		UserID:         userID,
		Name:           name,
		Type:           categoryType,
		ParentID:       parentID,
		Icon:           icon,
		CreatedAt:      now,
		UpdatedAt:      now,
		LastModifiedAt: now,
		Version:        1,
		IsDeleted:      false,
	}

	query := `
		INSERT INTO categories (id, user_id, name, type, parent_id, icon, created_at, updated_at, last_modified_at, version, is_deleted)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err := r.db.Exec(query,
		category.ID, category.UserID, category.Name, category.Type, category.ParentID,
		category.Icon, category.CreatedAt, category.UpdatedAt, category.LastModifiedAt,
		category.Version, category.IsDeleted,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	return category, nil
}

func (r *CategoryRepository) GetByID(id uuid.UUID, userID uuid.UUID) (*models.Category, error) {
	var category models.Category

	query := `
		SELECT id, user_id, name, type, parent_id, icon, created_at, updated_at, last_modified_at, version, is_deleted
		FROM categories
		WHERE id = $1 AND user_id = $2 AND is_deleted = false
	`

	err := r.db.Get(&category, query, id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCategoryNotFound
		}
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	return &category, nil
}

func (r *CategoryRepository) GetAll(userID uuid.UUID) ([]models.Category, error) {
	var categories []models.Category

	query := `
		SELECT id, user_id, name, type, parent_id, icon, created_at, updated_at, last_modified_at, version, is_deleted
		FROM categories
		WHERE user_id = $1 AND is_deleted = false
		ORDER BY type, name ASC
	`

	err := r.db.Select(&categories, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}

	return categories, nil
}

func (r *CategoryRepository) GetByType(userID uuid.UUID, categoryType models.CategoryType) ([]models.Category, error) {
	var categories []models.Category

	query := `
		SELECT id, user_id, name, type, parent_id, icon, created_at, updated_at, last_modified_at, version, is_deleted
		FROM categories
		WHERE user_id = $1 AND type = $2 AND is_deleted = false
		ORDER BY name ASC
	`

	err := r.db.Select(&categories, query, userID, categoryType)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories by type: %w", err)
	}

	return categories, nil
}

func (r *CategoryRepository) Update(category *models.Category, userID uuid.UUID) (*models.Category, error) {
	now := time.Now().UTC()
	category.UpdatedAt = now
	category.LastModifiedAt = now
	category.Version++

	query := `
		UPDATE categories
		SET name = $1, type = $2, parent_id = $3, icon = $4, updated_at = $5, last_modified_at = $6, version = $7
		WHERE id = $8 AND user_id = $9 AND is_deleted = false
	`

	result, err := r.db.Exec(query,
		category.Name, category.Type, category.ParentID, category.Icon,
		category.UpdatedAt, category.LastModifiedAt, category.Version,
		category.ID, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rows == 0 {
		return nil, ErrCategoryNotFound
	}

	return category, nil
}

func (r *CategoryRepository) Delete(id uuid.UUID, userID uuid.UUID) error {
	now := time.Now().UTC()

	query := `
		UPDATE categories
		SET is_deleted = true, updated_at = $1, last_modified_at = $2, version = version + 1
		WHERE id = $3 AND user_id = $4 AND is_deleted = false
	`

	result, err := r.db.Exec(query, now, now, id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rows == 0 {
		return ErrCategoryNotFound
	}

	return nil
}

func (r *CategoryRepository) GetModifiedSince(userID uuid.UUID, since time.Time) ([]models.Category, error) {
	var categories []models.Category

	query := `
		SELECT id, user_id, name, type, parent_id, icon, created_at, updated_at, last_modified_at, version, is_deleted
		FROM categories
		WHERE user_id = $1 AND last_modified_at > $2
	`

	err := r.db.Select(&categories, query, userID, since)
	if err != nil {
		return nil, fmt.Errorf("failed to get modified categories: %w", err)
	}

	return categories, nil
}

func (r *CategoryRepository) CreateMany(categories []models.Category) error {
	if len(categories) == 0 {
		return nil
	}

	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	query := `
		INSERT INTO categories (id, user_id, name, type, parent_id, icon, created_at, updated_at, last_modified_at, version, is_deleted)
		VALUES (:id, :user_id, :name, :type, :parent_id, :icon, :created_at, :updated_at, :last_modified_at, :version, :is_deleted)
		ON CONFLICT (id) DO UPDATE
		SET name = EXCLUDED.name,
		    type = EXCLUDED.type,
		    parent_id = EXCLUDED.parent_id,
		    icon = EXCLUDED.icon,
		    updated_at = EXCLUDED.updated_at,
		    last_modified_at = EXCLUDED.last_modified_at,
		    version = EXCLUDED.version,
		    is_deleted = EXCLUDED.is_deleted
		WHERE categories.last_modified_at <= EXCLUDED.last_modified_at
	`

	for _, category := range categories {
		_, err := tx.NamedExec(query, category)
		if err != nil {
			return fmt.Errorf("failed to insert category %s: %w", category.ID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *CategoryRepository) CreateDefaultCategories(userID uuid.UUID) error {
	now := time.Now().UTC()
	var categories []models.Category

	for _, name := range models.DefaultIncomeCategories {
		categories = append(categories, models.Category{
			ID:             uuid.New(),
			UserID:         userID,
			Name:           name,
			Type:           models.CategoryTypeIncome,
			ParentID:       nil,
			Icon:           "",
			CreatedAt:      now,
			UpdatedAt:      now,
			LastModifiedAt: now,
			Version:        1,
			IsDeleted:      false,
		})
	}

	for _, name := range models.DefaultExpenseCategories {
		categories = append(categories, models.Category{
			ID:             uuid.New(),
			UserID:         userID,
			Name:           name,
			Type:           models.CategoryTypeExpense,
			ParentID:       nil,
			Icon:           "",
			CreatedAt:      now,
			UpdatedAt:      now,
			LastModifiedAt: now,
			Version:        1,
			IsDeleted:      false,
		})
	}

	return r.CreateMany(categories)
}
