1. Do a git clone
2. Make appropriate changes to conn_params.go
3. sudo docker build -t scrape-app .
4. sudo docker run -p 8080:8080 --name scrape-app -it scrape-app
