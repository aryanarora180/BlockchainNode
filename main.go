package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

var nodeIdentifier = strings.Replace(uuid.New().String(), "-", "", -1)

type Transaction struct {
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
	Amount    int    `json:"amount"`
}

type Block struct {
	Index        int           `json:"index"`
	Timestamp    time.Time     `json:"timestamp"`
	Transactions []Transaction `json:"transactions"`
	Proof        int           `json:"proof"`
	PreviousHash string        `json:"previous_hash"`
}

var currentTransactions []Transaction
var chain []Block
var nodes []string

func main() {
	port, err := strconv.Atoi(os.Args[1])
	if err != nil {
		port = 5000
	}

	createNewBlock(100, "1")
	handleRequests(port)
}

func createNewBlock(proof int, previousHash string) Block {
	if previousHash == "" {
		previousHash = hash(chain[len(chain)-1])
	}

	block := Block{
		Index:        len(chain),
		Timestamp:    time.Now(),
		Transactions: currentTransactions,
		Proof:        proof,
		PreviousHash: previousHash,
	}
	currentTransactions = nil
	chain = append(chain, block)

	return block
}

func isValidChain(chain []Block) bool {
	lastBlock := chain[0]
	currentIndex := 1

	for currentIndex < len(chain) {
		block := chain[currentIndex]
		lastBlockHash := hash(lastBlock)
		if block.PreviousHash != lastBlockHash {
			return false
		}

		if !isValidProof(lastBlock.Proof, block.Proof, lastBlockHash) {
			return false
		}

		lastBlock = block
		currentIndex++
	}

	return true
}

func addNewTransaction(transaction Transaction) int {
	currentTransactions = append(currentTransactions, transaction)
	return chain[len(chain)-1].Index + 1
}

func isValidProof(lastProof int, currentProof int, lastHash string) bool {
	guess := strconv.Itoa(lastProof) + strconv.Itoa(currentProof) + lastHash
	guessHash := sha256.Sum256([]byte(guess))
	return hex.EncodeToString(guessHash[:])[:4] == "0000"
}

func hash(block Block) string {
	blockString, err := json.Marshal(block)
	if err != nil {
		panic(err)
	}
	hash := sha256.Sum256(blockString)
	return hex.EncodeToString(hash[:])
}

func calculateProofOfWork(lastBlock Block) int {
	lastProof := lastBlock.Proof
	lastHash := hash(lastBlock)

	proof := 0
	for !isValidProof(lastProof, proof, lastHash) {
		proof++
	}

	return proof
}

func registerNode(address string) {
	parsedUrl, err := url.Parse(address)
	if err != nil {
		panic(err)
	}

	if parsedUrl.Host != "" {
		nodes = append(nodes, parsedUrl.Host)
	} else if parsedUrl.Path != "" {
		nodes = append(nodes, parsedUrl.Path)
	} else {
		panic("Invalid URL")
	}
}

func resolveConflicts() bool {
	var newChain []Block
	maxLength := len(chain)

	for _, node := range nodes {
		response, err := http.Get("http://" + node + "/chain")
		if err == nil && response.StatusCode == 200 {
			responseBody, err := ioutil.ReadAll(response.Body)
			if err == nil {
				var conflictResolveResponse struct {
					length int
					chain  []Block
				}
				json.Unmarshal(responseBody, &conflictResolveResponse)
				if conflictResolveResponse.length > maxLength && isValidChain(conflictResolveResponse.chain) {
					maxLength = conflictResolveResponse.length
					newChain = conflictResolveResponse.chain
				}
			}
		}
	}

	if newChain != nil {
		chain = newChain
		return true
	}

	return false
}
