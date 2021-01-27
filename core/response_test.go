package core_test

import (
	"fmt"
	"testing"

	"github.com/schoeppi5/libts/core"
)

func TestUnmarshalResponseMapToStruct(t *testing.T) {
	// given
	m := []map[string]interface{}{{
		"Tag1": "test1",
		"Tag2": "test2",
	}}
	have := struct {
		Tag1 string
		Tag2 string
	}{}

	want := struct {
		Tag1 string
		Tag2 string
	}{
		Tag1: "test1",
		Tag2: "test2",
	}

	// when
	err := core.UnmarshalResponse(m, &have)

	// then
	if err != nil {
		LogTestError(have, want, t, fmt.Sprintf("Test failed with error: %+v", err))
	}
	if have != want {
		LogTestError(have, want, t)
	}
}

func TestUnmarshalResponseMapToStructError(t *testing.T) {
	// given
	m := []map[string]interface{}{{
		"Tag1": "test1",
		"Tag2": "test2",
	}}
	have := struct {
		Tag1 []bool
		Tag2 string
	}{}
	want := "1 error(s) decoding:\n\n* cannot parse 'Tag1[0]' as bool: strconv.ParseBool: parsing \"test1\": invalid syntax"

	// when
	err := core.UnmarshalResponse(m, &have)

	// then
	if err == nil {
		LogTestError(have, want, t, "expected unmarshal error for type []bool")
	}
	if err.Error() != want {
		LogTestError(err.Error(), want, t)
	}
}

func TestUnmarshalResponseMapToArray(t *testing.T) {
	// given
	m := []map[string]interface{}{{
		"Tag1": "test1_1",
		"Tag2": "test1_2",
	}, {
		"Tag1": "test2_1",
		"Tag2": "test2_2",
	}}
	have := [2]struct {
		Tag1 string
		Tag2 string
	}{}
	want := [2]struct {
		Tag1 string
		Tag2 string
	}{{
		Tag1: "test1_1",
		Tag2: "test1_2",
	}, {
		Tag1: "test2_1",
		Tag2: "test2_2",
	}}

	// when
	err := core.UnmarshalResponse(m, &have)

	// then
	if err != nil {
		LogTestError(have, want, t, fmt.Sprintf("test failed with error %s", err))
	}
	for i, v := range have {
		if v.Tag1 != want[i].Tag1 {
			LogTestError(v.Tag1, want[i].Tag1, t)
		}
		if v.Tag2 != want[i].Tag2 {
			LogTestError(v.Tag2, want[i].Tag2, t)
		}
	}
}

func TestUnmarshalResponseMapToShorterArray(t *testing.T) {
	// given
	m := []map[string]interface{}{{
		"Tag1": "test1_1",
		"Tag2": "test1_2",
	}, {
		"Tag1": "test2_1",
		"Tag2": "test2_2",
	}}
	have := [1]struct {
		Tag1 string
		Tag2 string
	}{}
	want := [1]struct {
		Tag1 string
		Tag2 string
	}{{
		Tag1: "test1_1",
		Tag2: "test1_2",
	}}

	// when
	err := core.UnmarshalResponse(m, &have)

	// then
	if err != nil {
		LogTestError(have, want, t, fmt.Sprintf("test failed with error %s", err))
	}
	for i, v := range have {
		if v.Tag1 != want[i].Tag1 {
			LogTestError(v.Tag1, want[i].Tag1, t)
		}
		if v.Tag2 != want[i].Tag2 {
			LogTestError(v.Tag2, want[i].Tag2, t)
		}
	}
}

func TestUnmarshalResponseMapToArrayError(t *testing.T) {
	// given
	m := []map[string]interface{}{{
		"Tag1": "test1_1",
		"Tag2": "test1_2",
	}, {
		"Tag1": "test2_1",
		"Tag2": "test2_2",
	}}
	have := [2]struct {
		Tag1 []bool
		Tag2 string
	}{}
	want := "1 error(s) decoding:\n\n* cannot parse 'Tag1[0]' as bool: strconv.ParseBool: parsing \"test1_1\": invalid syntax"

	// when
	err := core.UnmarshalResponse(m, &have)

	// then
	if err == nil {
		LogTestError(have, want, t, "expected error unmarshaling")
	}
	if err.Error() != want {
		LogTestError(err.Error(), want, t)
	}
}

