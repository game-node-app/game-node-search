# game-node-search  
GameNode's search system, made with ❤️. BLAZINGLY FAST.  

## Reasoning  
GameNode's search is powered by [ManticoreSearch](https://manticoresearch.com/).  
However, the data stored in the ManticoreSearch side and the way the query works wouldn't align with what we want to offer for our clients (web and mobile).  
We needed a service/API that wouldn't get in Manticore's way, and that's why Go was chosen. 
It features a very simple and powerful concurrency model, and is also very easy to learn.  


## Usage  
This API is very simple. It has a single endpoint, which receives a request, transforms it following `Manticore`'s query pattern, sends it to `Manticore`, and then returns a normalized version of the response (with formated date fields, camel case, etc.).  

### Important
If you are a developer trying to contribute to either `game-node-web` or any other frontend client, keep in mind that you don't need to host your own instance of `game-node-search`.  
You can just use our public version (which has all data we have available). Just use the following as the search system's URL:
`https://search.gamenode.com.br`.

## Installation

The easiest way to get started on development is to use Docker and Docker Compose.
Once you have those installed, you can just run:
```shell
make dev
```
**Alternatively**, If you don't have/want `make` installed, run this instead:
```shell
docker compose up --build -d
```  

Which will start a `Manticore` instance listening on ports `9306` (SQL) and `9308` (HTTP).  
`Manticore` is the search engine used by GameNode, and is the only dependency of this project.

Once the `Manticore` image is downloaded and the container started, you can start the server by running: 
```shell
go run main.go
```
This will start the game-node-search API on port `9000`. You can navigate to `localhost:9000/swagger/index.html` to see the API's documentation.  

### Where do i get data?
The Manticore's database is populated by entries in the GameNode's MySQL database.  
This means you need to set up `game-node-server` and `game-node-sync-igdb` to download and save IGDB entries locally.  

Please refer to the respective repositories on how to set them up (both are a `docker compose up` away).  
