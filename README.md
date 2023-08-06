#  zk-rollup-lite

This is a demo app of simple implementation of a layer2. We can see how layer2 works and what's the meaning of rollup.

## 0. Prerequisite

Make sure you have installed [NVM](https://github.com/nvm-sh/nvm#install--update-script), [Go](https://golang.google.cn/doc/install), [Circom](https://docs.circom.io/downloads/downloads/), [Docker](https://docs.docker.com/get-docker/) before you jump into the following steps.

## 1. Installation

```shell
git clone git@github.com:chentihe/zk-rollup-lite.git

cd zk-rollup-lite
```

## 2. Build the circuit files

### 2.0 Install packages

```shell
~/zk-rollup-lite: cd prover
~/prover: npm install
```

Since installing snarkjs will also install circom, but the version is very old, we need to delete it manually. Otherwise errors will be occurred when we execute the building script.
```shell
~/prover: sudo rm -rf node_modules/circom
```

### 2.1 Generate powersOfTau
You can see [powers of tau ceremony](https://github.com/iden3/snarkjs#1-start-a-new-powers-of-tau-ceremony) for details.
```shell
~/prover: ./powers-of-tau.sh 15
~/prover: ./powers-of-tau.sh 17
```

### 2.2 Compile circuits
The compile script will create build folder on operator directory because layer2 app need circuit files to generate zkp, we need to copy these files into the docker image.

It would take times to generate circuit related files.

- **Compile withdraw circuit**
```shell
~/prover: ./compile-circuit.sh circuits/withdraw.circom trusted_setup/pot15_final.ptau
```
- **Compile processTx circuit**
```shell
~/prover: ./compile-circuits.sh circuits/tx.circom trusted_setup/pot15_final.ptau
```
- **Compile rollupTx circuit**
```shell
~/prover: ./compile-circuits.sh circuits/rollupTx.circom trusted_setup/pot15_final.ptau
```
If you face the permission issue, execute below command to change the permission.
```shell
sudo chmod +x compile-circuits.sh
```
## 3. Start all services
docker compose up command starts containers including postgres, redis, anvil(ethereum node), layer2 app. Except for that, when layer2 app starts, it will deploy necessary contracts(**tx verifier**, **withdraw verifier**, **rollup**) at the same time.
```shell
~/prover: cd ..
~/zk-rollup-lite: docker compose up
```
We can navigate to postgres gui and redis gui to check the database.
- Pgweb: http://localhost:8085
- Redis commander: http://localhost:8081
    - username: root
    - password: qwerty
## 4. APIs
There are apis to query the account and the contract.
### Accounts
- **Query the account info:** http://localhost:8000/api/v1/accounts/{id}
    - id: 1, 2
### Rollup Contract
- **Query the user info:** http://localhost:8000/api/v1/contract/users/{id}
    - id: 1, 2
- **Query the contract balance:** http://localhost:8000/api/v1/contract/balance
- **Query the state root on the contract:** http://localhost:8000/api/v1/contract/root

## 5. Do a rollup
Before running the cli, we should install the packages first.
```shell
~/zk-rollup-lite: cd operator
~/operator: go mod download
```

There are two mock users on [env.yaml](https://github.com/chentihe/zk-rollup-lite/blob/main/operator/config/env.example.anvil.yaml#L1), and a rollup script to execute 2 rollup transactions.

Open a new terminal and navigate to operator folder to execute the rollup script.
```shell
~/operator: ./rollup.sh
```

When you see the terminal of docker containers indicates below message, you can check if the data is updated on both postgres and redis.
```shell
zk-rollup-lite-app-1 | 2023/08/04 03:31:27 Rollup success: 0x3327a5267c77fc272ad754881e9f608869b12fa185f9fb13d1ae170ccef4a18a
```

You will see that the value of **last-inserted** is 5 and the value of **rolluped-txs** is 4 on redis commander, because this app only rollup 2 transactions at once. The 5th transaction will be rolluped once the new transaction is created.
## 6. Try it yourself
First, navigate to cmd folder.
```shell
~/operator: cd cmd
```
There are some command you can try to interact with the demo app.
### 6.1 Deposit
There are only two mock users to choose, 0 for the first user, 1 for the second user.

The default balance for both users is 1000ETH.


```shell
~/cmd: go run main.go deposit --account 0 --amount 1 # 1 for 1ETH
```
### 6.2 Withdraw
Before you withdraw the ethers from layer2, you need to check the balance first, then type how many ethers you want to withdraw.
```shell
~/cmd: curl -s -X GET localhost:8000/api/v1/accounts/1

{
    "AccountIndex": 1,
    "PublicKey": "c433f7a696b7aa3a5224efb3993baf0ccd9e92eecee0c29a3f6c8208a9e81d9e",
    "L1Address": "0x70997970C51812dc3A010C7d01b50e0d17dc79C8",
    "Nonce": 0,
    "Balance": 10000000000000000000
}
```
```shell
~/cmd: go run main.go withdraw --account 0 --amount 0.2
```
### 6.3 Transfer ethers on layer2
In this demo app, you are only able to transfer ethers to the other user. You should notice that every transaction on layer2 may incur the transaction fee to the layer2 platform, hence, we will charge 0.05ETH on every transction.
```shell
~/cmd: go run main.go sendtx --account 0 --amount 0.5
```