package main

import (
	"io"
	"os"
	"reflect"
	"testing"
)

func TestReadDir(t *testing.T) {
	tests := []struct {
		name    string
		dir     string
		wantEnv Environment
		wantErr bool
	}{
		{
			name:    "empty dir",
			dir:     "testdata/empty",
			wantEnv: Environment{},
		},
		{
			name: "dir with files",
			dir:  "testdata/files",
			wantEnv: Environment{
				"FOO": EnvValue{Value: "bar"},
				"BAR": EnvValue{Value: "baz"},
			},
		},
		{
			name: "dir with empty file",
			dir:  "testdata/empty_file",
			wantEnv: Environment{
				"FOO": EnvValue{NeedRemove: true},
			},
		},
		{
			name:    "non-existent dir",
			dir:     "non-existent-dir",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env, err := ReadDir(tt.dir)
			if tt.wantErr && err == nil {
				t.Errorf("ReadDir() = nil, want error")
			} else if !tt.wantErr && err != nil {
				t.Errorf("ReadDir() = %v, want nil", err)
			}
			if !tt.wantErr && !reflect.DeepEqual(env, tt.wantEnv) {
				t.Errorf("ReadDir() = %v, want %v", env, tt.wantEnv)
			}
		})
	}
}

func setupTestDir(t *testing.T, dir string) {
	t.Helper()
	err := os.MkdirAll(dir, 0o755)
	if err != nil {
		t.Fatal(err)
	}
}

func teardownTestDir(t *testing.T, dir string) {
	t.Helper()
	err := os.RemoveAll(dir)
	if err != nil {
		t.Fatal(err)
	}
}

func createFile(t *testing.T, path string, content string) {
	t.Helper()
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	_, err = io.WriteString(f, content)
	if err != nil {
		t.Fatal(err)
	}
}
func TestMain(m *testing.M) {
	setupTestDir(nil, "testdata/empty")
	setupTestDir(nil, "testdata/files")
	setupTestDir(nil, "testdata/empty_file")
	createFile(nil, "testdata/files/FOO", "bar")
	createFile(nil, "testdata/files/BAR", "baz")
	createFile(nil, "testdata/empty_file/FOO", "")
	ret := m.Run()
	teardownTestDir(nil, "testdata/empty")
	teardownTestDir(nil, "testdata/files")
	teardownTestDir(nil, "testdata/empty_file")
	os.Exit(ret)
}
