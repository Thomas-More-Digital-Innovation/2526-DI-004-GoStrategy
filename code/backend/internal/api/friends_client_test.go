package api

import (
	"bytes"
	"digital-innovation/gostrategy/internal/db"
	"digital-innovation/gostrategy/internal/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

// friendsClient is a registered user talking to the real router over HTTP.
type friendsClient struct {
	t       *testing.T
	baseURL string
	cookies []*http.Cookie
	id      int
}

func registerFriendsClient(t *testing.T, ts *httptest.Server, name string) *friendsClient {
	t.Helper()
	// #nosec G117
	body, _ := json.Marshal(models.CreateUserRequest{Username: name, Password: db.TestPassword})
	resp, err := http.Post(ts.URL+"/users/register", contentTypeJSON, bytes.NewBuffer(body))
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var user models.UserDTO
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&user))
	return &friendsClient{t: t, baseURL: ts.URL, cookies: resp.Cookies(), id: user.ID}
}

// call sends a request and returns the status and raw body. Pass out to also decode JSON into it.
func (c *friendsClient) call(method, path string, body, out any) (int, string) {
	c.t.Helper()
	var reader *bytes.Reader
	if body != nil {
		raw, _ := json.Marshal(body)
		reader = bytes.NewReader(raw)
	} else {
		reader = bytes.NewReader(nil)
	}
	req, err := http.NewRequest(method, c.baseURL+path, reader)
	require.NoError(c.t, err)
	req.Header.Set("Content-Type", contentTypeJSON)
	for _, cookie := range c.cookies {
		req.AddCookie(cookie)
	}
	resp, err := http.DefaultClient.Do(req)
	require.NoError(c.t, err)
	defer func() { _ = resp.Body.Close() }()

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(resp.Body)
	if out != nil && resp.StatusCode < http.StatusBadRequest && buf.Len() > 0 {
		require.NoError(c.t, json.Unmarshal(buf.Bytes(), out))
	}
	return resp.StatusCode, buf.String()
}

func (c *friendsClient) friendCode() string {
	c.t.Helper()
	var res models.FriendCodeResponse
	status, body := c.call(http.MethodPost, "/users/me/friend-code", nil, &res)
	require.Equal(c.t, http.StatusOK, status, body)
	return res.Code
}

func (c *friendsClient) list() models.FriendListResponse {
	c.t.Helper()
	var res models.FriendListResponse
	status, body := c.call(http.MethodGet, "/friends", nil, &res)
	require.Equal(c.t, http.StatusOK, status, body)
	return res
}

func (c *friendsClient) sendRequest(code string, out any) int {
	status, _ := c.call(http.MethodPost, "/friends/requests", models.AddFriendRequest{Code: code}, out)
	return status
}
