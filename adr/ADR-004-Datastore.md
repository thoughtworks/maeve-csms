# Datastore

## Status

Accepted

## Context

The CSMS needs a persistent store to store data that is not transient. This includes:
* Charge station details
* Tokens
* Transactions

The datastore is used for the operational state of the CSMS. It is not designed to be used for analytics or reporting.
The access patterns make it suitable for a wide-variety of NoSQL stores, such as:
* Firestore
* MongoDB
* Cassandra
* DynamoDB

There is no standard API for writing to these stores.

## Decision

The datastore will be abstracted behind a set of interfaces that allow the CSMS to be written to without knowing the
specific implementation. This will allow the datastore to be swapped out for a different implementation in the future.

The initial implementations will be:
* Firestore: the development team currently deploy the CSMS to Google Cloud Platform and Firestore provides a simple
  and cost-effective datastore that is easy to use.
* In-memory: this implementation is designed for unit tests and must not be used in production.

## Consequences

* The current implementation must either be run with an in-memory datastore or a Google Cloud Platform project must be
  created and the Firestore implementation must be used.

### Alternatives considered

SQL databases have a moderately standard API (SQL) that would allow a single implementation to be used across a wide 
variety of providers. There is nothing that prevents us from implementing a SQL datastore in the future, but NoSQL
stores provide a more natural fit for the data that we are storing.