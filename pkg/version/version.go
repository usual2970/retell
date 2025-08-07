package version

import "fmt"

type Version struct {
	Major int
	Minor int
	Patch int
}

func FromString(versionStr string) (Version, error) {
	var v Version
	_, err := fmt.Sscanf(versionStr, "%d.%d.%d", &v.Major, &v.Minor, &v.Patch)
	if err != nil {
		return Version{}, fmt.Errorf("invalid version format: %s", versionStr)
	}
	return v, nil
}

func GreaterThan(v1, v2 string) (bool, error) {
	version1, err := FromString(v1)
	if err != nil {
		return false, err
	}
	version2, err := FromString(v2)
	if err != nil {
		return false, err
	}
	return version1.GreaterThan(version2), nil
}

func (v Version) GreaterThan(other Version) bool {
	if v.Major != other.Major {
		return v.Major > other.Major
	}
	if v.Minor != other.Minor {
		return v.Minor > other.Minor
	}
	return v.Patch > other.Patch
}

func (v Version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}
