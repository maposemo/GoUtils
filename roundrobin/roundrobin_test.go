package roundrobin

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

const MaxSliceLen int = 8

type testInf struct {
}

var testClient1 = testInf{}
var testClient2 = testInf{}

func TestAdd(t *testing.T) {
	element1 := "Add: {0}"
	element2 := "Add: {0, 1}"

	testCases := []struct {
		testCase string
		exp      int
	}{
		{testCase: element1, exp: 1},
		{testCase: element2, exp: 2},
	}

	for _, tc := range testCases {
		rr := new(RoundRobin)

		switch tc.testCase {
		case element1:
			rr.Add(&testInf{})
		case element2:
			rr.Add(&testInf{})
			rr.Add(&testInf{})
		}

		if len(rr.array) != tc.exp {
			t.Errorf("%s test failed(%v)", tc.testCase, rr.array)
		}
	}
}

func TestRemove(t *testing.T) {
	remove1st := "Remove 1st"
	remove2nd := "Remove 2nd"

	testcases := []struct {
		testCase string
		exp      *testInf
	}{
		{testCase: remove1st, exp: &testClient2},
		{testCase: remove2nd, exp: &testClient1},
	}

	for _, tc := range testcases {
		rr := new(RoundRobin)
		rr.Add(&testClient1)
		rr.Add(&testClient2)

		switch tc.testCase {
		case remove1st:
			rr.Remove(&testClient1)
		case remove2nd:
			rr.Remove(&testClient2)
		}

		if rr.array[0] != tc.exp {
			t.Errorf("%s test failed(%v)", tc.testCase, rr.array)
		}
	}
}

func TestNext(t *testing.T) {
	element0 := "Element: 0, CurrentIndex: 12"
	element1curIndex0 := "Element: 1, CurrentIndex: 0"
	element2curIndex0 := "Element: 2, CurrentIndex: 0"
	element2curIndex1 := "Element: 2, CurrentIndex: 1"

	testcases := []struct {
		testCase string
		exp      *testInf
	}{
		{testCase: element0},
		{testCase: element1curIndex0, exp: &testClient1},
		{testCase: element2curIndex0, exp: &testClient1},
		{testCase: element2curIndex1, exp: &testClient1},
	}

	for _, tc := range testcases {
		rr := new(RoundRobin)
		switch tc.testCase {
		case element0:
			rr.currentIdx = 12
		case element1curIndex0:
			rr.Add(&testClient1)
		case element2curIndex0:
			rr.Add(&testClient1)
			rr.Add(&testClient2)
		case element2curIndex1:
			rr.Add(&testClient1)
			rr.Add(&testClient2)
			rr.currentIdx = 1
		}

		result := rr.Next()

		if tc.testCase == element0 {
			if result != nil {
				t.Errorf("%s test failed(%v)", tc.testCase, result)
			}
		} else {
			if !cmp.Equal(result, tc.exp, cmp.AllowUnexported(testInf{})) {
				t.Errorf("%s test failed(%s)", tc.testCase, cmp.Diff(result, tc.exp))
			}
		}
	}
}
