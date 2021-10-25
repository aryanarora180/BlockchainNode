package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type ConsensusResponse struct {
	Message string        `json:"message"`
	Chain   []SignedBlock `json:"chain"`
}

type Message struct {
	Message string `json:"message"`
}

type RegisterNodeResponse struct {
	Message  string `json:"message"`
	AllNodes []Node `json:"all_nodes"`
}

/*
   Different functionalities implemented in Blockchain context
         :param port: port address
*/
func handleRequests(port int) {
	addr := ":" + strconv.Itoa(port)

	router := mux.NewRouter().StrictSlash(true)
	router.Use(jsonHeaderMiddleware)

	header := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST"})
	origins := handlers.AllowedOrigins([]string{"*"})

	router.HandleFunc("/api/transactions/unverified", getUnverifiedTransactions)
	router.HandleFunc("/api/transactions/verified", getVerifiedTransactions)
	router.HandleFunc("/api/mine", mine)
	router.HandleFunc("/api/chain", getChain)
	router.HandleFunc("/api/verify_chain", verifyCurrentChain)
	router.HandleFunc("/api/nodes", getNodes)
	router.HandleFunc("/api/nodes/resolve", getConsensus)

	router.HandleFunc("/api/transactions/new", postNewTransaction).Methods("POST")
	router.HandleFunc("/api/nodes/validatorify", makeNodeValidator).Methods("POST")

	fmt.Printf("Server attemping to listen on %v", addr)
	err := http.ListenAndServe(addr, handlers.CORS(header, methods, origins)(router))
	if err != nil {
		fmt.Printf("An error occured while setting up server on %v", addr)
	}
}

/*
   Function to mine
*/
func mine(w http.ResponseWriter, _ *http.Request) {
	if !checkIfValidator(getRsaPublicKeyAsBase64Str(publicKey)) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Message{Message: "Cannot mine block since you are not a validator"})
		return
	}

	lastBlock := chain[len(chain)-1]
	if lastBlock.SignerPublicKey == getRsaPublicKeyAsBase64Str(publicKey) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Message{Message: "Cannot mine block since you mined the last block"})
		return
	}

	addNewTransaction(Transaction{
		Sender:    "0",
		Recipient: getRsaPublicKeyAsBase64Str(publicKey),
		Amount:    1,
		Timestamp: time.Now(),
	})

	previousHash := hash(lastBlock.Data)

	block := createNewBlock(convertHashToString(previousHash), false)
	blockHash := hash(block)

	signature, err := rsa.SignPSS(rand.Reader, privateKey, crypto.SHA256, blockHash, nil)
	if err != nil {
		panic(err)
	}

	err = rsa.VerifyPSS(&privateKey.PublicKey, crypto.SHA256, blockHash, signature, nil)
	if err != nil {
		fmt.Println("Could not verify signature: ", err)
		return
	}

	newSignedBlock := SignedBlock{
		Data:            block,
		Signature:       base64.StdEncoding.EncodeToString(signature),
		SignerPublicKey: getRsaPublicKeyAsBase64Str(publicKey),
	}

	currentTransactions = nil
	chain = append(chain, newSignedBlock)

	json.NewEncoder(w).Encode(newSignedBlock)
}

/*
   Get the current chain
*/
func getChain(w http.ResponseWriter, _ *http.Request) {
	json.NewEncoder(w).Encode(chain)
}

/*
   Get the list of unverified transactions
*/
func getUnverifiedTransactions(w http.ResponseWriter, _ *http.Request) {
	json.NewEncoder(w).Encode(currentTransactions)
}

/*
   Get the list of verified transactions
		This is done by going over all the blocks in the chain and keeping track of all the transactions
*/
func getVerifiedTransactions(w http.ResponseWriter, _ *http.Request) {
	var verifiedTransactions []Transaction
	for _, block := range chain {
		verifiedTransactions = append(verifiedTransactions, block.Data.Transactions...)
	}

	json.NewEncoder(w).Encode(verifiedTransactions)
}

/*
   Get the list of all the registered nodes in the Blockchain network
*/
func getNodes(w http.ResponseWriter, _ *http.Request) {
	json.NewEncoder(w).Encode(getNodesList(false))
}

/*
   Resolve conflicts using Longest chain rule as the consensus parameter:
         Whether the chain is being replaced or not.
*/
func getConsensus(w http.ResponseWriter, _ *http.Request) {
	replaced := resolveConflicts()

	if replaced {
		json.NewEncoder(w).Encode(ConsensusResponse{
			Message: "Replaced",
			Chain:   chain,
		})
	} else {
		json.NewEncoder(w).Encode(ConsensusResponse{
			Message: "Authoritative",
			Chain:   chain,
		})
	}
}

/*
   Function to register a new transaction
*/
func postNewTransaction(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)

	var transaction struct {
		Sender    string `json:"sender"`
		Recipient string `json:"recipient"`
		Amount    int    `json:"amount"`
	}
	err := json.Unmarshal(body, &transaction)
	if err != nil {
		panic(err)
	}

	//adding the transactions to out current transactions list, which will be added on the block on being mined
	index := addNewTransaction(Transaction{
		Sender:    transaction.Sender,
		Recipient: transaction.Recipient,
		Amount:    transaction.Amount,
		Timestamp: time.Now(),
	})
	err = json.NewEncoder(w).Encode(Message{Message: fmt.Sprintf("Transaction will be added to block %v", index)})
	if err != nil {
		panic(err)
	}
}

func makeNodeValidator(w http.ResponseWriter, r *http.Request) {
	if checkIfValidator(getRsaPublicKeyAsBase64Str(publicKey)) {
		body, _ := ioutil.ReadAll(r.Body)

		var validator struct {
			Url       string `json:"url"`
			PublicKey string `json:"public_key"`
		}
		err := json.Unmarshal(body, &validator)
		if err != nil {
			panic(err)
		}

		addToList([]Node{
			{
				Url:         validator.Url,
				PublicKey:   validator.PublicKey,
				IsValidator: 1,
			},
		})

		json.NewEncoder(w).Encode(Message{Message: fmt.Sprintf("Made %v a validator", validator.PublicKey[0:9])})
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Message{Message: "Non-validators cannot make a node a validator"})
		return
	}
}

func verifyCurrentChain(w http.ResponseWriter, _ *http.Request) {
	valid := isValidChain(chain)

	json.NewEncoder(w).Encode(struct {
		Valid bool `json:"valid"`
	}{
		Valid: valid,
	})
}

func jsonHeaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
