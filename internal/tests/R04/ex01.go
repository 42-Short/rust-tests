package R04

import (
	"path/filepath"
	"rust-piscine/internal/alloweditems"

	Exercise "github.com/42-Short/shortinette/pkg/interfaces/exercise"
	"github.com/42-Short/shortinette/pkg/testutils"
)

var clippyTomlAsString01 = `
disallowed-macros = ["std::println"]
`

// cargo run > output.log 2>&1
// grep -i "panicked" output.log

// func testMissingPermissions(exercise Exercise.Exercise, workingDirectory string) Exercise.Result {
// 	if _, err := testutils.RunCommandLine(workingDirectory, "sh", []string{"-c", "mkdir ./foo && chmod 000 ./foo"}); err != nil {
// 		return Exercise.InternalError(err.Error())
// 	}
// 	if _, err := testutils.RunCommandLine(workingDirectory, "sh", []string{"-c", "cargo run -- ./foo/hello < a"}); err != nil {
// 		//
// 	}
// 	return Exercise.Passed("OK")
// }

func testRedirectionOneFile(workingDirectory string) Exercise.Result {
	if _, err := testutils.RunCommandLine(workingDirectory, "sh", []string{"-c", "echo 'Hello, World!' | cargo run -- a"}); err != nil {
		return Exercise.RuntimeError(err.Error())
	}
	if output, _ := testutils.RunCommandLine(workingDirectory, "cat", []string{"a"}); output != "Hello, World!\n" {
		return Exercise.AssertionError("Hello, World!\n", output)
	}
	return Exercise.Passed("OK")
}

func testStdout(workingDirectory string) Exercise.Result {
	output, err := testutils.RunCommandLine(workingDirectory, "sh", []string{"-c", "echo 'Hello, World!' | cargo run"})
	if err != nil {
		return Exercise.RuntimeError(err.Error())
	}
	if output != "Hello, World!\n" {
		return Exercise.AssertionError("Hello, World!", output)
	}
	return Exercise.Passed("OK")
}

func ex01Test(exercise *Exercise.Exercise) (result Exercise.Result) {
	if err := alloweditems.Check(*exercise, clippyTomlAsString01, nil); err != nil {
		return Exercise.CompilationError(err.Error())
	}
	workingDirectory := filepath.Join(exercise.CloneDirectory, exercise.TurnInDirectory)

	if result = testStdout(workingDirectory); !result.Passed {
		return result
	}
	if result = testRedirectionOneFile(workingDirectory); !result.Passed {
		return result
	}
	return Exercise.Passed("OK")
}

func ex01() Exercise.Exercise {
	return Exercise.NewExercise("01", "ex01", []string{"Cargo.toml", "src/main.rs"}, 10, ex01Test)
}
