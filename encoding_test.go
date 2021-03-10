package libovsdb

import (
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sort"
	"strings"
	"testing"
)


func getElemFromString(s string)(interface{} ,error){

	setPrefix := `["set",[`
	mapPrefix := `["map",[`
	var getElem  func(s string)(interface{},error)
	hasDoubleQuotationMark := func(s string)bool{
		return strings.HasPrefix(s,`"`) && strings.HasSuffix(s,`"`)
	}

	removeDoubleQuotationMark := func(s string)string{
		res := strings.TrimPrefix(s,`"`)
		res = strings.TrimSuffix(res,`"`)
		return res
	}
	split := func(s string)([]string,error){
		/*
		isSet :=false
		isMap :=false
		if strings.HasPrefix(s,setPrefix) {
			isSet = true
			s = strings.TrimPrefix(s,setPrefix)
		}else if strings.HasPrefix(s,mapPrefix) {
			isMap = true
			s = strings.TrimPrefix(s,mapPrefix)
		}else if strings.HasPrefix(s,mapPrefix) {
		}else {
			nil,errors.New("error split")
		}
		s=strings.TrimSuffix(s,`]]`)
		 */

		res := []string{}
		bracketCounter := 0
		l := len(s)
		for i:=0 ; i < l ; i++{
			if(string(s[i]) == `[`) {
				bracketCounter++
			} else if(string(s[i]) == `]`) {
				bracketCounter--
			} else if(string(s[i]) ==`,` && bracketCounter>0) {
				s= s[:i]+`~`+s[i+1:]
			}
		}
		res = strings.Split(s,`,`)
		for i := range(res) {
			if (strings.HasPrefix(res[i],setPrefix)||
				strings.HasPrefix(res[i],mapPrefix)) {
				res[i] = strings.Replace(res[i],`~`,`,`,-1)
			}
		}
		return res,nil
	}

	splitAndsort := func(s string)([]string,error){
		res,err := split(s)
		if err != nil{
			return nil, err
		}
		sort.Strings(res)
		return res,nil
	}

	getElemFromMap := func(s string)(map[interface{}]interface{},error){
		m := map[interface{}]interface{}{}
		if s ==""{
			return map[interface{}]interface{}{},nil
		}
		//assuming input from the form ["key1","value1"],["key2","value2"]
		if !(strings.HasPrefix(s,`[`) && strings.HasPrefix(s,`]`)){
			return map[interface{}]interface{}{},errors.New(s+"is not a map")
		}
		s= s[1:len(s)-1] //remove outer `[` `]`
		//list := strings.Split(s,`],[`)
		list,err := splitAndsort(s)
		if err!=nil{
			return nil, err
		}
		for _,elem := range(list){
			keyval := strings.Split(elem,`,`)
			if len(keyval) != 2{
				return map[interface{}]interface{}{},errors.New(s+"is not a map")
			}
			k , err := getElem(keyval[0])
			if err !=nil{
				return map[interface{}]interface{}{} , nil
			}
			v , err := getElem(keyval[1])
			if err !=nil{
				return map[interface{}]interface{}{} , nil
			}
			m[k]=v
		}
		return m,nil
	}
	getElemFromSet := func(s string)(map[interface{}]bool,error){
		//var r map[interface{}]bool
		r := map[interface{}]bool{}
		if s ==""{
			return r,nil
		}
		//assuming input from the form ["s","s2","s3","s4"]
		/*
		if !(strings.HasPrefix(s,`[`) && strings.HasPrefix(s,`]`)){
			return r,errors.New(s+"is not a set")
		}
		s= s[1:len(s)-1] //remove outer `[` `]`
		*/
		list,err := splitAndsort(s)
		if err!=nil{
			return nil, err
		}
		//list := strings.Split(s,`,`) //TODO think how to make it smarter split!
		//sort.Slice(list) //FIXME ADD SORT HERE!!
		sort.Strings(list)
		for _,val := range(list){
			v,err :=getElem(val)
			if err !=nil{
				return map[interface{}]bool{},err
			}
			r[v] = true
		}
		return r,nil
	}
	getElem = func(s string)(interface{},error){
		if(strings.HasPrefix(s,setPrefix)){
			return getElemFromSet(s[len(setPrefix):len(s)-2])
		}else if(strings.HasPrefix(s,mapPrefix)){
			return getElemFromMap(s[len(mapPrefix):len(s)-2])
		}else if(hasDoubleQuotationMark(s)){
			return removeDoubleQuotationMark(s) ,nil
		}else {
			return "", errors.New("err")
		}
	}




	return getElem(s)
}

