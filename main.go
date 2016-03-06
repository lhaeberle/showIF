// showIF project main.go

// Copyleft by Ludwig Haeberle 2016
/*
ShowIF is a command line utility to collect quickly basic information from your network interfaces.
It is useful if your server hosts a lot of interfaces. It supports flag parameters to select specific interfaces.
Logical flag compositions are supported.
*/

package main

import (
	"flag"
	"fmt"
	"net"
	"strconv"
	"strings"
)

func PrintExamples() {
	fmt.Println("\nshow all interfaces:" +
		"\nshowIF" +
		"\nshowIF -index 0" +
		"\n\nshow interface with index #5:" +
		"\nshowIF -index 5" +
		"\n\nshow only p2p interfaces which ip address contains string 192.168:" +
		"\nshowIF -IFflag pointtopoint -ip 192.168 -op and" +
		"\nshowIF -IFflag pointtopoint -ip 192.168" +
		"\n\nshow both all p2p interfaces and all interfaces with ip address string 10.101.101:" +
		"\nshowIF -IFflag pointtopoint -ip 10.101.101 -op or" +
		"\n\nshow interfaces with ip address string 10.101 and interface name contains vlan_:" +
		"\nshowIF -name vlan_ -ip 10.101" +
		"\n\nshow interfaces with ip address string 10.101 and interface name contains vlan_ and mtu > 1500:" +
		"\nshowIF -name vlan_ -ip 10.101 -mtu gt1500" +
		"\n\ngot it?" +
		"\nATTENTION: Flag - index >0 supports no composition with other flags. It overrules other parameters.")
	return
}

