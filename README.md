# json-to-mongo-pipeline

This is just a simple utility for converting a JSON representation of a MongoDB aggregation pipeline into a parameter that can be passed to the `mongosh` `aggregate` function.

Primarily this will convert _id strings to ObjectId(_id) and strings containing ISO dates to ISODate(date)


### Build

How to build this project:

```bash
make
```


### Usage

First, copy the JSON text that you want to convert to your clipboard. Then, run

```bash
./bin/json-to-mongo-pipeline
```

The results will be printed to `stdout` AND written back to your clipboard

