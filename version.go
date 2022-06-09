package ngrok

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Version struct {
	Major    int
	Minor    int
	Revision int
}

var versionPattern = regexp.MustCompile(`\d+(\.\d+)+`)

func ParseVersion(version string) *Version {
	if match := versionPattern.FindString(version); match != "" {
		parts := strings.Split(match, ".")
		nums := make([]int, len(parts))
		for i, str := range parts {
			nums[i], _ = strconv.Atoi(str)
		}
		result := Version{Major: nums[0], Minor: nums[1]}
		if len(nums) > 2 {
			result.Revision = nums[2]
		}
		// Other values are discarded
		return &result
	}
	return nil
}

func (v *Version) String() string { return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Revision) }
