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
        "sender": "Akhil",
        "recipient": "Aryan",
        "amount": 1,
        "timestamp": "2021-09-22T22:13:11.8489877+05:30"
    },
    {
        "sender": "Aryan",
        "recipient": "Rahul",
        "amount": 2,
        "timestamp": "2021-09-22T22:13:20.5471972+05:30"
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
    "timestamp": "2021-09-22T22:13:45.5501624+05:30",
    "transactions": [
        {
            "sender": "Akhil",
            "recipient": "Aryan",
            "amount": 1,
            "timestamp": "2021-09-22T22:13:11.8489877+05:30"
        },
        {
            "sender": "Aryan",
            "recipient": "Rahul",
            "amount": 2,
            "timestamp": "2021-09-22T22:13:20.5471972+05:30"
        },
        {
            "sender": "0",
            "recipient": "5e05a7e06b214cce9c516dd615817b60",
            "amount": 1,
            "timestamp": "2021-09-22T22:13:45.5501624+05:30"
        }
    ],
    "proof": 101161,
    "previous_hash": "c961c3aeb140ad3b42046b8e3faa1b16b542350550ce7017c7286defb424f7b8"
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
        "timestamp": "2021-09-22T22:13:06.3736145+05:30",
        "transactions": null,
        "proof": 100,
        "previous_hash": "1"
    },
    {
        "index": 1,
        "timestamp": "2021-09-22T22:13:45.5501624+05:30",
        "transactions": [
            {
                "sender": "Akhil",
                "recipient": "Aryan",
                "amount": 1,
                "timestamp": "2021-09-22T22:13:11.8489877+05:30"
            },
            {
                "sender": "Aryan",
                "recipient": "Rahul",
                "amount": 2,
                "timestamp": "2021-09-22T22:13:20.5471972+05:30"
            },
            {
                "sender": "0",
                "recipient": "5e05a7e06b214cce9c516dd615817b60",
                "amount": 1,
                "timestamp": "2021-09-22T22:13:45.5501624+05:30"
            }
        ],
        "proof": 101161,
        "previous_hash": "c961c3aeb140ad3b42046b8e3faa1b16b542350550ce7017c7286defb424f7b8"
    },
    {
        "index": 2,
        "timestamp": "2021-09-22T22:14:16.836531+05:30",
        "transactions": [
            {
                "sender": "Akhil",
                "recipient": "Rahul",
                "amount": 3,
                "timestamp": "2021-09-22T22:14:08.9713217+05:30"
            },
            {
                "sender": "0",
                "recipient": "5e05a7e06b214cce9c516dd615817b60",
                "amount": 1,
                "timestamp": "2021-09-22T22:14:16.836531+05:30"
            }
        ],
        "proof": 167273,
        "previous_hash": "96fd6d60e979fc5f65e4531b2ee1fc5c1cad0288a6d599fb9297a4f370dd5181"
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
        "sender": "Akhil",
        "recipient": "Aryan",
        "amount": 1,
        "timestamp": "2021-09-22T22:13:11.8489877+05:30"
    },
    {
        "sender": "Aryan",
        "recipient": "Rahul",
        "amount": 2,
        "timestamp": "2021-09-22T22:13:20.5471972+05:30"
    },
    {
        "sender": "0",
        "recipient": "5e05a7e06b214cce9c516dd615817b60",
        "amount": 1,
        "timestamp": "2021-09-22T22:13:45.5501624+05:30"
    },
    {
        "sender": "Akhil",
        "recipient": "Rahul",
        "amount": 3,
        "timestamp": "2021-09-22T22:14:08.9713217+05:30"
    },
    {
        "sender": "0",
        "recipient": "5e05a7e06b214cce9c516dd615817b60",
        "amount": 1,
        "timestamp": "2021-09-22T22:14:16.836531+05:30"
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
            "timestamp": "2021-09-22T22:13:06.3736145+05:30",
            "transactions": null,
            "proof": 100,
            "previous_hash": "1"
        },
        {
            "index": 1,
            "timestamp": "2021-09-22T22:13:45.5501624+05:30",
            "transactions": [
                {
                    "sender": "Akhil",
                    "recipient": "Aryan",
                    "amount": 1,
                    "timestamp": "2021-09-22T22:13:11.8489877+05:30"
                },
                {
                    "sender": "Aryan",
                    "recipient": "Rahul",
                    "amount": 2,
                    "timestamp": "2021-09-22T22:13:20.5471972+05:30"
                },
                {
                    "sender": "0",
                    "recipient": "5e05a7e06b214cce9c516dd615817b60",
                    "amount": 1,
                    "timestamp": "2021-09-22T22:13:45.5501624+05:30"
                }
            ],
            "proof": 101161,
            "previous_hash": "c961c3aeb140ad3b42046b8e3faa1b16b542350550ce7017c7286defb424f7b8"
        },
        {
            "index": 2,
            "timestamp": "2021-09-22T22:14:16.836531+05:30",
            "transactions": [
                {
                    "sender": "Akhil",
                    "recipient": "Rahul",
                    "amount": 3,
                    "timestamp": "2021-09-22T22:14:08.9713217+05:30"
                },
                {
                    "sender": "0",
                    "recipient": "5e05a7e06b214cce9c516dd615817b60",
                    "amount": 1,
                    "timestamp": "2021-09-22T22:14:16.836531+05:30"
                }
            ],
            "proof": 167273,
            "previous_hash": "96fd6d60e979fc5f65e4531b2ee1fc5c1cad0288a6d599fb9297a4f370dd5181"
        }
    ]
}
``` 

