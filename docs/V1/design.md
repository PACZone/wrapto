# Code design and structure

In this document, we will explain The Wrapto bridge V1 software design.

> [!NOTE]
> We suggest to read [protocol](./protocol(standard).md) document before this document.


> [!NOTE]
> The actor, side, or listener is used in this document interchangeably.

## Actor model

The first thing to be understood in the Wrapto is that each part of the software is responsible to listen a blockchain and its events to take correct action if needed. We implemented each of them as an actor in [actor model](https://en.wikipedia.org/wiki/Actor_model) paradigm.

> [!IMPORTANT]
> Currently, state of each actor is same and all of them use a single shared database which is not following actor model rules. In the next versions of the software we will launch and run a separate actor specified to the database and logs, other actors can later send anything they want to store in the database to this actor.

### Manager

All actors can only talk with an actor called the Manager and the Manager will forward the [message](../../types/message/message.go) to the final destination. Other actors are called sides (sides of a bridge).

All messages are fire-and-forget and forwarded. Managers have a channel with all actors.

### Sides

Each side is a listener and bridge for a blockchain. For example, the Pactus side listens on the Pactus lock address, and if it finds any bridge transaction (explained in [protocol](./protocol(standard).md) docs in detail.) it will make an [Order](../../types/order/order.go) and send it to destination actor to bridge it.

Let's explain it with an example, Pactus side finds a valid transaction that wants to bridge the PAC coin, in the TX memo there is an address and network ID (consider it's the Polygon network in our case), Pactus actor/side will wrap this info + transaction info in an Order structure and sends it to polygon actor/side, the polygon side will execute this Order and mint wPAC on Polygon network.

Other networks will check the smart contract event and transactions and send burn/bridge function calls to the Pactus actor.

## Fees

The Wrapto fee model is percentage-based and contains a min and max fee defined [here](../../types/params/params.go). This fee is used for bridge maintenance and operation and more importantly to prevent spam bridges and DDoS attacks.

The percentage is 0.5% of the amount that the user wants to bridge, if it is less than the minimum amount which is 1 PAC the fee will be set to 1 PAC and if it is more than 5 PAC the will be set to 5 PAC otherwise the 0.5% of the amount will be used. In this case, most of the normal usages need a fee of 1 PAC.
