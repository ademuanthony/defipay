package processors

import (
	"context"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestDollarToToken(t *testing.T) {
	amount, valid := common.Big1.SetString("1000000000000000", 10)
	if !valid {
		t.Fail()
	}

	p := dfcProcessor{}

	tokenAMount, err := p.DollarToToken(context.Background(), amount)
	fmt.Println(err)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if tokenAMount.String() == "1" {
		t.Log(tokenAMount)
		t.Fail()
	}
	//1000848565400000000
}

func TestDollarToTokenUsdt(t *testing.T) {
	amount, valid := common.Big1.SetString("1000000000000000", 10)
	if !valid {
		t.Fail()
	}

	p := usdtProcessor{}

	tokenAMount, err := p.DollarToToken(context.Background(), amount)
	fmt.Println(err)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if tokenAMount.String() == "1000848565400000000" {
		t.Log(tokenAMount)
		t.Fail()
	}
}
