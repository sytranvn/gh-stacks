package main

import (
	"github.com/cli/cli/api"
	"github.com/cli/cli/pkg/cmd/factory"
)

type Language struct {
	Name  string
	Color string
	ID    string
}

type Repository struct {
	Languages     []Language
	NameWithOwner string
}

type ContributionsPayload struct {
	Contributions   []Repository
	OwnRepositories []Repository
}

func Contributions(client *api.Client) (*ContributionsPayload, error) {
	type response struct {
		Viewer struct {
			RepositoriesContributedTo struct {
				TotalCount int
				Nodes      []Repository
			}
			Repositories struct {
				TotalCount int
				Nodes      []Repository
			}
		}
	}

	query := `
	query {
		viewer {
			repositoriesContributedTo(first: 100, contributionTypes: [COMMIT, PULL_REQUEST, REPOSITORY]) {
				totalCount
				nodes {
					nameWithOwner
					languages(first: 5, orderBy: {field: SIZE, direction: DESC}) {
						nodes {
							color
							name
						}
					}
				}
				pageInfo {
					endCursor
					hasNextPage
				}
			}
			repositories(first: 100, isFork: false) {
				totalCount
				nodes {
					nameWithOwner
					languages(first: 5, orderBy: {direction: DESC, field: SIZE}) {
						nodes {
							color
							name
						}
					}
				}
				pageInfo {
					endCursor
					hasNextPage
				}
			}
		}
	}`

	variables := map[string]interface{}{}

	var resp response
	err := client.GraphQL(repo.RepoHost(), query, variables, &resp)
	if err != nil {
		return nil, err
	}

	payload := ContributionsPayload{
		Contributions:   resp.Viewer.RepositoriesContributedTo.Nodes,
		OwnRepositories: resp.Viewer.Repositories.Nodes,
	}

	return &payload, nil
}

func main() {
	buildDate := build.Date
	buildVersion := build.Version
	cmdFactory := factory.New(buildVersion)
	apiClient := api.NewClientFromHTTP(cmdFactory.HttpClient)
	Contributions(apiClient)
}
