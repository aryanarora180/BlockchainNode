package main

import (
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
	Message string  `json:"message"`
	Chain   []Block `json:"chain"`
}

type Message struct {
	Message string `json:"message"`
}

type RegisterNodeResponse struct {
	Message  string   `json:"message"`
	AllNodes []string `json:"all_nodes"`
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
	router.HandleFunc("/api/nodes", getNodes)
	router.HandleFunc("/api/nodes/resolve", getConsensus)

	router.HandleFunc("/api/transactions/new", postNewTransaction).Methods("POST")
	router.HandleFunc("/api/nodes/register", postRegisterNode).Methods("POST")

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
	lastBlock := chain[len(chain)-1]
	proof := calculateProofOfWork(lastBlock)
	addNewTransaction(Transaction{
		Sender:    "0",
		Recipient: nodeIdentifier,
		Amount:    1,
		Timestamp: time.Now(),
	})

	previousHash := hash(lastBlock)
	block := createNewBlock(proof, previousHash)

	json.NewEncoder(w).Encode(block)
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
		verifiedTransactions = append(verifiedTransactions, block.Transactions...)
	}

	json.NewEncoder(w).Encode(verifiedTransactions)
}

/*
   Get the list of all the registered nodes in the Blockchain network
*/
func getNodes(w http.ResponseWriter, _ *http.Request) {
	json.NewEncoder(w).Encode(nodes)
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

/*
   Function to register a node to the network
*/
func postRegisterNode(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)

	var newNode struct {
		Node string `json:"node"`
	}
	json.Unmarshal(body, &newNode)
	registerNode(newNode.Node)

	json.NewEncoder(w).Encode(RegisterNodeResponse{
		Message:  "Node has been added to the network",
		AllNodes: nodes,
	})
}

func jsonHeaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
