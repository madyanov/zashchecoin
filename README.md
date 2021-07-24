**Zashchecoin** (Russian: Защекоин) – cryptocurrency, popular in the field of IT in the post-Soviet states.

*Currently in development, but I can give you some zashchecoins if you ask.*

## What's done

- [x] Blockchain
- [x] Mining
- [x] P2P
- [x] HTTP API
- [ ] Persistence
- [ ] Transactions
- [ ] Wallet

## How to use

CLI parameters:

```
  -api-port int
    	HTTP API port (default 8080)
  -peers string
    	Path to peers file (default "peers.txt")
  -srv-port int
    	P2P server port (default 8081)
```

HTTP API methods:

- `GET /mine` – mine a block and broadcast it to the peers.