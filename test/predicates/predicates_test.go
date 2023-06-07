package predicates

import "testing"

// Go1.7: The Method and NumMethod methods of Type and Value no longer return or count unexported methods.
// So cannot use reflect.TypeOf and MethodByName to test the implemented methods.
func TestPredicatesArgs(t *testing.T) {
	var cur any = &current{}
	_, ok := cur.(interface {
		onA5(any) (bool, error)
		onA9(any) (bool, error)
		onA13(any) (bool, error)
		onB9(any) (bool, error)
		onB10(any) (bool, error)
		onB11(any) (bool, error)
		onC1(any) (any, error)
	})
	if !ok {
		t.Errorf("want *current to have the expected methods")
	}
}
