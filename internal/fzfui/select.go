package fzfui

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/chege/azfind/internal/cache"
)

// trunc shortens a string to maxRunes characters, adding "..." if needed.
// It operates on runes to avoid breaking multi-byte characters mid-codepoint.
func trunc(s string, maxRunes int) string {
	if maxRunes <= 0 {
		return ""
	}

	runes := []rune(s)
	if len(runes) <= maxRunes {
		return s
	}
	if maxRunes <= 3 {
		return string(runes[:maxRunes])
	}
	return string(runes[:maxRunes-3]) + "..."
}

// shortType returns the last segment of the resource type, e.g. "microsoft.app/containerapps" → "containerapps".
func shortType(full string) string {
	if full == "" {
		return ""
	}
	if i := strings.LastIndex(full, "/"); i >= 0 && i+1 < len(full) {
		return full[i+1:]
	}
	return full
}

// SelectResource runs fzf on the given resources and returns the selected one.
func SelectResource(resources []cache.Resource, initialQuery string) (*cache.Resource, error) {
	if len(resources) == 0 {
		return nil, nil
	}

	// Fixed column widths for visible part; everything else goes into hidden fields.
	const (
		nameWidth = 40
		typeWidth = 24
		rgWidth   = 30
	)

	var buf bytes.Buffer
	for _, r := range resources {
		name := trunc(r.Name, nameWidth)
		typeShort := trunc(shortType(r.Type), typeWidth)
		rg := trunc(r.ResourceGroup, rgWidth)

		visible := fmt.Sprintf("%-*s | %-*s | %-*s",
			nameWidth, name,
			typeWidth, typeShort,
			rgWidth, rg,
		)

		// Hidden full fields after the first tab:
		// {2}=Name, {3}=Type, {4}=ResourceGroup, {5}=SubscriptionID, {6}=Location, {7}=ID
		line := fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\t%s",
			visible,
			r.Name,
			r.Type,
			r.ResourceGroup,
			r.SubscriptionID,
			r.Location,
			r.ID,
		)

		buf.WriteString(line)
		buf.WriteByte('\n')
	}

	// Create fzf command: show only column 1, but search across all fields.
	cmd := exec.Command("fzf",
		"--ansi",
		"--delimiter", "\t",
		"--with-nth", "1",
		"--nth", "1..7",
		"--preview", "echo -e \"Type:            {3}\\nName:            {2}\\nSubscription:    {5}\\nResource group:  {4}\\nLocation:        {6}\\nID:              {7}\"",
		"--preview-window", "right:40%",
		"--query="+initialQuery,
	)
	cmd.Stdin = &buf

	out, err := cmd.Output()
	if err != nil {
		// User pressed ESC or CTRL-C → no selection
		return nil, nil
	}

	selection := strings.TrimSpace(string(out))
	parts := strings.Split(selection, "\t")
	if len(parts) < 7 {
		return nil, nil
	}

	// Hidden full values
	nameFull := parts[1]
	typeFull := parts[2]
	rgFull := parts[3]
	sub := parts[4]
	id := parts[6]

	// Prefer matching by ID (should be unique); fall back to name/type/rg/sub if needed.
	for _, r := range resources {
		if r.ID == id && r.SubscriptionID == sub {
			return &r, nil
		}
	}
	for _, r := range resources {
		if r.Name == nameFull && r.Type == typeFull && r.ResourceGroup == rgFull && r.SubscriptionID == sub {
			return &r, nil
		}
	}

	return nil, nil
}
