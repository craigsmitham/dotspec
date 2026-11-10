package integration_test

import (
	"errors"
	"os"
	"os/exec"
	"strconv"
	"testing"

	"github.com/craigsmitham/dotspec/internal/cli"
	"github.com/rogpeppe/go-internal/testscript"
)

func TestMain(m *testing.M) {
	testscript.Main(m, map[string]func(){
		"dotspec": func() { os.Exit(cli.Run(os.Args[1:], os.Stdout, os.Stderr)) },
	})
}

func TestCLI(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir:                 "testdata",
		RequireExplicitExec: true,
		Setup: func(env *testscript.Env) error {
			env.Setenv("GIT_CONFIG_GLOBAL", os.DevNull)
			env.Setenv("GIT_CONFIG_NOSYSTEM", "1")
			return nil
		},
		Cmds: map[string]func(*testscript.TestScript, bool, []string){
			"status": func(ts *testscript.TestScript, neg bool, args []string) {
				if neg || len(args) < 2 {
					ts.Fatalf("usage: status CODE COMMAND [ARG...]")
				}
				want, err := strconv.Atoi(args[0])
				ts.Check(err)
				code := 0
				err = ts.Exec(args[1], args[2:]...)
				if err != nil {
					var exit *exec.ExitError
					if !errors.As(err, &exit) {
						ts.Fatalf("could not execute command: %v", err)
					}
					code = exit.ExitCode()
				}
				if code != want {
					ts.Fatalf("exit status = %d; want %d", code, want)
				}
			},
		},
	})
}
