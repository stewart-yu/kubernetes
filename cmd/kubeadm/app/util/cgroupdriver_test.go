/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package util

import (
	"testing"

	utilsexec "k8s.io/utils/exec"
	fakeexec "k8s.io/utils/exec/testing"
)

func TestGetCgroupDriverDocker(t *testing.T) {

	prefix := "Cgroup Driver: "
	testCases := []struct {
		name               string
		dockerInfo         string
		expectCgroupDriver string
		expectedError      bool
	}{
		{
			name: "valid: value is 'cgroupfs'",
			dockerInfo: "Logging Driver: json-file\n" +
				"Cgroup Driver: cgroupfs\n" +
				"Plugins:\n" +
				" Volume: local\n" +
				" Network: bridge host macvlan null overlay\n" +
				"Docker Root Dir: /var/lib/docker\n",
			expectCgroupDriver: `Cgroup Driver: cgroupfs`,
			expectedError:      false,
		},
		{
			name: "valid: value is 'systemd'",
			dockerInfo: "Logging Driver: json-file\n" +
				"Cgroup Driver: systemd\n" +
				"Plugins:\n" +
				" Volume: local\n" +
				" Network: bridge host macvlan null overlay\n" +
				"Docker Root Dir: /var/lib/docker\n",
			expectCgroupDriver: `Cgroup Driver: systemd`,
			expectedError:      false,
		},
		{
			name: "invalid: missing 'Cgroup Driver' key and value",
			dockerInfo: "Logging Driver: json-file\n" +
				"Plugins:\n" +
				" Volume: local\n" +
				" Network: bridge host macvlan null overlay\n" +
				"Docker Root Dir: /var/lib/docker\n",
			expectCgroupDriver: "",
			expectedError:      true,
		},
		{
			name: "invalid: only a 'Cgroup Driver' key is present",
			dockerInfo: "Logging Driver: json-file\n" +
				"Cgroup Driver\n" +
				"Plugins:\n" +
				" Volume: local\n" +
				" Network: bridge host macvlan null overlay\n" +
				"Docker Root Dir: /var/lib/docker\n",
			expectCgroupDriver: `Cgroup Driver`,
			expectedError:      true,
		},
		{
			name: "invalid: empty 'Cgroup Driver' value",
			dockerInfo: "Logging Driver: json-file\n" +
				"Cgroup Driver:\n" +
				"Plugins:\n" +
				" Volume: local\n" +
				" Network: bridge host macvlan null overlay\n" +
				"Docker Root Dir: /var/lib/docker\n",
			expectCgroupDriver: `Cgroup Driver: `,
			expectedError:      true,
		},
		{
			name: "invalid: unknown 'Cgroup Driver' value",
			dockerInfo: "Logging Driver: json-file\n" +
				"Cgroup Driver: invalid-value\n" +
				"Plugins:\n" +
				" Volume: local\n" +
				" Network: bridge host macvlan null overlay\n" +
				"Docker Root Dir: /var/lib/docker\n",
			expectCgroupDriver: `Cgroup Driver: invalid-value`,
			expectedError:      true,
		},
	}

	for i, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fcmd := fakeexec.FakeCmd{
				CombinedOutputScript: []fakeexec.FakeCombinedOutputAction{
					func() ([]byte, error) { return []byte(testCases[i].dockerInfo), nil }},
			}

			fexec := fakeexec.FakeExec{
				CommandScript: []fakeexec.FakeCommandAction{
					func(cmd string, args ...string) utilsexec.Cmd { return fakeexec.InitFakeCmd(&fcmd, cmd, args...) },
				},
			}

			cgroupDriver, err := GetCgroupDriverDocker(&fexec)

			if (err != nil) != tc.expectedError {
				t.Fatalf("case[%d]:expected error: %v, saw: %v, error: %v", i, tc.expectedError, (err != nil), err)
			}
			if (err == nil) && (prefix+cgroupDriver != tc.expectCgroupDriver) {
				t.Fatalf("case[%d]:expected cgroupDriver: %v, saw: %v", i, tc.expectCgroupDriver, prefix+cgroupDriver)
			}
		})
	}
}
