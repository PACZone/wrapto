# Wrapto Bridge User Documentation

This documentation contains all that a user needs to interact and work with Wrapto protocol.

## How to Bridge?

In this part we will talk about how to bridge PAC to other supported networks by Wrapto and how to burn wPAC
on other networks and get PAC on Pactus network.

The most straightforward way to do that is using [Wrapto](https://wrapto.app) website. But you are free to make
direct contract calls or sending transactions too.

### PAC to other networks (PAC to wPAC)

To bridge from Pactus to any other network, you have to just make a Pactus transaction (using any wallet) with a memo 
in this structure:

```
DESTINATION_NETWORK_ADDRESS@DESTINATION_NETWORK_ID
```

Example:

```
0xa6a9Def75CA1339Cb3514778948A1D67D826D89A@POLYGON
```

That means you want to bridge the PAC to this address on Polygon network.
(the website will help you to generate this)

You need to make this transaction to Wrapto lock address which is provided on the Wrapto website.

The amount of this transaction will be bridged to your destination address.

Fees are inclusive and explained [here](./protocol(standard).md)

### Other networks to PAC (wPAC to PAC)

To bridge your wPAC to Pactus network, you have to call bridge/burn function of our wPAC contract on other network and 
provide a Pactus address on the call.

You can find contract address on website and you are free to make direct call or use or website.
The EVM contracts ABI is available [here](../../abis/WrappedPac.json).


## Notes & Disclaimers

There is a minimum amount of 1 w/PAC for bridge. The address on contract calls and Pactus transaction MUST be valid.

Any bridge attempt with wrong info that we warned about will be counted as donations. if everything was following the
correct instructions and you got errors and failures, we will return your funds safely back to you!

The Wrapto protocol utilizes a multi-sig model and a warm and cold wallet system with the Pactus foundation. warm waller always keeps 10% of bridge liquidity in itself and the rest will be kept in cold wallet (a multisig wallet) so if you want to bridge more than 10% of whole bridge liquidity you probably need to contact email below and perform it manually.

> For this purpose you can contact us by [hi@dezh.tech](mailto:hi@dezh.tech), make sure you provide all data is needed to track the issue.
