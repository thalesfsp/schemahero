package mysql

import (
	"fmt"
	"regexp"
	"strings"
)

func unaliasUnparameterizedColumnType(requestedType string) string {
	switch requestedType {
	case "bool":
		return "tinyint (1)"
	case "boolean":
		return "tinyint (1)"
	}

	for _, unparameterizedColumnType := range unparameterizedColumnTypes {
		if unparameterizedColumnType == requestedType {
			return requestedType
		}
	}

	return ""
}

func unaliasParameterizedColumnType(requestedType string) string {
	if strings.HasPrefix(requestedType, "char ") || strings.HasPrefix(requestedType, "char(") ||
		requestedType == "character" || requestedType == "char" {
		r := regexp.MustCompile(`char\s*\((?P<len>\d*)\)`)

		matchGroups := r.FindStringSubmatch(requestedType)
		if len(matchGroups) == 0 {
			return "character"
		}

		return fmt.Sprintf("character (%s)", matchGroups[1])
	}
	if strings.HasPrefix(requestedType, "integer") {
		r := regexp.MustCompile(`integer\s*\(\s*(?P<max>\d*)\s*\)`)

		matchGroups := r.FindStringSubmatch(requestedType)
		if len(matchGroups) == 0 {
			return "int (11)"
		}

		return fmt.Sprintf("int (%s)", matchGroups[1])
	}
	if strings.HasPrefix(requestedType, "dec") {
		precisionAndScale := regexp.MustCompile(`dec\s*\(\s*(?P<precision>\d*),\s*(?P<scale>\d*)\s*\)`)
		precisionOnly := regexp.MustCompile(`dec\s*\(\s*(?P<precision>\d*)\s*\)`)

		precisionAndScaleMatchGroups := precisionAndScale.FindStringSubmatch(requestedType)
		precisionOnlyMatchGroups := precisionOnly.FindStringSubmatch(requestedType)

		if len(precisionAndScaleMatchGroups) == 3 {
			return fmt.Sprintf("decimal (%s, %s)", precisionAndScaleMatchGroups[1], precisionAndScaleMatchGroups[2])
		} else if len(precisionOnlyMatchGroups) == 2 {
			return fmt.Sprintf("decimal (%s, 0)", precisionOnlyMatchGroups[1])
		}

		return "decimal (10, 0)"
	}
	if strings.HasPrefix(requestedType, "double precision") {
		precisionAndScale := regexp.MustCompile(`double precision\s*\(\s*(?P<precision>\d*),\s*(?P<scale>\d*)\s*\)`)
		precisionAndScaleMatchGroups := precisionAndScale.FindStringSubmatch(requestedType)

		fmt.Printf("%#v\n", precisionAndScaleMatchGroups)
		if len(precisionAndScaleMatchGroups) == 3 {
			return fmt.Sprintf("double (%s, %s)", precisionAndScaleMatchGroups[1], precisionAndScaleMatchGroups[2])
		}

		return "double"
	}

	return ""
}