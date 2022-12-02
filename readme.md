# Pismo Transactions #

- Service to manage accounts and transactions.

## How to run application ? ##

- To run application, follow next steps:

```
$ cp .env-example .env
```

```
$ docker-compose up
```

## CURL to call api to create accounts ##

``` 
$ curl --location --request POST 'http://127.0.0.1:3000/accounts' \
--header 'Content-Type: application/json' \
--data-raw '{
    "document_number":"19345678910"
}'
``` 

## CURL to call api to create transactions ##

``` 
$ curl --location --request POST 'http://127.0.0.1:3000/transactions' \
--header 'Content-Type: application/json' \
--data-raw '{
    "account_id":1,
    "operation_type_id":4,
    "amount":1
}'
``` 

## CURL to call api to get account ##

``` 
$ curl --location --request GET 'http://127.0.0.1:3000/accounts/1' \
--header 'Content-Type: application/json'
``` 

## Tests ##

- To run tests, follow steps:

``` 
make test/cov
``` 