# Wrapto Protocol V1

In this document, we will talk about how Wrapto protocol V1 works.

## Semi-centralized model

The Wrapto V1 is a semi-centralized protocol because of the lack of smart contract support in the Pactus blockchain. Minting new wPAC on other networks is centralized same as unlocking locked PACs. It's planned to make it fully decentralized after Pactus updates.

## Workflow

Here is a simple workflow of how Wrapto works:

> [!IMPORTANT]
> The term listener is the same as an actor in this [document](./design.md).

```mermaid
stateDiagram-v2
    Pactus: Pactus network
    Polygon: Polygon network
    PactusListener: Pactus Listener
    PolygonListener: Polygon Listener

    note right of PactusListener: Listen bridge TXs, sending to Polygon Listener 
    note right of PolygonListener: Listen bridge Contract calls, sending to Pactus Listener 

    PactusListener --> Pactus
    PactusListener --> PolygonListener
    PolygonListener --> Polygon
    PolygonListener --> PactusListener
```

> [!NOTE]
> This mermaid diagram needs improvements.

Based on the chart we saw, a module listens to each transaction with this structure on the Pactus network:

The memo:
```
DESTINATION_NETWORK_ADDRESS@DESTINATION_NETWORK_ID
```

Example:

```
0xa6a9Def75CA1339Cb3514778948A1D67D826D89A@POLYGON
```

Amount:

more than 1 PAC

Destination:

The Wrapto lock address on the Pactus network. it's available at [Wrapto website](https://wrapto.app).

Later Pactus listener finds the destination based on the memo and sends a message to the specified listener ([actors](./design.md)). Invalid memos will be counted as donations.

Other networks will listen to token smart contracts and check burn/bridge functions, if there are any, they will send a message to the Pactus listener to unlock PAC on the Pactus network. The Pactus unlock address must be provided in the contract call otherwise its burning tokens and locked PACs on Pactus will be counted as fees.

### Rate and Health

The rate of PAC and wPAC is always 1:1 constantly. To check whether the bridge is healthy or not, you can check the `l >= (m + f)`, where the `l` is the balance of the Wrapto lock address in the Pactus network and `m` is the total wPAC minted in all of the networks and f is the fee needed to bridge all tokens on other networks to Pactus in the worst case. If it was true, the bridge is healthy.

> To calculate f: divide the total of wPAC by 2. `t = m / 2` and calculate the fee of t transactions with 1 PAC amount on the Pactus network.
That means if you make most bridges you can with 2 PACs the lock address has enough fee to pay for it and unlock tokens. which has a very low possibility of happening. more stuff like the fee for minting wPAC is needed in Wrapto service, but here we considered only Pactus. A decentralized version of Wrapto will make this health check much easier.  

Fees are exclusive for example consider fee 1 PAC (find the fee detail [here](./design.md)), if you bridge 10 PAC to polygon you will receive 9 wPAC then. 10 PAC is the balance of the lock address and can be collected by the team without hurting the health of bridge.

So, sending invalid TX in the lock address or burning wPAC means making more balance on the lock address with no token on the other networks and they are fees and donations.

## Structures

Here we explain data structures used in Wrapto.

## Order

Each request to unlock PAC or mint wPAC will be wrapped in a structure called Order (bridge order), listeners can execute this orders.

Here is Order example in Golang:

```go
type Order struct {
	// * unique ID on Wrapto system.
	ID string

	// * transaction or contract call that user made on source network.
	TxHash string

	// * address of receiver account on destination network.
	Receiver string

	// * address of sender on source network (account that made bridge transaction).
	Sender string

	// * amount of PAC to be bridged, **including fee**.
	amount amount.Amount

	// * status of order on Wrapto system.
	Status Status

	// * type of bridge.
	BridgeType BridgeType
}
```

You can find full detail [here](../../types/order).

## Message

Each Order needs to be transfer between listeners, we send Order and other data between actors in Message structure, here is Golang example:

```go
type Message struct {
	To      bypass.Name
	From    bypass.Name
	Payload *order.Order
}
```

You can find full detail [here](../../types/message).
