// Public domain.

package coord_test

import (
	"fmt"
	"math"

	"github.com/soniakeys/coord"
)

func ExampleCart_Add() {
	a1 := &coord.Cart{1, 2, 3}
	a2 := &coord.Cart{7, 2, 0}
	fmt.Printf("%+v\n", *new(coord.Cart).Add(a1, a2))
	// Output:
	// {X:8 Y:4 Z:3}
}

func ExampleCart_Cross() {
	a1 := &coord.Cart{1, 0, 0}
	a2 := &coord.Cart{0, 1, 0}
	fmt.Printf("%+v\n", *new(coord.Cart).Cross(a1, a2))
	// Output:
	// {X:0 Y:0 Z:1}
}

func ExampleCart_Dot() {
	a1 := &coord.Cart{1, 2, 3}
	a2 := &coord.Cart{7, 2, 0}
	fmt.Println(a1.Dot(a2))
	// Output:
	// 11
}

func ExampleCart_FromSphr() {
	c := new(coord.Cart)
	s := &coord.Sphr{Ra: 0, Dec: 30 * math.Pi / 180}
	fmt.Printf("%+.3v\n", *c.FromSphr(s))
	// Output:
	// {X:0.866 Y:0 Z:0.5}
}

func ExampleCart_MulScalar() {
	a := &coord.Cart{4, 1, 0}
	b := 3.
	fmt.Printf("%+v\n", *new(coord.Cart).MulScalar(a, b))
	// Output:
	// {X:12 Y:3 Z:0}
}

func ExampleCart_Mult3() {
	s, c := math.Sincos(30 * math.Pi / 180)
	rm := &coord.M3{ // rotate about X axis
		1, 0, 0,
		0, c, -s,
		0, s, c}
	a := &coord.Cart{0, 1, 0}
	fmt.Printf("%+.3v\n", *new(coord.Cart).Mult3(rm, a))
	// Output:
	// {X:0 Y:0.866 Z:0.5}
}

func ExampleCart_RotateX() {
	a := &coord.Cart{0, 1, 0}
	s, c := math.Sincos(30 * math.Pi / 180)
	fmt.Printf("%+.3v\n", *new(coord.Cart).RotateX(a, s, c))
	// Output:
	// {X:0 Y:0.866 Z:-0.5}
}

func ExampleCart_Square() {
	a := &coord.Cart{1, 2, 3}
	fmt.Println(a.Square())
	// Output:
	// 14
}

func ExampleCart_Sub() {
	a1 := &coord.Cart{8, 4, 3}
	a2 := &coord.Cart{1, 2, 3}
	fmt.Printf("%+v\n", *new(coord.Cart).Sub(a1, a2))
	// Output:
	// {X:7 Y:2 Z:0}
}

func ExampleCartS_FromSphrS() {
	s := coord.SphrS{
		{},
		{Ra: 30 * math.Pi / 180},
		{Dec: 30 * math.Pi / 180},
	}
	for _, c := range new(coord.CartS).FromSphrS(s) {
		fmt.Printf("%+.3v\n", c)
	}
	// Output:
	// {X:1 Y:0 Z:0}
	// {X:0.866 Y:0.5 Z:0}
	// {X:0.866 Y:0 Z:0.5}
}

func ExampleCartS_Mult3S() {
	s, c := math.Sincos(30 * math.Pi / 180)
	rm := &coord.M3{ // rotate about X axis
		1, 0, 0,
		0, c, -s,
		0, s, c}
	a := coord.CartS{
		{1, 0, 0},
		{0, 1, 0},
		{0, 0, 1},
	}
	for _, c := range new(coord.CartS).Mult3S(rm, a) {
		fmt.Printf("%+.3v\n", c)
	}
	// Output:
	// {X:1 Y:0 Z:0}
	// {X:0 Y:0.866 Z:0.5}
	// {X:0 Y:-0.5 Z:0.866}
}

func ExampleM3_Transpose() {
	m := &coord.M3{
		1, 2, 3,
		4, 5, 6,
		7, 8, 9}
	m.Transpose(m)
	fmt.Println(m[:3])
	fmt.Println(m[3:6])
	fmt.Println(m[6:])
	// Output:
	// [1 4 7]
	// [2 5 8]
	// [3 6 9]
}

func ExampleSphr_FromCart() {
	c := &coord.Cart{X: math.Sqrt(3) / 2, Z: 1. / 2}
	s := new(coord.Sphr).FromCart(c)
	fmt.Printf("RA:  %3.0f\n", s.Ra*180/math.Pi)
	fmt.Printf("Dec: %3.0f\n", s.Dec*180/math.Pi)
	// Output:
	// RA:    0
	// Dec:  30
}

func ExampleSphrS_FromCartS() {
	c := coord.CartS{
		{1, 0, 0},
		{math.Sqrt(3) / 2, 1. / 2, 0},
		{math.Sqrt(3) / 2, 0, 1. / 2},
	}
	for _, s := range new(coord.SphrS).FromCartS(c) {
		fmt.Printf("RA %3.0f, Dec %3.0f\n",
			s.Ra*180/math.Pi, s.Dec*180/math.Pi)
	}
	// Output:
	// RA   0, Dec   0
	// RA  30, Dec   0
	// RA   0, Dec  30
}
