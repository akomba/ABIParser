package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type ioRecord struct {
	Indexed bool
	Name    string
	Type    string
}
type abiRec struct {
	Constant        bool
	Inputs          []ioRecord
	Name            string
	Outputs         []ioRecord
	Payable         bool
	StateMutability string
	RecType         string `json:"type"`
}

func (ab *abiRec) inParams() (s string) {
	comma := ""
	for _, rec := range ab.Inputs {
		s += comma
		s += rec.Type
		s += " "
		s += rec.Name
		comma = ","
	}
	return
}

func (ab *abiRec) evParams() (s string) {
	comma := ""
	for _, rec := range ab.Inputs {
		s += comma
		s += rec.Type
		s += " "
		if rec.Indexed {
			s += "indexed"
			s += " "
		}
		s += rec.Name
		comma = ","
	}
	return
}

func (ab *abiRec) outParams() (s string) {
	if len(ab.Outputs) == 0 {
		return
	}
	s = " returns ("
	comma := ""
	for _, rec := range ab.Outputs {
		s += comma
		s += rec.Type
		if len(rec.Name) > 0 {
			s += " "
			s += rec.Name
		}
		comma = ","
	}
	s += ")"
	return
}

func (ab *abiRec) payable() (s string) {
	if strings.Compare(ab.StateMutability, "payable") == 0 {
		s = " payable"
	}
	return
}

func (ab *abiRec) string() string {
	switch ab.RecType {
	case "fallback":
		return fmt.Sprintf("func(%s) public%s%s;\n", ab.inParams(), ab.payable(), ab.outParams())
	case "event":
		return fmt.Sprintf("event %s(%s);\n", ab.Name, ab.evParams())
	case "constructor":
		return fmt.Sprintln("constructor() public;")
	case "function":
		return fmt.Sprintf("func %s(%s) public%s%s;", ab.Name, ab.inParams(), ab.payable(), ab.outParams())

	}
	return ""
}

var fName string
var oName string
var contractName string
var abi []abiRec

func main() {
	fmt.Println("ABI to Solididy Parser (c) David Appleton 2018")
	fmt.Println("contact : dave@akomba.com")
	fmt.Println("released under Apache 2.0 licence")
	var w *os.File
	var err error
	flag.StringVar(&fName, "input", "", "ABI file to process")
	flag.StringVar(&oName, "output", "", "Solidity file to produce")
	flag.StringVar(&contractName, "contract", "", "Name of contract")
	flag.Parse()
	if fName == "" {
		flag.Usage()
		os.Exit(0)
	}
	if oName == "" {
		w = os.Stdout
	} else {
		if _, err := os.Stat(oName + ".sol"); err == nil {
			fmt.Println("error : ", oName+".sol already exists")
			fmt.Println("we can't have you accidentally deleting files!!!")
			os.Exit(1)
		}

		w, err = os.Create(oName + ".sol")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	if contractName == "" {
		flag.Usage()
		os.Exit(0)
	}

	file, err := os.Open(fName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = json.NewDecoder(bufio.NewReader(file)).Decode(&abi)

	absoluteFile, _ := filepath.Abs(fName)

	fmt.Fprintln(w, "pragma solidity ^0.4.24")

	fmt.Fprintln(w, "// produced by the ABI to Solididy Parser (c) David Appleton 2018")
	fmt.Fprintln(w, "// contact : dave@akomba.com")
	fmt.Fprintln(w, "// released under Apache 2.0 licence")
	fmt.Fprintln(w, "// input ", absoluteFile)
	fmt.Fprintln(w, "// created : ", time.Now().UTC().Format(time.RFC850))

	fmt.Fprintln(w, "\n// produced by ABIParser from dave@akomba")
	fmt.Fprintln(w, "\n\ncontract ", contractName, "{")
	for _, rec := range abi {
		fmt.Fprintln(w, "    "+rec.string())
	}
	fmt.Fprintln(w, "}")
}
