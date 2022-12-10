package nanoarch

import (
	"bufio"
	"os"
	"strings"
	"log"
)

import "C"

type ConfigProperties map[string]*C.char

//mupen64plus_next_libretro.cfg 로드하는 것을 기준인가?
func ScanConfigFile(filename string) ConfigProperties {
	log.Println(filename)
	config := ConfigProperties{}

	if len(filename) == 0 {
		return config
	}
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
		return config
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if equal := strings.Index(line, "="); equal >= 0 {
			if key := strings.TrimSpace(line[:equal]); len(key) > 0 {
				value := ""
				if len(line) > equal {
					value = strings.TrimSpace(line[equal+1:])
				}
				config[key] = C.CString(value)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return config
}
