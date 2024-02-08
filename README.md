1. Do a git clone
2. Copy ipfs_cids.csv into current directory
3. Make appropriate changes to conn_params.go
4. sudo docker build -t scrape-app .
5. sudo docker run -p 8080:8080 --name scrape-app -it scrape-app
