package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type Mod struct {
	Name           string
	Version        string
	RequireModMap  map[string]*Mod
	RequiredModMap map[string]*Mod
}

func main() {
	args := os.Args

	if len(args) < 3 {
		log.Fatalln("args error")
	}

	graphFile := args[1]
	searchMod := args[2]

	lines, err := readline(graphFile)
	if err != nil {
		log.Fatalln(err)
	}

	modMap := map[string]*Mod{}

	getMod := func(modRaw string) *Mod {
		mod, ok := modMap[modRaw]
		if !ok {
			name, version, err := getModNameAndVersion(modRaw)
			if err != nil {
				log.Fatalln(err.Error())
			}

			mod = &Mod{
				Name:           name,
				Version:        version,
				RequireModMap:  map[string]*Mod{},
				RequiredModMap: map[string]*Mod{},
			}

			modMap[modRaw] = mod
		}

		return mod
	}

	for _, line := range lines {
		lineSplit := strings.Split(line, " ")
		if len(lineSplit) != 2 {
			log.Fatalln("error split line")
		}

		modRaw, requireModRaw := lineSplit[0], lineSplit[1]

		mod := getMod(modRaw)
		requireMod := getMod(requireModRaw)

		mod.RequireModMap[requireModRaw] = requireMod
		requireMod.RequiredModMap[modRaw] = mod
	}

	{
		name, version, err := getModNameAndVersion(searchMod)
		if err != nil {
			log.Fatalln(err.Error())
		}

		for _, mod := range modMap {
			if mod.Name != name {
				continue
			}
			if version != "" && mod.Version != version {
				continue
			}

			printMod(mod, 0)
		}
	}
}

func readline(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	lines := make([]string, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	err = scanner.Err()
	if err != nil {
		return nil, err
	}

	return lines, nil
}

func getModNameAndVersion(modRaw string) (string, string, error) {
	var name, version string

	modSplit := strings.Split(modRaw, "@")

	switch len(modSplit) {
	case 1:
		name = modSplit[0]
	case 2:
		name, version = modSplit[0], modSplit[1]
	default:
		return "", "", fmt.Errorf("modRaw error: %s", modRaw)
	}

	return name, version, nil
}

func printMod(mod *Mod, depth int) {
	fmtName := mod.Name
	if mod.Version != "" {
		fmtName = fmt.Sprintf("%s@%s", fmtName, mod.Version)
	}

	fmt.Printf("%s%s\n", strings.Repeat(" ", depth), fmtName)

	for _, subMod := range mod.RequiredModMap {
		printMod(subMod, depth+2)
	}

	// TODO:
	// for _, subMod := range mod.RequireModMap {
	//     printMod(subMod, depth+2)
	// }
}