func TestUnmarshalResponseMapToSlice(t *testing.T) {
	// given
	m := []map[string]interface{}{{
		"Tag1": "test1_1",
		"Tag2": "test1_2",
	}, {
		"Tag1": "test2_1",
		"Tag2": "test2_2",
	}}
	have := []struct {
		Tag1 string
		Tag2 string
	}{}
	want := []struct {
		Tag1 string
		Tag2 string
	}{{
		Tag1: "test1_1",
		Tag2: "test1_2",
	}, {
		Tag1: "test2_1",
		Tag2: "test2_2",
	}}

	// when
	err := core.UnmarshalResponse(m, &have)

	// then
	if err != nil {
		LogTestError(have, want, t, fmt.Sprintf("test failed with error %s", err))
	}
	for i, v := range have {
		if v.Tag1 != want[i].Tag1 {
			LogTestError(v.Tag1, want[i].Tag1, t)
		}
		if v.Tag2 != want[i].Tag2 {
			LogTestError(v.Tag2, want[i].Tag2, t)
		}
	}
}

func TestUnmarshalResponseMapToSliceWithContent(t *testing.T) {
	// given
	m := []map[string]interface{}{{
		"Tag1": "test1_1",
		"Tag2": "test1_2",
	}, {
		"Tag1": "test2_1",
		"Tag2": "test2_2",
	}}
	have := []struct {
		Tag1 string
		Tag2 string
	}{{
		Tag1: "do not",
		Tag2: "override",
	}}
	want := []struct {
		Tag1 string
		Tag2 string
	}{{
		Tag1: "do not",
		Tag2: "override",
	}, {
		Tag1: "test1_1",
		Tag2: "test1_2",
	}, {
		Tag1: "test2_1",
		Tag2: "test2_2",
	}}

	// when
	err := core.UnmarshalResponse(m, &have)

	// then
	if err != nil {
		LogTestError(have, want, t, fmt.Sprintf("test failed with error %s", err))
	}
	for i, v := range have {
		if v.Tag1 != want[i].Tag1 {
			LogTestError(v.Tag1, want[i].Tag1, t)
		}
		if v.Tag2 != want[i].Tag2 {
			LogTestError(v.Tag2, want[i].Tag2, t)
		}
	}
}

func TestUnmarshalResponseMapToSliceError(t *testing.T) {
	// given
	m := []map[string]interface{}{{
		"Tag1": "test1_1",
		"Tag2": "test1_2",
	}, {
		"Tag1": "test2_1",
		"Tag2": "test2_2",
	}}
	have := []struct {
		Tag1 []bool
		Tag2 string
	}{}
	want := "1 error(s) decoding:\n\n* cannot parse 'Tag1[0]' as bool: strconv.ParseBool: parsing \"test1_1\": invalid syntax"

	// when
	err := core.UnmarshalResponse(m, &have)

	// then
	if err == nil {
		LogTestError(have, want, t, "expected error unmarshaling")
	}
	if err.Error() != want {
		LogTestError(have, want, t)
	}
}

func TestUnmarshalResponseMapToUnsupportedTypeError(t *testing.T) {
	// given
	m := []map[string]interface{}{{
		"Tag1": "test",
	}}
	have := ""

	// when
	err := core.UnmarshalResponse(m, &have)

	// then
	if err == nil {
		LogTestError(have, "", t, "expected unmarshal error")
	}
	if err.Error() != "unsupported type string. Expected type struct, slice or array" {
		LogTestError(err.Error(), "unsupported type string. Expected type struct, slice or array", t)
	}
}

func TestDecodeMapToNonPtrError(t *testing.T) {
	// given
	m := map[string]interface{}{
		"Tag1": "test",
	}
	have := struct {
		Tag1 string
	}{}
	want := "expected pointer to value, not value"

	// when
	err := core.Decode(m, have)

	// then
	if err == nil {
		LogTestError(have, want, t, "expected error because of non pointer value")
	}
	if err.Error() != want {
		LogTestError(err.Error(), want, t)
	}
}
