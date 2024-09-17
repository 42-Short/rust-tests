package R03

import (
	"path/filepath"

	"github.com/42-Short/shortinette/pkg/logger"

	Exercise "github.com/42-Short/shortinette/pkg/interfaces/exercise"
	"github.com/42-Short/shortinette/pkg/testutils"
)

var Ex01TestMod = `
[cfg(test)]
mod shortinette_rust_test_module03_ex01_0001 {
    use super::*;

    #[test]
    fn it_works() {
        assert_eq!(min(12i32, -14i32), -14);
        assert_eq!(min(12f32, 14f32), 12f32);
        assert_eq!(min("abc", "def"), "abc");
        assert_eq!(min(String::from("abc"), String::from("def")), "abc");
        assert_eq!(min(0, 0), 0);
    }
}
`

func ex01Test(exercise *Exercise.Exercise) Exercise.Result {
	workingDirectory := filepath.Join(exercise.CloneDirectory, exercise.TurnInDirectory)

	if err := testutils.AppendStringToFile(Ex01TestMod, exercise.TurnInFiles[1]); err != nil {
		logger.Exercise.Printf("internal error: %v", err)
		return Exercise.InternalError(err.Error())
	}

	output, err := testutils.RunCommandLine(workingDirectory, "cargo", []string{"test", "--release", "shortinette_rust_test_module03_ex01_0001"})
	if err != nil {
		return Exercise.AssertionError("", output)
	}
	return Exercise.Passed("OK")
}

func ex01() Exercise.Exercise {
	return Exercise.NewExercise("01", "ex01", []string{"src/lib.rs", "Cargo.toml"}, 25, ex01Test)
}