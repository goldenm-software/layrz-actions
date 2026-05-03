package github

import (
	"fmt"
	"strings"
)

func (c *Client) PostOrUpdateComment(prNumber int, marker, body string) error {
	comments, err := c.ListIssueComments(prNumber)
	if err != nil {
		return fmt.Errorf("listing comments: %w", err)
	}

	for _, comment := range comments {
		if comment.User.Type == "Bot" && strings.Contains(comment.Body, marker) {
			return c.UpdateIssueComment(comment.ID, body)
		}
	}

	return c.CreateIssueComment(prNumber, body)
}
