hello:
	echo "Hello cint guys!"

build:
	docker build -t pawn_shop .

run:
	docker run -it -d --name=pawn_shop -p 8080:8080  --rm pawn_shop shop_number=10


stop:
	docker rm -f turbine

all: hello stop  build run
