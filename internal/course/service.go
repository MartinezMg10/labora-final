package course

import (
	"log"
	"time"

	"github.com/MartinezMg10/rest_golang/internal/domain"
)

type (
	Service interface {
		Create(name, starDate, endDate string) (*domain.Course, error)
		GetAll(filters Filters, offset, limit int) ([]domain.Course, error)
		Get(id string) (*domain.Course, error)
		Delete(id string) error
		Update(id string, Name *string, StartDate *string, EndDate *string) error
		Count(filter Filters) (int, error)
	}

	service struct {
		log  *log.Logger
		repo Repository
	}

	Filters struct {
		Name      string
		StartDate string
		EndDate   string
	}
)

func NewService(log *log.Logger, repo Repository) Service {
	return &service{
		log:  log,
		repo: repo,
	}
}

func (s service) Create(name, starDate, endDate string) (*domain.Course, error) {

	startDateParsed, err := time.Parse("2006-01-02", starDate)
	if err != nil {
		s.log.Println(err)
		return nil, err
	}

	endDateParsed, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		s.log.Println(err)
		return nil, err
	}

	course := &domain.Course{
		Name:      name,
		StartDate: startDateParsed.Format("2006-01-02 15:04:05"),
		EndDate:   endDateParsed.Format("2006-01-02 15:04:05"),
	}

	if err := s.repo.Create(course); err != nil {
		s.log.Println(err)
		return nil, err
	}

	return course, nil
}

func (s service) Get(id string) (*domain.Course, error) {
	course, err := s.repo.Get(id)
	if err != nil {
		return nil, err
	}
	return course, nil
}

func (s service) GetAll(filters Filters, offset, limit int) ([]domain.Course, error) {
	course, err := s.repo.GetAll(filters, offset, limit)
	if err != nil {
		return nil, err
	}
	return course, nil
}

func (s service) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s service) Update(id string, name *string, starDate *string, endDate *string) error {
	return s.repo.Update(id, name, starDate, endDate)
}

func (s service) Count(filters Filters) (int, error) {
	return s.repo.Count(filters)
}
