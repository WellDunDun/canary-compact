package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestContainsCanary(t *testing.T) {
	tests := []struct {
		name          string
		message       string
		canary        string
		caseSensitive bool
		wholeWord     bool
		want          bool
	}{
		{
			name:          "default match",
			message:       "ready CANARY_COMPACT",
			canary:        "CANARY_COMPACT",
			caseSensitive: true,
			wholeWord:     false,
			want:          true,
		},
		{
			name:          "case sensitive miss",
			message:       "ready canary_compact",
			canary:        "CANARY_COMPACT",
			caseSensitive: true,
			wholeWord:     false,
			want:          false,
		},
		{
			name:          "case insensitive match",
			message:       "ready canary_compact",
			canary:        "CANARY_COMPACT",
			caseSensitive: false,
			wholeWord:     false,
			want:          true,
		},
		{
			name:          "whole word miss inside identifier",
			message:       "ready XCANARY_COMPACTY",
			canary:        "CANARY_COMPACT",
			caseSensitive: true,
			wholeWord:     true,
			want:          false,
		},
		{
			name:          "whole word match with punctuation",
			message:       "ready (CANARY_COMPACT).",
			canary:        "CANARY_COMPACT",
			caseSensitive: true,
			wholeWord:     true,
			want:          true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := containsCanary(tt.message, tt.canary, tt.caseSensitive, tt.wholeWord)
			if got != tt.want {
				t.Fatalf("containsCanary() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnvBool(t *testing.T) {
	tests := []struct {
		name         string
		value        string
		set          bool
		defaultValue bool
		want         bool
	}{
		{name: "unset uses default true", defaultValue: true, want: true},
		{name: "unset uses default false", defaultValue: false, want: false},
		{name: "true string", value: "true", set: true, want: true},
		{name: "on string", value: "on", set: true, want: true},
		{name: "false string", value: "false", set: true, defaultValue: true, want: false},
		{name: "invalid uses default", value: "maybe", set: true, defaultValue: true, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			const envName = "CANARY_COMPACT_TEST_BOOL"
			if tt.set {
				t.Setenv(envName, tt.value)
			} else {
				_ = os.Unsetenv(envName)
			}
			if got := envBool(envName, tt.defaultValue); got != tt.want {
				t.Fatalf("envBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadPayload(t *testing.T) {
	tempDir := t.TempDir()
	stdinPath := filepath.Join(tempDir, "stdin.json")
	input := "{\"session_id\":\"s1\",\"last_assistant_message\":\"ready CANARY_COMPACT\",\"stop_hook_active\":true}"
	if err := os.WriteFile(stdinPath, []byte(input), 0o644); err != nil {
		t.Fatal(err)
	}

	file, err := os.Open(stdinPath)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	oldStdin := os.Stdin
	os.Stdin = file
	defer func() { os.Stdin = oldStdin }()

	payload, err := readPayload()
	if err != nil {
		t.Fatal(err)
	}
	if payload.SessionID != "s1" || !payload.StopHookActive {
		t.Fatalf("unexpected payload: %+v", payload)
	}
}

func TestBuildHookOutput(t *testing.T) {
	output := buildHookOutput("CUSTOM_CANARY")
	encoded, err := json.Marshal(output)
	if err != nil {
		t.Fatal(err)
	}
	text := string(encoded)
	if !strings.Contains(text, "CUSTOM_CANARY") {
		t.Fatalf("output missing canary word: %s", text)
	}
	if !strings.Contains(text, "/compact") {
		t.Fatalf("output missing compact command: %s", text)
	}
	if output.HookSpecificOutput.HookEventName != "Stop" {
		t.Fatalf("hook event = %q, want Stop", output.HookSpecificOutput.HookEventName)
	}
}

func TestRecordDetection(t *testing.T) {
	dataDir := t.TempDir()
	t.Setenv("CLAUDE_PLUGIN_DATA", dataDir)

	recordDetection(stopPayload{SessionID: "s1", TranscriptPath: "transcript.jsonl", CWD: "/work"}, "CANARY_COMPACT")

	contents, err := os.ReadFile(filepath.Join(dataDir, "canary-detections.jsonl"))
	if err != nil {
		t.Fatal(err)
	}
	var event detectionEvent
	if err := json.Unmarshal([]byte(strings.TrimSpace(string(contents))), &event); err != nil {
		t.Fatal(err)
	}
	if event.SessionID != "s1" || event.CanaryWord != "CANARY_COMPACT" || event.CWD != "/work" {
		t.Fatalf("unexpected detection event: %+v", event)
	}
}
