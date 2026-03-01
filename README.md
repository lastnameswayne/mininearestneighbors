# mininearestneighbors

HNSW (Hierarchical Navigable Small World) approximate nearest neighbor search implemented from scratch in Go. I use it to find clothes that fit me based on body measurements.

## What is HNSW

HNSW is a graph based algorithm for approximate nearest neighbor search. It builds a multi layer graph where the top layers have fewer, longer range connections and the bottom layers are dense. Search starts at the top layer and greedily descends, narrowing in on the nearest neighbors.

This is the same algorithm that Redis, Pinecone, pgvector, etc use under the hood.

## How it works

Given a query vector (e.g. chest=63, length=73, shoulder=55, sleeve=67), the algorithm:

1. Enters at the top layer through the entrance point
2. Greedily moves to the closest neighbor at each layer
3. Drops down to the next layer using that neighbor as the new entry point
4. At the bottom layer, runs a broader search to return the k nearest results

All data structures are implemented from scratch: graph (adjacency list), priority queue, min heap, set, vector operations.

## Structure

```
src/
  hnsw/           HNSW algorithm: insert, search, layer selection, serialization
  graph/          adjacency list graph with add/remove/neighbors
  priorityqueue/  priority queue backed by min heap
  heap/           min heap implementation
  vector/         euclidean distance, vector operations
  set/            generic set
  driver/         MongoDB integration to fetch product measurements
  main.go         entry point
```

## Example output

Searching for long sleeves with measurements chest=63, length=73, shoulder=55, sleeve=67:

```
item 0: 14067131S   distance: 2.83
item 1: 13845011XXL distance: 2.83
item 2: 14064871XL  distance: 2.83
...
item 47: 13834881XL distance: 1.41
item 48: 13834871XXL distance: 1.41
item 49: 13984121XL distance: 1.41
```

## Running

Requires a MongoDB instance with clothing measurement data.

```bash
go run src/main.go
```
