package enrollment

import (
	"fmt"
	"log"
	"strings"

	"github.com/MartinezMg10/rest_golang/internal/domain"
	"gorm.io/gorm"
)

type (
	Repository interface {
		Create(enroll *domain.Enrollment) error
		GetAll(filters Filters, offset, limit int) ([]domain.Enrollment, error)
		Get(id string) (*domain.Enrollment, error)
		Delete(id string) error
		Count(filter Filters) (int, error)
	}

	repo struct {
		db  *gorm.DB
		log *log.Logger
	}
)

func NewRepo(db *gorm.DB, l *log.Logger) Repository {
	return &repo{
		db:  db,
		log: l,
	}
}

func (repo *repo) Create(enroll *domain.Enrollment) error {
	if err := repo.db.Create(enroll).Error; err != nil {
		repo.log.Printf("error: %v", err)
		return err
	}

	repo.log.Println("Course created with id: ", enroll.ID)
	return nil
}

func (repo *repo) GetAll(filters Filters, offset, limit int) ([]domain.Enrollment, error) {
	var e []domain.Enrollment

	tx := repo.db.Model(&e)
	tx = applyFilter(tx, filters)
	tx = tx.Preload("User").Preload("Course")
	tx = tx.Limit(limit).Offset(offset)
	result := tx.Order("created_at DESC").Find(&e)
	if result.Error != nil {
		return nil, result.Error
	}

	return e, nil
}

func (repo *repo) Get(id string) (*domain.Enrollment, error) {
	var enrollment domain.Enrollment

	if err := repo.db.Preload("User").Preload("Course").First(&enrollment, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &enrollment, nil
}

func applyFilter(tx *gorm.DB, filters Filters) *gorm.DB {
	if filters.Name != "" {
		filters.Name = fmt.Sprintf("%%%s%%", strings.ToLower(filters.Name))
		tx = tx.Where("lower(name) like ?", filters.Name)
	}

	if filters.StartDate != "" {
		filters.StartDate = fmt.Sprintf("%%%s%%", strings.ToLower(filters.StartDate))
		tx = tx.Where("lower(start_date) like ?", filters.StartDate)
	}

	if filters.EndDate != "" {
		filters.EndDate = fmt.Sprintf("%%%s%%", strings.ToLower(filters.EndDate))
		tx = tx.Where("lower(end_date) like ?", filters.EndDate)
	}

	return tx
}

func (repo *repo) Count(filters Filters) (int, error) {
	var count int64
	tx := repo.db.Model(domain.Enrollment{})
	tx = applyFilter(tx, filters)
	if err := tx.Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count), nil

}

func (repo *repo) Delete(id string) error {
	enrollment := domain.Enrollment{ID: id}

	if err := repo.db.Delete(&enrollment).Error; err != nil {
		return err
	}

	return nil
}
