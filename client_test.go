package vectara_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/psyb0t/go-vectara"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testAPIKey = "EXAMPLE-DO-NOT-USE"
	jsonType   = "application/json"
	corpusKey  = "knowledge"
	corpusPath = "/v2/corpora/knowledge"
)

func newTestClient(
	t *testing.T,
	handler http.HandlerFunc,
) *vectara.ClientWithResponses {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	auth := vectara.WithRequestEditorFn(
		func(_ context.Context, req *http.Request) error {
			req.Header.Set("x-api-key", testAPIKey)

			return nil
		},
	)

	client, err := vectara.NewClientWithResponses(server.URL, auth)
	require.NoError(t, err)

	return client
}

func writeJSON(
	t *testing.T,
	w http.ResponseWriter,
	status int,
	body string,
) {
	t.Helper()

	w.Header().Set("Content-Type", jsonType)
	w.WriteHeader(status)

	_, err := w.Write([]byte(body))
	assert.NoError(t, err)
}

func TestListCorporaPagination(t *testing.T) {
	t.Parallel()

	const pageSize int32 = 2

	filter := "guide & caf\u00e9"
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/v2/corpora", r.URL.Path)
		assert.Equal(t, testAPIKey, r.Header.Get("x-api-key"))
		assert.Equal(t, "2", r.URL.Query().Get("limit"))
		assert.Equal(t, filter, r.URL.Query().Get("filter"))

		switch r.URL.Query().Get("page_key") {
		case "":
			writeJSON(t, w, http.StatusOK, `{
				"corpora":[
					{"id":"crp_1","key":"first","name":"First"},
					{"id":"crp_2","key":"second","name":"Second"}
				],
				"metadata":{"page_key":"opaque+/=cursor"}
			}`)
		case "opaque+/=cursor":
			writeJSON(t, w, http.StatusOK, `{
				"corpora":[
					{"id":"crp_3","key":"third","name":"Third"}
				],
				"metadata":{}
			}`)
		default:
			assert.Fail(t, "unexpected pagination cursor")
			writeJSON(t, w, http.StatusBadRequest, "{}")
		}
	})
	limit := pageSize
	params := &vectara.ListCorporaParams{Limit: &limit, Filter: &filter}
	keys := make([]string, 0)

	for range pageSize {
		res, err := client.ListCorporaWithResponse(context.Background(), params)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, res.StatusCode())
		require.NotNil(t, res.JSON200)
		require.NotNil(t, res.JSON200.Corpora)

		for _, corpus := range *res.JSON200.Corpora {
			keys = append(keys, corpus.Key)
		}

		require.NotNil(t, res.JSON200.Metadata)

		if res.JSON200.Metadata.PageKey == nil {
			break
		}

		params.PageKey = res.JSON200.Metadata.PageKey
	}

	assert.Equal(t, []string{"first", "second", "third"}, keys)
}

func TestQueryCorpus(t *testing.T) {
	t.Parallel()

	const interpolation float32 = 0.1234567

	const score float32 = 0.9876543

	const dimensionWeight float64 = 0.123456789012345

	const requestJSON = `{
		"query":"How does retrieval work?",
		"search":{
			"lexical_interpolation":0.1234567,"limit":2,
			"custom_dimensions":{"relevance":0.123456789012345}
		},
		"stream_response":false
	}`

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, corpusPath+"/query", r.URL.Path)
		assert.Equal(t, testAPIKey, r.Header.Get("x-api-key"))
		assert.Equal(t, jsonType, r.Header.Get("Content-Type"))

		var body vectara.QueryCorpusRequest

		if !assert.NoError(t, json.NewDecoder(r.Body).Decode(&body)) {
			return
		}

		assert.Equal(t, "How does retrieval work?", body.Query)

		if assert.NotNil(t, body.Search) &&
			assert.NotNil(t, body.Search.LexicalInterpolation) {
			assert.Equal(t, interpolation, *body.Search.LexicalInterpolation)
		}

		if assert.NotNil(t, body.Search) &&
			assert.NotNil(t, body.Search.CustomDimensions) {
			assert.Equal(t,
				dimensionWeight, (*body.Search.CustomDimensions)["relevance"],
			)
		}

		writeJSON(t, w, http.StatusOK, `{
			"summary":"Retrieval finds matching document parts.",
			"factual_consistency_score":0.9876543,
			"search_results":[]
		}`)
	})

	var body vectara.QueryCorpusRequest

	require.NoError(t, json.Unmarshal([]byte(requestJSON), &body))

	res, err := client.QueryCorpusWithResponse(
		context.Background(), corpusKey, nil, body,
	)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode())
	require.NotNil(t, res.JSON200)
	require.NotNil(t, res.JSON200.FactualConsistencyScore)
	require.NotNil(t, res.JSON200.Summary)
	assert.Equal(t, score, *res.JSON200.FactualConsistencyScore)
	assert.Equal(t,
		"Retrieval finds matching document parts.", *res.JSON200.Summary,
	)
}

