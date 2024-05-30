package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Account struct {
	Id   string `json:"Id"`
	Name string `json:"Name"`
	CNPJ string `json:"CNPJ__c"`
}

var accounts = []Account{}
var files = []string{}

func createAccountApi(c *gin.Context) {
	var newAccount Account
	if err := c.BindJSON(&newAccount); err != nil {
		return
	}
	accounts = append(accounts, newAccount)
	c.IndentedJSON(http.StatusCreated, newAccount)

}
func getAccounts(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, accounts)
}

func getAccountsByID(c *gin.Context) {
	id := c.Param("id")

	for _, a := range accounts {
		if a.Id == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "account not found"})
}

func putAccount(c *gin.Context) {
	id := c.Param("id")
	var newAccount Account
	if err := c.BindJSON(&newAccount); err != nil {
		return
	}
	for i, a := range accounts {
		if a.Id == id {
			accounts[i] = newAccount
			c.IndentedJSON(http.StatusOK, newAccount)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "account not found"})
}

func uploadFiles(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "file not found"})
		return
	}
	files = append(files, file.Filename)
	c.IndentedJSON(http.StatusCreated, gin.H{"message": "file uploaded", "file": file.Filename})
}
func getFiles(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, files)
}

func main() {
	router := gin.Default()
	router.POST("/accounts", createAccountApi)
	router.GET("/accounts", getAccounts)
	router.GET("/accounts/:id", getAccountsByID)
	router.POST("/upload", uploadFiles)
	router.GET("/files", getFiles)
	router.PUT("/accounts/:id", putAccount)
	err := router.Run("localhost:8080")
	if err != nil {
		return
	}
}
