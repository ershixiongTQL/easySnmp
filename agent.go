package easySnmp

import (
	"fmt"
	"io"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/posteo/go-agentx"
	"github.com/posteo/go-agentx/value"
)

type EnterpriseInfoType struct {
	ID          uint32            //Enterprise ID registered to IANA
	BriefName   string            //Short, letters only name of the organization
	FullName    string            //Fullname of the organization
	ContactInfo map[string]string //Contact informations in key-value data structure
	ProductName string            //Readble name of the product, letters only
}

func (info *EnterpriseInfoType) defaultFill() *EnterpriseInfoType {

	if info.ID == 0 {
		info.ID = 33369
	}

	if info.BriefName == "" {
		info.BriefName = "undefined"
	}

	if info.FullName == "" {
		info.FullName = "undefined"
	}

	if info.ProductName == "" {
		info.ProductName = "undefined"
	}

	return info
}

type SNMP_VAL_TYPE int

const (
	SNMP_VAL_UINT64 SNMP_VAL_TYPE = iota
	SNMP_VAL_STRING
	// SNMP_VAL_FLOAT
)

type StandardTableEntry struct {
	Name string
	Type SNMP_VAL_TYPE
	Desc string
	// ParentName *string
}

type standardTable struct {
	Name       string
	ParentNode *string
	IndexName  string
	Entries    []StandardTableEntry
	getCell    func(col uint, row uint) (string, error)
}

type serviceInfoType struct {
	masterAgent   string
	agentxClient  *agentx.Client
	agentxSession *agentx.Session
}

type MiniSnmpAgent struct {
	EnterpriseInfo EnterpriseInfoType

	//standard structure: Tables
	tables []*standardTable

	//standard structure: Lists
	//TODO

	//standard structure: Controllers
	//TODO

	serviceInfo serviceInfoType

	lock sync.Mutex
}

func (agent *MiniSnmpAgent) AddTable(name string, entries []StandardTableEntry, getCell func(col uint, row uint) (string, error)) {

	agent.lock.Lock()
	defer agent.lock.Unlock()

	for _, tb := range agent.tables {
		if tb.Name == name {
			return
		}
	}

	table := new(standardTable)

	table.Name = name
	table.Entries = entries
	table.getCell = getCell
	table.ParentNode = &agent.EnterpriseInfo.ProductName
	table.IndexName = entries[0].Name

	// for i := range entries {
	// 	entries[i].ParentName = &table.Name
	// }

	agent.tables = append(agent.tables, table)
}

func (agent *MiniSnmpAgent) rootOID() string {
	return fmt.Sprintf("1.3.6.1.4.1.%d.1", agent.EnterpriseInfo.ID)
}

func (agent *MiniSnmpAgent) MibExport(output io.Writer) error {

	agent.lock.Lock()
	defer agent.lock.Unlock()

	output.Write([]byte(strings.ToUpper(agent.EnterpriseInfo.ProductName) + "-MIB DEFINITIONS ::= BEGIN\n"))

	tempMibImports, _ :=
		template.New("mibImports").Funcs(helperFuncs).Parse(templateMibImports)
	tempMibEnterpriseInfo, _ :=
		template.New("mibEnterpriseInfo").Funcs(helperFuncs).Parse(templateMibEnterpriseInfo)
	tempMibProductNode, _ :=
		template.New("mibProductNode").Funcs(helperFuncs).Parse(templateMibProductNode)
	tempMibDataTable, _ :=
		template.New("mibDataTable").Funcs(helperFuncs).Parse(templateMibDataTable)

	enterpriseInfo := agent.EnterpriseInfo.defaultFill()

	tempMibImports.Execute(output, mibImports)
	tempMibEnterpriseInfo.Execute(output, enterpriseInfo)
	tempMibProductNode.Execute(output, enterpriseInfo)

	for i, tb := range agent.tables {
		tempMibDataTable.Execute(output, templateMibDataTableStruct{
			Index: i + 1,
			Table: tb,
		})
	}

	output.Write([]byte("\nEND"))

	return nil
}

func (agent *MiniSnmpAgent) SetMasterAgent(addrPort string) {
	agent.serviceInfo.masterAgent = strings.TrimSpace(addrPort)
}

func (agent *MiniSnmpAgent) start(network string, address string) {

	for agent.serviceInfo.agentxClient == nil {
		agent.serviceInfo.agentxClient, _ = agentx.Dial(network, address)
	}

	agent.serviceInfo.agentxClient.Timeout = 1 * time.Minute
	agent.serviceInfo.agentxClient.ReconnectInterval = 5 * time.Second

	for agent.serviceInfo.agentxSession == nil {
		agent.serviceInfo.agentxSession, _ = agent.serviceInfo.agentxClient.Session()
	}

	agent.serviceInfo.agentxSession.Handler = agent

	agent.serviceInfo.agentxSession.Register(127, value.MustParseOID(agent.rootOID()))
}

func CreateAgent(basicInfo *EnterpriseInfoType, masterAgent string) (agent *MiniSnmpAgent) {

	agent = new(MiniSnmpAgent)

	if basicInfo != nil {
		agent.EnterpriseInfo = *basicInfo
	}

	if masterAgent != "" {
		agent.serviceInfo.masterAgent = masterAgent
		go agent.start("tcp", masterAgent)
	}

	return
}
