package main

import (
	"encoding/csv"
	"os"
	"strconv"
)

type Node struct {
	Url         string `json:"url"`
	PublicKey   string `json:"public_key"`
	IsValidator int    `json:"is_validator"`
}

func addToList(newNodes []Node) {
	allNodes := getNodesList(false)

	var pubKeys []string
	for _, line := range allNodes {
		pubKeys = append(pubKeys, line.PublicKey)
	}

	for _, node := range newNodes {
		if !stringInSlice(node.PublicKey, pubKeys) {
			allNodes = append(allNodes, node)
		} else {
			for index, oldNode := range allNodes {
				if oldNode.PublicKey == node.PublicKey {
					allNodes[index] = node
					break
				}
			}
		}
	}

	csvFile, err := os.Create("nodes.csv")
	defer csvFile.Close()
	if err != nil {
		panic(err)
	}

	csvWriter := csv.NewWriter(csvFile)
	defer csvWriter.Flush()

	for _, node := range allNodes {
		err = csvWriter.Write([]string{node.PublicKey, node.Url, strconv.Itoa(node.IsValidator)})
		if err != nil {
			panic(err)
		}
	}
}

func addCurrentNode(path string) {
	if _, err := os.Stat("nodes.csv"); os.IsNotExist(err) {
		addToList([]Node{
			{
				Url:         path,
				PublicKey:   getRsaPublicKeyAsBase64Str(publicKey),
				IsValidator: 1,
			},
		})
	} else {
		addToList([]Node{
			{
				Url:         path,
				PublicKey:   getRsaPublicKeyAsBase64Str(publicKey),
				IsValidator: 0,
			},
		})
	}
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func getNodesList(validatorsOnly bool) []Node {
	if _, err := os.Stat("nodes.csv"); os.IsNotExist(err) {
		os.Create("nodes.csv")
		return nil
	}

	csvFile, err := os.Open("nodes.csv")
	if err != nil {
		panic(err)
	}

	csvLines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		panic(err)
	}

	var nodes []Node
	for _, line := range csvLines {
		var node Node
		node.PublicKey = line[0]
		node.Url = line[1]
		node.IsValidator, _ = strconv.Atoi(line[2])

		if validatorsOnly {
			if node.IsValidator == 1 {
				nodes = append(nodes, node)
			}
		} else {
			nodes = append(nodes, node)
		}
	}

	csvFile.Close()

	return nodes
}

func checkPubkeyExists(pubKey string) bool {
	nodes := getNodesList(false)

	for _, node := range nodes {
		if pubKey == node.PublicKey {
			return true
		}
	}

	return false
}

func checkIfValidator(pubkey string) bool {
	nodes := getNodesList(false)

	for _, node := range nodes {
		if pubkey == node.PublicKey {
			if node.IsValidator == 1 {
				return true
			}
		}
	}

	return false
}
