package libovsdb

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sort"
	"strings"
	"testing"
)

/*
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

func emptyOvsSet()OvsSet{
	return newOvsSetOrPanic([]string{})
}

func emptyOvsMap()OvsMap{
	return newOvsMapOrPanic(map[string]string{})
}

type tuple struct {
	obj interface{}
	json string
}

var tests = []tuple{
	{emptyOvsSet(),`["set",[]]`},
	{,`["set",[]]`},
}
*/

func jsonEq(t *testing.T,str1 string,str2 string) {
	_,eq := jsonEqReq(t,str1,str2,false,false)
	if !eq{
		require.Fail(t,fmt.Sprintf("jsons \n%s\n%s\n are not equal",str1,str2))
	}
	//return res
}

func jsonEqReq(t *testing.T,str1 string,str2 string,isMap bool,isSet bool)(string,bool){
	// TODO better names
	findTheCloseIdxForElement :=func(s string)int{
		squareBracketsBalance := 2
		l := len(s)
		for i:=0 ; i < l ; i++ {
			switch string(s[i]){
			case `[`:
				squareBracketsBalance++
			case `]`:
				squareBracketsBalance--
			}
			if squareBracketsBalance==0{
				return i
			}
		}
		require.Error(t,errors.New(`Unbalanced square brackets (there are no too extra "]" as expected`))
		return -1 //we should not get here
	}
	sortList := func(s string)string{
		//recive a json list and sort it ass
		if string(s[0]) != `[` || string(s[len(s)-1]) != `]`{
			require.Error(t,errors.New(`not a valid json list`))
			return "" //we should not get here
		}
		innerListStr := s[1:len(s)-1]
		list := strings.Split(innerListStr,`,`)
		/*
		if len(list)==1{
			return `[` +list[0]+ `]`
		}
		*/
		sort.Strings(list)
		res:=list[0]
		for _ ,elem :=range list[1:]{
			res += `,`+elem
		}
		return `[` +res+ `]`
	}

	//uncomment after debugg
	/*
	const setPrefix = `["set",[`
	const mapPrefix = `["map",[`
	const shouldEqualMsg = "they should be equal\n"
	const setPrefixLen = len(setPrefix)
	const mapPrefixLen = len(mapPrefix)
	*/
	//non const for debug
	setPrefix := `["set",[`
	mapPrefix := `["map",[`
	shouldEqualMsg := "they should be equal\n"
	setPrefixLen := len(setPrefix)
	mapPrefixLen := len(mapPrefix)

	inMapOrSet := isSet || isMap

	if isMap && isSet{
		require.Error(t,errors.New(`cannot be set and map on the same time`)) // this is bug of in my program
		return "",false
	}//TODO remove the panic in the future and find better alternative

	if len(str1) != len(str2){return "",false}
	if inMapOrSet{ //TODO make code more elegant
		//buid
		e1 := findTheCloseIdxForElement(str1)
		e2 := findTheCloseIdxForElement(str2)
		if e1 != e2{return "",false}
		str1 = str1[:e1-1]
		str2 = str2[:e2-1]
		//now we will take care inner sets or maps
	}
	l := len(str1)
	for i:=0 ; i < l ; i++{
		if !inMapOrSet{
			if str1[i] != str2[i]{
				return "",false
			}
		}
		if i <setPrefixLen || i< mapPrefixLen{
			continue
		}
		isEqual := true
		var resStr string
		if str1[i-setPrefixLen:i]==setPrefix {
			resStr, isEqual = jsonEqReq(t,string(str1[i:]), string(str2[i:]), false, true)
		}else if str1[i-mapPrefixLen:i]==mapPrefix{
			resStr, isEqual = jsonEqReq(t,string(str1[i:]), string(str2[i:]), true, false)
		}
		if resStr != "" && isEqual{ //maybe find better condition
			str1 = str1[0:i-setPrefixLen] + resStr+ str1[i+len(resStr):]
			str2 = str2[0:i-setPrefixLen] + resStr+ str2[i+len(resStr):]
			l = len(str1)
			i=i+len(resStr)
		}
	}
	if inMapOrSet{
		//todo (maybe we need to convert string to proper json at this point.it dependes on implementation of JSONEqf
		if isSet{
			str1 = sortList(`[` +str1 +`]`)
			str2 = sortList(`[` +str2 +`]`)
			assert.JSONEqf(t, str1, str2,shouldEqualMsg)
			return str1,true
		}else {
		//is map
			str1 = `{` +str1 +`}`
			str2 = `{` +str2 +`}`
			//TODO need to add more string proccessing to be equal to map
			assert.JSONEqf(t, str1, str2,shouldEqualMsg)
			return str1,true

		}
	}
	return "",true
}


//omer test tmp
func TestOmerTmp2(t *testing.T) {
	//ovsSet, err := NewOvsSet([]string{"aa", "bb"})
	//ovsSet, err := NewOvsSet([]string{"bb", "aa"})
	//assert.Nil(t, err)
	//jsonStr, err := json.Marshal(ovsSet)
	//assert.Nil(t, err)
	//expected := "[\"set\",[\"aa\",\"bb\"]]"
	s1 := `["set",["aa","bb"]]`
	s2 := `["set",["bb","aa"]]`

	jsonEq(t,s1,s2)

	s1 = `["set",["ac","aa","bb"]]`
	s2 = `["set",["bb","ac","aa"]]`

	//require.JSONEqf(t, expected, string(jsonStr), "they should be equal\n")
	jsonEq(t,s1,s2)


	s1 = `["set",["ac",["set",["aa","bb"]],"aa","bb"]]`
	s2 = `["set",[["set",["bb","aa"]],"bb","ac","aa"]]`
	//TODO this should pass

	//require.JSONEqf(t, expected, string(jsonStr), "they should be equal\n")
	jsonEq(t,s1,s2)
}
/*
//omer test tmp
func TestOmerTmp(t *testing.T) {
	//ovsSet, err := NewOvsSet([]string{"aa", "bb"})
	ovsSet, err := NewOvsSet([]string{"bb", "aa"})
	assert.Nil(t, err)
	jsonStr, err := json.Marshal(ovsSet)
	assert.Nil(t, err)
	expected := "[\"set\",[\"aa\",\"bb\"]]"
	require.JSONEqf(t, expected, string(jsonStr), "they should be equal\n")
}

 */

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