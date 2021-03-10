package libovsdb

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type marshalSetTestTuple struct {
	objInput           interface{}
	jsonExpectedOutput string
}

type marshalMapsTestTuple struct {
	objInput           map[string]string
	jsonExpectedOutput string
}

var validUuidStr0 = `00000000-0000-0000-0000-000000000000`
var validUuidStr1 = `11111111-1111-1111-1111-111111111111`
var validUuid0 = UUID{GoUUID: validUuidStr0}
var validUuid1 = UUID{GoUUID: validUuidStr1}

var setTestList = []marshalSetTestTuple{
	{
		objInput:           []string{},
		jsonExpectedOutput: `["set",[]]`,
	},

	{
		objInput:           `aa`,
		jsonExpectedOutput: `"aa"`,
	},

	{
		objInput:           []string{`aa`},
		jsonExpectedOutput: `"aa"`,
	},

	{
		objInput:           []string{`aa`, `bb`},
		jsonExpectedOutput: `["set",["aa","bb"]]`,
	},

	{
		objInput:           []UUID{},
		jsonExpectedOutput: `["set",[]]`,
	},

	{
		objInput:           UUID{GoUUID: `aa`},
		jsonExpectedOutput: `["named-uuid","aa"]`,
	},

	{
		objInput:           []UUID{UUID{GoUUID: `aa`}},
		jsonExpectedOutput: `["named-uuid","aa"]`,
	},

	{
		objInput:           []UUID{UUID{GoUUID: `aa`}, UUID{GoUUID: `bb`}},
		jsonExpectedOutput: `["set",[["named-uuid","aa"],["named-uuid","bb"]]]`,
	},

	{
		objInput:           validUuid0,
		jsonExpectedOutput: fmt.Sprintf(`["uuid","%v"]`, validUuidStr0),
	},

	{
		objInput:           []UUID{validUuid0},
		jsonExpectedOutput: fmt.Sprintf(`["uuid","%v"]`, validUuidStr0),
	},

	{
		objInput:           []UUID{validUuid0, validUuid1},
		jsonExpectedOutput: fmt.Sprintf(`["set",[["uuid","%v"],["uuid","%v"]]]`, validUuidStr0, validUuidStr1),
	},
}

var mapTestList = []marshalMapsTestTuple{
	{
		objInput:           map[string]string{},
		jsonExpectedOutput: `["map",[]]`,
	},

	{
		objInput:           map[string]string{`v0`: `k0`},
		jsonExpectedOutput: `["map",[["v0","k0"]]]`,
	},

	{
		objInput:           map[string]string{`v0`: `k0`, `v1`: `k1`},
		jsonExpectedOutput: `["map",[["v0","k0"],["v1","k1"]]]`,
	},
}

func setsAreEqual(t *testing.T, set1 OvsSet, set2 OvsSet) {
	res1 := map[interface{}]bool{}
	for _, elem := range set1.GoSet {
		res1[elem] = true
	}

	res2 := map[interface{}]bool{}
	for _, elem := range set2.GoSet {
		res2[elem] = true
	}

	assert.Equal(t, res1, res2, "they should be equal\n")
}

func TestMarshalSet(t *testing.T) {

	for _, e := range setTestList {
		set, err := NewOvsSet(e.objInput)
		assert.Nil(t, err)
		jsonStr, err := json.Marshal(set)
		assert.Nil(t, err)
		assert.JSONEqf(t, e.jsonExpectedOutput, string(jsonStr), "they should be equal\n")
	}

}

func TestMarshalMap(t *testing.T) {

	for _, e := range mapTestList {
		m, err := NewOvsMap(e.objInput)
		assert.Nil(t, err)
		jsonStr, err := json.Marshal(m)
		assert.Nil(t, err)
		assert.JSONEqf(t, e.jsonExpectedOutput, string(jsonStr), "they should be equal\n")
	}

}

func TestUnmarshalSet(t *testing.T) {

	for _, e := range setTestList {
		set, err := NewOvsSet(e.objInput)
		assert.Nil(t, err)
		jsonStr, err := json.Marshal(set)
		assert.Nil(t, err)
		var res OvsSet
		err = json.Unmarshal(jsonStr, &res)
		assert.Nil(t, err)
		setsAreEqual(t, *set, res)
	}

}

func TestUnmarshalMap(t *testing.T) {

	for _, e := range mapTestList {
		m, err := NewOvsMap(e.objInput)
		assert.Nil(t, err)
		jsonStr, err := json.Marshal(m)
		assert.Nil(t, err)
		var res OvsMap
		err = json.Unmarshal(jsonStr, &res)
		assert.Nil(t, err)
		assert.Equal(t, *m, res, "they should be equal\n")
	}

}
