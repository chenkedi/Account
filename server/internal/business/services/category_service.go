package services

import (
	"account/internal/business/models"
	"account/internal/data/repository"
	"fmt"

	"github.com/google/uuid"
)

type CategoryService struct {
	categoryRepo *repository.CategoryRepository
}

func NewCategoryService(categoryRepo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{categoryRepo: categoryRepo}
}

type CreateCategoryRequest struct {
	Name     string                `json:"name" binding:"required"`
	Type     models.CategoryType   `json:"type" binding:"required"`
	ParentID *uuid.UUID            `json:"parent_id"`
	Icon     string                `json:"icon"`
}

type UpdateCategoryRequest struct {
	Name     string                `json:"name"`
	Type     models.CategoryType   `json:"type"`
	ParentID *uuid.UUID            `json:"parent_id"`
	Icon     string                `json:"icon"`
}

func (s *CategoryService) CreateCategory(userID uuid.UUID, req *CreateCategoryRequest) (*models.Category, error) {
	if !isValidCategoryType(req.Type) {
		return nil, fmt.Errorf("invalid category type: %s", req.Type)
	}

	category, err := s.categoryRepo.Create(userID, req.Name, req.Type, req.ParentID, req.Icon)
	if err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	return category, nil
}

func (s *CategoryService) GetCategory(id uuid.UUID, userID uuid.UUID) (*models.Category, error) {
	category, err := s.categoryRepo.GetByID(id, userID)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (s *CategoryService) GetAllCategories(userID uuid.UUID) ([]models.Category, error) {
	categories, err := s.categoryRepo.GetAll(userID)
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (s *CategoryService) GetCategoriesByType(userID uuid.UUID, categoryType models.CategoryType) ([]models.Category, error) {
	if !isValidCategoryType(categoryType) {
		return nil, fmt.Errorf("invalid category type: %s", categoryType)
	}

	categories, err := s.categoryRepo.GetByType(userID, categoryType)
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (s *CategoryService) UpdateCategory(id uuid.UUID, userID uuid.UUID, req *UpdateCategoryRequest) (*models.Category, error) {
	category, err := s.categoryRepo.GetByID(id, userID)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		category.Name = req.Name
	}
	if req.Type != "" {
		if !isValidCategoryType(req.Type) {
			return nil, fmt.Errorf("invalid category type: %s", req.Type)
		}
		category.Type = req.Type
	}
	if req.ParentID != nil {
		category.ParentID = req.ParentID
	}
	if req.Icon != "" {
		category.Icon = req.Icon
	}

	updatedCategory, err := s.categoryRepo.Update(category, userID)
	if err != nil {
		return nil, err
	}

	return updatedCategory, nil
}

func (s *CategoryService) DeleteCategory(id uuid.UUID, userID uuid.UUID) error {
	return s.categoryRepo.Delete(id, userID)
}

func isValidCategoryType(t models.CategoryType) bool {
	return t == models.CategoryTypeIncome || t == models.CategoryTypeExpense
}
