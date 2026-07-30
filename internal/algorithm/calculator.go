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

// CalculateNextVersion calculates the next version based on release type and last release.
// If initialVersion is non-empty it is used for first releases; otherwise FirstRelease is the default.
func (c *Calculator) CalculateNextVersion(lastRelease *Release, releaseType ReleaseType, branch *Branch, initialVersion string) (string, error) {
	if lastRelease == nil || lastRelease.Version == "" {
		base := FirstRelease
		if initialVersion != "" {
			base = initialVersion
		}
		if branch != nil && branch.Type == BranchTypePrerelease {
			return fmt.Sprintf("%s-%s.%s", base, branch.Prerelease, FirstPrerelease), nil
		}
		return base, nil
	}

	// Parse last version
	lastVersion, err := semver.NewVersion(lastRelease.Version)
	if err != nil {
		return "", fmt.Errorf("invalid last version %q: %w", lastRelease.Version, err)
	}

	// Calculate next version based on release type.
	// Default to release branch if nil (defense-in-depth for public API callers).
	if branch == nil {
		branch = &Branch{Type: BranchTypeRelease}
	}

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
		parts := splitPrerelease(lastVersion.Prerelease())
		if len(parts) > 0 && parts[0] == preid {
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
	default:
		return nil, fmt.Errorf("cannot start prerelease series for release type %q", releaseType)
	}

	next, err := semver.NewVersion(fmt.Sprintf("%s-%s.%s", baseVersion.String(), preid, FirstPrerelease))
	if err != nil {
		return nil, err
	}
	return next, nil
}

// calculateMaintenance calculates the next maintenance version.
// Unlike regular releases, maintenance branches respect the release type
// for minor/major bumps but always stay within the same major.minor line
// for patches (BUG-calculator-maintenance-ignores-type).
func (c *Calculator) calculateMaintenance(lastVersion *semver.Version, releaseType ReleaseType) (*semver.Version, error) {
	switch releaseType {
	case ReleaseTypeMajor:
		// Major bump on maintenance: bump major, reset minor+patch.
		nextVersion := *semver.MustParse(fmt.Sprintf("%d.0.0", lastVersion.Major()+1))
		return &nextVersion, nil
	case ReleaseTypeMinor:
		// Minor bump on maintenance: bump minor, reset patch.
		nextVersion := *semver.MustParse(fmt.Sprintf("%d.%d.0", lastVersion.Major(), lastVersion.Minor()+1))
		return &nextVersion, nil
	default:
		// Patch (default for maintenance): bump patch within same major.minor.
		nextVersion := *semver.MustParse(fmt.Sprintf("%d.%d.%d", lastVersion.Major(), lastVersion.Minor(), lastVersion.Patch()+1))
		return &nextVersion, nil
	}
}

// incrementPrerelease increments the prerelease number
func (c *Calculator) incrementPrerelease(version *semver.Version) *semver.Version {
	parts := splitPrerelease(version.Prerelease())
	if len(parts) < 2 {
		return version
	}

	preid := parts[0]
	var num int
	_, _ = fmt.Sscanf(parts[1], "%d", &num)

	return semver.MustParse(fmt.Sprintf("%d.%d.%d-%s.%d", version.Major(), version.Minor(), version.Patch(), preid, num+1))
}

// splitPrerelease splits a prerelease string like "alpha.1" into ["alpha", "1"]
func splitPrerelease(prerelease string) []string {
	if prerelease == "" {
		return nil
	}
	result := make([]string, 0, 2)
	current := ""
	for i := 0; i < len(prerelease); i++ {
		if prerelease[i] == '.' {
			if current != "" {
				result = append(result, current)
			}
			current = ""
		} else {
			current += string(prerelease[i])
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}
