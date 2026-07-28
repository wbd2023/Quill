// packhelper is a test-only external Pack runtime. It reads one JSON request from standard input
// and switches on check_id to exercise the subprocess protocol's success and failure paths. It is
// compiled by the external driver tests and is never part of the Quill build.
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type request struct {
	CheckID        string         `json:"check_id"`
	RepositoryRoot string         `json:"repository_root"`
	Files          []string       `json:"files"`
	Configuration  map[string]any `json:"configuration"`
}

func main() {
	raw, _ := io.ReadAll(os.Stdin)

	var req request
	_ = json.Unmarshal(raw, &req)

	switch req.CheckID {
	case "diagnostic":
		fmt.Println(`{"type":"diagnostic","code":"found","message":"direct database access","file":"internal/service/users.go","start":{"line":42,"column":5},"end":{"line":42,"column":18}}`)
		fmt.Println(`{"type":"complete","success":true}`)

	case "fail-completion":
		fmt.Println(`{"type":"complete","success":false,"error":"configuration field is required"}`)

	case "malformed":
		fmt.Println("this is not json")
		fmt.Println(`{"type":"complete","success":true}`)

	case "no-completion":
		fmt.Println(`{"type":"diagnostic","message":"orphaned finding"}`)

	case "bad-range":
		fmt.Println(`{"type":"diagnostic","message":"bad location","file":"/etc/passwd"}`)
		fmt.Println(`{"type":"complete","success":true}`)

	case "nonzero":
		os.Exit(1)

	case "timeout":
		time.Sleep(30 * time.Second)
		fmt.Println(`{"type":"complete","success":true}`)

	case "truncate":
		for range 200000 {
			fmt.Println(`{"type":"diagnostic","message":"padding-padding-padding-padding-padding"}`)
		}
		fmt.Println(`{"type":"complete","success":true}`)

	case "stderr-debug":
		fmt.Fprintln(os.Stderr, "starting check")
		fmt.Println(`{"type":"complete","success":true}`)

	case "marker":
		_ = os.WriteFile(req.RepositoryRoot+string(os.PathSeparator)+".pack-ran-marker", []byte("ran"), 0o644)
		fmt.Println(`{"type":"complete","success":true}`)

	case "inspect-request":
		encoded, _ := json.Marshal(req.Configuration)
		escaped := strings.ReplaceAll(strings.ReplaceAll(string(encoded), `\`, `\\`), `"`, `'`)
		fileNil := "false"
		if req.Files == nil {
			fileNil = "true"
		}
		fmt.Printf(`{"type":"diagnostic","code":"inspect","message":"config=%s filecount=%d filenil=%s files=%s"}`+"\n",
			escaped, len(req.Files), fileNil, strings.Join(req.Files, ","))
		fmt.Println(`{"type":"complete","success":true}`)
	default:
		fmt.Println(`{"type":"complete","success":true}`)
	}
}
