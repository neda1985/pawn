package apis

import (
	"dvb_pawn_shop/pkg/input"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestShop_checkOffer(t *testing.T) {
	one:=1
	two:=2
	three:=3
	five:=5
	minuseTwo:=-2
	type fields struct {
		Inventory []int
		mu        sync.Mutex
	}
	type args struct {
		w       http.ResponseWriter
		pawnReq input.PawnRequest
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "Accept with Value 1",
			fields: fields{
				Inventory: []int{1,1 , 1, 1,1},
			},
			args: args{
				w: httptest.NewRecorder(),
				pawnReq: input.PawnRequest{
					Demand: &minuseTwo,
					Offer:  &five,
				},
			},
			want: true,
		},
		{
			name: "Sanity failed, Demand should be less than Offer",
			fields: fields{
				Inventory: []int{1, 1, 1, 1,5},
			},
			args: args{
				w: httptest.NewRecorder(),
				pawnReq: input.PawnRequest{
					Demand: &three,
					Offer:  &two,
				},
			},
			want: false,
		},
		{
			name: "Accept with Value 1",
			fields: fields{
				Inventory: []int{1,1,1,1,5},
			},
			args: args{
				w: httptest.NewRecorder(),
				pawnReq: input.PawnRequest{
					Demand: &one,
					Offer:  &two,
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Shop{
				Inventory: tt.fields.Inventory,
				mu:        tt.fields.mu,
			}
			if got := s.checkOffer(tt.args.w, tt.args.pawnReq); got != tt.want {
				t.Errorf("checkOffer() = %v, want %v", got, tt.want)
			}
		})
	}
}
