package keyvalues

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

const (
	defaultContentType     = "application/vnd.microsoft.appconfig.kv+json"
	keyVaultRefContentType = "application/vnd.microsoft.appconfig.keyvaultref+json;charset=utf-8"
	apiVersion             = "1.0"
)

// Client is an interface with all methods to
// manage App Configuration Key Values
type Client interface {
	// ListKeyValues returns an array of App Configuration KeyValues. The list
	// of KeyValues are filtered by the provided Key and/or Label.
	//
	// Optional: Key; Label (if not specified, it implies any Key/Label).
	ListKeyValues(ListKeyValuesArgs) (KeyValues, error)

	// GetKeyValue gets an App Configuration Key-Value.
	GetKeyValue(key, label string) (KeyValue, error)

	// CreateOrUpdateKeyValue create/update an App Configuration Key-Value.
	//
	// Required parameters: Key; Value
	//
	// Optional parameters: Label; ContentType; Tags; IsSecret
	CreateOrUpdateKeyValue(CreateOrUpdateKeyValueArgs) (KeyValue, error)

	// DeleteKeyValue deletes an App Configuration Key-Value.
	DeleteKeyValue(key, label string) error
}

// ClientImpl implements the Client interface
type ClientImpl struct {
	autorest.Client
	Endpoint string
}

// NewClientAzureAD creates a Client configured from Azure AD credentials.
func NewClientAzureAD(args NewClientAzureADArgs) (Client, error) {
	creds := auth.NewClientCredentialsConfig(args.ClientID, args.ClientSecret, args.TenantID)
	creds.Resource = args.ResourceEndpoint
	if args.AADEndpoint != "" {
		creds.AADEndpoint = args.AADEndpoint
	}

	auth, err := creds.Authorizer()
	if err != nil {
		return nil, err
	}

	return newClient(args.ResourceEndpoint, auth), nil
}

// NewClientCli creates a Client configured from Azure CLI 2.0.
func NewClientCli(endpoint string) (Client, error) {
	auth, err := auth.NewAuthorizerFromCLIWithResource(endpoint)
	if err != nil {
		return nil, err
	}
	return newClient(endpoint, auth), nil
}

func newClient(endpoint string, authorizer autorest.Authorizer) Client {
	client := autorest.NewClientWithUserAgent(autorest.UserAgent())
	client.Authorizer = authorizer

	return &ClientImpl{
		Client:   client,
		Endpoint: endpoint,
	}
}

// ListKeyValues returns an array of App Configuration KeyValues. The list
// of KeyValues are filtered by the provided Key and/or Label.
//
// Optional: Key; Label (if not specified, it implies any Key/Label).
func (client *ClientImpl) ListKeyValues(args ListKeyValuesArgs) (KeyValues, error) {
	if args.Key == "" {
		args.Key = "*"
	}
	if args.Label == "" {
		args.Label = "*"
	}

	response, err := client.sendRequest(args.Label, args.Key, true, autorest.AsGet())
	if err != nil {
		return KeyValues{}, err
	}

	var result KeyValues
	if err = getJSON(response, &result); err != nil {
		return KeyValues{}, err
	}
	return result, nil
}

// GetKeyValue gets an App Configuration Key-Value.
func (client *ClientImpl) GetKeyValue(key, label string) (KeyValue, error) {
	result := KeyValue{}
	response, err := client.sendRequest(label, url.QueryEscape(key), false, autorest.AsGet())
	if err != nil {
		return result, err
	}

	if err = getJSON(response, &result); err != nil {
		return result, err
	}

	return result, nil
}

// CreateOrUpdateKeyValue create/update an App Configuration Key-Value.
//
// Required parameters: Key; Value
// Optional parameters: Label; ContentType; Tags; IsSecret
func (client *ClientImpl) CreateOrUpdateKeyValue(args CreateOrUpdateKeyValueArgs) (KeyValue, error) {
	result := KeyValue{}

	if args.IsSecret {
		args.Value = fmt.Sprintf("{\"uri\":\"%s\"}", args.Value)
		if args.ContentType == "" {
			args.ContentType = keyVaultRefContentType
		}
	}

	response, err := client.sendRequest(
		args.Label,
		args.Key,
		false,
		autorest.AsContentType(defaultContentType),
		autorest.AsPut(),
		autorest.WithJSON(args),
	)
	if err != nil {
		return result, err
	}

	if err = getJSON(response, &result); err != nil {
		return result, err
	}

	return result, nil
}

// DeleteKeyValue deletes an App Configuration Key-Value.
func (client *ClientImpl) DeleteKeyValue(key, label string) error {
	_, err := client.sendRequest(
		label,
		url.QueryEscape(key),
		false,
		autorest.AsDelete(),
	)
	if err != nil {
		return err
	}
	return nil
}

func (client *ClientImpl) sendRequest(label string, key string, isList bool, additionalDecorator ...autorest.PrepareDecorator) (*http.Response, error) {
	var req *http.Request
	var err error
	queryParameters := map[string]interface{}{
		"label":       label,
		"api-version": apiVersion,
	}

	if !isList {
		req, err = client.preparer(
			label,
			key,
			queryParameters,
			additionalDecorator...,
		).Prepare(&http.Request{})
	} else {
		req, err = client.listPreparer(
			label,
			key,
			queryParameters,
			additionalDecorator...,
		).Prepare(&http.Request{})
	}
	if err != nil {
		return nil, err
	}

	var resp *http.Response
	resp, err = client.Send(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		s, _ := ioutil.ReadAll(resp.Body)
		body := string(s)
		return resp, fmt.Errorf("ERROR: %s - Response Body: %s", resp.Status, body)
	}

	return resp, err
}

func (client *ClientImpl) preparer(label, key string, query map[string]interface{}, additionalDecorators ...autorest.PrepareDecorator) autorest.Preparer {
	pathParameters := map[string]interface{}{
		"key": key,
	}
	decorators := []autorest.PrepareDecorator{
		autorest.WithBaseURL(client.Endpoint),
		autorest.WithPathParameters("/kv/{key}", pathParameters),
		autorest.WithQueryParameters(query),
		client.Client.WithAuthorization(),
	}
	decorators = append(decorators, additionalDecorators...)

	return autorest.CreatePreparer(decorators...)
}

func (client *ClientImpl) listPreparer(label string, key string, query map[string]interface{}, additionalDecorators ...autorest.PrepareDecorator) autorest.Preparer {
	query["key"] = key
	decorators := []autorest.PrepareDecorator{
		autorest.WithBaseURL(fmt.Sprintf("%s/kv", client.Endpoint)),
		autorest.WithQueryParameters(query),
		client.Client.WithAuthorization(),
	}
	decorators = append(decorators, additionalDecorators...)

	return autorest.CreatePreparer(decorators...)
}

func getJSON(response *http.Response, target interface{}) error {
	defer response.Body.Close()
	return json.NewDecoder(response.Body).Decode(target)
}
