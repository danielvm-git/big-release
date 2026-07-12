package algorithm

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
)

const (
	// FirstRelease is the default version for the first release
	FirstRelease = "1.0.0"
	
	// FirstPrerelease is the default prerelease number
	FirstPrerelease = "1"
)

// Calculator calculates the next version
type Calculator struct{}

// NewCalculator creates a new Calculator
func NewCalculator() *Calculator {
	return &Calculator{}
}

// CalculateNextVersion calculates the next version based on release type and last release
func (c *Calculator) CalculateNextVersion(lastRelease *Release, releaseType ReleaseType, branch *Branch) (string, error) {
	// First release
	if lastRelease == nil || lastRelease.Version == "" {
		if branch != nil && branch.Type == BranchTypePrerelease {
			return fmt.Sprintf("%s-%s.%s", FirstRelease, branch.Prerelease, FirstPrerelease), nil
		}
		return FirstRelease, nil
	}

	// Parse last version
	lastVersion, err := semver.NewVersion(lastRelease.Version)
	if err != nil {
		return "", fmt.Errorf("invalid last version %q: %w", lastRelease.Version, err)
	}

	// Calculate next version based on release type
	var nextVersion *semver.Version

	switch branch.Type {
	case BranchTypePrerelease:
		nextVersion, err = c.calculatePrerelease(lastVersion, releaseType, branch)
	case BranchTypeMaintenance:
		nextVersion, err = c.calculateMaintenance(lastVersion, releaseType)
	default:
		nextVersion, err = c.calculateRegular(lastVersion, releaseType)
	}

	if err != nil {
		return "", err
	}

	return nextVersion.String(), nil
}

// calculateRegular calculates the next regular version
func (c *Calculator) calculateRegular(lastVersion *semver.Version, releaseType ReleaseType) (*semver.Version, error) {
	var nextVersion semver.Version

	switch releaseType {
	case ReleaseTypeMajor:
		nextVersion = *semver.MustParse(fmt.Sprintf("%d.0.0", lastVersion.Major()+1))
	case ReleaseTypeMinor:
		nextVersion = *semver.MustParse(fmt.Sprintf("%d.%d.0", lastVersion.Major(), lastVersion.Minor()+1))
	case ReleaseTypePatch:
		nextVersion = *semver.MustParse(fmt.Sprintf("%d.%d.%d", lastVersion.Major(), lastVersion.Minor(), lastVersion.Patch()+1))
	default:
		return nil, fmt.Errorf("no releasable commits found")
	}

	return &nextVersion, nil
}

// calculatePrerelease calculates the next prerelease version
func (c *Calculator) calculatePrerelease(lastVersion *semver.Version, releaseType ReleaseType, branch *Branch) (*semver.Version, error) {
	preid := branch.Prerelease
	if preid == "" || preid == "true" {
		preid = branch.Name
	}

	// If current version is already a prerelease with same preid
	if lastVersion.Prerelease() != "" {
	 prereleases := lastVersion.Prerelease()
		if len(prereleases) > 0 && prereleases[0] == preid {
			// Increment prerelease number
			return c.incrementPrerelease(lastVersion), nil
		}
	}

	// Start new prerelease series
	var baseVersion *semver.Version
	switch releaseType {
	case ReleaseTypeMajor:
		baseVersion = semver.MustParse(fmt.Sprintf("%d.0.0", lastVersion.Major()+1))
	case ReleaseTypeMinor:
		baseVersion = semver.MustParse(fmt.Sprintf("%d.%d.0", lastVersion.Major(), lastVersion.Minor()+1))
	case ReleaseTypePatch:
		baseVersion = semver.MustParse(fmt.Sprintf("%d.%d.%d", lastVersion.Major(), lastVersion.Minor(), lastVersion.Patch()+1))
	}

	return semver.NewVersion(fmt.Sprintf("%s-%s.%s", baseVersion.String(), preid, FirstPrerelease)), nil
}

// calculateMaintenance calculates the next maintenance version
func (c *Calculator) calculateMaintenance(lastVersion *semver.Version, releaseType ReleaseType) (*semver.Version, error) {
	// Maintenance releases are typically patches
	var nextVersion semver.Version
	nextVersion = *semver.MustParse(fmt.Sprintf("%d.%d.%d", lastVersion.Major(), lastVersion.Minor(), lastVersion.Patch()+1))
	return &nextVersion, nil
}

// incrementPrerelease increments the prerelease number
func (c *Calculator) incrementPrerelease(version *semver.Version) *semver.Version {
	prerelease := version.Prerelease()
	if len(prerelease) < 2 {
		// Can't increment, return as-is
		return version
	}

	// Parse preid.N format
	preid := prerelease[0]
	var num int
	fmt.Sscanf(prerelease[1], "%d", &num)

	return semver.MustParse(fmt.Sprintf("%s-%s.%d", version.Original(), preid, num+1))
}
