package registry

import (
	"testing"
)

func Test_RegistryKindDetect(t *testing.T) {
	type test struct {
		repo      string
		expect    string
		expectErr bool
	}

	tests := []test{
		{
			repo:      "gcr.io/dcvisor-162009/alcide/dcvisor/cp-kafka",
			expect:    "gcr",
			expectErr: false,
		},
		{
			repo:      "us.gcr.io/dcvisor-162009/alcide/dcvisor/cp-kafka",
			expect:    "gcr",
			expectErr: false,
		},
		{
			repo:      "666.dkr.ecr.us-west-2.amazonaws.com/kaudit:latest",
			expect:    "ecr",
			expectErr: false,
		},
	}

	for _, tst := range tests {
		kind, err := DetectRegistryKind(tst.repo)
		if tst.expectErr && err == nil {
			t.Errorf("Expected error in '%v'", tst.repo)
		}
		if kind != tst.expect {
			t.Errorf("Expected '%v' to have kind '%v' got '%v'", tst.repo, tst.expect, kind)
		}
	}
}
