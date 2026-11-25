package utils

// OrganizationProjectLabel returns "<organizationSlug>/<projectName>" for a project of an organization.
// When the organization slug is unknown, it only returns the project name, which may be empty.
func OrganizationProjectLabel(organizationSlug, projectName string) string {
	if organizationSlug == "" || projectName == "" {
		return projectName
	}
	return organizationSlug + "/" + projectName
}
