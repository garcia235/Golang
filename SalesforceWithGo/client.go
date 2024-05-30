package main

import (
	"bufio"
	"fmt"
	"github.com/simpleforce/simpleforce"
	"os"
)

var (
	sfURL      = os.Getenv("SF_URL")
	sfUser     = os.Getenv("SF_USER")
	sfPassword = os.Getenv("SF_PASSWORD")
	sfToken    = os.Getenv("SF_TOKEN")
	reader     = bufio.NewReader(os.Stdin)
)

func newAccount(name, cnpj string) Account {
	return Account{
		Name: name,
		CNPJ: cnpj,
	}
}

func createClient() *simpleforce.Client {
	client := simpleforce.NewClient(sfURL, simpleforce.DefaultClientID, simpleforce.DefaultAPIVersion)
	if client == nil {
		return nil
	}
	err := client.LoginPassword(sfUser, sfPassword, sfToken)
	if err != nil {
		return nil
	}
	return client
}

func Query(client *simpleforce.Client) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter a CNPJ: ")
	cnpj, _ := reader.ReadString('\n')
	cnpj = cnpj[:len(cnpj)-1]

	q := fmt.Sprintf("select Id, Name from Account where CNPJ__c = '%s' limit 10", cnpj)

	result, err := client.Query(q)
	if err != nil {
		return
	}

	for _, record := range result.Records {
		fmt.Println(record)
	}
}

func WorkWithRecords() {
	client := simpleforce.NewClient(sfURL, simpleforce.DefaultClientID, simpleforce.DefaultAPIVersion)
	client.LoginPassword(sfUser, sfPassword, sfToken)

	// Get an SObject with given type and external ID
	obj := client.SObject("Case").Get("__ID__")
	if obj == nil {
		return
	}
	attrs := obj.AttributesField()
	if attrs != nil {
		fmt.Println(attrs.Type)
		fmt.Println(attrs.URL)
	}

	userObj := obj.SObjectField("UserCSV", "CreatedById")
	if userObj == nil {
		return
	}
	fmt.Println(userObj.StringField("Name"))

	userObj.Get()
	fmt.Println(userObj.StringField("Name"))

	updateObj := client.SObject("Contact").
		Set("Id", "001D300000pdybHIAQ").
		Set("FirstName", "New Name").
		Update()
	fmt.Println(updateObj)

	upsertObj := client.SObject("Contact").
		Set("ExternalIDField", "customExtIdField__c").
		Set("customExtIdField__c", "__ExtID__").
		Set("FirstName", "New Name").
		Upsert()
	fmt.Println(upsertObj)

	err := client.SObject("Case").
		Set("Subject", "Case created by simpleforce").
		Set("Comments", "Case commented by simpleforce").
		Create().
		Get().
		Delete()

	fmt.Println(err)
}

func executeAnonymous(client *simpleforce.Client) {
	code := `
		System.debug('test anonymous apex');
	`
	result, err := client.ExecuteAnonymous(code)
	if err != nil {
		// handle the error
		return
	}

	fmt.Println(result)
}

func main() {
	client := createClient()
	if client == nil {
		fmt.Println("Failed to create client")
		return
	}

	// Query(client)
	// WorkWithRecords()
	// executeAnonymous(client)
	// createAccount(client)
	updateAccount(client)
}

func updateAccount(client *simpleforce.Client) {
	updateObj := client.SObject("Account").
		Set("Id", "001D300000pdybHIAQ").
		Set("Name", "Teste Go 2").
		Update()

	fmt.Println(updateObj)
}
func createAccount(client *simpleforce.Client) {
	fmt.Println("Enter a Name of Account: ")
	name, _ := reader.ReadString('\n')
	name = name[:len(name)-1]
	fmt.Println("Enter a CNPJ of Account: ")
	cnpj, _ := reader.ReadString('\n')
	cnpj = cnpj[:len(cnpj)-1]

	account := newAccount(name, cnpj)
	createObj := client.SObject("Account").
		Set("Name", account.Name).
		Set("CNPJ__c", account.CNPJ).
		Create()

	fmt.Println(createObj)
}
