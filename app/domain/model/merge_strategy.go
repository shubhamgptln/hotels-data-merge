package model

import (
	"fmt"
	"reflect"
	"strings"
)

type MergeStrategy string

const (
	MajorityNonZero MergeStrategy = "majority-filled"
	Longest         MergeStrategy = "longest"
	AppendUnique    MergeStrategy = "append-unique"
	Append          MergeStrategy = "append"
)

func (ms *MergeStrategy) LongestString(data []string) string {
	maxLen := 0
	var result string
	for _, str := range data {
		if len(str) > maxLen {
			maxLen = len(str)
			result = str
		}
	}
	return result
}

// Append only unique string irrespective of lowercase or uppercase
func (ms *MergeStrategy) AppendUniqueEntries(data []string) []string {
	freq := make(map[string]int, 0)
	var result []string
	for _, str := range data {
		_, ok := freq[strings.ToLower(str)]
		if ok {
			continue
		}
		result = append(result, str)
		freq[str]++
	}
	return result
}

func (ms *MergeStrategy) MajorityNonZeroField(data []interface{}) interface{} {
	maxMajority := 0
	var winner interface{}
	for _, st := range data {
		val := ReadStruct(st)
		if val > maxMajority {
			maxMajority = val
			winner = st
		}
	}
	fmt.Printf("debug Location %v", winner)
	return winner
}

func (ms *MergeStrategy) AssignFirstNonZeroField(data []interface{}, fieldName string) interface{} {
	for _, st := range data {
		val := reflect.ValueOf(st)
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		for i := 0; i < val.NumField(); i++ {
			if val.Type().Field(i).Name == fieldName && val.Field(i).IsZero() {
				return st
			}
		}
	}
	return nil
}

func ReadStruct(st interface{}) int {
	cnt := 0
	val := reflect.ValueOf(st)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	for i := 0; i < val.NumField(); i++ {
		f := val.Field(i)
		switch f.Kind() {
		case reflect.Struct:
			cnt += ReadStruct(f.Interface())
		default:
			if !f.IsZero() {
				cnt++
			}
		}
	}
	return cnt
}
