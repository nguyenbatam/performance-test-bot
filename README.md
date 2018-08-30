## Performance Test Bot with version Test.
## Config masternode
../build/bin/tomo --datadir node1/ --syncmode 'full' --port 30311 --rpc --rpcaddr 'localhost' --rpcport 8501 --rpcapi 'personal,db,eth,net,web3,txpool,miner' --bootnodes 'enode://7e7ae37981dbf8de1a7206eb8bed4066b8f4afe1dd4eeb74dc155f2ab69d8a3635e2f8a2293f942af910484cb73e58cff2c57ba5f28411bdec966e4ce1a6c732@127.0.0.1:30310' --networkid 89 --gasprice '1' --unlock '0x21292d56e2a8de3cc4672db039aaa27f9190b1f6' --password node1/password.txt --mine  --txpool.globalqueue 10000000 --txpool.globalslots 10000000 --txpool.accountqueue 100000000 --txpool.accountslots 10000000 --targetgaslimit 1000000000 > logs/node1.txt 2>&1
## or with new version use toml config file
../build/bin/tomo --config node2/config.toml  --txpool.globalqueue 10000000 --txpool.globalslots 10000000 --txpool.accountqueue 100000000 --txpool.accountslots 10000000 --targetgaslimit 1000000000 > logs/node2.txt
##cmd run bot
./bot -password 123456 -n 10 -req 10000 -url "http://localhost:8501" -key-file /home/tamnb/_projects/src/github.com/ethereum/go-ethereum/devnet/node1/keystore/UTC--2018-07-04T07-32-41.529884508Z--21292d56e2a8de3cc4672db039aaa27f9190b1f6