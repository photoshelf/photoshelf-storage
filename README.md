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
go build
```

and run
```bash
./photoshelf-storage
```

### Using Docker
```bash
docker build -t photoshelf/photoshelf-storage .
docker run -p 1323:1323 photoshelf/photoshelf-storage
```

## License
MIT License

Copyright (c) 2017 Shunsuke Maeda

See [LICENSE](./LICENSE) file
