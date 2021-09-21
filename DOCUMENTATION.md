# API Documentation

## Structure and basic protocols followed

- **BASE_URL:** The hostname like http://localhost:5000
- **API_PATH:** BASE_URL/api

- GET: Used to get some data from the server. 
- POST: To send any data to the server.

## Endpoints

### Get unverified transactions
Get a list of transactions that have not been mined (unverified) yet

**Endpoint:** API_PATH/transactions/unverified

**Method:** GET

**Parameters:** None

**Response:**
- Status: 200
- Body: A list of transactions.
```
[
    {
        "sender": "Rahul",
        "recipient": "Akhil",
        "amount": 1000
    },
    {
        "sender": "Aryan",
        "recipient": "Rahul",
        "amount": 2000
    }
]
```

### Add new transaction
Adds a transaction to the list of unverified transactions, which when mined get added to a new block and become verified.

**Endpoint:** API_PATH/transactions/new

**Method:** POST

**Parameters:**
- "sender": string, the name of the sender
- "recipient": string, the name of the recipient
- "amount": integer, the amount of coins the sender is sending to the recipient

**Response:**
- Status: 200
- Body: A message denoting the index of the block the transaction will be added to once mined.
```
{
    "message": "Transaction will be added to block 1"
}
``` 

### Mine
Adds all unverified transactions to a new block, which verifies the transactions.

**Endpoint:** API_PATH/mine

**Method:** GET

**Parameters:** None

**Response:**
- Status: 200
- Body: The new block that's been mined.
```
{
    "index": 1,
    "timestamp": "2021-09-21T11:34:14.3950913+05:30",
    "transactions": [
        {
            "sender": "Rahul",
            "recipient": "Akhil",
            "amount": 2100
        },
        {
            "sender": "Akhil",
            "recipient": "Aryan",
            "amount": 1100
        },
        {
            "sender": "Aryan",
            "recipient": "Akhil",
            "amount": 2000
        },
        {
            "sender": "0",
            "recipient": "d168b64f829745358867c3dd5257d2fc",
            "amount": 1
        }
    ],
    "proof": 18504,
    "previous_hash": "4debf4626a20b7b2bb6f54496286017e2a957c4eea3af3139ec881b10ad38fda"
}
``` 

### Chain
Gets a list of all the blocks that have been mined on the network.

**Endpoint:** API_PATH/chain

**Method:** GET

**Parameters:** None

**Response:**
- Status: 200
- Body: The list of all blocks.
```
[
    {
        "index": 0,
        "timestamp": "2021-09-21T11:31:22.0185028+05:30",
        "transactions": null,
        "proof": 100,
        "previous_hash": "1"
    },
    {
        "index": 1,
        "timestamp": "2021-09-21T11:34:14.3950913+05:30",
        "transactions": [
            {
                "sender": "Rahul",
                "recipient": "Akhil",
                "amount": 2100
            },
            {
                "sender": "Akhil",
                "recipient": "Aryan",
                "amount": 1100
            },
            {
                "sender": "Aryan",
                "recipient": "Akhil",
                "amount": 2000
            },
            {
                "sender": "0",
                "recipient": "d168b64f829745358867c3dd5257d2fc",
                "amount": 1
            }
        ],
        "proof": 18504,
        "previous_hash": "4debf4626a20b7b2bb6f54496286017e2a957c4eea3af3139ec881b10ad38fda"
    }
]
``` 

### Get verified transactions
Get a list of transactions from all the blocks that have been mined

**Endpoint:** API_PATH/transactions/verified

**Method:** GET

**Parameters:** None

**Response:**
- Status: 200
- Body: A list of transactions.
```
[
    {
        "sender": "Aryan",
        "recipient": "Rahul",
        "amount": 2000
    },
    {
        "sender": "Aryan",
        "recipient": "Rahul",
        "amount": 2000
    },
    {
        "sender": "0",
        "recipient": "ef4bd7dffe264867b5694c0a0936a95f",
        "amount": 1
    },
    {
        "sender": "Rahul",
        "recipient": "Aryan",
        "amount": 1000
    },
    {
        "sender": "Rahul",
        "recipient": "Aryan",
        "amount": 1000
    },
    {
        "sender": "Akhil",
        "recipient": "Aryan",
        "amount": 2000
    },
    {
        "sender": "0",
        "recipient": "ef4bd7dffe264867b5694c0a0936a95f",
        "amount": 1
    }
]
```

### Get nodes
Gets a list of all nodes registered on the network.

**Endpoint:** API_PATH/nodes

**Method:** GET

**Parameters:** None

**Response:**
- Status: 200
- Body: A list of all the hostnames of the nodes
```
[
    "localhost:5001",
    "localhost:5002",
    "localhost:5003"
]
``` 

### Register a new node
Register a new node onto the network. To get the new node's blocks upto date, you need to send another request to resolve conflicts.

**Endpoint:** API_PATH/nodes/register

**Method:** POST

**Parameters:** 
- "node": The full URL of the node that needs to be registered onto the network. Eg: "http://localhost:5001"

**Response:**
- Status: 200
- Body: A list of all the hostnames of the nodes
```
{
    "message": "Node has been added to the network",
    "all_nodes": [
        "localhost:5001",
        "localhost:5002",
        "localhost:5003"
    ]
}
``` 

### Resolve conflicts
Updates the node's chain
**Endpoint:** API_PATH/nodes/resolve

**Method:** GET

**Parameters:** None

**Response:**
- Status: 200
- Body: 
    - "message":
        - "Replaced": Your node was outdated, and the chain has been updated
        - "Authoritative":  Your node was up to date, and there was no need to update the node
    - "chain": The new, up to date chain
```
{
    "message": "Authoritative",
    "chain": [
        {
            "index": 0,
            "timestamp": "2021-09-21T12:34:24.0753126+05:30",
            "transactions": null,
            "proof": 100,
            "previous_hash": "1"
        },
        {
            "index": 1,
            "timestamp": "2021-09-21T12:43:27.441507+05:30",
            "transactions": [
                {
                    "sender": "Aryan",
                    "recipient": "Akhil",
                    "amount": 2000
                },
                {
                    "sender": "0",
                    "recipient": "dafd9c4b73ee42d19398490245c54b88",
                    "amount": 1
                }
            ],
            "proof": 6786,
            "previous_hash": "fa7f9a6c544a0bc316070a9578521e9638511608912ebd71c0caf80e315502d6"
        }
    ]
}
``` 

