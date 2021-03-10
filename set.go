package libovsdb

import (
	"encoding/json"
	"errors"
	"reflect"
)

// OvsSet is an OVSDB style set
// RFC 7047 has a wierd (but understandable) notation for set as described as :
// Either an <atom>, representing a set with exactly one element, or
// a 2-element JSON array that represents a database set value.  The
// first element of the array must be the string "set", and the
// second element must be an array of zero or more <atom>s giving the
// values in the set.  All of the <atom>s must have the same type.
type OvsSet struct {
	GoSet []interface{}
}

// NewOvsSet creates a new OVSDB style set from a Go slice
func NewOvsSet(i interface{}) (*OvsSet, error) {
	v := reflect.ValueOf(i)
	var ovsSet []interface{}
	switch v.Kind() {
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			ovsSet = append(ovsSet, v.Index(i).Interface())
		}
	case reflect.String:
		ovsSet = append(ovsSet, v.Interface())
	case reflect.ValueOf(UUID{}).Kind():
		ovsSet = append(ovsSet, v.Interface())
	default:
		return nil, errors.New("OvsSet supports only Go Slice/string/uuid types")
	}
	return &OvsSet{ovsSet}, nil
}

// MarshalJSON wil marshal an OVSDB style set in to a JSON byte array
func (o OvsSet) MarshalJSON() ([]byte, error) {
	if len(o.GoSet) == 1 {
		return json.Marshal(o.GoSet[0])
	} else if len(o.GoSet) > 1 {
		var oSet []interface{}
		oSet = append(oSet, "set")
		oSet = append(oSet, o.GoSet)
		return json.Marshal(oSet)
	}
	return []byte("[\"set\",[]]"), nil
}

// UnmarshalJSON will unmarshal a JSON byte array to an OVSDB style set
func (o *OvsSet) UnmarshalJSON(b []byte) (err error) {
	addToSet := func(o *OvsSet,v interface{}) {
		goVal, err := ovsSliceToGoNotation(v)
		if err == nil {
			o.GoSet = append(o.GoSet, goVal)
		}
	}

	var inter interface{}
	if err = json.Unmarshal(b, &inter); err != nil {
		return err
	}
	switch inter.(type) {
	case []interface{}:
		var oSet []interface{}
		oSet = inter.([]interface{})
		switch  oSet[1].(type){
		case []interface{}:
			innerSet := oSet[1].([]interface{})
			for _, val := range innerSet {
				switch  val.(type){
				case []interface{}:
					// it is a uuid object
					addToSet(o,UUID{GoUUID: val.([]interface{})[1].(string)})
				default:
					addToSet(o,val)
				}
			}
		default:
			// it is a single uuid object
			addToSet(o, UUID{GoUUID: oSet[1].(string)})
		}
	default:
		// it is a single object
		addToSet(o,inter)
	}
	return err
}
