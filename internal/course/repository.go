package course

import (
	"fmt"
	"log"
	"strings"

	"github.com/MartinezMg10/rest_golang/internal/domain"
	"gorm.io/gorm"
)

type (
	Repository interface {
		Create(course *domain.Course) error
		GetAll(filters Filters, offset, limit int) ([]domain.Course, error)
		Get(id string) (*domain.Course, error)
		Delete(id string) error
		Update(id string, Name *string, StartDate *string, EndDate *string) error
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

func (repo *repo) Create(course *domain.Course) error {
	if err := repo.db.Create(course).Error; err != nil {
		repo.log.Printf("error: %v", err)
		return err
	}

	repo.log.Println("Course created with id: ", course.ID)
	return nil
}

func (repo *repo) GetAll(filters Filters, offset, limit int) ([]domain.Course, error) {
	var c []domain.Course

	tx := repo.db.Model(&c)
	tx = applyFilter(tx, filters)
	tx = tx.Limit(limit).Offset(offset)
	result := tx.Order("created_at DESC").Find(&c)
	if result.Error != nil {
		return nil, result.Error
	}

	return c, nil
}

func (repo *repo) Get(id string) (*domain.Course, error) {
	course := domain.Course{ID: id}

	if err := repo.db.First(&course).Error; err != nil {
		return nil, err
	}

	return &course, nil
}

func (repo *repo) Delete(id string) error {
	course := domain.Course{ID: id}

	if err := repo.db.Delete(&course).Error; err != nil {
		return err
	}

	return nil
}

func (repo *repo) Update(id string, Name *string, StartDate *string, EndDate *string) error {
	values := make(map[string]interface{})

	if Name != nil {
		values["name"] = *Name
	}

	if StartDate != nil {
		values["start_date"] = *StartDate
	}

	if EndDate != nil {
		values["end_date"] = *EndDate
	}

	if err := repo.db.Model(&domain.Course{}).Where("id = ?", id).Updates(values).Error; err != nil {
		return err
	}

	return nil
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
	tx := repo.db.Model(domain.Course{})
	tx = applyFilter(tx, filters)
	if err := tx.Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count), nil

}
