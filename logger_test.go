package log

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	F1 string
	F2 int
}

const (
	testJSON          = `{"a":"apple","b":"banana","c":{"d":"e"},1:3}`
	testInt           = 1
	testFloat         = 1.1
	testBool          = true
	testStr           = "1"
	testJSONResult    = `{"a":"apple","b":"banana"}`
	testStruct0Result = `{"F1":"hello","F2":2}`
)

var (
	testStruct0 = testStruct{
		"hello",
		2,
	}
	testMap = map[string]interface{}{"a": "apple", "b": "banana"}
)

func TestNormalizeArgs(t *testing.T) {
	res := normalizeArgs([]interface{}{testJSON, testStr, testInt, testBool, testFloat, testMap, testStruct0})
	// Validate map is converted to json string testJSONResult
	assert.Equal(t, res, []interface{}{testJSON, testStr, testInt, testBool, testFloat, testJSONResult, testStruct0Result})
}

func TestLogging(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "requestId", "request-id")
	ctx = context.WithValue(ctx, "userId", "user-id")

	Init(SimpleFormatter, logrus.DebugLevel, "requestId", "userId")
	Debug(ctx, "Debug Message 1")
	Debugf(ctx, "Debug Message %d", 2)
	Info(ctx, "Informational Message 1")
	Infof(ctx, "Informational Message %d", 2)
	Warn(ctx, "Warning Message 1")
	Warnf(ctx, "Warning Message %d", 2)
	Error(ctx, "Error Message 1")
	Errorf(ctx, "Error Message %d", 2)

	Init(TextFormatter, logrus.InfoLevel, "requestId", "userId")
	Info(ctx, "Informational Message 1")
	Info(ctx, "Informational Message 2", Field("field1", "value1"))
	Infof(ctx, "Informational Message %d", 3)
	Warn(ctx, "Warning Message 1")
	Warnf(ctx, "Warning Message %d", 2)
	Error(ctx, "Error Message 1")
	Errorf(ctx, "Error Message %d", 2)

	Init(JSONFormatter, logrus.DebugLevel, "requestId", "userId")
	Debug(ctx, "Debug Message 1")
	Debugf(ctx, "Debug Message %d", 2)
	Info(ctx, "Informational Message 1")
	Info(ctx, "Informational Message 2", Field("field1", "value1"))
	Infof(ctx, "Informational Message %d", 3)
	Warn(ctx, "Warning Message 1")
	Warnf(ctx, "Warning Message %d", 2)
	Error(ctx, "Error Message 1")
	Errorf(ctx, "Error Message %d", 2)

	ctx = context.TODO()
	Init(JSONFormatter, logrus.InfoLevel, "requestId", "userId")
	Info(ctx, "Informational Message 1")
	Infof(ctx, "Informational Message %d", 2)
	Warn(ctx, "Warning Message 1")
	Warnf(ctx, "Warning Message %d", 2)
	Error(ctx, "Error Message 1")
	Errorf(ctx, "Error Message %d", 2)
}
