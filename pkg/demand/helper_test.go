package demand

import (
	"fmt"
	"reflect"
)

// This is gomock matcher for the DemandExtParams type

func (p *DemandExtParams) Matches(x interface{}) bool {
	reflectedValue := reflect.ValueOf(x).Elem()

	if string(p.Body) != string(reflectedValue.FieldByName("Body").Bytes()) {
		return false
	} else if p.SamFlag != reflectedValue.FieldByName("SamFlag").String() {
		return false
	} else if p.SamHbTag != reflectedValue.FieldByName("SamHbTag").String() {
		return false
	}
	return true
}

func (p *DemandExtParams) String() string {
	return fmt.Sprintf("{DemandExtParams - SamFlag:%s SamHbTag:%s Body:%s}", p.SamFlag, p.SamHbTag, string(p.Body))
}
