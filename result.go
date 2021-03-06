package runscope

import (
	"errors"
	"fmt"
	"time"
)

// Result represents the outcome of a test run
type Result struct {
	AssertionsDefined int       `json:"assertions_defined"`
	AssertionsFailed  int       `json:"assertions_failed"`
	AssertionsPassed  int       `json:"assertions_passed"`
	BucketKey         string    `json:"bucket_key"`
	FinishedAt        float64   `json:"finished_at"`
	Region            string    `json:"region"`
	RequestsExecuted  int       `json:"requests_executed"`
	Result            string    `json:"result"`
	ScriptsDefined    int       `json:"scripts_defined"`
	ScriptsFailed     int       `json:"scripts_failed"`
	ScriptsPassed     int       `json:"scripts_passed"`
	StartedAt         float64   `json:"started_at"`
	TestRunID         string    `json:"test_run_id"`
	TestRunURL        string    `json:"test_run_url"`
	TestID            string    `json:"test_id"`
	VariablesDefined  int       `json:"variables_defined"`
	VariablesFailed   int       `json:"variables_failed"`
	VariablesPassed   int       `json:"variables_passed"`
	EnvironmentID     string    `json:"environment_id"`
	EnvironmentName   string    `json:"environment_name"`
	Requests          []Request `json:"requests"`
}

// Request represents the result of a request made by a given test
type Request struct {
	Result            string      `json:"result"`
	URL               string      `json:"url"`
	Method            string      `json:"method"`
	AssertionsDefined int         `json:"assertions_defined"`
	AssertionsFailed  int         `json:"assertions_failed"`
	AssertionsPassed  int         `json:"assertions_passed"`
	ScriptsDefined    int         `json:"scripts_defined"`
	ScriptsFailed     int         `json:"scripts_failed"`
	ScriptsPassed     int         `json:"scripts_passed"`
	VariablesDefined  int         `json:"variables_defined"`
	VariablesFailed   int         `json:"variables_failed"`
	VariablesPassed   int         `json:"variables_passed"`
	Assertions        []Assertion `json:"assertions"`
	Scripts           []Script    `json:"scripts"`
	Variables         []Variable  `json:"variables"`
}

// ListResults returns all results for a given test
func (client *Client) ListResults(bucketKey string, testID string) ([]Result, error) {
	var results = []Result{}

	path := fmt.Sprintf("buckets/%s/tests/%s/results", bucketKey, testID)
	content, err := client.Get(path)
	if err != nil {
		return results, err
	}

	err = unmarshal(content, &results)
	return results, err
}

// FilterResults returns results for a given test and supplied conditions
func (client *Client) FilterResults(bucketKey, testID string, count int64, since, before *time.Time) ([]Result, error) {
	var results = []Result{}

	filterQs, err := client.buildFilterQS(count, since, before)
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("buckets/%s/tests/%s/results%s", bucketKey, testID, filterQs)
	content, err := client.Get(path)
	if err != nil {
		return results, err
	}

	err = unmarshal(content, &results)
	return results, err
}

// GetResult returns a more detail result for a result ID
func (client *Client) GetResult(bucketKey string, testID string, testRunID string) (Result, error) {
	var result = Result{}

	path := fmt.Sprintf("buckets/%s/tests/%s/results/%s", bucketKey, testID, testRunID)
	content, err := client.Get(path)
	if err != nil {
		return result, err
	}

	err = unmarshal(content, &result)
	return result, err
}

// GetResultLatest returns the last known result for a given test
func (client *Client) GetResultLatest(bucketKey string, testID string) (Result, error) {
	return client.GetResult(bucketKey, testID, "latest")
}

// Builds filter for list results (count maximum is 50 and since/before are exclusive!)
func (client *Client) buildFilterQS(count int64, since, before *time.Time) (string, error) {
	if since != nil && before != nil {
		return "", errors.New("Filter parameters 'since' and 'before' are exclusive!")
	}

	if count > 50 {
		return "", errors.New("Filter parameter 'count' has maximal value of 50!")
	}

	var qs = fmt.Sprintf("?count=%v", count)
	if since != nil {
		qs += fmt.Sprintf("&since=%f", unixTimestampToFloat(*since))
	}
	if before != nil {
		qs += fmt.Sprintf("&before=%f", unixTimestampToFloat(*before))
	}
	return qs, nil
}
