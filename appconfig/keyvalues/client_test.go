package keyvalues

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/stone-payments/appconfig-go-sdk/appconfig/test"
)

var (
	fakeIDs            = "fake"
	fakeKey            = "fakeKey"
	fakeLabel          = "fakeLabel"
	fakeValue          = "fakeValue"
	fakeSecretResponse = "{\"uri\":\"fakeValue\"}"
	errString          = "ERROR: 500 Internal Server Error - Response Body: "
)

type token struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	ExpiresIn    json.Number `json:"expires_in"`
	ExpiresOn    string      `json:"expires_on"`
	NotBefore    json.Number `json:"not_before"`
	Resource     string      `json:"resource"`
	Type         string      `json:"token_type"`
}

func TestNewClientAzureAD(t *testing.T) {
	type args struct {
		NewClientAzureADArgs
	}
	type want struct {
		err error
	}

	cases := map[string]struct {
		reason  string
		handler http.Handler
		args    args
		want    want
	}{
		"CreateNewClientSuccess": {
			reason: "Should create a new Client sucessfully",
			args: args{
				NewClientAzureADArgs: NewClientAzureADArgs{
					ClientID:         "fakeID",
					ClientSecret:     "fakeSecret",
					TenantID:         "fakeID",
					ResourceEndpoint: "fake.com",
				},
			},
			want: want{
				err: nil,
			},
		},
		"CreateNewClientError": {
			reason: "Should return an error if the Client creation fails",
			args: args{
				NewClientAzureADArgs: NewClientAzureADArgs{
					ClientID:         "fakeID",
					ClientSecret:     "fakeSecret",
					TenantID:         "fakeID",
					ResourceEndpoint: "",
				},
			},
			want: want{
				err: errors.New("failed to get SPT from client credentials: parameter 'resource' cannot be empty"),
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := NewClientAzureAD(tc.args.NewClientAzureADArgs)
			if diff := cmp.Diff(tc.want.err, err, test.EquateErrors()); diff != "" {
				t.Errorf("NewClientAzureAD(...): -want error, +got error:\n%s", diff)
			}
		})
	}
}

func TestNewClientCli(t *testing.T) {
	type args struct {
		endpoint string
	}
	type want struct {
		err error
	}

	cases := map[string]struct {
		reason  string
		handler http.Handler
		args    args
		want    want
	}{
		"CreateNewClientError": {
			reason: "Should return an error if the Client creation fails",
			args: args{
				endpoint: "",
			},
			want: want{
				err: errors.New("Resource  is not in expected format. Only alphanumeric characters, [dot], [colon], [hyphen], and [forward slash] are allowed."),
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := NewClientCli(tc.args.endpoint)
			if diff := cmp.Diff(tc.want.err, err, test.EquateErrors()); diff != "" {
				t.Errorf("NewClientCli(...): -want error, +got error:\n%s", diff)
			}
		})
	}
}

func TestListKeyValues(t *testing.T) {
	type args struct {
		ListKeyValuesArgs
	}
	type want struct {
		kvs KeyValues
		err error
	}

	cases := map[string]struct {
		reason  string
		handler http.Handler
		args    args
		want    want
	}{
		"ListKeyValuesSucessfully": {
			reason: "Should return a list of KeyValues",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				r.Body.Close()
				w.WriteHeader(http.StatusOK)
				if strings.Contains(r.URL.String(), "oauth") {
					_ = json.NewEncoder(w).Encode(&token{})
				}
				if strings.Contains(r.URL.String(), "kv") {
					response := KeyValues{Items: []KeyValue{{}, {}}}
					_ = json.NewEncoder(w).Encode(&response)
				}
			}),
			args: args{
				ListKeyValuesArgs: ListKeyValuesArgs{},
			},
			want: want{
				kvs: KeyValues{Items: []KeyValue{{}, {}}},
				err: nil,
			},
		},
		"ListKeyValuesInternalError": {
			reason: "Should return an error if the list returns Status Code greater than 399",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				r.Body.Close()
				if strings.Contains(r.URL.String(), "oauth") {
					w.WriteHeader(http.StatusOK)
					_ = json.NewEncoder(w).Encode(&token{})
				}
				if strings.Contains(r.URL.String(), "kv") {
					w.WriteHeader(http.StatusInternalServerError)
				}
			}),
			args: args{
				ListKeyValuesArgs: ListKeyValuesArgs{},
			},
			want: want{
				kvs: KeyValues{},
				err: errors.New(errString),
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			server := httptest.NewServer(tc.handler)
			defer server.Close()
			args := NewClientAzureADArgs{ClientID: fakeIDs, TenantID: fakeIDs, ClientSecret: fakeIDs, ResourceEndpoint: server.URL, AADEndpoint: server.URL}
			c, _ := NewClientAzureAD(args)

			got, err := c.ListKeyValues(tc.args.ListKeyValuesArgs)

			if diff := cmp.Diff(tc.want.kvs, got); diff != "" {
				t.Errorf("ListKeyValues(...): -want, +got:\n%s", diff)
			}
			if diff := cmp.Diff(tc.want.err, err, test.EquateErrors()); diff != "" {
				t.Errorf("ListKeyValues(...): -want error, +got error:\n%s", diff)
			}
		})
	}
}

