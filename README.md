# Shopping cart kata


Most businesses sell their products online by means of an e-commerce site which provides a real-life experience taking advantage of the shopping cart capability.


## Overview
A company wants to sell online the following products:

```
Code         | Name                     |  Price
--------------------------------------------------
VOUCHER      | AcME Voucher             |   5.00 €
TSHIRT       | AcME T-Shirt             |  20.00 €
MUG          | AcME Coffee Mug          |   7.50 €
```

Two different offers will be used to boost sales:
 - buy one get one free (aka two for the price of one)
 - multibuy discount (buying at least a quantity of a product its unit price is discounted)

The promotional strategy planned by the company is to apply:
 - buy one get one free to `VOUCHER` items
 - multibuy discount to `TSHIRT` items for a quantity of 3 or more and reducing the unit price to 19.00 €

**Examples**
```
Items: VOUCHER, TSHIRT, MUG
Subtotal: 32.50€

Items: VOUCHER, TSHIRT, VOUCHER
Subtotal: 25.00€

Items: TSHIRT, TSHIRT, TSHIRT, VOUCHER, TSHIRT
Subtotal: 81.00€

Items: VOUCHER, TSHIRT, VOUCHER, VOUCHER, MUG, TSHIRT, TSHIRT
Subtotal: 74.50€
```


## Goal

Design and implement a client/server solution supporting the following operations:
 - create a new cart
 - add a product to a cart
 - get the cart subtotal
 - remove a cart


## Requirements
 - Do not use any external databases of any kind
 - Support multiple concurrent remote invocations and operations on the same cart and on different carts
 - Use the protocol exposed by the server in the client to translate user input to remote calls
 - Allow easy and frequent changes to the set of rules in the promotional strategy
 - Write production-ready code that is easy to maintain and evolve


## Deliverables
 - Description of the solution justifying choices
 - Proof of concept of the proposed solution that can be built and executed on a UNIX/Linux environment