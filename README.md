# Simple Go Wallet Service

## Overview

The **Go Wallet Service** is a simple RESTful API that allows users to manage their wallet transactions. Users can deposit, withdraw, and transfer money, as well as check their wallet balance and view transaction history.

## User Stories

- **As a user**, I want to deposit money into my wallet.
- **As a user**, I want to withdraw money from my wallet.
- **As a user**, I want to send money to another user.
- **As a user**, I want to check my wallet balance.
- **As a user**, I want to view my transaction history.

## RESTful API Endpoints

### 1. Deposit Money

- **Endpoint**: `POST /wallet/{userId}/deposit`
- **Description**: Deposit money into the specified user's wallet.
- **Response**:
    ```json
    {
    "message": "Transfer Successful",
    "status": "ok"
    }
    ```

### 2. Withdraw Money

- **Endpoint**: `POST /wallet/{userId}/withdraw`
- **Description**: Withdraw money from the specified user's wallet.
- **Response**:
    ```json
    {
    "message": "Transfer Successful",
    "status": "ok"
    }
    ```

### 3. Transfer Money

- **Endpoint**: `POST /user/2/transfer`
- **Description**: Transfer money from one user to another.
- **Response Body**:
    ```json
    {
    "message": "Transfer Successful",
    "status": "ok"
    }
    ```

### 4. Get Wallet Balance

- **Endpoint**: `GET /wallet/{userId}/balance`
- **Description**: Get the specified user's wallet balance.
- **Response**:
    ```json
    {
    "balance": 28930
    }
    ```

### 5. Get Transaction History

- **Endpoint**: `GET user/{userId}/transactions`
- **Description**: Get the specified user's transaction history.
- **Response Body**:
    ```json
    [{"Type":"transfer","Amount":570,"Receiver":"matt","TransactionDate":"2024-12-20T09:58:58.755195Z"},{"Type":"withdraw","Amount":500,"Receiver":"clare"}]
    ```

## Example API Requests (Postman)

You can use Postman to test the API endpoints. Here are some example requests:

### 1. Deposit Money
- **Method**: POST
- **URL**: `http://localhost:8080/wallet/userId/deposit`
- **Body** (JSON):
    ```json
    {
        "userId":2,
        "amount":2000
    }
    ```

### 2. Withdraw Money
- **Method**: POST
- **URL**: `http://localhost:8080/wallet/userId/withdraw`
- **Body** (JSON):
    ```json
    {
        "userId":2,
        "amount": 500.00
    }
    ```

### 3. Transfer Money
- **Method**: POST
- **URL**: `http://localhost:8080/user/userId/transfer`
- **Body** (JSON):
    ```json
     {
        "userId":2,
        "receiverUserId":1,
        "amount": 570.00
    }
    ```

### 4. Get Wallet Balance
- **Method**: GET
- **URL**: `http://localhost:8080/wallet/user1/balance`
- **Body** (JSON):
    ```json
    {
        "userId": "user1"
    }
    ```

### 5. Get Transaction History
- **Method**: GET
- **URL**: `http://localhost:8080/user/userId/transactions`

## Setup

To set up the project, please ensure you have the following prerequisites installed:

### Requirements