func TestGetKeyValue(t *testing.T) {
	type args struct {
		key   string
		label string
	}
	type want struct {
		kv  KeyValue
		err error
	}

	cases := map[string]struct {
		reason  string
		handler http.Handler
		args    args
		want    want
	}{
		"GetKeyValueSucessfully": {
			reason: "Should return a KeyValue",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				r.Body.Close()
				w.WriteHeader(http.StatusOK)
				if strings.Contains(r.URL.String(), "oauth") {
					_ = json.NewEncoder(w).Encode(&token{})
				}
				if strings.Contains(r.URL.String(), "kv") {
					_ = json.NewEncoder(w).Encode(KeyValue{Key: &fakeKey, Label: &fakeLabel})
				}
			}),
			args: args{
				key:   fakeKey,
				label: fakeLabel,
			},
			want: want{
				kv: KeyValue{
					Key:   &fakeKey,
					Label: &fakeLabel,
				},
				err: nil,
			},
		},
		"GetKeyValueInternalError": {
			reason: "Should return an error if the GET returns Status Code greater than 399",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				r.Body.Close()
				if strings.Contains(r.URL.String(), "oauth") {
					w.WriteHeader(http.StatusOK)
					_ = json.NewEncoder(w).Encode(&token{})
				}
				if strings.Contains(r.URL.String(), "kv") {
					w.WriteHeader(http.StatusInternalServerError)
				}
			}),
			args: args{
				key:   fakeKey,
				label: fakeLabel,
			},
			want: want{
				kv:  KeyValue{},
				err: errors.New(errString),
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			server := httptest.NewServer(tc.handler)
			defer server.Close()
			args := NewClientAzureADArgs{ClientID: fakeIDs, TenantID: fakeIDs, ClientSecret: fakeIDs, ResourceEndpoint: server.URL, AADEndpoint: server.URL}
			c, _ := NewClientAzureAD(args)

			got, err := c.GetKeyValue(tc.args.key, tc.args.label)

			if diff := cmp.Diff(tc.want.kv, got); diff != "" {
				t.Errorf("GetKeyValue(...): -want, +got:\n%s", diff)
			}
			if diff := cmp.Diff(tc.want.err, err, test.EquateErrors()); diff != "" {
				t.Errorf("GetKeyValue(...): -want error, +got error:\n%s", diff)
			}
		})
	}
}

