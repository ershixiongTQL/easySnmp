package easySnmp

const ERROR_INSERTION = "!!!ERR!!!"

var mibImports = map[string][]string{

	"SNMPv2-SMI": {
		"MODULE-IDENTITY",
		"OBJECT-TYPE",
		"NOTIFICATION-TYPE",
		"Unsigned32",
		"Counter32",
		"Counter64",
		"OCTET STRING",
		"TimeTicks",
		"Gauge32",
		"Integer32",
	},

	"SNMPv2-CONF": {
		"MODULE-COMPLIANCE",
		"NOTIFICATION-GROUP",
		"OBJECT-GROUP",
	},

	"SNMPv2-TC": {
		"DisplayString",
		"TruthValue",
		"MacAddress",
		"TimeInterval",
		"RowStatus",
		"TimeStamp",
		"DateAndTime",
	},

	"SNMP-FRAMEWORK-MIB": {
		"SnmpAdminString",
	},

	"ENTITY-MIB": {
		"PhysicalIndex",
	},

	"INET-ADDRESS-MIB": {
		"InetAddressType",
		"InetAddress",
		"InetPortNumber",
	},

	"IANAifType-MIB": {
		"IANAifType",
	},
}

const (
	templateMibImports = `
IMPORTS
{{- range $name, $item := .}}
{{- range $item}}{{- printf "\t%s,\n" .}}{{- end}}
{{- printf "\t\t%s %s\n" "FROM" $name}}
{{- end}};
`

	templateMibEnterpriseInfo = `
{{.BriefName}} MODULE-IDENTITY
LAST-UPDATED "{{lastUpdate}}Z"
ORGANIZATION "{{.FullName}}"
CONTACT-INFO
	"
	{{- range $name, $item := .ContactInfo}}
	{{- printf "\n\t%s: %s" $name $item}}
	{{- end}}
	"
DESCRIPTION
	"None"
REVISION
	"202202090000Z"
::= { enterprises {{.ID}} }
`

	templateMibProductNode = `
{{.ProductName}} OBJECT IDENTIFIER ::= { {{.BriefName}} 1 }	
`

	templateMibDataTable = `
{{with $tb := .Table}}
{{notationFormat .Name}} OBJECT-TYPE
SYNTAX SEQUENCE OF syntax{{notationFormat .Name}}Entry
MAX-ACCESS not-accessible
STATUS current
::= { {{.ParentNode}} {{$.Index}} }

{{notationFormat .Name}}Entry OBJECT-TYPE
	SYNTAX syntax{{notationFormat .Name}}Entry
	MAX-ACCESS not-accessible
	STATUS current
	INDEX {{notationFormat .IndexName | printf "{%s}"}}
	::= {{notationFormat .Name | printf "{%s 1}"}}

syntax{{notationFormat .Name}}Entry ::= SEQUENCE {
	{{- range $item := .Entries}}
	{{- printf "\n\t%s %s," (notationFormat .Name) (valueTypeToSmiType .Type)}}
	{{- end}}
}
{{range $id, $item := .Entries}}
{{notationFormat $item.Name}} OBJECT-TYPE
	SYNTAX      {{valueTypeToSmiType $item.Type}}
	MAX-ACCESS  read-only
	STATUS      current
	DESCRIPTION "{{.Desc}}"
		::= {{printf "{ %sEntry %d }" (notationFormat $tb.Name) (intInc $id)}}
{{end}}
{{end}}
`
)

type templateMibDataTableStruct struct {
	Index int
	Table *standardTable
}
