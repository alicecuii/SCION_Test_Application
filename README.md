# SCION_Test_Application
SCION client/server application using PAN library 

First, run ```scion address``` to find out your machine's SCION address;

Open two terminals, one is for client, the other is for server;

On the server side, run ```go run hello.go -listen 127.0.0.1:1234```

One the client side, run ```go run hello.go -remote *scion_address*,[127.0.0.1]:1234```, then input the message you want to send to the server.
