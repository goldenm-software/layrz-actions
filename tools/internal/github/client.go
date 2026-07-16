package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const baseURL = "https://api.github.com"

type Client struct {
	token string
	repo  string
	http  *http.Client
}

func NewClient(token, repo string) *Client {
	return &Client{
		token: token,
		repo:  repo,
		http:  &http.Client{Timeout: 15 * time.Second},
	}
}

const maxAttempts = 3

func (c *Client) do(method, path string, body any) ([]byte, int, error) {
	var jsonBody []byte
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, 0, err
		}
		jsonBody = b
	}

	var data []byte
	var status int
	var err error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if attempt > 1 {
			time.Sleep(time.Duration(attempt-1) * 2 * time.Second)
		}
		data, status, err = c.doOnce(method, path, jsonBody)
		// Retry on network errors and 5xx responses; anything else is final.
		if err == nil && status < 500 {
			return data, status, nil
		}
	}
	return data, status, err
}

func (c *Client) doOnce(method, path string, jsonBody []byte) ([]byte, int, error) {
	var reqBody io.Reader
	if jsonBody != nil {
		reqBody = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequest(method, baseURL+path, reqBody)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	if jsonBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	return data, resp.StatusCode, err
}

func (c *Client) get(path string, out any) error {
	data, status, err := c.do("GET", path, nil)
	if err != nil {
		return err
	}
	if status < 200 || status >= 300 {
		return fmt.Errorf("GET %s: status %d: %s", path, status, string(data))
	}
	return json.Unmarshal(data, out)
}

func (c *Client) post(path string, body, out any) error {
	data, status, err := c.do("POST", path, body)
	if err != nil {
		return err
	}
	if status < 200 || status >= 300 {
		return fmt.Errorf("POST %s: status %d: %s", path, status, string(data))
	}
	if out != nil {
		return json.Unmarshal(data, out)
	}
	return nil
}

func (c *Client) patch(path string, body any) error {
	data, status, err := c.do("PATCH", path, body)
	if err != nil {
		return err
	}
	if status < 200 || status >= 300 {
		return fmt.Errorf("PATCH %s: status %d: %s", path, status, string(data))
	}
	return nil
}

type CommitResponse struct {
	Author *struct {
		Login string `json:"login"`
	} `json:"author"`
}

func (c *Client) GetCommitAuthorLogin(sha string) (string, error) {
	var resp CommitResponse
	err := c.get(fmt.Sprintf("/repos/%s/commits/%s", c.repo, sha), &resp)
	if err != nil {
		return "", err
	}
	if resp.Author == nil {
		return "", nil
	}
	return resp.Author.Login, nil
}

type UserSearchResponse struct {
	Items []struct {
		Login string `json:"login"`
	} `json:"items"`
}

func (c *Client) SearchUserByEmail(email string) (string, error) {
	var resp UserSearchResponse
	err := c.get(fmt.Sprintf("/search/users?q=%s+in:email", email), &resp)
	if err != nil {
		return "", err
	}
	if len(resp.Items) == 0 {
		return "", nil
	}
	return resp.Items[0].Login, nil
}

type IssueComment struct {
	ID   int64  `json:"id"`
	Body string `json:"body"`
	User struct {
		Type string `json:"type"`
	} `json:"user"`
}

func (c *Client) ListIssueComments(number int) ([]IssueComment, error) {
	var comments []IssueComment
	err := c.get(fmt.Sprintf("/repos/%s/issues/%d/comments?per_page=100", c.repo, number), &comments)
	return comments, err
}

func (c *Client) CreateIssueComment(number int, body string) error {
	return c.post(fmt.Sprintf("/repos/%s/issues/%d/comments", c.repo, number), map[string]string{"body": body}, nil)
}

func (c *Client) UpdateIssueComment(commentID int64, body string) error {
	return c.patch(fmt.Sprintf("/repos/%s/issues/comments/%d", c.repo, commentID), map[string]string{"body": body})
}
