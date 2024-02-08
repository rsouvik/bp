1. Do a git clone
2. Make appropriate changes to conn_params.go
3. sudo docker build -t scrape-app .
4. sudo docker run -p 8080:8080 --name scrape-app6 -it scrape-app

```
```
APIs: (e.g)
1) http://ec2-44-203-160-144.compute-1.amazonaws.com:8080/tokens
2) http://ec2-44-203-160-144.compute-1.amazonaws.com:8080/tokens/bafkreicpcdl32e5l4kusphczsswo3wcrjo7fyt4iktgnhctdnhim6o3xwe
