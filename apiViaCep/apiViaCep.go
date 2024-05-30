package main

import (
	"encoding/json"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"time"
)

type Endereco struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	DDD         string `json:"ddd"`
	Siafi       string `json:"siafi"`

	Erro     bool   `json:"erro"`
	Mensagem string `json:"mensagem"`
}

var enderecoList []Endereco

func viaCep(c *gin.Context) {
	cep := c.Param("cep")
	url := "https://viacep.com.br/ws/" + cep + "/json/"
	resp, err := http.Get(url)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "error"})
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "error"})
		return
	}

	var endereco Endereco
	err = json.Unmarshal(body, &endereco)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "error"})
		return
	}
	enderecoList = append(enderecoList, endereco)
	c.IndentedJSON(http.StatusOK, endereco)
}

func getEnderecos(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, enderecoList)
}

func main() {
	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"POST", "GET", "PUT", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Accept", "User-Agent", "Cache-Control", "Pragma"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true
	config.MaxAge = 12 * time.Hour

	router.Use(cors.New(config))
	router.GET("/viacep/:cep", viaCep)
	router.GET("/enderecos", getEnderecos)
	err := router.Run(":8080")
	if err != nil {
		return
	}
}
