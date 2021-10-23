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
- Body: A list of transactions. If no unverified transactions exist, it returns null.
```
[
  {
    "sender": "Rahul",
    "recipient": "Akhil",
    "amount": 2,
    "timestamp": "2021-10-25T15:34:01.7255311+05:30"
  },
  {
    "sender": "Aryan",
    "recipient": "Akhil",
    "amount": 1,
    "timestamp": "2021-10-25T15:34:20.3794288+05:30"
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
You can only mine if:
- Your node is a validator
- You haven't mined the last block.

**Endpoint:** API_PATH/mine

**Method:** GET

**Parameters:** None

**Response if allowed to mine:**
- Status: 200
- Body: The new block that's been mined.
```
{
  "data": {
    "index": 1,
    "timestamp": "2021-10-25T15:38:05.7972542+05:30",
    "transactions": [
      {
        "sender": "Rahul",
        "recipient": "Akhil",
        "amount": 2,
        "timestamp": "2021-10-25T15:34:01.7255311+05:30"
      },
      {
        "sender": "Aryan",
        "recipient": "Akhil",
        "amount": 1,
        "timestamp": "2021-10-25T15:34:20.3794288+05:30"
      },
      {
        "sender": "0",
        "recipient": "b1e8793688be425d9af1b0682cfb96ca",
        "amount": 1,
        "timestamp": "2021-10-25T15:38:05.7971658+05:30"
      }
    ],
    "previous_hash": "6e94b60c1a194342a5d526c6ffa61daa0208ae620fdb411ba3d81bf250da7b6a"
  },
  "signature": "k+uZG3XbxO9CA9XcVLKJ8NjurDkh9OMsduHp7dHh97rh+E6EyyclTb4ikDF/cEnAwVdY7znCayw6gUs//pNPbuP35vlsDw7IG78kMb0PBEXN4bdrtMKUVqCT/G/Q/T8kTdHc6SONlfjELegBmu5Ki3ZtzDDR4TCqsfmJt8sphoklKMMbyxvxkYrdKD2f83pHSPYoMFLrwIxqFCABPGSad47xd7ukwM2TXkowqPquq/3DGztxEjSvU/0N0UcZsS9IN1S8MSXz9toA9TrXA6/zB+EtZh58hHA/cs/57oto9jF/Z+Ps1qyQg3GvmwOoWmM3iWRpE3v0TbzyHIuGnYWLXj7xSnqAY1MCc9ebwSi9rtA5W+mNJSYZGhbm7GiOfFxdIIrbcAkHf6LYe6Bp2JOyd35cYvyOpRRPwtN92muOZuCYpVLlp9YxfLtMy3aMKxBe/mgIzyul72TdmUSiR41fvoBXU2UAWdJegcSHeL/aJJHoobqhVfg75wQExRpFnL/d0+lQKNbJ1htQcXny/tUbkwbeUEgk09Y1WbpBt5OR5Z7JnsUSi8FUj7OxZIr8tpCBOoa+dRpa3aCp7dswWwnjB5SkoV4aRj/iXskpQ/G/sQlQWP+JngKaNrIQkB4HAYEvXlNhS8jeViNWegD8MWkQpmnQq0LZHPui/ttSBJcWu5g=",
  "signer_public_key": "LS0tLS1CRUdJTiBSU0EgUFVCTElDIEtFWS0tLS0tCk1JSUNJakFOQmdrcWhraUc5dzBCQVFFRkFBT0NBZzhBTUlJQ0NnS0NBZ0VBdHpHcDRTUEh5Q0k5TUtBWU9TdmsKR0trVlRjbVJVYnFpZ0x0QVZ1Vm9KcDRvOFU3Z1h6Z21RNitNb1RWTy9VN2Y4dENQeHp4MEUrWWlxZWlVOE1wTgpuc2wxS0lHYi9ndUdrK2ExMDhKNzgxK3B6aEt1ekY0RlBRcXYvNmM3WFI3N1l2NUN6YTJkUTRONHlxT3p3TjZKCkNtSjRPNU9hVEowcFIzeTBUd0dHTjlQdmdHVFF0WnFjZnFoL2NKMDNBMVU3VlBuQTlMdDFkQnpydnFIZ3BjSzEKU1Q1TldodW1HR2lBNnFWakQ3QlpHQmJTaXdaSW1mVmRiUVRFRUJsa09BZ1BUcDViak0xUWpMK0FmVjlKd1BVVApsNkYvQVVrNEFKVmlyRjlMZS9FSmpzdzJtcmNDN0swNmtLaXN4ZWp1MURSaElML1RYaXJzNTRnUWI3dGQ5c2xMClkyalhOWEVnU1RPVWs2K04zejk3cVJqMHRBS1RlbWJGb2dPUHlISWVPdktYQ3FuYkJ0VTRpSVNpU2lzeGdnQWYKV1NzWGFGZ210SHVBcEpicVl5RFZLVG9BSmF4Zklvbm8vd21yNUJOT1dBK29jcHRyV0xlTmVWQnpMVzJsMlFRUgpXTjR2NVQ0WnF6bjdLQ2FJdU4vWG9iV3M3NUp1UHd2bFFsUUR6Z1BjUnp6UTRVbFlmKzdiTCtHamdRSFFXL1R3CjlmUlJ6M2dHV3VqeDdvcWYyRk56UysydnhtSmxGZ1dFdVh4ekNqeU1SMWl2cjFlbVQ3akZudVRiYzQrZ0NMSkUKUXhqQUl3b0YrRXlwK2dSOFdQQW1ScmVOcU9kSXlSdEF3TkhYRDJJVTc4NHJxZTBkQktoUFVvcDFYYW1OWitvQgpocWNvanAxWlB5ZS9oTXJkUWpPNzJwMENBd0VBQVE9PQotLS0tLUVORCBSU0EgUFVCTElDIEtFWS0tLS0tCg=="
}
``` 

**Response if not allowed to mine:**
- Status: 400
- Body: A message denoting why you're not allowed to mine.
``` 
{
  "message": "Cannot mine block since you mined the last block"
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
    "data": {
      "index": 0,
      "timestamp": "2021-10-25T15:33:58.1437614+05:30",
      "transactions": null,
      "previous_hash": "1"
    },
    "signature": "",
    "signer_public_key": ""
  },
  {
    "data": {
      "index": 1,
      "timestamp": "2021-10-25T15:38:05.7972542+05:30",
      "transactions": [
        {
          "sender": "Rahul",
          "recipient": "Akhil",
          "amount": 2,
          "timestamp": "2021-10-25T15:34:01.7255311+05:30"
        },
        {
          "sender": "Aryan",
          "recipient": "Akhil",
          "amount": 1,
          "timestamp": "2021-10-25T15:34:20.3794288+05:30"
        },
        {
          "sender": "0",
          "recipient": "b1e8793688be425d9af1b0682cfb96ca",
          "amount": 1,
          "timestamp": "2021-10-25T15:38:05.7971658+05:30"
        }
      ],
      "previous_hash": "6e94b60c1a194342a5d526c6ffa61daa0208ae620fdb411ba3d81bf250da7b6a"
    },
    "signature": "k+uZG3XbxO9CA9XcVLKJ8NjurDkh9OMsduHp7dHh97rh+E6EyyclTb4ikDF/cEnAwVdY7znCayw6gUs//pNPbuP35vlsDw7IG78kMb0PBEXN4bdrtMKUVqCT/G/Q/T8kTdHc6SONlfjELegBmu5Ki3ZtzDDR4TCqsfmJt8sphoklKMMbyxvxkYrdKD2f83pHSPYoMFLrwIxqFCABPGSad47xd7ukwM2TXkowqPquq/3DGztxEjSvU/0N0UcZsS9IN1S8MSXz9toA9TrXA6/zB+EtZh58hHA/cs/57oto9jF/Z+Ps1qyQg3GvmwOoWmM3iWRpE3v0TbzyHIuGnYWLXj7xSnqAY1MCc9ebwSi9rtA5W+mNJSYZGhbm7GiOfFxdIIrbcAkHf6LYe6Bp2JOyd35cYvyOpRRPwtN92muOZuCYpVLlp9YxfLtMy3aMKxBe/mgIzyul72TdmUSiR41fvoBXU2UAWdJegcSHeL/aJJHoobqhVfg75wQExRpFnL/d0+lQKNbJ1htQcXny/tUbkwbeUEgk09Y1WbpBt5OR5Z7JnsUSi8FUj7OxZIr8tpCBOoa+dRpa3aCp7dswWwnjB5SkoV4aRj/iXskpQ/G/sQlQWP+JngKaNrIQkB4HAYEvXlNhS8jeViNWegD8MWkQpmnQq0LZHPui/ttSBJcWu5g=",
    "signer_public_key": "LS0tLS1CRUdJTiBSU0EgUFVCTElDIEtFWS0tLS0tCk1JSUNJakFOQmdrcWhraUc5dzBCQVFFRkFBT0NBZzhBTUlJQ0NnS0NBZ0VBdHpHcDRTUEh5Q0k5TUtBWU9TdmsKR0trVlRjbVJVYnFpZ0x0QVZ1Vm9KcDRvOFU3Z1h6Z21RNitNb1RWTy9VN2Y4dENQeHp4MEUrWWlxZWlVOE1wTgpuc2wxS0lHYi9ndUdrK2ExMDhKNzgxK3B6aEt1ekY0RlBRcXYvNmM3WFI3N1l2NUN6YTJkUTRONHlxT3p3TjZKCkNtSjRPNU9hVEowcFIzeTBUd0dHTjlQdmdHVFF0WnFjZnFoL2NKMDNBMVU3VlBuQTlMdDFkQnpydnFIZ3BjSzEKU1Q1TldodW1HR2lBNnFWakQ3QlpHQmJTaXdaSW1mVmRiUVRFRUJsa09BZ1BUcDViak0xUWpMK0FmVjlKd1BVVApsNkYvQVVrNEFKVmlyRjlMZS9FSmpzdzJtcmNDN0swNmtLaXN4ZWp1MURSaElML1RYaXJzNTRnUWI3dGQ5c2xMClkyalhOWEVnU1RPVWs2K04zejk3cVJqMHRBS1RlbWJGb2dPUHlISWVPdktYQ3FuYkJ0VTRpSVNpU2lzeGdnQWYKV1NzWGFGZ210SHVBcEpicVl5RFZLVG9BSmF4Zklvbm8vd21yNUJOT1dBK29jcHRyV0xlTmVWQnpMVzJsMlFRUgpXTjR2NVQ0WnF6bjdLQ2FJdU4vWG9iV3M3NUp1UHd2bFFsUUR6Z1BjUnp6UTRVbFlmKzdiTCtHamdRSFFXL1R3CjlmUlJ6M2dHV3VqeDdvcWYyRk56UysydnhtSmxGZ1dFdVh4ekNqeU1SMWl2cjFlbVQ3akZudVRiYzQrZ0NMSkUKUXhqQUl3b0YrRXlwK2dSOFdQQW1ScmVOcU9kSXlSdEF3TkhYRDJJVTc4NHJxZTBkQktoUFVvcDFYYW1OWitvQgpocWNvanAxWlB5ZS9oTXJkUWpPNzJwMENBd0VBQVE9PQotLS0tLUVORCBSU0EgUFVCTElDIEtFWS0tLS0tCg=="
  }
]
``` 

### Verify chain
Goes over all the blocks in the chain and checks whether the chain is valid.
Criteria for a chain to be valid:
- The block's previous hash should actually equal the previous block's data's hash.
- The block's signer shouldn't have signed the previous block.
- The RSA signature can be verified.

**Endpoint:** API_PATH/verify_chain

**Method:** GET

**Parameters:** None

**Response:**
- Status: 200
- Body: A true/false value denoting whether the chain is valid or not.
``` 
{
  "valid": true
}
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
    "sender": "Rahul",
    "recipient": "Akhil",
    "amount": 2,
    "timestamp": "2021-10-25T15:34:01.7255311+05:30"
  },
  {
    "sender": "Aryan",
    "recipient": "Akhil",
    "amount": 1,
    "timestamp": "2021-10-25T15:34:20.3794288+05:30"
  },
  {
    "sender": "0",
    "recipient": "b1e8793688be425d9af1b0682cfb96ca",
    "amount": 1,
    "timestamp": "2021-10-25T15:38:05.7971658+05:30"
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

