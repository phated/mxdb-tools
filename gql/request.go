package gql

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	graphqlURL = "https://api.graph.cool/simple/v1/metaxdb"
)

var token string

// SetToken sets the token within the package to make requests
func SetToken(t string) {
	token = t
}

// Request makes a graphql request
func Request(query []byte, variables interface{}) ([]byte, error) {
	reqBody, err := queryToRequest(query, variables)
	if err != nil {
		return nil, err
	}

	respBody, err := makeRequest(reqBody)
	if err != nil {
		return nil, err
	}

	return bodyToResponse(respBody)
}

/* request utils */

func queryToRequest(queryString []byte, variables interface{}) (*bytes.Buffer, error) {
	type payload struct {
		Query     string      `json:"query"`
		Variables interface{} `json:"variables,omitempty"`
	}

	replacer := strings.NewReplacer("\n", "")
	compactQuery := replacer.Replace(string(queryString))

	body, err := json.Marshal(payload{
		Query:     compactQuery,
		Variables: variables,
	})

	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(body), nil
}

func bodyToResponse(body []byte) ([]byte, error) {
	type GraphqlError struct {
		Message string `json:"message"`
	}

	type GraphQlResponse struct {
		Data   json.RawMessage `json:"data"`
		Errors []*GraphqlError `json:"errors"`
	}

	jsonResp := GraphQlResponse{}
	if err := json.Unmarshal(body, &jsonResp); err != nil {
		return nil, err
	}

	if len(jsonResp.Errors) != 0 {
		return nil, errors.New(jsonResp.Errors[0].Message)
	}

	return jsonResp.Data, nil
}

func makeRequest(body *bytes.Buffer) ([]byte, error) {
	req, err := http.NewRequest("POST", graphqlURL, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
