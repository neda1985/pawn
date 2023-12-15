package apis

import (
	"dvb_pawn_shop/pkg/input"
	"dvb_pawn_shop/pkg/output"
	"encoding/json"
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
	fmt.Printf("Inventory : %v \n",s.Inventory)

	decoder := json.NewDecoder(r.Body)
	var pawnReq input.PawnRequest
	err := decoder.Decode(&pawnReq)
	if err != nil {
		response := output.PawnResponse{Code: "REJECT",Message: "failed to pars Jason"}
		fmt.Printf("failed to decode request %s", err.Error())
		sendJSONResponse(w, response, http.StatusBadRequest)
		return
	}
	if err := input.ValidateRequest(pawnReq); err != nil {
		response := output.PawnResponse{Code: "REJECT",Message: "invalid json provided"}
		fmt.Printf("failed to validate request %s", err.Error())
		sendJSONResponse(w, response, http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.checkOffer(w, pawnReq)

}

func (s *Shop) checkOffer(w http.ResponseWriter, pawnReq input.PawnRequest) bool {
	fmt.Printf("Demand : %d , Offer : %d \n", pawnReq.Demand, pawnReq.Offer)
	if *pawnReq.Demand > *pawnReq.Offer {
		response := output.PawnResponse{Code: "REJECT", Message: "Sanity failed, Demand should be less than Offer"}
		sendJSONResponse(w, response, http.StatusBadRequest)
		return false
	}

	for i, invValue := range s.Inventory {
		if invValue >= *pawnReq.Demand && invValue < *pawnReq.Offer {
			s.Inventory[i] = *pawnReq.Offer
			response := output.PawnResponse{Code: "ACCEPT",Value: invValue}
			sort.Ints(s.Inventory)
			fmt.Println("modified",s.Inventory)
			sendJSONResponse(w, response, http.StatusOK)
			return true
		}
		fmt.Printf("index : %d , inventory value of %d was not bigger than Demand (%d) and smaller than Offer(%d) at the same time\n",i,invValue,pawnReq.Demand,pawnReq.Offer)
	}
	response := output.PawnResponse{Code: "REJECT",Message: "no item found in inventory to be bigger than Demand and less than offer at the same time"}
	sendJSONResponse(w, response, http.StatusBadRequest)
	return false
}
func sendJSONResponse(w http.ResponseWriter, response interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), status)
	}
}
