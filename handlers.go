package easySnmp

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/posteo/go-agentx/pdu"
	"github.com/posteo/go-agentx/value"
)

func getLocalOidArray(agent *MiniSnmpAgent, oid string) []int {

	localOidsRaw := strings.TrimPrefix(oid, agent.rootOID())

	localOids := []int{}

	for _, seg := range strings.Split(localOidsRaw, ".") {

		if seg == "" {
			continue
		}

		if id, err := strconv.Atoi(seg); err != nil {
			return []int{}
		} else {
			localOids = append(localOids, id)
		}
	}

	return localOids
}

func tableGet(agent *MiniSnmpAgent, objNum uint, tableId uint, subOids []int, getNext bool) (retOid string, retType pdu.VariableType, retValue interface{}, err error) {

	if len(subOids) < 2 || subOids[0] != 1 {
		return "", pdu.VariableTypeNull, nil, nil
	}

	var colNum = subOids[1]

	thisTable := agent.tables[tableId]

	retType = pdu.VariableTypeOctetString

	if getNext {

		if len(subOids) == 2 {

			retOid = agent.rootOID() + "." + strconv.Itoa(int(objNum)) + ".1." + strconv.Itoa(colNum) + ".1"
			retValue, err = thisTable.getCell(thisTable.Name, uint(colNum-1), 0)

		} else if len(subOids) == 3 {

			rowNum := subOids[2]
			retOid = fmt.Sprintf(agent.rootOID()+".%d.1.%d.%d", int(objNum), colNum, rowNum+1)
			retValue, err = thisTable.getCell(thisTable.Name, uint(colNum-1), uint(rowNum-1+1))
		}

	} else {

		if len(subOids) != 3 {
			return
		}

		rowNum := subOids[2]
		retOid = fmt.Sprintf(agent.rootOID()+".%d.1.%d.%d", int(tableId+1), colNum, rowNum)
		retValue, err = thisTable.getCell(thisTable.Name, uint(colNum-1), uint(rowNum-1))
	}

	switch thisTable.Entries[colNum-1].Type {
	case SNMP_VAL_UINT64:
		retValue, err = strconv.Atoi(retValue.(string))
		retValue = uint64(retValue.(int))
		retType = pdu.VariableTypeCounter64
	default:
		break
	}

	return
}

func (agent *MiniSnmpAgent) Get(oid value.OID) (value.OID, pdu.VariableType, interface{}, error) {

	// fmt.Printf("get oid %s\n", oid.String())

	if !strings.HasPrefix(oid.String(), agent.rootOID()) {
		return oid, pdu.VariableTypeNull, nil, nil
	}

	localIds := getLocalOidArray(agent, oid.String())

	if len(localIds) == 0 {
		return nil, pdu.VariableTypeNull, nil, nil
	}

	if len(agent.tables) != 0 && localIds[0] <= len(agent.tables) {
		//type Table

		retOid, retType, retValue, err := tableGet(agent, uint(localIds[0]), uint(localIds[0]-1), localIds[1:], false)

		if retOid == "" || err != nil {
			return nil, pdu.VariableTypeNull, nil, nil
		} else {
			return value.MustParseOID(retOid), retType, retValue, nil
		}

	} else {
		//TODO: other type of objects
	}

	return nil, pdu.VariableTypeNull, nil, nil
}

func (agent *MiniSnmpAgent) GetNext(oid value.OID, what bool, noid value.OID) (value.OID, pdu.VariableType, interface{}, error) {

	// fmt.Printf("get next %s\n", oid.String())

	if !strings.HasPrefix(oid.String(), agent.rootOID()) {
		return noid, pdu.VariableTypeNull, nil, nil
	}

	localIds := getLocalOidArray(agent, oid.String())

	if len(localIds) == 0 {
		return noid, pdu.VariableTypeNull, nil, nil
	}

	if len(agent.tables) != 0 && localIds[0] <= len(agent.tables) {
		//type Table

		retOid, retType, retValue, err := tableGet(agent, uint(localIds[0]), uint(localIds[0]-1), localIds[1:], true)

		if retOid == "" || err != nil {
			return nil, pdu.VariableTypeNull, nil, nil
		} else {
			return value.MustParseOID(retOid), retType, retValue, nil
		}

	} else {
		//TODO: other type of objects
	}

	return nil, pdu.VariableTypeNull, nil, nil
}
