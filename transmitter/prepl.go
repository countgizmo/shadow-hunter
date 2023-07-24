package transmitter

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"

	"ziggytwister.com/shadow-hunter/ast"
	"ziggytwister.com/shadow-hunter/parser"
)


func GetAppDB(host string, port string) *ast.EDN {
	address := host + ":" + port
	conn, err := net.Dial("tcp", address)

	if err != nil {
		log.Fatalf("Cannot create TCP connection with %s", address)
	}

	fmt.Fprintf(conn, "@re-frame.db/app-db" + "\n")

	message, _ := bufio.NewReader(conn).ReadString('\n')

	edn := parser.ParseString(message)

	mapElement := edn.Elements[0].(*ast.MapElement)
	var appDBString string 

	for i, key := range mapElement.Keys {
		if key.String() == ":val" {
			appDBString = mapElement.Values[i].String()
			appDBString = strings.ReplaceAll(appDBString, "\\", "")
			break
		}
	}

	return parser.ParseString(appDBString)
}