func SelectInterfaces(i int, n, f, maca, ipa, mtu, o string) { // selects interfaces according to flag intuts

	ifs, err := net.Interfaces()

	if err != nil { // check interface errors
		fmt.Println(err)
	}

	var isOR bool
	if !((strings.ToUpper(o) != "OR") || (strings.ToUpper(o) != "AND")) {

		fmt.Println("No valid input for parameter op. Usage: showIF -op or. Default is AND. Try -h for help...\n")
		return
	} else {
		if strings.ToUpper(o) == "OR" {
			isOR = true
		}
	}
	if i != 0 { // Index not 0 so we print the selected Index only

		if (i < 0) || (i > len(ifs)) { // Error check, Index must not be negativ & has to exist
			fmt.Println("\nNo such index... please try showIF -h for help.\n")
			return
		}

		fmt.Println(strings.Repeat("-", 55))

		for _, inf := range ifs {
			if inf.Index == i {
				PrintInterfaces(inf)
			}
		}

		fmt.Println()
		return
	} else { // Index is 0 (also default value for Index) --> shows all Interfaces

		if (n != "") || (f != "") || (maca != "") || (ipa != "") || (mtu != "") { // parameter (other than Intex) entered?

			var mtuOperator string
			var mtuI int

			if mtu != "" {
				if !(strings.HasPrefix(mtu, "lt") || strings.HasPrefix(mtu, "eq") || strings.HasPrefix(mtu, "gt")) {
					fmt.Println("Please use prefix lt | eq | gt  followed by an integer... -mtu eq1500. Please try shoIF -h for help\n")
					return
				}

				if strings.HasPrefix(mtu, "lt") {
					mtuOperator = "lt"
				}
				if strings.HasPrefix(mtu, "gt") {
					mtuOperator = "gt"
				}
				if strings.HasPrefix(mtu, "eq") {
					mtuOperator = "eq"
				}

				r := strings.NewReplacer("lt", "", "gt", "", "eq", "")
				mtu = r.Replace(mtu)
				mtuI, err = strconv.Atoi(mtu)
				if err != nil {
					fmt.Println("Please enter an integer (up to 6 digits) right after prefix... -mtu eq1500. Please try -h for help.\n")
					return
				}
			}

			IFNames := []string{}
			IPString := ""
			for _, inf := range ifs {
				z, _ := inf.Addrs()
				if z != nil {
					IPString = ""
					i := 0
					for _, zz := range z { // !! unfortunately neccessary to create []string type from []net.addr type - didn't find another way...
						if zz != nil { // it does nothing; just to avoid compiler message "zz declared and not used"; grrrrr desperate bid...
							fmt.Print("") // let me know if you know a more elegant solution
						}
						zs := z[i].String()
						IPString = IPString + " " + zs
						i = i + 1
					}

				}
				if isOR == false {
					switch { // needed for MTU cases
					case mtuOperator == "eq":
						if strings.Contains(inf.Name, n) && strings.Contains(inf.Flags.String(), f) && strings.Contains(inf.HardwareAddr.String(), maca) && strings.Contains(IPString, ipa) && (inf.MTU == mtuI) {
							PrintInterfaces(inf)
							IFNames = append(IFNames, inf.Name) // used for func FoundInterfaces
							IPString = ""                       // needed to clean IPString
						}

					case mtuOperator == "gt":
						if strings.Contains(inf.Name, n) && strings.Contains(inf.Flags.String(), f) && strings.Contains(inf.HardwareAddr.String(), maca) && strings.Contains(IPString, ipa) && (inf.MTU > mtuI) {
							PrintInterfaces(inf)
							IFNames = append(IFNames, inf.Name) // used for func FoundInterfaces
							IPString = ""                       // needed to clean IPString
						}

					case mtuOperator == "lt":
						if strings.Contains(inf.Name, n) && strings.Contains(inf.Flags.String(), f) && strings.Contains(inf.HardwareAddr.String(), maca) && strings.Contains(IPString, ipa) && (inf.MTU < mtuI) {
							PrintInterfaces(inf)
							IFNames = append(IFNames, inf.Name) // used for func FoundInterfaces
							IPString = ""                       // needed to clean IPString
						}
					default: // no MTU input
						if strings.Contains(inf.Name, n) && strings.Contains(inf.Flags.String(), f) && strings.Contains(inf.HardwareAddr.String(), maca) && strings.Contains(IPString, ipa) {
							PrintInterfaces(inf)
							IFNames = append(IFNames, inf.Name) // used for func FoundInterfaces
							IPString = ""                       // needed
						}
					}
				} else { // ifOR true (not false :-) ) --> Parameter -or set to OR

					const noEmptyString = "8bits.or.8bids?64!t&d" // need the following to avoid showing intefaces with empty fields...
					if n == "" {                                  // if the above const hits your interface name(s) change one of them :-)
						n = noEmptyString
					}
					if f == "" {
						f = noEmptyString
					}
					if maca == "" {
						maca = noEmptyString
					}
					if ipa == "" {
						ipa = noEmptyString
					}

					switch { // needed for MTU cases
					case mtuOperator == "eq":
						if strings.Contains(inf.Name, n) || strings.Contains(inf.Flags.String(), f) || strings.Contains(inf.HardwareAddr.String(), maca) || strings.Contains(IPString, ipa) || (inf.MTU == mtuI) {
							PrintInterfaces(inf)
							IFNames = append(IFNames, inf.Name) // used for func FoundInterfaces
							IPString = ""                       // needed to clean IPString
						}

					case mtuOperator == "gt":
						if strings.Contains(inf.Name, n) || strings.Contains(inf.Flags.String(), f) || strings.Contains(inf.HardwareAddr.String(), maca) || strings.Contains(IPString, ipa) || (inf.MTU > mtuI) {
							PrintInterfaces(inf)
							IFNames = append(IFNames, inf.Name) // used for func FoundInterfaces
							IPString = ""                       // needed to clean IPString
						}

					case mtuOperator == "lt":
						if strings.Contains(inf.Name, n) || strings.Contains(inf.Flags.String(), f) || strings.Contains(inf.HardwareAddr.String(), maca) || strings.Contains(IPString, ipa) || (inf.MTU < mtuI) {
							PrintInterfaces(inf)
							IFNames = append(IFNames, inf.Name) // used for func FoundInterfaces
							IPString = ""                       // needed to clean IPString
						}
					default: // no MTU input
						if strings.Contains(inf.Name, n) || strings.Contains(inf.Flags.String(), f) || strings.Contains(inf.HardwareAddr.String(), maca) || strings.Contains(IPString, ipa) {
							PrintInterfaces(inf)
							IFNames = append(IFNames, inf.Name) // used for func FoundInterfaces
							IPString = ""                       // needed to clean IPString
						}
					}

				}
			}
			FoundInterfaces(IFNames)

			fmt.Println(n)
			return
		}

		fmt.Println(strings.Repeat("-", 55))
		IFNames := []string{}
		for _, inf := range ifs {

			PrintInterfaces(inf)
			IFNames = append(IFNames, inf.Name)
		}

		FoundInterfaces(IFNames)
	}

	fmt.Println()
	return
}

