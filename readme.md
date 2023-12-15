# Neda Divbandi cint Code Challenge
Inventory status can be checked from stdout by each request 
## requirements
we need to have docker installed or just run the project by `go run main.go -shop_number=<number of desired shop>`
to run with docker simply run `docker build -t your-image-name .` and `docker run -p 8080:8080 your-image-name shop_number=<number of desired shop>`
or use `make build` and `make run` the application will be up and running on port 8080 , the shop number can be modified in make file 
### End points
* "/offer"  
* you can send request with data like ``{"code": "PAWN", "offer": 5, "demand": -2} `` and get response from the inventory
### run tests
go test ./... 
