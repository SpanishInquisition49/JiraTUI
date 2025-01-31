package jira

import (
	"fmt"
	"log"

	jira "github.com/andygrunwald/go-jira"
)

type Jira struct {
	client *jira.Client
}

/**
  * Create a new Jira client
  * @param url string - The URL of the Jira server
  * @param username string - The username to authenticate with
  * @param password string - The password to authenticate with
  * @return Jira - A new Jira client
  */
func Connect(auth jira.BasicAuthTransport, url string) *Jira {
  client, err := jira.NewClient(auth.Client(), url)
  if err != nil {
    log.Println(fmt.Sprintf("Error creating Jira client: %s", err))
    return nil
  }
  return &Jira{client}
}

/**
 * Search for issues in Jira
 * @param jql string - The JQL query to search for issues
 * @return []jira.Issue - A list of issues matching the JQL query or nil if an error occurred
 */
func (j Jira) searchIssues(jql string) []jira.Issue {
	issues, _, err := j.client.Issue.Search(jql, nil)
	if err != nil {
		log.Println(fmt.Sprintf("Error searching for issues: %s", err))
		return nil
	}
	return issues
}

/**
 * Add a comment to a Jira issue
 * @param issueKey string - The key of the issue to add the comment to
 * @param comment string - The comment to add to the issue
 * @return bool - True if the comment was added successfully, false otherwise
 */
func (j Jira) addComment(issueKey string, comment string) bool {
	_, _, err := j.client.Issue.AddComment(issueKey, &jira.Comment{Body: comment})
	if err != nil {
		log.Println(fmt.Sprintf("Error adding comment to issue %s: %s", issueKey, err))
		return false
	}
	return true
}
