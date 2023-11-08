# game-node-search
GameNode's search system, made with ❤️. BLAZINGLY FAST.

## Reasoning
GameNode's search is powered by [ManticoreSearch](https://manticoresearch.com/).  
However, the data stored in the ManticoreSearch side and the way the query works wouldn't align with what we want to offer for our clients (web and mobile).  
That's why we need a service/API that:  

1. Wouldn't get in Manticore's way
Manticore is fast. Very fast. We needed a service that wouldn't get in it's way, by increasing latency or doing expensive computations.
That's why Go was chosen. It has an incredible (and simple) concurrency model, which allows our frontend clients to send many requests/sec to this API.

2. Allows our's clients to send requests following a specific query model (DTO)

## Usage
This API is very simple. It has a single endpoint, which receives a request, transforms it following Manticore's query pattern, sends it to Manticore, and then returns a normalized version of the response (with formated date fields, camel case, etc.).
