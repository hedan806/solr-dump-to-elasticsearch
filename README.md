# solr-dump-to-elasticsearch
Export SOLR documents efficiently to Elasticsearch with cursors.

### 1. build
```bash
make build-linux
```

### 2. run
```bash
./solrdump2es_unix dump -s {src_index} -t {target_host}-i {target_index} -c {consumer_size} -f {unique_field}
```

> EXAMPLE
```bash
./solrdump2es_unix dump -s http://127.0.0.1:8983/solr/foo -t http://127.0.0.1:9200 -i foo -c 1 -f id
```
