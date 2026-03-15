package repositories

import (
	"action-tag-expense-tracker/backend/models"

	"gorm.io/gorm"
)

type ActionTagRepository struct {
	DB *gorm.DB
}

func (r *ActionTagRepository) FindAll() ([]models.ActionTag, error) {
	var tags []models.ActionTag
	err := r.DB.Find(&tags).Error
	return tags, err
}

func (r *ActionTagRepository) Create(tag *models.ActionTag) error {
	return r.DB.Create(tag).Error
}

func (r *ActionTagRepository) Update(tag *models.ActionTag) error {
	return r.DB.Save(tag).Error
}

func (r *ActionTagRepository) Delete(id uint) error {
	return r.DB.Delete(&models.ActionTag{}, id).Error
}

func (r *ActionTagRepository) FindByID(id uint) (*models.ActionTag, error) {
	var tag models.ActionTag
	err := r.DB.First(&tag, id).Error
	return &tag, err
}

func (r *ActionTagRepository) FindOrCreateByName(name string) (*models.ActionTag, error) {
	var tag models.ActionTag
	err := r.DB.Where("name = ?", name).First(&tag).Error
	if err == gorm.ErrRecordNotFound {
		tag = models.ActionTag{Name: name}
		if createErr := r.DB.Create(&tag).Error; createErr != nil {
			return nil, createErr
		}
		return &tag, nil
	}
	return &tag, err
}
