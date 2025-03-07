package ctx

import (
	"github.com/cli/go-gh"
	"github.com/cli/go-gh/pkg/api"
	graphql "github.com/cli/shurcooL-graphql"
	"github.com/prnk28/gh-task/internal/ghc"
)

// listOrgs returns a list of organization names the authenticated user is a member of
func listOrgs() ([]string, error) {
	// Create a GraphQL client
	client, err := gh.GQLClient(nil)
	if err != nil {
		return nil, err
	}

	// Define the query to get organizations
	var query struct {
		Viewer struct {
			Organizations struct {
				Nodes []struct {
					Login string
				}
				PageInfo struct {
					HasNextPage bool
					EndCursor   string
				}
			} `graphql:"organizations(first: 100)"`
		}
	}

	// Execute the query
	err = client.Query("UserOrganizations", &query, nil)
	if err != nil {
		return nil, err
	}

	// Extract organization names
	orgs := make([]string, 0, len(query.Viewer.Organizations.Nodes))
	for _, org := range query.Viewer.Organizations.Nodes {
		orgs = append(orgs, org.Login)
	}

	// If there are more pages, fetch them
	if query.Viewer.Organizations.PageInfo.HasNextPage {
		err = fetchRemainingOrgs(client, &orgs, query.Viewer.Organizations.PageInfo.EndCursor)
		if err != nil {
			return nil, err
		}
	}

	return orgs, nil
}

// fetchRemainingOrgs recursively fetches additional pages of organizations
func fetchRemainingOrgs(client api.GQLClient, orgs *[]string, cursor string) error {
	var query struct {
		Viewer struct {
			Organizations struct {
				Nodes []struct {
					Login string
				}
				PageInfo struct {
					HasNextPage bool
					EndCursor   string
				}
			} `graphql:"organizations(first: 100, after: $cursor)"`
		}
	}

	variables := map[string]interface{}{
		"cursor": graphql.String(cursor),
	}

	err := client.Query("UserOrganizationsPaginated", &query, variables)
	if err != nil {
		return err
	}

	// Add organizations from this page
	for _, org := range query.Viewer.Organizations.Nodes {
		*orgs = append(*orgs, org.Login)
	}

	// If there are more pages, fetch them recursively
	if query.Viewer.Organizations.PageInfo.HasNextPage {
		return fetchRemainingOrgs(client, orgs, query.Viewer.Organizations.PageInfo.EndCursor)
	}

	return nil
}

func filterActiveOrgs() ([]string, error) {
	// Create new context if it doesn't exist
	orgs, err := listOrgs()
	if err != nil {
		return nil, err
	}
	activeOrgs := make([]string, 0)
	for _, org := range orgs {
		if ghc.OrgHasRepo(org, ".github") {
			activeOrgs = append(activeOrgs, org)
		}
	}
	return activeOrgs, nil
}
