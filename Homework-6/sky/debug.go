package sky

import (
	"fmt"
	"strings"
)

func debugPrint(info string,value ...interface{}) {
	if isDebugging(){
		if !strings.HasSuffix(info,"\n"){
			info += "\n"
		}
		fmt.Printf("[Sky-Debug] "+info,value...)
	}
}

func debugPrintError(err error) {
	if err != nil {
		fmt.Printf("[Sky-Debug] [Error] %v\n",err)
	}
}
