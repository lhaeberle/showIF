ShowIF readme file
===================
Version v0.203. Copyleft by Ludwig Haeberle 2016

Files:
readme.txt 	->this file
main.go		->code
showIF		->Mach-O 64-bit executable x86_64 (Apple OSX)
showIF.bin	->ELF 64-bit LSB executable, x86-64, version 1 (SYSV), statically linked, not stripped
showIF.exe	->PE32+ executable for MS Windows (console) Mono/.Net assembly
*.info		->file info 		

ShowIF is a command line utility to collect quickly basic information from your network interfaces. It is useful if your server hosts a lot of interfaces. it supports flag parameters to select specific interfaces. Logical flag compositions are supported.

Some usage examples ------------------------------------------------------------------------------
show all interfaces:
shoIF
showIF -index 0

show interface with index #5
showIF -index 5

show only p2p interfaces which ip address contains string 192.168
showIF -IFflag pointtopoint -ip 192.168 -op and
showIF -IFflag pointtopoint -ip 192.168

ahow all p2p interfaces and all interfaces with ip address string 10.101
showIF -IFflag pointtopoint -ip 10.101 -op or

interfaces with both ip address string 10.101 and interface name contains vlan_
showIF -name vlan_ -ip 10.101 

got it?

ATTENTION: Flag - index >0 supports no composition with other flags. It overrules other parameters.
---------------------------------------------------------------------------------------------------

Usage of ./showIF:
  -IFflag string
    	Values: up | broadcast | multicast | loopback |pointtopoint. Substrings are supported. E.g. showIF -IFflag up
  -help string
    	 usage: showIF -help ex; provides a couple of exampels
  -index int
    	Values reflect interface index. Use 0 for all (default). ATTENTION: Flag -index >0 supports no composition with other flags. It overrules any other flags.
  -ip string
    	A string value is used to look-up IP address matches. Substrings are suppotred. E.g. showIF -ip 192.168
  -mac string
    	A string value is used to look-up hardware address (MAC) matches. Substrings are suppotred. E.g. showIF -mac 3c:07
  -mtu string
    	usage: showIF -mtu <prefix><integer>; prefix values lt, gt, eq; e.g. -mtu gt1500, -mtu eq1500, -mtu lt1200
  -name string
    	String value used to look-up interface names matches. Substrings are supported. E.g. showIF -name vlan_; showIF -name eth
  -op string
    	 usage: showIF -op or; -op and; defines logical parameter compositions. Default is AND. (default "AND")

----------------------------------------------------------------------------------------------------
Known limitations / bugs:

Windows Server 2012: Flag -index ignored for args > 0

