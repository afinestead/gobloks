package utilities

import (
	"fmt"
	"testing"
)

func TestColumnGetter(t *testing.T) {
	test := func(n uint64, col, exp uint8) {
		res := getColumn(n, col)
		if res != exp {
			t.Errorf("error getting column %v. expected %b, got %b", col, exp, res)
		}
	}

	test(0x0000000000000000, 0, 0x00)
	test(0x0101010101010101, 0, 0xff)
	test(0x0101010101010101, 1, 0x00)
	test(0x0202020202020202, 0, 0x00)
	test(0x0202020202020202, 1, 0xff)
	test(0x0404040404040201, 0, 0x01)
	test(0x0404040403040202, 1, 0b1011)
}

func TestRowGetter(t *testing.T) {
	test := func(n uint64, row, exp uint8) {
		res := getRow(n, row)
		if res != exp {
			t.Errorf("error getting column %v. expected %b, got %b", row, exp, res)
		}
	}

	test(0x0000000000000000, 0, 0x00)
	test(0x0101010101010101, 0, 0x01)
	test(0x0101010101010101, 1, 0x01)
	test(0x0202020202020202, 0, 0x02)
	test(0x0202020202020202, 1, 0x02)
	test(0x0404040404040201, 2, 0x04)
	test(0x0404040403040202, 1, 0x02)
}

func TestRotator(t *testing.T) {
	testRot := func(n, exp uint64) {
		res := rotate64(n)
		if res != exp {
			t.Errorf("error rotating %b. expected %b, got %b", n, exp, res)
		}
	}

	testRot(0x0101010101010101, 0x00000000000000ff)
	testRot(0x00000000000000ff, 0x0101010101010101)
}

func TestReflector(t *testing.T) {
	testReflect := func(ax Axis, n, exp uint64) {
		var result uint64
		if ax == X {
			result = reflectX64(n)
		} else {
			result = reflectY64(n)
		}
		if result != exp {
			t.Errorf("error reflecting %b. expected %b, got %b", n, exp, result)
		}
	}

	testReflect(X, 0x0101010101010101, 0x8080808080808080)
	testReflect(X, 0x8080808080808080, 0x0101010101010101)
	testReflect(X, 0x00000000000000ff, 0x00000000000000ff)
	testReflect(
		X,
		0b100000011100000100,
		0b010000001110000000100000,
	)

	testReflect(Y, 0x0101010101010101, 0x0101010101010101)
	testReflect(Y, 0xff00000000000000, 0x00000000000000ff)
}

func TestStringifier(t *testing.T) {
	test := func(n uint64, exp string) {
		res := stringify64(n)
		if res != exp {
			t.Errorf("error stringifying %b\nexpected %v\n got %v\n", n, exp, res)
		}
	}

	test(
		0x01010101010101ff,
		"\n########\n#       \n#       \n#       \n#       \n#       \n#       \n#       ",
	)

}

func TestNormalizer(t *testing.T) {
	test := func(n, exp uint64) {
		res := normalize64(n)
		fmt.Printf("%b\n%b\n", n, res)
		if res != exp {
			t.Errorf("error normalizing %b. expected %b, got %b", n, exp, res)
		}
	}

	test(
		0b1000000111000001000000000000000000000000000000000000000000,
		0b0000000000000000000000000000000000000000100000011100000100,
	)
}

func TestRowFlipper(t *testing.T) {
	test := func(n, exp uint8) {
		res := flipRow(n)
		if res != exp {
			t.Errorf("error flipping row %b. expected %b, got %b", n, exp, res)
		}
	}

	test(0b00000001, 0b10000000)
	test(0b00000011, 0b11000000)
}
