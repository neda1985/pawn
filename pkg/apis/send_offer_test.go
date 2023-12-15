package apis

import (
	"dvb_pawn_shop/pkg/input"

	"reflect"
	"sync"
	"testing"
)

func TestShop_checkOffer(t *testing.T) {
	one := 1
	two := 2
	three := 3
	four:=4
	five := 5
	minusTwo := -2
	type fields struct {
		Inventory []int
		mu        sync.Mutex
	}
	type args struct {
		pawnReq input.PawnRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *int
		wantErr bool
	}{
		{
			name: "Accept with return Value 1 and minus demand",
			fields: fields{
				Inventory: []int{1, 1, 1, 1, 1},
			},
			args: args{

				pawnReq: input.PawnRequest{
					Demand: &minusTwo,
					Offer:  &five,
				},
			},
			want:    &one,
			wantErr: false,
		},
		{
			name: "Sanity failed, Demand should be less than Offer",
			fields: fields{
				Inventory: []int{1, 1, 1, 1, 5},
			},
			args: args{
				pawnReq: input.PawnRequest{
					Demand: &three,
					Offer:  &two,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Accept with return Value 1 and demand two",
			fields: fields{
				Inventory: []int{1, 1, 1, 1, 5},
			},
			args: args{
				pawnReq: input.PawnRequest{
					Demand: &one,
					Offer:  &two,
				},
			},
			wantErr: false,
			want:    &one,
		},
		{
			name: "no item found in inventory to be bigger than Demand and less than offer at the same time",
			fields: fields{
				Inventory: []int{1, 1, 1, 1, 5},
			},
			args: args{
				pawnReq: input.PawnRequest{
					Demand: &three,
					Offer:  &four,
				},
			},
			wantErr: true,
			want:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Shop{
				Inventory: tt.fields.Inventory,
				mu:        tt.fields.mu,
			}
			got, err := s.checkOffer(tt.args.pawnReq)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkOffer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("checkOffer() got = %v, want %v", got, tt.want)
			}
		})
	}
}
