package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const defaultCanaryWord = "CANARY_COMPACT"

type stopPayload struct {
	SessionID            string `json:"session_id"`
	TranscriptPath       string `json:"transcript_path"`
	CWD                  string `json:"cwd"`
	LastAssistantMessage string `json:"last_assistant_message"`
	StopHookActive       bool   `json:"stop_hook_active"`
}

type hookOutput struct {
	HookSpecificOutput stopHookOutput `json:"hookSpecificOutput"`
}

type stopHookOutput struct {
	HookEventName     string `json:"hookEventName"`
	AdditionalContext string `json:"additionalContext"`
}

type detectionEvent struct {
	DetectedAt     string `json:"detected_at"`
	CanaryWord     string `json:"canary_word"`
	SessionID      string `json:"session_id,omitempty"`
	TranscriptPath string `json:"transcript_path,omitempty"`
	CWD            string `json:"cwd,omitempty"`
}

func main() {
	payload, err := readPayload()
	if err != nil || payload.LastAssistantMessage == "" {
		return
	}

	canaryWord := getenvDefault("CANARY_COMPACT_WORD", defaultCanaryWord)
	caseSensitive := envBool("CANARY_COMPACT_CASE_SENSITIVE", true)
	wholeWord := envBool("CANARY_COMPACT_WHOLE_WORD", false)

	if !containsCanary(payload.LastAssistantMessage, canaryWord, caseSensitive, wholeWord) {
		return
	}

	recordDetection(payload, canaryWord)

	if payload.StopHookActive {
		return
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetEscapeHTML(false)
	_ = encoder.Encode(buildHookOutput(canaryWord))
}

func buildHookOutput(canaryWord string) hookOutput {
	return hookOutput{
		HookSpecificOutput: stopHookOutput{
			HookEventName: "Stop",
			AdditionalContext: fmt.Sprintf(
				"The assistant reply contained the canary word `%s`. Treat this as a request to compact the session. Produce a concise handoff summary that preserves current task state, decisions, changed files, tests, blockers, and next steps. Then show the exact `/compact` command the user should submit with that summary as focus instructions. Do not continue unrelated implementation work in this continuation.",
				canaryWord,
			),
		},
	}
}

func readPayload() (stopPayload, error) {
	var payload stopPayload
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		return payload, err
	}
	input = []byte(strings.TrimSpace(string(input)))
	if len(input) == 0 {
		return payload, nil
	}
	err = json.Unmarshal(input, &payload)
	return payload, err
}

func containsCanary(message, canaryWord string, caseSensitive, wholeWord bool) bool {
	if canaryWord == "" {
		return false
	}

	pattern := regexp.QuoteMeta(canaryWord)
	if wholeWord {
		pattern = `(^|[^A-Za-z0-9_])` + pattern + `([^A-Za-z0-9_]|$)`
	}
	if !caseSensitive {
		pattern = `(?i)` + pattern
	}

	matched, err := regexp.MatchString(pattern, message)
	return err == nil && matched
}

func envBool(name string, defaultValue bool) bool {
	value, ok := os.LookupEnv(name)
	if !ok {
		return defaultValue
	}
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	default:
		return defaultValue
	}
}

func getenvDefault(name, defaultValue string) string {
	value := os.Getenv(name)
	if value == "" {
		return defaultValue
	}
	return value
}

func recordDetection(payload stopPayload, canaryWord string) {
	dataDir := os.Getenv("CLAUDE_PLUGIN_DATA")
	if dataDir == "" {
		return
	}
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return
	}

	event := detectionEvent{
		DetectedAt:     time.Now().UTC().Format(time.RFC3339Nano),
		CanaryWord:     canaryWord,
		SessionID:      payload.SessionID,
		TranscriptPath: payload.TranscriptPath,
		CWD:            payload.CWD,
	}
	encoded, err := json.Marshal(event)
	if err != nil {
		return
	}

	path := filepath.Join(dataDir, "canary-detections.jsonl")
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return
	}
	defer file.Close()

	_, _ = file.Write(append(encoded, '\n'))
}
