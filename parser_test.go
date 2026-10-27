package httpx

import (
	"net/http/httptest"
	"testing"
	"uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseValueAsString(t *testing.T) {
	got, gotErr := ParseValueAs[string]("e7990bab-0957-4ee0-9d03-dfd5fd4cd529")

	require.NoError(t, gotErr)
	assert.Equal(t, "e7990bab-0957-4ee0-9d03-dfd5fd4cd529", got)
}

func TestParseValueAsString_Empty(t *testing.T) {
	got, gotErr := ParseValueAs[string]("")

	require.NoError(t, gotErr)
	assert.Empty(t, got)
}

func TestParseValueAsUUID(t *testing.T) {
	got, gotErr := ParseValueAs[uuid.UUID]("e7990bab-0957-4ee0-9d03-dfd5fd4cd529")

	require.NoError(t, gotErr)
	assert.Equal(t, uuid.MustParse("e7990bab-0957-4ee0-9d03-dfd5fd4cd529"), got)
}

func TestParseValueAsInt(t *testing.T) {
	got, gotErr := ParseValueAs[int]("123")

	require.NoError(t, gotErr)
	assert.Equal(t, int(123), got)
}

func TestParseValueAsInt64(t *testing.T) {
	got, gotErr := ParseValueAs[int64]("123")

	require.NoError(t, gotErr)
	assert.Equal(t, int64(123), got)
}

func TestParseValueAsUUID_Error(t *testing.T) {
	got, gotErr := ParseValueAs[uuid.UUID]("e7990bab-0957-4ee0-9d03-dfd5fd4cd")

	require.Error(t, gotErr)
	assert.Zero(t, got)
}

func TestParseValueAsInt_Error(t *testing.T) {
	got, gotErr := ParseValueAs[int]("a")

	require.Error(t, gotErr)
	assert.Zero(t, got)
}

func TestParseValueAsInt64_Error(t *testing.T) {
	got, gotErr := ParseValueAs[int64]("a")

	require.Error(t, gotErr)
	assert.Zero(t, got)
}

func TestParseValueAs_ErrorUnsupportedType(t *testing.T) {
	type mytype string

	got, gotErr := ParseValueAs[mytype]("abc")

	require.Error(t, gotErr)
	assert.Zero(t, got)
}

func TestParsePathParamAs(t *testing.T) {
	paramName := "param_1"
	paramVal := "param_value"

	req := httptest.NewRequest("GET", "/api/v1/smth/"+paramVal, nil)
	req.SetPathValue(paramName, paramVal)

	gotVal, gotErr := ParsePathParamAs[string](req, paramName)

	require.NoError(t, gotErr)
	assert.Equal(t, paramVal, gotVal)
}

func TestParsePathParamAs_Error(t *testing.T) {
	paramName := "param_1"
	paramVal := "param_value"

	req := httptest.NewRequest("GET", "/api/v1/smth/"+paramVal, nil)
	req.SetPathValue(paramName, paramVal)

	gotVal, gotErr := ParsePathParamAs[uuid.UUID](req, paramName)

	var err *ErrInvalidParam
	require.ErrorAs(t, gotErr, &err)
	assert.Equal(t, "path", err.ParamType)
	assert.Equal(t, paramName, err.ParamName)
	assert.Zero(t, gotVal)
}

func TestParsePathParamOrAs(t *testing.T) {
	paramName := "param_1"
	paramVal := "param_value"
	defaultVal := "default_value"

	req := httptest.NewRequest("GET", "/api/v1/smth/"+paramVal, nil)
	req.SetPathValue(paramName, paramVal)

	gotVal, gotErr := ParsePathParamOrAs(req, paramName, defaultVal)

	require.NoError(t, gotErr)
	assert.Equal(t, paramVal, gotVal)
}

func TestParsePathParamOrAs_Default(t *testing.T) {
	paramName := "param_1"
	paramVal := "param_value"
	defaultVal := "default_value"

	req := httptest.NewRequest("GET", "/api/v1/smth/"+paramVal, nil)
	req.SetPathValue(paramName, paramVal)

	gotVal, gotErr := ParsePathParamOrAs(req, "not_found", defaultVal)

	require.NoError(t, gotErr)
	assert.Equal(t, defaultVal, gotVal)
}

func TestParsePathParamOrAs_Error(t *testing.T) {
	paramName := "param_1"
	paramVal := "param_value"
	defaultVal := uuid.New()

	req := httptest.NewRequest("GET", "/api/v1/smth/"+paramVal, nil)
	req.SetPathValue(paramName, paramVal)

	gotVal, gotErr := ParsePathParamOrAs(req, paramName, defaultVal)

	var err *ErrInvalidParam
	require.ErrorAs(t, gotErr, &err)
	assert.Equal(t, paramName, err.ParamName)
	assert.Equal(t, "path", err.ParamType)
	assert.Zero(t, gotVal)
}

func TestParseQueryParamAs(t *testing.T) {
	paramName := "param_1"
	paramVal := "param_value"

	req := httptest.NewRequest("GET", "/api/v1/links?"+paramName+"="+paramVal+"", nil)

	gotVal, gotErr := ParseQueryParamAs[string](req, paramName)

	require.NoError(t, gotErr)
	assert.Equal(t, paramVal, gotVal)
}

func TestParseQueryParamAs_Error(t *testing.T) {
	paramName := "param_1"
	paramVal := "param_value"

	req := httptest.NewRequest("GET", "/api/v1/links?"+paramName+"="+paramVal+"", nil)

	gotVal, gotErr := ParseQueryParamAs[uuid.UUID](req, paramName)

	var err *ErrInvalidParam
	require.ErrorAs(t, gotErr, &err)
	assert.Equal(t, paramName, err.ParamName)
	assert.Equal(t, "query", err.ParamType)
	assert.Zero(t, gotVal)
}

func TestParseQueryParamOrAs(t *testing.T) {
	paramName := "param_1"
	paramVal := "param_value"
	defaultVal := "default_value"

	req := httptest.NewRequest("GET", "/api/v1/links?"+paramName+"="+paramVal+"", nil)

	gotVal, gotErr := ParseQueryParamOrAs(req, paramName, defaultVal)

	require.NoError(t, gotErr)
	assert.Equal(t, paramVal, gotVal)
}

func TestParseQueryParamOrAs_Default(t *testing.T) {
	paramName := "param_1"
	paramVal := "param_value"
	defaultVal := "default_value"

	req := httptest.NewRequest("GET", "/api/v1/links?"+paramName+"="+paramVal+"", nil)

	gotVal, gotErr := ParseQueryParamOrAs(req, "not_found", defaultVal)

	require.NoError(t, gotErr)
	assert.Equal(t, defaultVal, gotVal)
}

func TestParseQueryParamOrAs_Error(t *testing.T) {
	paramName := "param_1"
	paramVal := "param_value"
	defaultVal := uuid.New()

	req := httptest.NewRequest("GET", "/api/v1/links?"+paramName+"="+paramVal+"", nil)

	gotVal, gotErr := ParseQueryParamOrAs(req, paramName, defaultVal)

	var err *ErrInvalidParam
	require.ErrorAs(t, gotErr, &err)
	assert.Equal(t, paramName, err.ParamName)
	assert.Equal(t, "query", err.ParamType)
	assert.Zero(t, gotVal)
}
