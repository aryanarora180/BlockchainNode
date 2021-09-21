package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
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

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.Use(jsonHeaderMiddleware)

	router.HandleFunc("/api/transactions/unverified", getUnverifiedTransactions)
	router.HandleFunc("/api/transactions/verified", getVerifiedTransactions)
	router.HandleFunc("/api/mine", mine)
	router.HandleFunc("/api/chain", getChain)
	router.HandleFunc("/api/nodes", getNodes)
	router.HandleFunc("/api/nodes/resolve", getConsensus)

	router.HandleFunc("/api/transactions/new", postNewTransaction).Methods("POST")
	router.HandleFunc("/api/nodes/register", postRegisterNode).Methods("POST")

	http.ListenAndServe(":5000", router)
}

func mine(w http.ResponseWriter, _ *http.Request) {
	lastBlock := chain[len(chain)-1]
	proof := calculateProofOfWork(lastBlock)
	addNewTransaction(Transaction{
		Sender:    "0",
		Recipient: nodeIdentifier,
		Amount:    1,
	})

	previousHash := hash(lastBlock)
	block := createNewBlock(proof, previousHash)

	json.NewEncoder(w).Encode(block)
}

func getChain(w http.ResponseWriter, _ *http.Request) {
	json.NewEncoder(w).Encode(chain)
}

func getUnverifiedTransactions(w http.ResponseWriter, _ *http.Request) {
	json.NewEncoder(w).Encode(currentTransactions)
}

func getVerifiedTransactions(w http.ResponseWriter, _ *http.Request) {
	var verifiedTransactions []Transaction
	for _, block := range chain {
		verifiedTransactions = append(verifiedTransactions, block.Transactions...)
	}

	json.NewEncoder(w).Encode(verifiedTransactions)
}

func getNodes(w http.ResponseWriter, _ *http.Request) {
	json.NewEncoder(w).Encode(nodes)
}

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

func postNewTransaction(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)

	var transaction Transaction
	err := json.Unmarshal(body, &transaction)
	if err != nil {
		panic(err)
	}

	index := addNewTransaction(transaction)
	err = json.NewEncoder(w).Encode(Message{Message: fmt.Sprintf("Transaction will be added to block %v", index)})
	if err != nil {
		panic(err)
	}
}

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
