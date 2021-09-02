package main

import (
	"fmt"

	"github.com/stone-payments/appconfig-go-sdk/appconfig/keyvalues"
)

func main() {
	endpoint := "https://my-config.azconfig.io"
	client, err := keyvalues.NewClientCli(endpoint)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	kv, err := createKeyValues(client)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("KeyValue created. Key: %v\n", *kv.Key)

	list, err := listKeyValues(client)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	for i, item := range list.Items {
		fmt.Printf("KeyValue %v. Key: %v\n", i+1, *item.Key)
	}
}

func createKeyValues(client keyvalues.Client) (keyvalues.KeyValue, error) {
	args := keyvalues.CreateOrUpdateKeyValueArgs{
		Key:   "mykey",
		Label: "mylabel",
		Value: "myvalue",
	}
	kv, err := client.CreateOrUpdateKeyValue(args)
	if err != nil {
		return keyvalues.KeyValue{}, err
	}
	return kv, nil
}

func listKeyValues(client keyvalues.Client) (keyvalues.KeyValues, error) {
	args := keyvalues.ListKeyValuesArgs{}
	kvs, err := client.ListKeyValues(args)
	if err != nil {
		return keyvalues.KeyValues{}, err
	}
	return kvs, nil
}
