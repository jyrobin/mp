package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/jyrobin/goutil"
)

func mj(buf []byte, attrs, ints, strVal, intVal string) error {
	if strVal != "" {
		var v string
		if err := json.Unmarshal(buf, &v); err != nil {
			return err
		}
		if v != strVal {
			return fmt.Errorf("Got %v, expected string %s", v, strVal)
		}
		return nil
	}

	if intVal != "" {
		var v int
		if err := json.Unmarshal(buf, &v); err != nil {
			return err
		}
		if strconv.Itoa(v) != intVal {
			return fmt.Errorf("Got %v, expected int %s", v, intVal)
		}
		return nil
	}

	var val map[string]interface{}
	if err := json.Unmarshal(buf, &val); err != nil {
		return err
	}

	if attrs != "" {
		params := goutil.ParseParams(attrs, ",")
		for k, v := range params {
			if val[k] != v {
				return fmt.Errorf("Got attr %s=%v, expected %s=%v", k, val[k], k, v)
			}
		}
	}

	if ints != "" {
		params := goutil.ParseParams(ints, ",")
		for k, v := range params {
			if fmt.Sprintf("%v", v) != val[k] {
				return fmt.Errorf("Got int %s=%v, expected %s=%s", k, v, k, val[k])
			}
		}
	}

	return nil
}

func main() {
	finFlag := flag.Bool("fin", false, "Last one")
	attrsFlag := flag.String("attrs", "", "String attrs")
	intsFlag := flag.String("ints", "", "Int attrs")
	intFlag := flag.String("int", "", "Int value")
	strFlag := flag.String("str", "", "String value")
	flag.Parse()

	buf, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := mj(buf, *attrsFlag, *intsFlag, *strFlag, *intFlag); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if !*finFlag {
		if _, err := os.Stdout.Write(buf); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
