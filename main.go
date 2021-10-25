package main

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

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
	PreviousHash string        `json:"previous_hash"`
}

type SignedBlock struct {
	Data            Block  `json:"data"`
	Signature       string `json:"signature"`
	SignerPublicKey string `json:"signer_public_key"`
}

var currentTransactions []Transaction
var chain []SignedBlock

var privateKey *rsa.PrivateKey
var publicKey *rsa.PublicKey

/*
Create the genesis Block in the Blockchain
        As it is the first block, the previous hash v is set to 1.
*/
func main() {
	privateKey, publicKey = getRsaKeyPair()

	port, err := strconv.Atoi(os.Args[1])
	if err != nil {
		port = 5000
	}

	parsedUrl, err := url.Parse("http://localhost:" + strconv.Itoa(port))
	if err != nil {
		panic(err)
	}

	var path string
	if parsedUrl.Host != "" {
		path = parsedUrl.Host
	} else if parsedUrl.Path != "" {
		path = parsedUrl.Path
	} else {
		path = ""
	}
	addCurrentNode(path)

	createNewBlock("1", true)
	handleRequests(port)
}

/*
Create a new Block in the Blockchain
        :param proof: The proof given by the Proof of Work algorithm
        :param previousHash: hash of previous Block
        :return: Newly created Block
*/
func createNewBlock(previousHash string, addToChain bool) Block {
	if previousHash == "" {
		previousHash = convertHashToString(hash(chain[len(chain)-1].Data))
	}

	block := Block{
		Index:        len(chain),
		Timestamp:    time.Now(),
		Transactions: currentTransactions,
		PreviousHash: previousHash,
	}

	if addToChain {
		currentTransactions = nil
		chain = append(chain, SignedBlock{
			Data:            block,
			Signature:       "",
			SignerPublicKey: "",
		})
	}

	return block
}

/*
   Determine if the given blockchain is valid
       :param chain: A blockchain
       :return: True if valid, False if not
*/
func isValidChain(chain []SignedBlock) bool {
	lastBlock := chain[0]
	currentIndex := 1

	for currentIndex < len(chain) {
		block := chain[currentIndex]

		// verifying hash of the Block
		lastBlockHash := hash(lastBlock.Data)
		if block.Data.PreviousHash != convertHashToString(lastBlockHash) {
			return false
		}

		if lastBlock.SignerPublicKey == block.SignerPublicKey {
			return false
		}

		if !checkIfValidator(block.SignerPublicKey) {
			return false
		}

		signature, err := base64.StdEncoding.DecodeString(block.Signature)
		signerPubKeyString, err := base64.StdEncoding.DecodeString(block.SignerPublicKey)
		signerPubKey, err := parseRsaPublicKeyFromPemStr(string(signerPubKeyString))

		err = rsa.VerifyPSS(signerPubKey, crypto.SHA256, hash(block.Data), signature, nil)
		if err != nil {
			fmt.Println("Could not verify signature: ", err)
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
	return chain[len(chain)-1].Data.Index + 1
}

/*
   Creates the hash of a Block (using SHA-256)
       :param block: Block
*/
func hash(block Block) []byte {
	blockString, err := json.Marshal(block)
	if err != nil {
		panic(err)
	}
	hash := sha256.Sum256(blockString)
	return hash[:]
}

func convertHashToString(hash []byte) string {
	return hex.EncodeToString(hash)
}

/*
   Resolve conflicts using the longest chain rule.
       :return: true if the chain is being replaced, false if not.
*/
func resolveConflicts() bool {
	var newChain []SignedBlock

	// interested in the longest chain, which we will find.
	maxLength := len(chain)

	// check the chains for all the nodes and replace the current chain with the longest chain (if it isn't already the longest)
	for _, node := range getNodesList(true) {
		response, err := http.Get("http://" + node.Url + "/api/chain")
		if err == nil && response.StatusCode == 200 {
			responseBody, err := ioutil.ReadAll(response.Body)
			if err == nil {
				var chain []SignedBlock
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
