package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"testing"
)

func TestRunCmd(t *testing.T) {
	tests := []struct {
		name    string
		cmd     []string
		env     Environment
		wantErr bool
	}{
		{
			name: "simple command",
			cmd:  []string{"cmd.exe", "/c", "echo", "hello"},
			env:  Environment{},
		},
		{
			name: "command with env",
			cmd:  []string{"cmd.exe", "/c", "echo", "%FOO%"},
			env:  Environment{"FOO": EnvValue{Value: "bar"}},
		},
		{
			name: "command with removed env",
			cmd:  []string{"cmd.exe", "/c", "echo", "%FOO%"},
			env:  Environment{"FOO": EnvValue{NeedRemove: true}},
		},
		{
			name:    "non-existent command",
			cmd:     []string{"non-existent-command"},
			env:     Environment{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			cmd := exec.Command(tt.cmd[0], tt.cmd[1:]...)
			cmd.Stdout = &buf
			cmd.Env = append(cmd.Env, "PATH=C:\\Windows\\System32;C:\\Windows")
			for k, v := range tt.env {
				if v.NeedRemove {
					cmd.Env = append(cmd.Env, fmt.Sprintf("%s=", k))
				} else {
					cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v.Value))
				}
			}
			err := cmd.Run()
			if tt.wantErr && err == nil {
				t.Errorf("RunCmd() = nil, want error")
			} else if !tt.wantErr && err != nil {
				t.Errorf("RunCmd() = %v, want nil", err)
			}
			if buf.String() != "" {
				t.Logf("output: %s", buf.String())
			}
		})
	}
}
