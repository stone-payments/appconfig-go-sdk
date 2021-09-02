[![Documentation](https://godoc.org/github.com/<username>/<library>?status.svg)](http://godoc.org/github.com/stone-payments/appconfig-go-sdk)


# appconfig-go-sdk
appconfig-go-sdk is a Go client library for accessing the [Azure App Configuration REST API](https://docs.microsoft.com/en-us/azure/azure-app-configuration/rest-api).

## Installation
First, you need to get the library:
```
go get github.com/stone-payments/appconfig-go-sdk
```

Then you can import the library into your code:
```golang
import "github.com/stone-payments/appconfig-go-sdk"
```

## Usage
```golang
import "github.com/stone-payments/appconfig-go-sdk"
```

Construct a new client for the desired resource. In this example, we will use the KeyValue Client.
```golang
args := keyvalues.NewClientAzureADArgs{
		ClientID:         clientID,
		ClientSecret:     clientSecret,
		TenantID:         tenantID,
		ResourceEndpoint: endpoint,
	}
client, err := keyvalues.NewClientAzureAD(args)
```

Then you can use the various methods on the client to access the App Configuration API. For Example:
```golang
list, err := client.ListKeyValues(keyvalues.ListKeyValuesArgs{})
kv, err := client.GetKeyValue("mykey", "mylabel")
```
For more sample code snippets, head over to the [example](example/) directory.
### Testing code that uses appconfig-go-sdk
All clients provide interfaces to SDK calls to improve testability, so you can create a mock struct that implements the methods that you need to test.

Head over to the [KeyValues client](appconfig/keyvalues/client.go) to see an interface example.

## Contributing
Contributions are always welcome. Please, submit pull requests and open new issues!

## License
This library is distributed under the MIT license found in the [LICENSE](LICENSE) file.
