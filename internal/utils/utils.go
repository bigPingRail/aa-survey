package utils

import (
	"log"
	"path/filepath"
	"strconv"
)

func ToAbsPath(ans interface{}) interface{} {
	p := ans.(string)
	if len(p) == 0 {
		return nil
	}
	abs_path, _ := filepath.Abs(p)
	return abs_path
}

func Contains(s []string, x string) bool {
	for _, v := range s {
		if v == x {
			return true
		}
	}
	return false
}

func ConvertBoolToStr(s string) bool {
	if s != "" {
		parsed, err := strconv.ParseBool(s)
		if err != nil {
			log.Fatalf(
				`error parsing value: %s
Available values is: 
	"1", "t", "T", "true", "TRUE", "True"
	"0", "f", "F", "false", "FALSE", "False"`, s)
		}
		return parsed
	}
	return true

}
