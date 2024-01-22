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
	}

	service struct {
		log       *log.Logger
		repo      Repository
		courseSrv course.Service
		userSrv   user.Service
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
