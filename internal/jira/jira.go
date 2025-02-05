package jira

import (
	"fmt"
	"log"

	jira "github.com/andygrunwald/go-jira"
	"github.com/google/go-cmp/cmp"
)

type Client struct {
	client *jira.Client
}

type Issue struct {
	Key         string
	Summary     string
	Assignee    string
	Status      string
	Reporter    string
	Description string
}

/**
 * Create a new Jira client
 * @param email string - The email address of the user to authenticate as
 * @param api_token string - The API token of the user to authenticate as
 * @param url string - The URL of the Jira instance to connect to
 * @return Jira - A new Jira client
 */
func CreateClient(email string, api_token string, url string) *Client {
	tp := jira.BasicAuthTransport{
		Username: email,
		Password: api_token,
	}

	client, err := jira.NewClient(tp.Client(), url)
	if err != nil {
		log.Println(fmt.Sprintf("Error creating Jira client: %s", err))
		return nil
	}
	return &Client{client}
}

/**
 * Search for issues in Jira
 * @param jql string - The JQL query to search for issues
 * @return []jira.Issue - A list of issues matching the JQL query or nil if an error occurred
 */
func (j Client) SearchIssues(jql string) []Issue {
	issues, _, err := j.client.Issue.Search(jql, nil)
	if err != nil {
		log.Println(fmt.Sprintf("Error searching for issues: %s", err))
		return nil
	}

	var result []Issue
	for _, issue := range issues {
		// Check if fields are nil to avoid panics using the cmp package
		i := Issue{
			Key:         issue.Key,
			Assignee:    "",
			Reporter:    "",
			Description: "",
			Status:      "",
		}
		if issue.Fields.Assignee != nil {
			i.Assignee = issue.Fields.Assignee.DisplayName
		}
		if issue.Fields.Reporter != nil {
			i.Reporter = issue.Fields.Reporter.DisplayName
		}
		if issue.Fields.Description != "" {
			i.Description = issue.Fields.Description
		}
		if issue.Fields.Status != nil {
			i.Status = issue.Fields.Status.Name
		}
		result = append(result, i)
	}
	return result
}

/**
* Add a comment to a Jira issue
 * @param issue Issue - The issue to add the comment to
* @param comment string - The comment to add to the issue
* @return bool - True if the comment was added successfully, false otherwise
*/
func (j Client) AddComment(issue Issue, comment string) bool {
	_, _, err := j.client.Issue.AddComment(issue.Key, &jira.Comment{Body: comment})
	if err != nil {
		log.Println(fmt.Sprintf("Error adding comment to issue %s: %s", issue.Key, err))
		return false
	}
	return true
}
