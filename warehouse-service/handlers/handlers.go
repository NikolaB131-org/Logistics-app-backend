package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/NikolaB131-org/logistics-backend/warehouse-service/models"
	"github.com/NikolaB131-org/logistics-backend/warehouse-service/rabbitmq"
)

func HandleProducts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGetProducts(w, r)
	case http.MethodPost:
		handlePostProduct(w, r)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func handleGetProducts(w http.ResponseWriter, r *http.Request) {
	jsonResp, err := json.Marshal(models.GetProducts())
	if err != nil {
		fmt.Println("Error happened in JSON marshal: ", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}

func handleGetProduct(w http.ResponseWriter, r *http.Request) {

}

func handlePostProduct(w http.ResponseWriter, r *http.Request) {
	headerContentTtype := r.Header.Get("Content-Type")
	if headerContentTtype != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte("Content-Type is not application/json"))
		return
	}
	var p models.Product
	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&p)
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad Request. Wrong Type provided for field " + unmarshalErr.Field))
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad Request"))
		}
		return
	}

	models.CreateProduct(p)
	rabbitmq.SendNotification(fmt.Sprintf("Product {name: %v, quantity: %v} was created", p.Name, p.Quantity))
	w.WriteHeader(http.StatusOK)
}