func TestGetCorpusBearerAuth(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, corpusPath, r.URL.Path)
		assert.Empty(t, r.Header.Get("x-api-key"))
		assert.Equal(t, "Bearer "+testAPIKey, r.Header.Get("Authorization"))
		writeJSON(t, w, http.StatusOK, `{
			"id":"crp_1","key":"knowledge","name":"Knowledge",
			"enabled":true
		}`)
	})
	res, err := client.GetCorpusWithResponse(
		context.Background(), corpusKey, nil,
		func(_ context.Context, req *http.Request) error {
			req.Header.Del("x-api-key")
			req.Header.Set("Authorization", "Bearer "+testAPIKey)

			return nil
		},
	)
	require.NoError(t, err)
	require.NotNil(t, res.JSON200)
	assert.Equal(t, corpusKey, res.JSON200.Key)
	assert.Equal(t, "Knowledge", res.JSON200.Name)
	require.NotNil(t, res.JSON200.Enabled)
	assert.True(t, *res.JSON200.Enabled)
}

func TestHTTPFailures(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name   string
		status int
	}{
		{name: "unauthorized", status: http.StatusUnauthorized},
		{name: "forbidden", status: http.StatusForbidden},
		{name: "rate_limited", status: http.StatusTooManyRequests},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			client := newTestClient(t,
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, http.MethodGet, r.Method)
					assert.Equal(t, corpusPath, r.URL.Path)
					writeJSON(t, w, tc.status, `{
						"messages":["Request rejected"],
						"request_id":"request-example"
					}`)
				},
			)
			res, err := client.GetCorpusWithResponse(
				context.Background(), corpusKey, nil,
			)
			require.ErrorIs(t, err, nil)
			require.Equal(t, tc.status, res.StatusCode())
			assert.Nil(t, res.JSON200)

			var body vectara.Error

			require.NoError(t, json.Unmarshal(res.Body, &body))
			require.NotNil(t, body.Messages)
			assert.Equal(t, []string{"Request rejected"}, *body.Messages)
		})
	}
}

func TestMalformedJSON(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(t, w, http.StatusOK, `{"corpora":invalid}`)
	})
	res, err := client.ListCorporaWithResponse(context.Background(), nil)

	var syntaxError *json.SyntaxError

	require.ErrorAs(t, err, &syntaxError)
	assert.Nil(t, res)
}

func TestCanceledRequest(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(_ http.ResponseWriter, _ *http.Request) {
		assert.Fail(t, "canceled request reached the server")
	})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	res, err := client.ListCorporaWithResponse(ctx, nil)
	require.ErrorIs(t, err, context.Canceled)
	assert.Nil(t, res)
}

func TestCorrectHallucinationsWithoutQuery(t *testing.T) {
	t.Parallel()

	const original = "The handbook allows pets."

	const corrected = "The handbook excludes pets."

	const requestJSON = `{
		"generated_text":"The handbook allows pets.",
		"documents":[{"text":"The handbook excludes pets."}],
		"model_name":"corrector-example"
	}`

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t,
			"/v2/hallucination_correctors/correct_hallucinations",
			r.URL.Path,
		)
		assert.Equal(t, testAPIKey, r.Header.Get("x-api-key"))
		assert.Equal(t, jsonType, r.Header.Get("Content-Type"))

		raw, err := io.ReadAll(r.Body)
		if !assert.NoError(t, err) {
			return
		}

		var body map[string]json.RawMessage

		if !assert.NoError(t, json.Unmarshal(raw, &body)) {
			return
		}

		assert.NotContains(t, body, "query")
		assert.JSONEq(t, requestJSON, string(raw))
		writeJSON(t, w, http.StatusOK, `{
			"corrected_text":"The handbook excludes pets.",
			"corrections":[{
				"original_text":"The handbook allows pets.",
				"corrected_text":"The handbook excludes pets.",
				"explanation":"The source excludes pets."
			}]
		}`)
	})

	var body vectara.HallucinationCorrectionRequest

	require.NoError(t, json.Unmarshal([]byte(requestJSON), &body))
	require.Nil(t, body.Query)

	res, err := client.CorrectHallucinationsWithResponse(
		context.Background(), nil, body,
	)
	require.NoError(t, err)
	require.NotNil(t, res.JSON200)
	assert.Equal(t, corrected, res.JSON200.CorrectedText)
	require.Len(t, res.JSON200.Corrections, 1)
	assert.Equal(t, original, res.JSON200.Corrections[0].OriginalText)
}

func TestJsonSchemaObjectAdditionalProperties(t *testing.T) {
	t.Parallel()

	const original = `{
		"type":"object",
		"additionalProperties":false,
		"default":null,
		"x-provider-option":"kept"
	}`

	var schema vectara.JsonSchemaObject

	require.NoError(t, json.Unmarshal([]byte(original), &schema))
	require.NotNil(t, schema.AllowsAdditionalProperties)
	assert.False(t, *schema.AllowsAdditionalProperties)
	assert.Equal(t, "kept", schema.AdditionalProperties["x-provider-option"])

	body, err := json.Marshal(schema)
	require.NoError(t, err)
	assert.JSONEq(t, original, string(body))
}