//omer test tmp
func TestOmerTmp2(t *testing.T) {
	//ovsSet, err := NewOvsSet([]string{"aa", "bb"})
	//ovsSet, err := NewOvsSet([]string{"bb", "aa"})
	//assert.Nil(t, err)
	//jsonStr, err := json.Marshal(ovsSet)
	//assert.Nil(t, err)
	//expected := "[\"set\",[\"aa\",\"bb\"]]"

	//s1 := `["set",["aa","bb"]]`
	//s2 := `["set",["bb","aa"]]`

	//jsonEq(t,s1,s2)

	//s1 = `["set",["ac","aa","bb"]]`
	//s2 = `["set",["bb","ac","aa"]]`

	////require.JSONEqf(t, expected, string(jsonStr), "they should be equal\n")
	//jsonEq(t,s1,s2)
/*
	f0 := `["map",[["one","first"],["two","second"]]]`
	f1 := `["map",[["two","second"],["one","first"]]]`

	m0,err:=toCanonicalForm(f0)
	if err!=nil{
		print(err)
	}

	m1,err:=toCanonicalForm(f1)
	if err!=nil{
		print(err)
	}
	print(m0)
	print(m1)

	s0 := `["set",["ac","aa","bb",["set",["aa","bb"]]]]`
	s1 := `["set",["ac",["set",["aa","bb"]],"aa","bb"]]`
	s2 := `["set",[["set",["bb","aa"]],"bb","ac","aa"]]`
	//TODO this should pass
	a0,err:=toCanonicalForm(s0)
	if err!=nil{
		print(err)
	}
	a1,err:=toCanonicalForm(s1)
	if err!=nil{
		print(err)
	}
	a2,err:=toCanonicalForm(s2)
	if err!=nil{
		print(err)
	}
*/
	//s3 := `["set",[["set",[["map",[["one","first"],["two","second"]]],"aa"]],"bb","ac","aa"]]`
	//s4 := `["set",[["set",["aa",["map",[["one","first"],["two","second"]]]]],"bb","ac","aa"]]`

	//s3 := `["set",[["set",[["set",["one","first","two","second"]],"aa"]],"bb","ac","aa"]]`
	//s4 := `["set",[["set",["aa",["set",["one","first","two","second"]]]],"bb","ac","aa"]]`

	//s3 := `["set",["ac",["set",[["set",["ac","aa","bb"]],"bb","ac","aa"]],"bb","y"]]`
	s3 :=`["set",["aa","bb"]]`
	//s4 := `["set",["bb","y","ac",["set",["ac","aa","bb",["set",["bb","ac","aa"]]]]]]`
	s4 := `["set",["bb","y","ac",["set",["ac","aa","bb"]]]]`
	//["set",["ac","aa","bb"]]
	//["set",["bb","ac","aa"]]
	a3,err:=getElemFromString(s3)
	if err!=nil{
		print(err)
	}
	a4,err:=getElemFromString(s4)
	if err!=nil{
		print(err)
	}
	//print(a0)
	//print(a1)
	//print(a2)
	print(a3)
	print(a4)
	//print(a2)
	//require.JSONEqf(t, expected, string(jsonStr), "they should be equal\n")
	//jsonEq(t,s1,s2)
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