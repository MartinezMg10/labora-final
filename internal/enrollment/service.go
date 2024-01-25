package enrollment

import (
	"errors"
	"log"

	"github.com/MartinezMg10/rest_golang/internal/course"
	"github.com/MartinezMg10/rest_golang/internal/domain"
	"github.com/MartinezMg10/rest_golang/internal/user"
)

type (
	Service interface {
		Create(userId, courseId string) (*domain.Enrollment, error)
		GetAll(filters Filters, offset, limit int) ([]domain.Enrollment, error)
		Get(id string) (*domain.Enrollment, error)
		Delete(id string) error
		Count(filter Filters) (int, error)
	}

	service struct {
		log       *log.Logger
		repo      Repository
		courseSrv course.Service
		userSrv   user.Service
	}

	Filters struct {
		Name      string
		StartDate string
		EndDate   string
	}
)

func NewService(log *log.Logger, repo Repository, courseSrv course.Service, userSrv user.Service) Service {
	return &service{
		log:       log,
		repo:      repo,
		courseSrv: courseSrv,
		userSrv:   userSrv,
	}
}

func (s service) Create(userId, courseId string) (*domain.Enrollment, error) {

	enroll := &domain.Enrollment{
		UserID:   userId,
		CourseID: courseId,
		Status:   "P",
	}

	if _, err := s.userSrv.Get(enroll.UserID); err != nil {
		return nil, errors.New("user id doesn't exists")
	}

	if _, err := s.courseSrv.Get(enroll.CourseID); err != nil {
		return nil, errors.New("course id doesn't exists")
	}

	if err := s.repo.Create(enroll); err != nil {
		s.log.Println(err)
		return nil, err
	}

	return enroll, nil
}

func (s service) Get(id string) (*domain.Enrollment, error) {
	enrollment, err := s.repo.Get(id)
	if err != nil {
		return nil, err
	}
	return enrollment, nil
}

func (s service) GetAll(filters Filters, offset, limit int) ([]domain.Enrollment, error) {
	enrollments, err := s.repo.GetAll(filters, offset, limit)
	if err != nil {
		return nil, err
	}
	return enrollments, nil
}

func (s service) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s service) Count(filters Filters) (int, error) {
	return s.repo.Count(filters)
}
