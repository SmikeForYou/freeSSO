package utils

import (
	"github.com/stretchr/testify/assert"
	"os"
	"strconv"
	"testing"
)

func Test_mapEnvToStruct(t *testing.T) {
	testStringValue := "TESTSTRING"
	var testUintValue uint16 = 16555
	type testStruct struct {
		TestString string `env:"TEST_STRING"`
		TestInt    uint16 `env:"TEST_UINT16"`
		TestBool   bool   `env:"TEST_BOOL"`
		TestSkip   string
	}
	_ = os.Setenv("TEST_STRING", testStringValue)
	_ = os.Setenv("TEST_UINT16", strconv.Itoa(int(testUintValue)))
	_ = os.Setenv("TEST_BOOL", "1")
	_ = os.Setenv("TEST_SKIP", "skipped")
	ts := testStruct{}
	_ = mapEnvToStruct(&ts)
	assert.Equal(t, testStringValue, ts.TestString)
	assert.Equal(t, testUintValue, ts.TestInt)
	assert.True(t, ts.TestBool)
	assert.Equal(t, "", ts.TestSkip)
}

func Test_mapEnvToStruct_wrong_bool(t *testing.T) {
	type testStruct struct {
		TestWrongBool bool `env:"TEST_WRONG_BOOL"`
	}
	_ = os.Setenv("TEST_WRONG_BOOL", "wrong")
	ts := testStruct{}
	err := mapEnvToStruct(&ts)
	assert.Error(t, err)
}

func Test_initStructFromEnv(t *testing.T) {
	testStringValue := "TESTSTRING"
	var testUintValue uint16 = 16555
	type testStruct struct {
		TestString string `env:"TEST_STRING"`
		TestInt    uint16 `env:"TEST_UINT16"`
	}
	_ = os.Setenv("TEST_STRING", testStringValue)
	_ = os.Setenv("TEST_UINT16", strconv.Itoa(int(testUintValue)))
	ts := testStruct{}
	err := InitStructFromEnv(&ts)
	if err != nil {
		t.Fatal(err)
	}
	if ts.TestString != testStringValue {
		t.Errorf("%s  !=  %s \n", ts.TestString, testStringValue)
	}
	if ts.TestInt != testUintValue {
		t.Errorf("%d  !=  %d \n", ts.TestInt, testUintValue)
	}
}
