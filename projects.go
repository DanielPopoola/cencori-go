package cencori

import (
	"context"
	"fmt"
)

// ProjectsService provides methods for managing project-related operations.
// It uses a Client to communicate with the projects API endpoints.
type ProjectsService struct {
	client *Client
}

// List retrieves all projects for the specified organization.
//
// Parameters:
//   - ctx: context.Context for request cancellation and timeouts
//   - orgSlug: the organization slug identifier
//
// Returns:
//   - []Project: a slice of projects belonging to the organization
//   - error: an error if the request fails
func (s *ProjectsService) List(ctx context.Context, orgSlug string) ([]Project, error) {
	path := fmt.Sprintf("/api/organizations/%s/projects", orgSlug)

	type response struct {
		Projects []Project `json:"projects"`
	}

	resp, err := doRequest[any, response](s.client, ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	return resp.Projects, nil
}

// Create creates a new project within the specified organization.
// It sends a POST request to the organization's projects endpoint with the provided parameters.
// Returns the created Project or an error if the request fails.
func (s *ProjectsService) Create(ctx context.Context, orgSlug string, params CreateProjectParams) (*Project, error) {
	path := fmt.Sprintf("/api/organizations/%s/projects", orgSlug)
	return doRequest[CreateProjectParams, Project](s.client, ctx, "POST", path, &params)
}

// Get retrieves a project by its organization slug and project slug.
// It returns the project details or an error if the request fails.
func (s *ProjectsService) Get(ctx context.Context, orgSlug, projectSlug string) (*Project, error) {
	path := fmt.Sprintf("/api/organizations/%s/projects/%s", orgSlug, projectSlug)
	return doRequest[any, Project](s.client, ctx, "GET", path, nil)
}

// Update updates a project in the specified organization.
// It sends a PATCH request to the projects API endpoint with the given organization and project slugs.
// Returns an error if the request fails.
func (s *ProjectsService) Update(ctx context.Context, orgSlug, projectSlug string, params CreateProjectParams) error {
	path := fmt.Sprintf("/api/organizations/%s/projects/%s", orgSlug, projectSlug)
	_, err := doRequest[CreateProjectParams, any](s.client, ctx, "PATCH", path, &params)
	return err
}

// Delete deletes a project in the specified organization.
// It sends a DELETE request to the projects API endpoint with the given organization and project slugs.
// Returns an error if the request fails.
func (s *ProjectsService) Delete(ctx context.Context, orgSlug, projectSlug string) error {
	path := fmt.Sprintf("/api/organizations/%s/projects/%s", orgSlug, projectSlug)
	_, err := doRequest[any, any](s.client, ctx, "DELETE", path, nil)
	return err
}
