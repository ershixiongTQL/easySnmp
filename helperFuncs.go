package easySnmp

import (
	"strings"
	"text/template"
	"time"
)

var helperFuncs = template.FuncMap{
	"intInc":             intInc,
	"lastUpdate":         lastUpdate,
	"notationFormat":     notationFormat,
	"valueTypeToSmiType": valueTypeToSmiType,
}

func intInc(num int) int {
	return num + 1
}

func lastUpdate() string {
	return time.Now().Format("200601021504")
}

func notationFormat(orig string) string {
	return strings.Join(strings.Fields(strings.Title(orig)), "")
}

func valueTypeToSmiType(t SNMP_VAL_TYPE) string {
	switch t {
	case SNMP_VAL_UINT64:
		return "Counter64"
	case SNMP_VAL_STRING:
		return "OCTECT STRING"
	default:
		return ERROR_INSERTION
	}
}
