package cli

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/alecthomas/kong"

	"github.com/wbd2023/quill/internal/engine"
	"github.com/wbd2023/quill/internal/report"
)

/* -------------------------------------------- Types ------------------------------------------- */

// errHelpRequested is the private sentinel returned by the help flag hook and wrapped by Kong in a
// ParseError. Runner.Run detects it via errors.Is.
var errHelpRequested = errors.New("help requested")

// command is the one behavior every top-level Quill command implements after Kong selects and
// populates its command struct.
type command interface {
	run(context.Context, Runner) (exitCode int)
}

// repoFlags is the grammar shared by commands that operate on a repository and support text or
// JSON output. It is not a command base class: it contains only parsed inputs and creates the
// operation's engine after parsing succeeds.
type repoFlags struct {
	Root string `kong:"name=repository-root,help=repository root (auto-detected when omitted)"`

	Format report.OutputFormat `kong:"default=text,enum='text,json',help=format: text|json"`
}

func (flags repoFlags) newEngine(options ...engine.Option) (application *engine.Engine, err error) {
	root, err := resolveRepositoryRoot(flags.Root)
	if err != nil {
		return nil, err
	}

	return engine.New(root, options...)
}

// ruleRef is the `rule:<id>` positional grammar. It validates syntax during Kong parsing so an
// invalid subject receives the same deterministic usage behavior as every other malformed value.
type ruleRef string

func (reference *ruleRef) Decode(context *kong.DecodeContext) (err error) {
	var subject string
	if err := context.Scan.PopValueInto("subject", &subject); err != nil {
		return err
	}

	kind, id, found := strings.Cut(subject, ":")
	if !found {
		return fmt.Errorf("invalid subject %q: expected rule:<id>", subject)
	}

	if kind != "rule" {
		return fmt.Errorf("unsupported subject %q: only rule:<id> is supported", subject)
	}

	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("invalid subject %q: rule id must not be empty", subject)
	}

	*reference = ruleRef(id)
	return nil
}

func (reference ruleRef) ruleID() (id string) {
	return string(reference)
}

/* ------------------------------------------ Help Flag ----------------------------------------- */

// quillHelpFlag is the sole help surface. Kong's default --help/-h are disabled via NoDefaultHelp;
// this flag's BeforeReset hook returns a private sentinel whenever --help or -h is present, so
// Runner.Run can render deterministic help to stdout and return 0 before preparation or execution.
// The hook only fires when the flag is actually parsed: Kong skips BeforeReset for unset flags
// that carry no default, so the sentinel never surfaces on ordinary invocations.
type quillHelpFlag bool

func (quillHelpFlag) BeforeReset(*kong.Context) (err error) {
	return errHelpRequested
}

/* ---------------------------------------- Command Model --------------------------------------- */

// commandLine is Quill's single Kong model. Every supported command is declared here; help is a
// root flag inherited by all commands. The declaration order is the deterministic order in which
// commands appear in root usage.
type commandLine struct {
	Help     quillHelpFlag `kong:"name=help,short=h,help=show help"`
	Check    checkCmd      `kong:"cmd,help=run STYLE.md checks"`
	Fix      fixCmd        `kong:"cmd,help=run safe style auto-fixes"`
	Doctor   doctorCmd     `kong:"cmd,help=check pinned style tools"`
	Coverage coverageCmd   `kong:"cmd,help=show STYLE.md automation coverage"`
	Install  installCmd    `kong:"cmd,help=install pinned style tools"`
	Lock     lockCmd       `kong:"cmd,help=resolve archive-tool hashes to quill.lock"`
	Version  versionCmd    `kong:"cmd,help=print the Quill version"`
	Init     initCmd       `kong:"cmd,help=create a minimal STYLE.md and quill.toml"`
	List     listCmd       `kong:"cmd,help='list packs, rules, tools, or scopes'"`
	Explain  explainCmd    `kong:"cmd,help=explain an active rule"`
}

// lookup returns the command implementation backing the parsed node, or nil for the root.
func (cl *commandLine) lookup(node *kong.Node) (cmd command) {
	if node == nil {
		return nil
	}
	switch node.Name {
	case "check":
		return &cl.Check

	case "fix":
		return &cl.Fix

	case "doctor":
		return &cl.Doctor

	case "coverage":
		return &cl.Coverage

	case "install":
		return &cl.Install

	case "lock":
		return &cl.Lock

	case "version":
		return &cl.Version

	case "init":
		return &cl.Init

	case "list":
		return &cl.List

	case "explain":
		return &cl.Explain
	}
	return nil
}

/* ---------------------------------------- Construction ---------------------------------------- */

// newParser builds a Kong parser over a fresh command model. NoDefaultHelp suppresses Kong's
// built-in help flag (and its os.Exit behaviour); the commandLine help flag is the only help path.
func newParser() (parser *kong.Kong, model *commandLine, err error) {
	model = &commandLine{}
	parser, err = kong.New(model, kong.NoDefaultHelp(), kong.Name("quill"))
	return parser, model, err
}
