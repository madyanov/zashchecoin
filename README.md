**Zashchecoin** (Russian: Защекоин) – cryptocurrency, popular in the field of IT in the post-Soviet states.

*Currently in development, but I can give you some zashchecoins if you ask.*

## What's done

- [x] Blockchain
- [x] Mining
- [x] P2P
- [x] HTTP API
- [x] Persistence
- [ ] Transactions
- [ ] Wallet

## How to use

CLI parameters:

```
  -api-port int
    	HTTP API port (default 8080)
  -bc string
    	Path to the blockchain file (default "bc.dat")
  -p2p-port int
    	P2P server port (default 8081)
  -peers string
    	Path to the peers file (default "peers.txt")
```

HTTP API methods:

- `GET /mine` – mine a block and broadcast it to the peers.