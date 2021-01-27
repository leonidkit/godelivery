package api

import (
	"encoding/json"
	"fmt"
	"godelivery/internal/converter"
	"godelivery/internal/storage"
	"godelivery/pkg/logger"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type APIServer struct {
	Router  *mux.Router
	storage storage.ItemRepository
	conv    converter.Interface
	logger  logger.Interface
}

func New(converter converter.Interface, storage storage.ItemRepository, logger logger.Interface) *APIServer {
	apiservice := &APIServer{
		storage: storage,
		logger:  logger,
		conv:    converter,
		Router:  mux.NewRouter(),
	}

	apiservice.endpoints()

	return apiservice
}

func (as *APIServer) Run(host, port string) error {
	return http.ListenAndServe(host+":"+port, as.Router)
}

func (as *APIServer) endpoints() {
	as.Router.HandleFunc("/create", as.Create).Methods(http.MethodPost)
	as.Router.HandleFunc("/read/{format}/{id}", as.Read).Methods(http.MethodGet)
	as.Router.HandleFunc("/delete/{format}/{id}", as.Delete).Methods(http.MethodDelete)
}

func (as *APIServer) writeError(w http.ResponseWriter, code int, err string) {
	w.WriteHeader(code)

	errJSON := Response{
		Error: err,
	}
	json.NewEncoder(w).Encode(&errJSON)
}

func (as *APIServer) Create(w http.ResponseWriter, r *http.Request) {
	var req Request

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		msg := fmt.Sprintf("convert request error: %v", err)

		as.logger.Error(msg)
		as.writeError(w, http.StatusBadRequest, msg)
		return
	}

	err = req.validate()
	if err != nil {
		msg := fmt.Sprintf("request error: %v", err)

		as.logger.Error(msg)
		as.writeError(w, http.StatusBadRequest, msg)
		return
	}

	convertedFormat, err := as.conv.Convert(req.Format, req.FormatType)
	if err != nil {
		msg := fmt.Sprintf("converting format error: %v", err)

		as.logger.Error(msg)
		as.writeError(w, http.StatusInternalServerError, msg)
		return
	}

	item := storage.Item{
		ID:         req.ID,
		FormatType: req.FormatType,
		Format:     convertedFormat,
	}

	err = as.storage.Create(r.Context(), item)
	if err != nil {
		msg := fmt.Sprintf("inserting item error: %v", err)

		as.logger.Error(msg)
		as.writeError(w, http.StatusInternalServerError, msg)
		return
	}

	resp := Response{
		Message: "OK",
	}
	json.NewEncoder(w).Encode(&resp)
}

func (as *APIServer) Read(w http.ResponseWriter, r *http.Request) {
	queryMap := mux.Vars(r)

	if _, ok := queryMap["id"]; !ok {
		as.logger.Error("request error")
		as.writeError(w, http.StatusInternalServerError, "field id not exists")
		return
	}
	if _, ok := queryMap["format"]; !ok {
		as.logger.Error("request error")
		as.writeError(w, http.StatusInternalServerError, "field format not exists")
		return
	}

	id, err := strconv.ParseInt(queryMap["id"], 10, 64)
	if err != nil {
		msg := fmt.Sprintf("request error: %v", err)

		as.logger.Error(msg)
		as.writeError(w, http.StatusInternalServerError, msg)
		return
	}

	item, err := as.storage.Get(r.Context(), id, queryMap["format"])
	if err != nil {
		msg := fmt.Sprintf("getting item error: %v", err)

		as.logger.Errorf(msg)
		as.writeError(w, http.StatusInternalServerError, msg)
		return
	}

	resp := Response{
		Message: "OK",
		Format:  string(item.Format),
	}
	json.NewEncoder(w).Encode(&resp)
}

func (as *APIServer) Delete(w http.ResponseWriter, r *http.Request) {
	queryMap := mux.Vars(r)

	if _, ok := queryMap["id"]; !ok {
		as.logger.Error("request error")
		as.writeError(w, http.StatusInternalServerError, "field id not exists")
		return
	}
	if _, ok := queryMap["format"]; !ok {
		as.logger.Error("request error")
		as.writeError(w, http.StatusInternalServerError, "field format not exists")
		return
	}

	id, err := strconv.ParseInt(queryMap["id"], 10, 64)
	if err != nil {
		msg := fmt.Sprintf("request error: %v", err)

		as.logger.Error(msg)
		as.writeError(w, http.StatusInternalServerError, msg)
		return
	}

	err = as.storage.Delete(r.Context(), id, queryMap["format"])
	if err != nil {
		msg := fmt.Sprintf("getting item error: %v", err)

		as.logger.Error(msg)
		as.writeError(w, http.StatusInternalServerError, msg)
		return
	}

	resp := Response{
		Message: "OK",
	}
	json.NewEncoder(w).Encode(&resp)
}
