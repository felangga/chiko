package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const githubLatestReleaseURL = "https://api.github.com/repos/felangga/chiko/releases/latest"

type githubRelease struct {
	TagName string `json:"tag_name"`
}

func CheckLatestVersion(currentVersion string) (latest string, isNewer bool, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, githubLatestReleaseURL, nil)
	if err != nil {
		return "", false, err
	}
	req.Header.Set("User-Agent", fmt.Sprintf("chiko/%s", currentVersion))
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", false, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release githubRelease
	if err = json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", false, err
	}

	latest = strings.TrimPrefix(release.TagName, "v")
	isNewer, err = compareVersions(latest, currentVersion)
	if err != nil {
		return "", false, err
	}
	return latest, isNewer, nil
}

// compareVersions returns true if a is newer than b, false if equal or older.
// It compares all available segments and reports malformed non-numeric parts.
func compareVersions(a, b string) (bool, error) {
	aParts := strings.Split(a, ".")
	bParts := strings.Split(b, ".")
	max := len(aParts)
	if len(bParts) > max {
		max = len(bParts)
	}
	for i := 0; i < max; i++ {
		var av, bv int
		var err error
		if i < len(aParts) {
			av, err = strconv.Atoi(aParts[i])
			if err != nil {
				return false, fmt.Errorf("invalid version segment %q in %q", aParts[i], a)
			}
		}
		if i < len(bParts) {
			bv, err = strconv.Atoi(bParts[i])
			if err != nil {
				return false, fmt.Errorf("invalid version segment %q in %q", bParts[i], b)
			}
		}
		if av != bv {
			return av > bv, nil
		}
	}
	return false, nil
}
