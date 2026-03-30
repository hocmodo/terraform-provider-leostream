// pool_http_helpers.go is a temporary shim for API fields missing from the
// upstream leostream-client-go PoolAwsCenter struct (confirmed gap in v0.1.7).
//
// Go silently drops undeclared struct fields during JSON serialisation, so
// fields like launch_template_version would never reach the API through the
// normal typed client calls. These helpers work around that by injecting extra
// fields into the JSON payload after serialisation.
//
// TODO: delete this file once leostream-client-go adds the missing fields to
// PoolAwsCenter. Steps: wire fields directly through centerConfig, revert
// CreateNested/UpdateNested to r.client.CreatePool/UpdatePool, read fields
// directly from poolConfig.Provision.Center in Read.

package leostream

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	leo "gitlab.com/hocmodo/leostream-client-go"
)

// poolReadExtra reads provision.center as a free-form map to extract fields absent from leo.PoolAwsCenter.
type poolReadExtra struct {
	Provision *struct {
		Center map[string]interface{} `json:"center,omitempty"`
	} `json:"provision,omitempty"`
}

// addPoolCenterExtraFields injects extra key/value pairs into provision.center after JSON serialisation.
func addPoolCenterExtraFields(poolConfig leo.Pool, centerExtraFields map[string]string) ([]byte, error) {
	rawPool, err := json.Marshal(poolConfig)
	if err != nil {
		return nil, err
	}

	if len(centerExtraFields) == 0 {
		return rawPool, nil
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(rawPool, &payload); err != nil {
		return nil, err
	}

	provision, ok := payload["provision"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("pool payload missing provision object")
	}

	center, ok := provision["center"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("pool payload missing center object")
	}

	for key, value := range centerExtraFields {
		if value == "" {
			continue
		}
		center[key] = value
	}

	return json.Marshal(payload)
}

// doAuthorizedRequest executes an authenticated HTTP request, mirroring the upstream client auth setup.
func doAuthorizedRequest(client *leo.Client, req *http.Request, authToken *string) ([]byte, error) {
	token := client.Token
	if authToken != nil {
		token = *authToken
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	u, _ := req.URL.Parse(client.HostURL)
	for _, cookie := range client.CookieJar.Cookies(u) {
		req.AddCookie(cookie)
	}

	res, err := client.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	for _, cookie := range res.Cookies() {
		client.CookieJar.SetCookies(u, []*http.Cookie{cookie})
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, nil
}

// createPoolWithCenterExtraFields replaces client.CreatePool, injecting extra center fields before POST.
func createPoolWithCenterExtraFields(client *leo.Client, poolConfig leo.Pool, centerExtraFields map[string]string, authToken *string) (*leo.PoolsStored, error) {
	payload, err := addPoolCenterExtraFields(poolConfig, centerExtraFields)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/rest/v1/pools", client.HostURL), strings.NewReader(string(payload)))
	if err != nil {
		return nil, err
	}

	body, err := doAuthorizedRequest(client, req, authToken)
	if err != nil {
		return nil, err
	}

	stored := leo.PoolsStored{}
	if err := json.Unmarshal(body, &stored); err != nil {
		return nil, err
	}

	return &stored, nil
}

// updatePoolWithCenterExtraFields replaces client.UpdatePool, injecting extra center fields before PUT.
func updatePoolWithCenterExtraFields(client *leo.Client, poolID string, poolConfig leo.Pool, centerExtraFields map[string]string, authToken *string) (*leo.PoolsStored, error) {
	payload, err := addPoolCenterExtraFields(poolConfig, centerExtraFields)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/rest/v1/pools/%s", client.HostURL, poolID), strings.NewReader(string(payload)))
	if err != nil {
		return nil, err
	}

	body, err := doAuthorizedRequest(client, req, authToken)
	if err != nil {
		return nil, err
	}

	stored := leo.PoolsStored{}
	if err := json.Unmarshal(body, &stored); err != nil {
		return nil, err
	}

	return &stored, nil
}

// getPoolCenterStringField reads a named string field from provision.center that is absent from leo.PoolAwsCenter.
func getPoolCenterStringField(client leo.Client, poolID string, field string) (string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/rest/v1/pools/%s", client.HostURL, poolID), nil)
	if err != nil {
		return "", err
	}

	body, err := doAuthorizedRequest(&client, req, nil)
	if err != nil {
		return "", err
	}

	pool := poolReadExtra{}
	if err := json.Unmarshal(body, &pool); err != nil {
		return "", err
	}

	if pool.Provision == nil || pool.Provision.Center == nil {
		return "", nil
	}

	value, ok := pool.Provision.Center[field]
	if !ok {
		return "", nil
	}

	stringValue, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("pool center field %q is not a string", field)
	}

	return stringValue, nil
}
