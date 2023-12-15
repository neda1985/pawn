package apis

import (
	"dvb_pawn_shop/pkg/input"
	"dvb_pawn_shop/pkg/output"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"sync"
)

type Shop struct {
	Inventory []int
	mu        sync.Mutex
}

func NewShop(size int) *Shop {
	shop := &Shop{
		Inventory: make([]int, size),
	}
	for i := range shop.Inventory {
		shop.Inventory[i] = 1
	}
	return shop
}

func (s *Shop) HandleOffer(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Initial Inventory : %v \n", s.Inventory)
	decoder := json.NewDecoder(r.Body)
	var pawnReq input.PawnRequest
	err := decoder.Decode(&pawnReq)
	if err != nil {
		response := output.PawnResponse{Code: "REJECT", Message: "failed to pars Jason"}
		sendJSONResponse(w, response, http.StatusBadRequest)
		return
	}
	if err := input.ValidateRequest(pawnReq); err != nil {
		response := output.PawnResponse{Code: "REJECT", Message: "invalid json provided"}
		sendJSONResponse(w, response, http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	inventoryValue, err := s.checkOffer(pawnReq)
	if err != nil {
		response := output.PawnResponse{Code: "REJECT", Message: err.Error()}
		sendJSONResponse(w, response, http.StatusBadRequest)

		return
	}
	response := output.PawnResponse{Code: "ACCEPT", Value: *inventoryValue}
	sendJSONResponse(w, response, http.StatusOK)

}

func (s *Shop) checkOffer(pawnReq input.PawnRequest) (*int, error) {
	fmt.Printf("Demand : %d , Offer : %d \n", pawnReq.Demand, pawnReq.Offer)
	if *pawnReq.Demand > *pawnReq.Offer {
		return nil, errors.New("sanity failed, Demand should be less than Offer")
	}
	for i, invValue := range s.Inventory {
		if invValue >= *pawnReq.Demand && invValue < *pawnReq.Offer {
			s.Inventory[i] = *pawnReq.Offer
			sort.Ints(s.Inventory)
			fmt.Println("modified inventory",s.Inventory)
			return &invValue,nil
		}
		fmt.Printf("index : %d , inventory value of %d was not bigger than Demand (%d) and smaller than Offer(%d) at the same time\n", i, invValue, pawnReq.Demand, pawnReq.Offer)
	}
	return nil, errors.New("no item found in inventory to be bigger than Demand and less than offer at the same time")

}
func sendJSONResponse(w http.ResponseWriter, response interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), status)
	}
}
