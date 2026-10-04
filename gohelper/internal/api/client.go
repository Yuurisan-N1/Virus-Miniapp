package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
)

const graphqlURL = "https://virusgift.pro/api/graphql/query"

const authQuery = "mutation authTelegramInitData($initData: String!, $refCode: String) { authTelegramInitData(initData: $initData, refCode: $refCode) { token success __typename } }"
const meQuery = "query me { me { firstName balance starsBalance nextFreeSpin nextCaseFreeSpin __typename } }"
const spinQuery = "mutation startRouletteSpin($input: StartRouletteSpinInput!) { startRouletteSpin(input: $input) { success prize { name __typename } storyReward __typename } }"
const openCaseQuery = "mutation openCase($id: ID!, $demo: Boolean!) { openCase(id: $id, demo: $demo) { success prize { name __typename } __typename } }"

type gqlRequest struct {
	OperationName string      `json:"operationName"`
	Variables     interface{} `json:"variables"`
	Query         string      `json:"query"`
}

type Client struct {
	Token         string
	Client        *http.Client
	ClientVersion string
	Timezone      string
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randomClientVersion() string {
	chars := "0123456789abcdef"
	out := make([]byte, 8)
	for i := range out {
		out[i] = chars[rand.Intn(len(chars))]
	}
	return string(out)
}

func randomTimezone() string {
	zones := []string{
		"Asia/Jakarta",
		"Asia/Bangkok",
		"Asia/Singapore",
		"Asia/Tokyo",
		"Europe/London",
		"Europe/Berlin",
		"America/New_York",
		"America/Chicago",
		"UTC",
	}
	return zones[rand.Intn(len(zones))]
}

func NewClient(token string) *Client {
	return &Client{
		Token:         token,
		Client:        &http.Client{Timeout: 30 * time.Second},
		ClientVersion: randomClientVersion(),
		Timezone:      randomTimezone(),
	}
}

func (c *Client) postGraphQL(reqBody gqlRequest, extraHeaders map[string]string) ([]byte, error) {
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", graphqlURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("accept", "*/*")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("Origin", "https://virusgift.pro")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36")
	req.Header.Set("x-client-version", c.ClientVersion)
	req.Header.Set("x-timezone", c.Timezone)
	req.Header.Set("apollo-require-preflight", "*")
	for k, v := range extraHeaders {
		req.Header.Set(k, v)
	}
	if c.Token != "" {
		req.Header.Set("authorization", "Bearer "+c.Token)
	}
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("status %d %s", resp.StatusCode, string(body))
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func decodeAuth(body []byte) (*AuthResponse, error) {
	var single AuthResponse
	if err := json.Unmarshal(body, &single); err == nil && single.Data.AuthTelegramInitData.Token != "" {
		return &single, nil
	}
	var batch []AuthResponse
	if err := json.Unmarshal(body, &batch); err != nil {
		return nil, err
	}
	if len(batch) == 0 {
		return nil, fmt.Errorf("empty auth batch response")
	}
	return &batch[0], nil
}

func decodeMe(body []byte) (*MeResponse, error) {
	var single MeResponse
	if err := json.Unmarshal(body, &single); err == nil && single.Data.Me.FirstName != "" {
		return &single, nil
	}
	var batch []MeResponse
	if err := json.Unmarshal(body, &batch); err != nil {
		return nil, err
	}
	if len(batch) == 0 {
		return nil, fmt.Errorf("empty me batch response")
	}
	return &batch[0], nil
}

func decodeOpenCase(body []byte) (*OpenCaseResponse, error) {
	var single OpenCaseResponse
	if err := json.Unmarshal(body, &single); err == nil && single.Data.OpenCase.Prize.Name != "" {
		return &single, nil
	}
	var batch []OpenCaseResponse
	if err := json.Unmarshal(body, &batch); err != nil {
		return nil, err
	}
	if len(batch) == 0 {
		return nil, fmt.Errorf("empty openCase batch response")
	}
	return &batch[0], nil
}

func (c *Client) Auth(initData, refCode string) (string, error) {
	variables := map[string]interface{}{
		"initData": initData,
		"refCode":  refCode,
	}
	reqBody := gqlRequest{
		OperationName: "authTelegramInitData",
		Variables:     variables,
		Query:         authQuery,
	}
	extra := map[string]string{
		"x-batch": "true",
		"Referer": "https://virusgift.pro/?tgWebAppStartParam=openRoulette",
	}
	body, err := c.postGraphQL(reqBody, extra)
	if err != nil {
		return "", err
	}
	resp, err := decodeAuth(body)
	if err != nil {
		return "", err
	}
	auth := resp.Data.AuthTelegramInitData
	if !auth.Success {
		return "", fmt.Errorf("auth not successful")
	}
	if auth.Token == "" {
		return "", fmt.Errorf("token is empty")
	}
	return auth.Token, nil
}

func (c *Client) Me() (*MeBlock, error) {
	reqBody := gqlRequest{
		OperationName: "me",
		Variables:     map[string]interface{}{},
		Query:         meQuery,
	}
	extra := map[string]string{
		"x-batch": "true",
		"Referer": "https://virusgift.pro/roulette",
	}
	body, err := c.postGraphQL(reqBody, extra)
	if err != nil {
		return nil, err
	}
	resp, err := decodeMe(body)
	if err != nil {
		return nil, err
	}
	return &resp.Data.Me, nil
}

func (c *Client) Spin() (*SpinInner, error) {
	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"type": "X1",
			"demo": false,
		},
	}
	reqBody := gqlRequest{
		OperationName: "startRouletteSpin",
		Variables:     variables,
		Query:         spinQuery,
	}
	extra := map[string]string{
		"x-batch": "false",
		"Referer": "https://virusgift.pro/roulette/main",
	}
	body, err := c.postGraphQL(reqBody, extra)
	if err != nil {
		return nil, err
	}
	var resp SpinResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	result := resp.Data.StartRouletteSpin
	return &result, nil
}

func (c *Client) OpenCase(caseID int) (*CaseInner, error) {
	variables := map[string]interface{}{
		"id":   caseID,
		"demo": false,
	}
	reqBody := gqlRequest{
		OperationName: "openCase",
		Variables:     variables,
		Query:         openCaseQuery,
	}
	extra := map[string]string{
		"x-batch": "true",
		"Referer": "https://virusgift.pro/roulette/cases",
	}
	body, err := c.postGraphQL(reqBody, extra)
	if err != nil {
		return nil, err
	}
	resp, err := decodeOpenCase(body)
	if err != nil {
		return nil, err
	}
	result := resp.Data.OpenCase
	return &result, nil
}