func TestCreateOrUpdateKeyValue(t *testing.T) {
	type args struct {
		CreateOrUpdateKeyValueArgs
	}
	type want struct {
		kv  KeyValue
		err error
	}

	cases := map[string]struct {
		reason  string
		handler http.Handler
		args    args
		want    want
	}{
		"CreateOrUpdateKeyValueSucessfully": {
			reason: "Should return the created KeyValue",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				r.Body.Close()
				if strings.Contains(r.URL.String(), "oauth") {
					w.WriteHeader(http.StatusOK)
					_ = json.NewEncoder(w).Encode(&token{})
				}
				if strings.Contains(r.URL.String(), "kv") {
					w.WriteHeader(http.StatusCreated)
					_ = json.NewEncoder(w).Encode(KeyValue{Key: &fakeKey, Label: &fakeLabel})
				}
			}),
			args: args{
				CreateOrUpdateKeyValueArgs: CreateOrUpdateKeyValueArgs{
					Key:   fakeKey,
					Label: fakeLabel,
				},
			},
			want: want{
				kv: KeyValue{
					Key:   &fakeKey,
					Label: &fakeLabel,
				},
				err: nil,
			},
		},
		"CreateOrUpdateSecretKeyValue": {
			reason: "Should create a secret key value",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				r.Body.Close()
				if strings.Contains(r.URL.String(), "oauth") {
					w.WriteHeader(http.StatusOK)
					_ = json.NewEncoder(w).Encode(&token{})
				}
				if strings.Contains(r.URL.String(), "kv") {
					w.WriteHeader(http.StatusCreated)
					_ = json.NewEncoder(w).Encode(KeyValue{Key: &fakeKey, Label: &fakeLabel, Value: &fakeSecretResponse})
				}
			}),
			args: args{
				CreateOrUpdateKeyValueArgs: CreateOrUpdateKeyValueArgs{
					Key:      fakeKey,
					Label:    fakeLabel,
					Value:    fakeValue,
					IsSecret: true,
				},
			},
			want: want{
				kv: KeyValue{
					Key:   &fakeKey,
					Label: &fakeLabel,
					Value: &fakeSecretResponse,
				},
				err: nil,
			},
		},
		"CreateOrUpdateKeyValueInternalError": {
			reason: "Should return an error if the request returns Status Code greater than 399",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				r.Body.Close()
				if strings.Contains(r.URL.String(), "oauth") {
					w.WriteHeader(http.StatusOK)
					_ = json.NewEncoder(w).Encode(&token{})
				}
				if strings.Contains(r.URL.String(), "kv") {
					w.WriteHeader(http.StatusInternalServerError)
				}
			}),
			args: args{
				CreateOrUpdateKeyValueArgs: CreateOrUpdateKeyValueArgs{
					Key:   fakeKey,
					Label: fakeLabel,
				},
			},
			want: want{
				kv:  KeyValue{},
				err: errors.New(errString),
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			server := httptest.NewServer(tc.handler)
			defer server.Close()
			args := NewClientAzureADArgs{ClientID: fakeIDs, TenantID: fakeIDs, ClientSecret: fakeIDs, ResourceEndpoint: server.URL, AADEndpoint: server.URL}
			c, _ := NewClientAzureAD(args)

			got, err := c.CreateOrUpdateKeyValue(tc.args.CreateOrUpdateKeyValueArgs)

			if diff := cmp.Diff(tc.want.kv, got); diff != "" {
				t.Errorf("CreateOrUpdateKeyValue(...): -want, +got:\n%s", diff)
			}
			if diff := cmp.Diff(tc.want.err, err, test.EquateErrors()); diff != "" {
				t.Errorf("CreateOrUpdateKeyValue(...): -want error, +got error:\n%s", diff)
			}
		})
	}
}

func TestDeleteKeyValue(t *testing.T) {
	type args struct {
		key   string
		label string
	}
	type want struct {
		err error
	}

	cases := map[string]struct {
		reason  string
		handler http.Handler
		args    args
		want    want
	}{
		"DeleteKeyValueSucessfully": {
			reason: "Should delete the KeyValue successfully",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				r.Body.Close()
				if strings.Contains(r.URL.String(), "oauth") {
					w.WriteHeader(http.StatusOK)
					_ = json.NewEncoder(w).Encode(&token{})
				}
				if strings.Contains(r.URL.String(), "kv") {
					w.WriteHeader(http.StatusNoContent)
				}
			}),
			args: args{
				key:   fakeKey,
				label: fakeLabel,
			},
			want: want{
				err: nil,
			},
		},
		"DeleteKeyValueInternalError": {
			reason: "Should return an error if the request returns Status Code greater than 399",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				r.Body.Close()
				if strings.Contains(r.URL.String(), "oauth") {
					w.WriteHeader(http.StatusOK)
					_ = json.NewEncoder(w).Encode(&token{})
				}
				if strings.Contains(r.URL.String(), "kv") {
					w.WriteHeader(http.StatusInternalServerError)
				}
			}),
			args: args{
				key:   fakeKey,
				label: fakeLabel,
			},
			want: want{
				err: errors.New(errString),
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			server := httptest.NewServer(tc.handler)
			defer server.Close()
			args := NewClientAzureADArgs{ClientID: fakeIDs, TenantID: fakeIDs, ClientSecret: fakeIDs, ResourceEndpoint: server.URL, AADEndpoint: server.URL}
			c, _ := NewClientAzureAD(args)

			err := c.DeleteKeyValue(tc.args.key, tc.args.key)

			if diff := cmp.Diff(tc.want.err, err, test.EquateErrors()); diff != "" {
				t.Errorf("DeleteKeyValue(...): -want error, +got error:\n%s", diff)
			}
		})
	}
}