1. **Go v1.23**
   - Download and install Go from the official website: [golang.org](https://golang.org/dl/).
   - Verify the installation by running:
     ```bash
     go version
     ```

2. **Docker**
   - Install Docker by following the instructions at: [docker.com](https://docs.docker.com/get-docker/).
   - Check that Docker is running correctly:
     ```bash
     docker --version
     ```

3. **PostgreSQL**
   - Install PostgreSQL using the instructions provided on the official site: [postgresql.org](https://www.postgresql.org/download/).
   - After installation, ensure PostgreSQL is running:
     ```bash
     psql --version
     ```
### **Docker** Setup

This will help set up the Docker container where the project will be stored.

## Setup Steps

### 1. Navigate to the folder of the project

```bash
cd {path}\go-wallet-service
```

### 2. Build the Docker container using the Dockerfile and docker-compose.yml file

```bash
docker-compose up --build
```

You should see the server being started if LOG_LEVEL is in DEBUG mode
```
{"level":"debug","time":"2024-12-22T16:14:17Z","message":"Connected to Server"}
```

### 3. Complete

Check if the docker containers have been created, you should see two docker containers

```bash
docker ps
```

### **PostgreSQL** Setup

This will help set up a PostgreSQL database to be used.

## Setup Steps

### 1. Create tables in the container

Open your terminal and create the PostgreSQL tables in go-wallet database, using the container id for 

```bash
docker exec -it {containerid} psql -U postgres -d go-wallet
```

### 2. Run the queries

#### Create tables

You can now run your SQL queries. For example, to create a new database:

***Create users***
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    creation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

***Add users***
```sql
INSERT INTO  users  (username) VALUES
('user1'),
('user2'),
('user3');
```

***Create wallets***
```sql
CREATE TABLE wallets (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users (id) ON DELETE CASCADE,
    balance INTEGER NOT NULL,
    creation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

***Add wallets***
```sql
INSERT INTO wallets (user_id, balance) VALUES 
(1,100), 
(2,200), 
(3,0);
```

***Create oauth***
```sql
CREATE TABLE oauth (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users (id),
    token VARCHAR(255) NOT NULL
);
```

***Add oauth***
```sql
INSERT INTO oauth (user_id, token) VALUES 
(1,'G5gkR3FEQacVfjkB8tLzO8JiUhhIVv6qq'), 
(2,'w7Iyg4TxMhO3PXMFQi0Hpp5sZcTjF6o6y'), 
(3,'xGc7WAGcYZmXoBijkzvXwFy4XLCPIPFKm');
```

### 3. Exit the command prompt

Once the tables have been created, exit the postgres terminal 

```bash
exit
```

## Test the APIs

### 1. Navigate to the go-wallet-service container
```bash
docker exec -it {containerid} bash
```

### 2. GET balance API

```bash
curl -X GET "http://localhost:8080/wallet/2/balance" -H "Authorization: Bearer w7Iyg4TxMhO3PXMFQi0Hpp5sZcTjF6o6y" -H "Content-Type: application/json" -d '{"userId":2}'
```
Response:

```json
{"balance":200}
```

### 2. POST Deposit API

```bash
curl -X POST "http://localhost:8080/wallet/2/deposit" -H "Authorization: Bearer w7Iyg4TxMhO3PXMFQi0Hpp5sZcTjF6o6y" -d '{"user_id":2,"amount":500}'
```
Response:

```json
{"message":"Deposit Successful","status":"ok"}
```

### 3. POST Withdraw API

```bash
curl -X POST "http://localhost:8080/wallet/2/withdraw" -H "Authorization: Bearer w7Iyg4TxMhO3PXMFQi0Hpp5sZcTjF6o6y" -d '{"user_id":2,"amount":50}'
```
Response:

```json
{"message":"Withdraw Successful","status":"ok"}
```

### 4. POST Transfer API

```bash
curl -X POST "http://localhost:8080/user/2/transfer" -H "Authorization: Bearer w7Iyg4TxMhO3PXMFQi0Hpp5sZcTjF6o6y" -d '{"user_id":2,"receiver_user_id":3,"amount":50}'
```
Response:

```json
{"message":"Transfer Successful","status":"ok"}
```
Sending an amount greater than balance
```bash
curl -X POST "http://localhost:8080/user/2/transfer" -H "Authorization: Bearer w7Iyg4TxMhO3PXMFQi0Hpp5sZcTjF6o6y" -d '{"user_id":2,"receiver_user_id":3,"amount":10000}'
```
Response:

```json
{"error":"operation_not_permitted","error_description":"The transfer amount is greater than wallet balance","code":403}
```

### 5. GET Transactions API

```bash
curl -X GET "http://localhost:8080/user/2/transactions" -H "Authorization: Bearer w7Iyg4TxMhO3PXMFQi0Hpp5sZcTjF6o6y"
```
Response:

```json
[{"Type":"transfer","Amount":10000,"Receiver":"user3","TransactionDate":"2024-12-22T16:48:18.101066Z"},{"Type":"transfer","Amount":50,"Receiver":"user3","TransactionDate":"2024-12-22T16:47:50.448901Z"},{"Type":"withdraw","Amount":50,"Receiver":"user2","TransactionDate":"2024-12-22T16:44:52.587845Z"},{"Type":"deposit","Amount":500,"Receiver":"user2","TransactionDate":"2024-12-22T16:43:02.193361Z"}]
```


## Further Development

| Task   | Description   |
|------------|------------|
| Add more OAuth features | Fully implement an automatic token refresher and update the token for a the given user when accessed, upon expiry |
| Add more User data | User data can be enriched with various statuses, such as active, inactive, deleted, and so on. |
| Optimize performance | Improve performance with caching strategies using REDIS |
| Code practices | Optimize code structure to follow best practices |

## Others

| Notes   | Comment   |
|------------|------------|
| Hours spent | Due to my work commitments and setup challenges, I was only able to dedicate approximately 25 hours to the project. |
| As simple as possible | Given the limited hours commitable, I aimed to keep the project setup as straightforward as possible. |
| Lesson learned | We use Go for creating microservices, but not for REST APIs. Regardless of the outcome, this has been an enjoyable experience. |