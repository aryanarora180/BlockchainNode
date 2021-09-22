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
	Sender    string    `json:"sender"`
	Recipient string    `json:"recipient"`
	Amount    int       `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
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

/*
Create the genesis Block in the Blockchain
        As it is the first block, the previous hash is set to 1.
*/
func main() {
	port, err := strconv.Atoi(os.Args[1])
	if err != nil {
		port = 5000
	}

	createNewBlock(100, "1")
	handleRequests(port)
}

/*
Create a new Block in the Blockchain
        :param proof: The proof given by the Proof of Work algorithm
        :param previousHash: hash of previous Block
        :return: Newly created Block
*/
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

	// the current transaction list is reset after being added to the block
	currentTransactions = nil
	chain = append(chain, block)

	return block
}

/*
   Determine if the given blockchain is valid
       :param chain: A blockchain
       :return: True if valid, False if not
*/
func isValidChain(chain []Block) bool {
	lastBlock := chain[0]
	currentIndex := 1

	for currentIndex < len(chain) {
		block := chain[currentIndex]

		// verifying hash of the Block
		lastBlockHash := hash(lastBlock)
		if block.PreviousHash != lastBlockHash {
			return false
		}

		// validating the proof of work
		if !isValidProof(lastBlock.Proof, block.Proof, lastBlockHash) {
			return false
		}

		lastBlock = block
		currentIndex++
	}

	return true
}

/*
   Add a Transaction to the list of unverified Transactions [ Will be added to the next mined Block ]
       :param transaction: Details of the transaction
       :return: The index of the Block that will hold this transaction
*/
func addNewTransaction(transaction Transaction) int {
	currentTransactions = append(currentTransactions, transaction)
	return chain[len(chain)-1].Index + 1
}

/*
   Check if the proof is valid
       :param lastProof: Previous proof
       :param proof: Current proof
       :param lastHash: The hash of the Previous Block
       :return: True if valid, False if not.
*/
func isValidProof(lastProof int, currentProof int, lastHash string) bool {
	guess := strconv.Itoa(lastProof) + strconv.Itoa(currentProof) + lastHash
	guessHash := sha256.Sum256([]byte(guess))

	// checking if the hash has 4 leading zeroes.
	return hex.EncodeToString(guessHash[:])[:4] == "0000"
}

/*
   Creates the hash of a Block (using SHA-256)
       :param block: Block
*/
func hash(block Block) string {
	blockString, err := json.Marshal(block)
	if err != nil {
		panic(err)
	}
	hash := sha256.Sum256(blockString)
	return hex.EncodeToString(hash[:])
}

/*
   Proof of Work Algorithm implemented is:
        - Find a number p' such that hash(pp') contains leading 4 zeroes
        - Where p is the previous proof, and p' is the new proof

       :param lastBlock: last Block
       :return: the proof value [ p' ]
*/
func calculateProofOfWork(lastBlock Block) int {
	lastProof := lastBlock.Proof
	lastHash := hash(lastBlock)

	proof := 0
	for !isValidProof(lastProof, proof, lastHash) {
		proof++
	}

	return proof
}

/*
   Registering a Node:
       :param address: address of the Node
*/
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

/*
   Resolve conflicts using the longest chain rule.
       :return: true if the chain is being replaced, false if not.
*/
func resolveConflicts() bool {
	var newChain []Block

	// interested in the longest chain, which we will find.
	maxLength := len(chain)

	// check the chains for all the nodes and replace the current chain with the longest chain (if it isn't already the longest)
	for _, node := range nodes {
		response, err := http.Get("http://" + node + "/api/chain")
		if err == nil && response.StatusCode == 200 {
			responseBody, err := ioutil.ReadAll(response.Body)
			if err == nil {
				var chain []Block
				json.Unmarshal(responseBody, &chain)
				if len(chain) > maxLength && isValidChain(chain) {
					maxLength = len(chain)
					newChain = chain
				}
			}
		}
	}

	// if a longer chain id found replace current chain with it
	if newChain != nil {
		chain = newChain
		return true
	}

	return false
}
