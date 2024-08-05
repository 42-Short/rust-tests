package R00

import (
	"github.com/42-Short/shortinette/pkg/logger"

	Module "github.com/42-Short/shortinette/pkg/interfaces/module"

	Exercise "github.com/42-Short/shortinette/pkg/interfaces/exercise"
)

func R00() *Module.Module {
	exercises := map[string]Exercise.Exercise{
		"00": ex00(),
		"01": ex01(),
		"02": ex02(),
		"03": ex03(),
		"04": ex04(),
		"05": ex05(),
		"06": ex06(),
		"07": ex07(),
	}
	r00, err := Module.NewModule("00", 70, exercises)
	if err != nil {
		logger.Error.Printf("internal error: %v", err)
		return nil
	}
	return &r00
}