func PrintInterfaces(inf net.Interface) { // prints the selected stuff

	x, _ := inf.Addrs()
	y, _ := inf.MulticastAddrs()

	fmt.Print("Index		: ", inf.Index, "\n")
	fmt.Print("Name		: ", inf.Name, "\n")
	fmt.Print("MTU		: ", inf.MTU, "\n")
	fmt.Print("MAC		: ", inf.HardwareAddr, "\n")
	fmt.Print("IF Flags	: ", inf.Flags, "\n")
	fmt.Print("IP Addr		: ", x, "\n")
	fmt.Print("Multicast	: ", y, "\n")
	fmt.Println(strings.Repeat("-", 55))

	return
}

func FoundInterfaces(n []string) { // tells number of found interfaces & prints name of found interaces

	fmt.Print("\nNunber of Interfaces found: ", len(n), " --> ", n, "\n")

}

func main() {

	// Flag definition

	flagNAME := flag.String("name", "", "String value used to look-up interface names matches. Substrings are supported. E.g. showIF -name vlan_; showIF -name eth")
	flagINDEX := flag.Int("index", 0, "Values reflect interface index. Use 0 for all (default). ATTENTION: Flag -index >0 supports no composition with other flags. It overrules any other flags.")
	flagIFFLAG := flag.String("IFflag", "", "Values: up | broadcast | multicast | loopback |pointtopoint. Substrings are supported. E.g. showIF -IFflag up")
	flagMAC := flag.String("mac", "", "A string value is used to look-up hardware address (MAC) matches. Substrings are suppotred. E.g. showIF -mac 3c:07")
	flagIP := flag.String("ip", "", "A string value is used to look-up IP address matches. Substrings are suppotred. E.g. showIF -ip 192.168")
	flagMTU := flag.String("mtu", "", "usage: showIF -mtu <prefix><integer>; prefix values lt, gt, eq; e.g. -mtu gt1500, -mtu eq1500, -mtu lt1200")
	flagOPERATOR := flag.String("op", "AND", " usage: showIF -op or; -op and; defines logical parameter compositions. Default is AND.")
	flagHELP := flag.String("help", "", " usage: showIF -help ex; provides a couple of exampels")

	flag.Parse()

	fmt.Println("\n\nshowIF version v0.203. Copyleft by Ludwig Haeberle 2016")

	if len(flag.Args()) > 0 { // i support one parameter per flag (but multiple flags, of course)
		fmt.Println("\n You can use one parameter per flag. Multiple flags are supported. Please use showIF -h for help.\n")
		return
	}

	if strings.HasPrefix(*flagHELP, "ex") {
		PrintExamples()
		return
	}

	SelectInterfaces(*flagINDEX, *flagNAME, *flagIFFLAG, *flagMAC, *flagIP, *flagMTU, *flagOPERATOR) //do the job

	if *flagINDEX > 0 { // print input parameter & tell -index >0 overrules...
		fmt.Println("You requested: showIF -index", *flagINDEX, "\n")
		fmt.Println("Be aware: flag -index >0 shows always one interface. No flag composition with -index >0 supported. Try -h for help\n")
	} else { // print input parameters
		fmt.Println("You requested: showIF -index", *flagINDEX, "-name", *flagNAME, "-IFflag", *flagIFFLAG, "-mac", *flagMAC, "-ip", *flagIP, "-mtu", *flagMTU, "-op", *flagOPERATOR, ";(-index 0 and -op AND are defaults.)\n")
	}

}
