# photoshelf-storage
[![Build Status](https://travis-ci.org/photoshelf/photoshelf-storage.svg?branch=master)](https://travis-ci.org/photoshlef/photoshelf-storage)
[![Coverage Status](http://coveralls.io/repos/github/photoshelf/photoshelf-storage/badge.svg?branch=master)](https://coveralls.io/github/photoshelf/photoshelf-storage?branch=master)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)  

Image upload server with REST API.

## Run Server
### Locally
build executable file.
```bash
git clone https://github.com/photoshelf/photoshelf-storage.git
cd photoshelf-storage
dep ensure
go build -o photoshelf-storage
```

and run
```bash
./photoshelf-storage
```

If you have already set $GOPATH, you can install it with the following command.
```bash
go get -u github.com/photoshelf/photoshelf-storage
```

### Running options
|flag|description            |default |
|----|-----------------------|--------|
|c   |configuration file path|        |
|p   |port number            |1323    |
|t   |storage type           |leveldb |
|s   |storage path           |./photos|

#### configuration file
photoshelf-storage can recognized external file.  
It needs YAML format
```yaml
server:
  port: 8080
storage:
  type: file
  path: /path/to/storage
```
#### storage type
You can use `file` or embedded kvs (`leveldb` or `boltdb`) to store photos.

### Using Docker
```bash
git clone https://github.com/photoshelf/photoshelf-storage.git
cd photoshelf-storage
docker build -t photoshelf/photoshelf-storage .
docker run -p 1323:1323 -v $PWD/photos:/photoshelf/photos photoshelf/photoshelf-storage
```

## CRUD photo
### Create
```bash
curl -X POST http://localhost:1323/photos/ -F "photo=@/path/to/photo"
```

returns 
```json
{
  "Id": "identifier"
}
```

### Read
```bash
curl -X GET http://localhost:1323/photos/:id
```
  
or  
  
Access with browser to `http://localhost:1323/photos/:id`


### Update
```bash
curl -X POST http://localhost:1323/photos/:id -F "photo=@/path/to/new_photo"
```

### Delete
```bash
curl -X DELETE http://localhost:1323/photos/:id
```

## License
MIT License

Copyright (c) 2017 Shunsuke Maeda

See [LICENSE](./LICENSE) file
