package pdf

import (
	"bytes"
	"cpm-rad-backend/domain/config"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type DitoRequest struct {
	TemplateProjectPath string            `json:"templateProjectPath"`
	TemplateName        string            `json:"templateName"`
	Properties          map[string]string `json:"properties"`
	Data                any               `json:"data"`
}

type DitoResponse struct {
	Data string `json:"data"`
}

func ToResponse(data string) DitoResponse {
	return DitoResponse{
		Data: data,
	}
}

func GetReport(path string, data any) ([]byte, error) {
	ditoReq := DitoRequest{
		TemplateProjectPath: path,
		TemplateName:        "output",
		Properties: map[string]string{
			"pdfVersion": "2.0",
		},
		Data: data,
	}

	jsonBody, err := json.Marshal(ditoReq)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err

	}

	url := config.DitoApi
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept", "application/pdf")

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return body, nil
}
