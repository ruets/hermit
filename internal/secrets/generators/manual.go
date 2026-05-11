package generators

import (
	"fmt"
)

type ManualGenerator struct {
	name string
}

func (g *ManualGenerator) Generate(path string) error {
	fmt.Printf("  Enter value for %s: ", g.name)
	var value string
	fmt.Scanln(&value)
	if value == "" {
		return fmt.Errorf("empty value for %s", g.name)
	}
	return writeSecret(path, value)
}
