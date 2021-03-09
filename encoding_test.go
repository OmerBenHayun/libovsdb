package libovsdb

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)
// newOvsSetOrPanic is a convenience wrapper over NewOvsSet
func newOvsSetOrPanic(s interface{})OvsSet{
	// TODO add better panic explanations

	//make sure that the input is a slice
	if reflect.TypeOf(s).Kind() != reflect.Slice{
		panic(s)
	}
	res, err := NewOvsSet(s)
	if err !=nil{
		panic(s)
	}
	return *res
}

// newOvsMapOrPanic is a convenience wrapper over NewOvsMap
func newOvsMapOrPanic(v interface{})OvsMap{
	// TODO add better panic explanations
	//make sure that the input is a map
	if reflect.TypeOf(v).Kind() != reflect.Map{
		panic(v)
	}
	res, err := NewOvsMap(v)
	if err !=nil{
		panic(v)
	}
	return *res
}

// empty Set test
func TestEmptySet(t *testing.T) {
	emptySet, err := NewOvsSet([]string{})
	assert.Nil(t, err)
	jsonStr, err := json.Marshal(emptySet)
	assert.Nil(t, err)
	expected := "[\"set\",[]]"
	assert.JSONEqf(t, expected, string(jsonStr), "they should be equal\n")
}

// test Set
func TestSet(t *testing.T) {
	ovsSet, err := NewOvsSet([]string{"aa", "bb"})
	assert.Nil(t, err)
	jsonStr, err := json.Marshal(ovsSet)
	assert.Nil(t, err)
	expected := "[\"set\",[\"aa\",\"bb\"]]"
	require.JSONEqf(t, expected, string(jsonStr), "they should be equal\n")
}

// empty Map test
func TestEmptyMap(t *testing.T) {
	emptyMap, err := NewOvsMap(map[string]string{})
	assert.Nil(t, err)
	jsonStr, err := json.Marshal(emptyMap)
	assert.Nil(t, err)
	expected := "[\"map\",[]]"
	assert.JSONEqf(t, expected, string(jsonStr), "they should be equal\n")
}

// test Map
func TestMap(t *testing.T) {
	ovsMap, err := NewOvsMap(map[string]string{"one": "first", "two": "second"})
	assert.Nil(t, err)
	jsonStr, err := json.Marshal(ovsMap)
	assert.Nil(t, err)
	expected := "[\"map\",[[\"one\",\"first\"],[\"two\",\"second\"]]]"
	assert.JSONEqf(t, expected, string(jsonStr), "they should be equal\n")
}
/*
type tuple struct {
	obj interface{}
	json string
}
var flagtests = []tuple {
	{"%a", "[%a]"},
	{"%-a", "[%-a]"},
	{"%+a", "[%+a]"},
	{"%#a", "[%#a]"},
	{"% a", "[% a]"},
	{"%0a", "[%0a]"},
	{"%1.2a", "[%1.2a]"},
	{"%-1.2a", "[%-1.2a]"},
	{"%+1.2a", "[%+1.2a]"},
	{"%-+1.2a", "[%+-1.2a]"},
	{"%-+1.2abc", "[%+-1.2a]bc"},
	{"%-1.2abc", "[%-1.2a]bc"},
}

 */

//func parserTest(t *testing.T)