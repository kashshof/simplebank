package api

import (
	"testing"

	"github.com/golang/mock/gomock"
	//mockdb "github.com/kashshof/simplebank/db/mock"
	db "github.com/kashshof/simplebank/db/sqlc"
	"github.com/kashshof/simplebank/util"
)

func TestGetAccountAPI(t *testing.T) {
	//account := randomAccount()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	//mockdb.MockStore(ctrl)

}

func randomAccount() db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    util.RandomOwner(),
		Balance:  util.RandomBalance(),
		Currency: util.RandomCurrency(),
	}
}
