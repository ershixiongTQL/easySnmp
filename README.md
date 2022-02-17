# easySnmp

This go module is working on to export some kinds of data inside any application to standard SNMP interface smoothly. It works alone with AgentX protocol by acting as a sub-agent, which also means an master agent will be required within the final application(system).

The major concept of this module is to hide any complexity of the SNMP protocol and makes it easy enough for an application to support SNMP feature(at some level), by using some pre-defined common data structures and few custom callbacks.

Also, a good news is with this module, no need to write MIB files by hand, but just simply use the MibExport() function to generate the whole MIB file for use. (this feature perfectly suitable for "no given MIB" project)


# background

This module initial for embedded network system development.
This module is under a extreamly unstable status, any imports must after carefully evaluations.

# example
None

