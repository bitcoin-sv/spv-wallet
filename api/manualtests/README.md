## Manual tests helper tool

### Instruction

### First time run

1. Start spv-wallet from code locally
   1. Ensure that in your config.yaml you have `new_transaction_flow_enabled: true` which means that you have v2 experimental enabled.
2. Create user with [admin_create_user_test.go](adminapi/admin_create_user_test.go)
    1. The config/state file [state.yaml](state.yaml) will be created at the first run
    2. At the first start you will get error about missing config options with the link to config/state file
3. I strongly encourage you to create a user at least twice to have sender and recipient
4. After creating user you can play arround with outher tests 
   1. for example, you can try to receive and make some transactions [handle_transactions_test.go](userapi/handle_transactions_test.go)
   2. Receiving and making paymail transactions requires spv-wallet to be exposed on domain.
