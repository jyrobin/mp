package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/jyrobin/goutil"
	"github.com/jyrobin/mp"
)

func mt(buf []byte, kind, tags, attrs string) error {
	m, err := mp.ParseMeta(buf)
	if err != nil {
		return err
	}
	if kind != "" && m.Kind() != kind {
		return fmt.Errorf("Got kind %s, expected %s", m.Kind(), kind)
	}
	if tags != "" {
		params := goutil.ParseParams(tags, ",")
		for k, v := range params {
			if !m.HasTag(k) || m.Tag(k) != v {
				return fmt.Errorf("Got tag %s=%s, expected %s=%s", k, m.Tag(k), k, v)
			}
		}
	}

	if attrs != "" {
		params := goutil.ParseParams(attrs, ",")
		for k, v := range params {
			if !m.HasAttr(k) || m.Attr(k) != v {
				return fmt.Errorf("Got attr %s=%s, expected %s=%s", k, m.Attr(k), k, v)
			}
		}
	}

	return nil
}

func main() {
	finFlag := flag.Bool("fin", false, "Last one")
	kindFlag := flag.String("kind", "", "Kind")
	tagsFlag := flag.String("tags", "", "Tags")
	attrsFlag := flag.String("attrs", "", "Attrs")
	flag.Parse()

	buf, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := mt(buf, *kindFlag, *tagsFlag, *attrsFlag); err != nil {
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
