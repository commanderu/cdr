package cdrutil_test

import (
	"fmt"
	"math"

	"github.com/commanderu/cdrd/cdrutil"
)

func ExampleAmount() {

	a := cdrutil.Amount(0)
	fmt.Println("Zero Atom:", a)

	a = cdrutil.Amount(1e8)
	fmt.Println("100,000,000 Atoms:", a)

	a = cdrutil.Amount(1e5)
	fmt.Println("100,000 Atoms:", a)
	// Output:
	// Zero Atom: 0 cdr
	// 100,000,000 Atoms: 1 cdr
	// 100,000 Atoms: 0.001 cdr
}

func ExampleNewAmount() {
	amountOne, err := cdrutil.NewAmount(1)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(amountOne) //Output 1

	amountFraction, err := cdrutil.NewAmount(0.01234567)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(amountFraction) //Output 2

	amountZero, err := cdrutil.NewAmount(0)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(amountZero) //Output 3

	amountNaN, err := cdrutil.NewAmount(math.NaN())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(amountNaN) //Output 4

	// Output: 1 cdr
	// 0.01234567 cdr
	// 0 cdr
	// invalid coin amount
}

func ExampleAmount_unitConversions() {
	amount := cdrutil.Amount(44433322211100)

	fmt.Println("Atom to kCoin:", amount.Format(cdrutil.AmountKiloCoin))
	fmt.Println("Atom to Coin:", amount)
	fmt.Println("Atom to MilliCoin:", amount.Format(cdrutil.AmountMilliCoin))
	fmt.Println("Atom to MicroCoin:", amount.Format(cdrutil.AmountMicroCoin))
	fmt.Println("Atom to Atom:", amount.Format(cdrutil.AmountAtom))

	// Output:
	// Atom to kCoin: 444.333222111 kcdr
	// Atom to Coin: 444333.222111 cdr
	// Atom to MilliCoin: 444333222.111 mcdr
	// Atom to MicroCoin: 444333222111 μcdr
	// Atom to Atom: 44433322211100 Atom
}
