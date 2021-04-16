package simple

import (
	"fmt"
	"github.com/xunull/goc/commonx"
	"testing"
)

type TOne struct {
	Name string
	Age  int
}

type TTwo struct {
	Name   string
	Height int
	Age    int
}

func Test1(t *testing.T) {
	a := TOne{Name: "a", Age: 18}
	b := TTwo{Name: "b", Height: 20}

	err := Merge(&b, &a)
	commonx.CheckErrOrFatal(err)
	fmt.Printf("%+v\n", b)
}

func Test2(t *testing.T) {
	a := TOne{Name: "a", Age: 0}
	b := TTwo{Name: "b", Height: 20}

	err := Merge(&b, &a)
	commonx.CheckErrOrFatal(err)
	fmt.Printf("%+v\n", b)
}

func Test3(t *testing.T) {
	a := TOne{Name: "a"}
	b := TTwo{Name: "b", Height: 20}

	err := Merge(&b, &a)
	commonx.CheckErrOrFatal(err)
	fmt.Printf("%+v\n", b)
}

func Test4(t *testing.T) {
	a := TOne{Name: "a"}
	b := TTwo{Name: "b", Height: 20, Age: 10}

	err := Merge(&b, &a)
	commonx.CheckErrOrFatal(err)
	fmt.Printf("%+v\n", b)
}

type TOneOne struct {
	Name *string
	Age  *int
}

type TTwoTwo struct {
	Name   *string
	Height *int
	Age    *int
}

func Test44(t *testing.T) {

}

func Test55(t *testing.T) {

}
