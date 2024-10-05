package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"go.temporal.io/sdk/client"
	"log"
	"net/http"
)

type handlers struct {
	temporalClient client.Client
}

func (h *handlers) handleMenuFetch(w http.ResponseWriter, r *http.Request) {
	menu := Menu{
		Items: []MenuItem{
			{Name: "Coffee", Type: "beverage", Price: 300},
			{Name: "Latte", Type: "beverage", Price: 350},
			{Name: "Milkshake", Type: "beverage", Price: 450},
			{Name: "Bagel", Type: "food", Price: 500},
			{Name: "Sandwich", Type: "food", Price: 600},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(menu)
}

func (h *handlers) handleOrdersCreate(w http.ResponseWriter, r *http.Request) {
	var input Order

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//var items []*proto.OrderLineItem
	//for i := range input.Items {
	//	item, err := convertItemAPIToProto(&input.Items[i])
	//	if err != nil {
	//		http.Error(w, err.Error(), http.StatusInternalServerError)
	//	}
	//	items = append(items, item)
	//}
	//
	//_, err = h.temporalClient.ExecuteWorkflow(
	//	r.Context(),
	//	client.StartWorkflowOptions{
	//		TaskQueue: "cafe",
	//	},
	//	workflows.Order,
	//	&proto.OrderInput{
	//		Name:         input.Name,
	//		Email:        input.Email,
	//		PaymentToken: "fake",
	//		Items:        items,
	//	},
	//)

	if err != nil {
		log.Printf("failed to start workflow: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func Router(c client.Client) *mux.Router {
	r := mux.NewRouter()

	h := handlers{temporalClient: c}

	r.HandleFunc("/menu", h.handleMenuFetch).Methods("GET").Name("menu_fetch")

	r.HandleFunc("/orders", h.handleOrdersCreate).Methods("POST").Name("orders_create")

	return r
}
