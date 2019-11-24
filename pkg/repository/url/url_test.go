package url

import (
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name        string
		url         string
		repo, owner string
		shouldErr   bool
	}{
		{
			name:  "https github URL",
			url:   "https://github.com/corneliusweig/rakkess",
			repo:  "rakkess",
			owner: "corneliusweig",
		},
		{
			name:  "http github URL",
			url:   "http://github.com/corneliusweig/rakkess",
			repo:  "rakkess",
			owner: "corneliusweig",
		},
		{
			name:  "github URL with anchor",
			url:   "https://github.com/corneliusweig/rakkess#howto",
			repo:  "rakkess",
			owner: "corneliusweig",
		},
		{
			name:  "github URL with query param",
			url:   "https://github.com/corneliusweig/rakkess?howto=not",
			repo:  "rakkess",
			owner: "corneliusweig",
		},
		{
			name:  "github URL with subpath",
			url:   "https://github.com/corneliusweig/rakkess/nested",
			repo:  "rakkess",
			owner: "corneliusweig",
		},
		{
			name:  "github URL without protocol",
			url:   "github.com/corneliusweig/rakkess",
			repo:  "rakkess",
			owner: "corneliusweig",
		},
		{
			name:      "gitLAB URL with subpath",
			url:       "https://gitlab.com/corneliusweig/rakkess",
			shouldErr: true,
		},
		{
			name:      "any url",
			url:       "https://kudo.org",
			shouldErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actualOwner, actualRepo, err := Parse(test.url)

			if test.shouldErr && err == nil {
				t.Errorf("wanted %s to fail", test.url)
			} else if !test.shouldErr && err != nil {
				t.Errorf("wanted %s not to fail", test.url)
			}
			if actualOwner != test.owner {
				t.Errorf("wanted %s, got %s", test.owner, actualOwner)
			}
			if actualRepo != test.repo {
				t.Errorf("wanted %s, got %s", test.repo, actualRepo)
			}
		})
	}
}
