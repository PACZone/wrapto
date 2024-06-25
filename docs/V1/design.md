# Code design and structure

In this document we will explain The Wrapto bridge V1 software design.

[!NOTE]
> We suggest to read [protocol](./protocol(standard).md) document before this document.


[!NOTE]
> The actor, side or listener is used in this document interchangeably.

## Actor model

The first thing to be understand in the Wrapto is that each part of software is responsible to listen a blockchain and its events to take correct action if needed. We implemented each of them as an actor in [actor model](https://en.wikipedia.org/wiki/Actor_model) paradigm.

[!IMPORTANT]
> Currently, state of each actor is same and all of them using a single shared database which is not following actor model rules. In the next versions of the software we will launch and run a separated actor specified to database and logs, other actors can later send anything they want to store in database to this actor.

### Manager

All actors can only talk with an actor called Manager and Manager will forward the [message](../../types/message/message.go) to final destination. Other actors is called sides (sides of a bridge).

All messages are fire-and-forget and forwarded. Manager have channel with all actors.

### Sides

Each side is a listener and bridge for a blockchain. For example Pactus side listen on Pactus lock address, if it find any bridge transaction (explained in [protocol](./protocol(standard).md) docs in detail.) it will make an [Order](../../types/order/order.go) and send it to destination actor to bridge it.

Lets explain it with an example, Pactus side finds a valid transaction that wants to bridge PAC coin, in the TX memo there is an address and network ID (consider it's polygon network in our case), Pactus actor/side will wrap this info + transaction info in an Order structure and sends it to polygon actor/side, the polygon side will execute this Order and mint wPAC on Polygon network.

Other networks will check the smart contract event and transactions and send burn/bridge function calls to Pactus actor.

## Fees

The Wrapto fee model is percentage base and contains a min and max fee defined [here](../../types/params/params.go). This fee is used for bridge maintain and operation and more importantly to prevent spam bridges and DDoS attacks.

The percentage is 0.5% of the amount that user wants to bridge, if it was less than min amount which is 1 PAC the fee will be set to 1 PAC and if it was more than 5 PAC the will be set to 5 PAC otherwise the 0.5% of the amount will be used. In this case most of normal usages need a fee of 1 PAC.
