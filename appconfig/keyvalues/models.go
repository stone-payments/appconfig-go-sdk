package keyvalues

// KeyValue represents a Key Value response
type KeyValue struct {
	Etag         *string            `json:"etag,omitempty"`
	Key          *string            `json:"key,omitempty"`
	Label        *string            `json:"label,omitempty"`
	ContentType  *string            `json:"content_type,omitempty"`
	Value        *string            `json:"value,omitempty"`
	LastModified *string            `json:"last_modified,omitempty"`
	Locked       *bool              `json:"locked,omitempty"`
	Tags         *map[string]string `json:"tags,omitempty"`
}

// ListKeyValuesArgs represents the argument for the
// ListKeyValues SDK method.
type ListKeyValuesArgs struct {
	Key   string
	Label string
}

// CreateOrUpdateKeyValueArgs represents the argument for the
// CreateOrUpdateKeyValue SDK method.
//
// IsSecret by default is false and must be true if the KeyValue is a
// reference to Azure Key Vault.
//
// The Value field must be a Secret Identifier when the KeyValue
// is an Azure Key Vault reference.
// Example:
// https://my-vault.vault.azure.net/secrets/mysecret
type CreateOrUpdateKeyValueArgs struct {
	Key         string            `json:"key,omitempty"`
	Label       string            `json:"label,omitempty"`
	ContentType string            `json:"content_type,omitempty"`
	Value       string            `json:"value,omitempty"`
	Tags        map[string]string `json:"tags,omitempty"`
	IsSecret    bool              `json:"isSecret,omitempty"`
}

// NewClientAzureADArgs represents the argument for the
// NewClientAzureAD SDK method.
//
// Required: ClientID, ClientSecret, TenantID, ResourceEndpoint
// Optional: AADEndpoint
type NewClientAzureADArgs struct {
	ClientID         string
	ClientSecret     string
	TenantID         string
	AADEndpoint      string
	ResourceEndpoint string
}

// KeyValues represents the response of the
// ListKeyValues SDK method.
type KeyValues struct {
	Items []KeyValue `json:"items"`
}
