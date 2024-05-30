package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
}

func getToken() (string, error) {
	url := os.Getenv("SF_URL_TOKEN")
	method := "POST"

	payload := strings.NewReader(os.Getenv("CLIENT"))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var tokenResponse TokenResponse
	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		return "", err
	}

	return tokenResponse.AccessToken, nil
}

type FileUpload struct {
	CSVFile *multipart.FileHeader `form:"csvfile" binding:"required"`
}

type UserCSV struct {
	ID                     string `json:"id"`
	Cargo                  string `json:"cargo"`
	FlagComprador          string `json:"flag_comprador"`
	CPF                    string `json:"cpf"`
	PermissaoDemonstrativo string `json:"permissao_demonstrativo"`
	Email                  string `json:"email"`
	Nome                   string `json:"nome"`
	DocumentoNumero        string `json:"documento_numero"`
	FlagRecebeNotaFiscal   string `json:"flag_recebe_nota_fiscal"`
	FlagRepresentanteLegal string `json:"flag_representante_legal"`
	TelefoneCelular        string `json:"telefone_celular"`
	TelefoneFixo           string `json:"telefone_fixo"`
	DocumentoID            string `json:"documento_id"`
	EmpresaID              string `json:"empresa_id"`
	EstadoCivilID          string `json:"estado_civil_id"`
	WBCContato             string `json:"wbc_contato"`
	DepartamentoID         string `json:"departamento_id"`
	CNPJ                   string `json:"cnpj"`
}

func GetCSV(c *gin.Context) {
	var csvfile FileUpload
	if err := c.ShouldBind(&csvfile); err != nil {
		log.Fatalln("Error in JSON binding: ", err)
	}

	if csvfile.CSVFile == nil {
		log.Fatalln("File is missing")
	}

	token, err := getToken()
	if err != nil {
		log.Fatalln("Error getting token:", err)
	}

	file, err := csvfile.CSVFile.Open()
	if err != nil {
		log.Fatalln("Error in opening file")
	}
	defer file.Close()

	reader := csv.NewReader(file)

	if _, err := reader.Read(); err != nil {
		log.Fatalln("Error in reading header line")
	}

	var users []UserCSV
	for {
		record, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			log.Fatalln("Error in reading file", err)
		}

		user := UserCSV{
			ID:                     record[0],
			Cargo:                  record[1],
			FlagComprador:          record[2],
			CPF:                    record[3],
			PermissaoDemonstrativo: record[4],
			Email:                  record[5],
			Nome:                   record[6],
			DocumentoNumero:        record[7],
			FlagRecebeNotaFiscal:   record[8],
			FlagRepresentanteLegal: record[9],
			TelefoneCelular:        record[10],
			TelefoneFixo:           record[11],
			DocumentoID:            record[12],
			EmpresaID:              record[13],
			EstadoCivilID:          record[14],
			WBCContato:             record[15],
			DepartamentoID:         record[16],
			CNPJ:                   record[17],
		}
		users = append(users, user)
	}

	// Processamento em batches
	batchSize := 100 // Define o tamanho do batch
	var wg sync.WaitGroup
	userCh := make(chan []UserCSV, len(users)/batchSize+1)

	go func() {
		defer close(userCh)
		for i := 0; i < len(users); i += batchSize {
			end := i + batchSize
			if end > len(users) {
				end = len(users)
			}
			userCh <- users[i:end]
		}
	}()

	for batch := range userCh {
		wg.Add(1)
		go func(batch []UserCSV) {
			defer wg.Done()
			createContactsInSalesforce(token, batch)
		}(batch)
	}

	wg.Wait()
	c.IndentedJSON(http.StatusOK, users)
}

func createContactsInSalesforce(token string, users []UserCSV) {
	url := os.Getenv("SF_URL_CONTACT")
	method := "POST"

	records := make([]map[string]interface{}, len(users))
	for i, user := range users {
		records[i] = map[string]interface{}{
			"attributes": map[string]string{"type": "Contact"},
			"LastName":   user.Nome,
			"Email":      user.Email,
			"Phone":      user.TelefoneCelular,
		}
	}

	payload := map[string]interface{}{
		"allOrNone": true,
		"records":   records,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, strings.NewReader(string(jsonData)))
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(body))
}

func main() {
	router := gin.Default()
	router.POST("/upload", GetCSV)
	router.Run(":8080")
}
