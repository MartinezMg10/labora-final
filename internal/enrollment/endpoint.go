package enrollment

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/MartinezMg10/rest_golang/pkg/meta"
	"github.com/gorilla/mux"
)

type (
	Controller func(w http.ResponseWriter, r *http.Request)

	Endpoints struct {
		Create Controller
		Get    Controller
		GetAll Controller
		Delete Controller
	}

	CreateReq struct {
		UserID   string `json:"user_id"`
		CourseID string `json:"course_id"`
	}

	Response struct {
		Status int         `json:"status"`
		Data   interface{} `json:"data,omitempty"`
		Err    string      `json:"error,omitempty"`
		Meta   *meta.Meta  `json:"meta,omitempty"`
	}
)

func MakeEndpoints(s Service) Endpoints {

	return Endpoints{
		Create: makeCreateEndpoint(s),
		Get:    makeGetEndpoint(s),
		GetAll: makeGetAllEndpoint(s),
		Delete: makeDeleteEndpoint(s),
	}

}

func makeCreateEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {

		var req CreateReq

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "Invalid request format"})
			return
		}

		if req.UserID == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "User id is required"})
			return
		}

		if req.CourseID == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "Course id is required"})
			return
		}

		enroll, err := s.Create(req.UserID, req.CourseID)
		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: err.Error()})
			return
		}

		json.NewEncoder(w).Encode(&Response{Status: 200, Data: enroll})
	}
}

func makeGetAllEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {

		v := r.URL.Query()

		filters := Filters{
			Name:      v.Get("name"),
			StartDate: v.Get("start_date"),
			EndDate:   v.Get("end_date"),
		}

		limit, _ := strconv.Atoi(v.Get("limit"))
		page, _ := strconv.Atoi(v.Get("page"))

		count, err := s.Count(filters)
		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: err.Error()})
		}

		meta, err := meta.New(page, limit, count)
		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: err.Error()})
		}

		enrollments, err := s.GetAll(filters, meta.Offset(), meta.Limit())
		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: err.Error()})
		}

		json.NewEncoder(w).Encode(&Response{Status: 200, Data: enrollments, Meta: meta})
	}
}

func makeGetEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		path := mux.Vars(r)
		id := path["id"]
		course, err := s.Get(id)
		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: err.Error()})
		}
		json.NewEncoder(w).Encode(&Response{Status: 200, Data: course})
	}
}

func makeDeleteEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		path := mux.Vars(r)
		id := path["id"]

		if err := s.Delete(id); err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Data: err.Error()})
			return
		}
		json.NewEncoder(w).Encode(&Response{Status: 200, Data: "sucess"})
	}
}
