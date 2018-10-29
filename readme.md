# ABI to Soldity Stub Generator

Halfway through the ethstatus hackathon just prior to Devcon IV, I needed my smart contract to interact with another contract but I only had the ABI.

Naturally I did what any self respectin coder would do when unable to find one, I built it.

## installation

`ABIParser` is written in Go and needs to be compiled. Put the executable in the path

## Usage

`ABIParser -input inputfile -contract contract_name -output output_file`

example using gx.json as provided

file : `ABIParser -input gx.json -contract myGX -output gx`

This reads the file gx.json, creates gx.sol

```solidity
pragma solidity ^0.4.24


contract  myGX {
    event OwnershipTransferred(address indexed previousOwner,address indexed newOwner);

    func receiveToken(address _from,uint256 _value) public;
    event GoldPurchased(address buyer,uint256 fiatValue,uint256 goldAmount);

    event GoldSold(address seller,uint256 goldAmount,uint256 fiatAmount);

    func setExchangeRates(uint256 f2g,uint256 g2f) public;
    func setFiatToken(address fiat_) public;
    func setGoldOracle(address _goldOracle) public;
    func setGOLDX(address goldx_) public;
    func transferOwnership(address newOwner) public;
    func fiat_to_gold() public returns (uint256);
    func fiatToken() public returns (address);
    func gold_to_fiat() public returns (uint256);
    func goldOracle() public returns (address);
    func goldx() public returns (address);
    func owner() public returns (address);
    func testBuyGold(uint256 fiatValue) public returns (uint256 goldx_amount);
    func testSellGold(uint256 goldValue) public returns (uint256 fiat_amount);
}
```
