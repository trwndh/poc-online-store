# POC about online store problem

## What I Think

I think the system didn't implement locking stock per item for every request, 
so the stock will keep going reduced and become negative everytime user checkout the item.
In checkout stage, the system did not check or reserve stock per item in the cart. 
Or if it did, it reserving wrong value of stock.
and then finally, they can paid without knowing the availability of their purchased item stock.
Row locking in database is a must for a row that can be edited by many users at (nearly) same time. 


Why? 
Because it will block every transaction before it commited. 
This will prevent newer transaction read uncommited data proceed by previous transaction.

### How stock can be a negative number?
```
Let's say I have Product A, with stock of 100.
Transaction A is modifying Product A, to reduce the stock by 50. 
And in the same time, transaction B is coming and try to modify Item A too, reducing stock by 100.
```
The ideal results are: 
- final stock for Product A must be 0
- transaction B must be failed, so the user can't move to payment page
 
Without Locking, here is what will happen in timeline:

| Sequence | Transaction A                                       | Transaction B                                       |
|----------|-----------------------------------------------------|-----------------------------------------------------|
|    1     | START TRANSACTION                                   |                                                     |
|    2     | SELECT stock FROM product WHERE id = 3, (return 100)  | START TRANSACTION                                   |
|    3     | UPDATE product SET stock = stock - 50 WHERE id = 3  | SELECT stock FROM product WHERE id = 3, (return 100)  |
|    4     | COMMIT, (real stock = 50)                           | UPDATE product SET stock = stock - 100 WHERE id = 3 |
|    5     |                                                     | COMMIT, (real stock = -50)                            |

Assuming the service already has stock reducing logic: 
```
IF stock still eligible to take (request <= stock available), then transaction is eligible to update stock.
```
What happen after Transaction B is committed? 

The answer is, stock become ```-50```. This happens because transaction B read stock before stock A committed update, look seq 3. Stock value
from Read operation in Transaction B is still ``100``, and in service logic, it will defined as eligible to take the stock.

In seq 4, Transaction A is committed stock update: ```100 - 50 = 50 ```. At this point, the real stock is ```50```

Then in seq 5 Transaction B committed update, try reducing 100 from stock, which we knew it already reduced after Transaction A committed the change.
So the result is ``-50``, from ``` 50 - 100 ```.
 
---

## What I purpose
To avoid wrong value from reading stock, I will implement row locking to every row that needs to be updated. 

I will use ```FOR UPDATE``` after ```SELECT``` statement to lock selected row.
then to unlock, just use ```COMMIT``` statement.
Then the timeline will be like this:

| Sequence | Transaction A    (request 50)                                   | Transaction B   (request 100)                       |
|----------|-----------------------------------------------------------------|-----------------------------------------------------|
|    1     | START TRANSACTION                                               |                                                     |
|    2     | SELECT stock FROM product WHERE id = 3 FOR UPDATE, (return 100)   | START TRANSACTION                                   |
|    3     | UPDATE product SET stock = stock - 50 WHERE id = 3              | SELECT stock FROM product WHERE id = 3 FOR UPDATE, (not returning here..)             |
|    4     | COMMIT, (real stock = 50)                                       | waiting..         |
|    5     |                                                                 | return here. get stock = 50, not eligible ( <= 100 ).|
|    5     |                                                                 | ROLLBACK                                  |                                              |

I will Create 2 endpoints to prove this proposal.
```
1. Endpoint WITH locking row, which will return error when stock is already 0
2. Endpoint without locking row, which will leads to negative number of stock
3-... . Endpoint to get product(s), add product and change product stock
```

### Endpoints
Please refer here for full endpoint available and example : https://documenter.getpostman.com/view/9258280/TzsZr8Js

### Scope
This poc is only for proving how to prevent negative stock, and return error in checkout if stock unavailable.
No cart service and no payment service provided in this repo.

### Run service
- Make sure you have docker and docker-compose installed.
- then run these commands:
```
    $ git clone https://github.com/trwndh/poc-online-store.git
    $ cd poc-online-store
    $ docker-compose up
```

### Testing
- To run concurrent simulation, look at ```poc``` folder, and open folder you want to test.
- There are 2 folders inside ```poc```  folder.
- If you want to see negative number case, run ```drain_stock``` binary or ```go run drain_stock.go``` in folder ```drain_stock_without_lock```
- If you want to see my solution case, run ```drain_stock```binary or ```go run drain_stock.go``` in folder ```drain_stock```
- Inside each folder, i've made binary file ready to execute for windows, linux, and osx (darwin). There's also go file if you want to see.
- To change request payload you can edit file ```request.json``` in each folder
- Hit ```localhost:9999/api/v1/product/``` before and after drain_stock execution to compare stock change.

### Tech used

```
- Jaeger as opentracing tracer. open http://localhost:16686/ to access Jaeger UI. and select
- MySQL
- Golang 
```

