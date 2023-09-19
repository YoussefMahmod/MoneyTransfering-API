# MoneyTransfering-API
This is an API for transfer money between two accounts using GOLang.

1. **Database Design(Used in-memory DS to mimic database)**
  - `account` table: id (int, PK), name (string), balance (decimal), created_at (timestamp), updated_at (timestamp).
  - `transaction` table: id (int, PK), from_account_id (int, FK -> account.id), to_account_id (int, FK -> account.id), amount (decimal), created_at (timestamp).

2. **Models**
  - Built Structs for `Account` and `Transaction` corresponding to the database tables.
  - Expose Only interfaces instead of the actual structs.

3. **Data Store**
  - Mimic models and indices using DSA(Datastructures and Algorithms).
  - Built `shardedMap` for more concurrency performance.
  - Handled concurrency (read and write locks)

4. **Services**
  - Splited bulks into chunks and assign goroutine to increase the performance.
  - Used goroutine for inserting each element in chunk(As it now sharded).

5. **APIs**
  - Split the API into two resources `/accounts` and `/transactions`.
  - Each one handle it's own flow.
  - Wrote Integration tests for all endpoints.
    - you can run it using `go test -v ./api/api_tests`
    - you can run the server using `go run main.go`
  - Documented the API [HERE](https://documenter.getpostman.com/view/25231966/2s9YC8upsN)


Suggestion:
  - If we want to add getting transactions between two dates(Range Search) we can do that with 3 different options.
    1. AVL, Balanced BSTs.
    2. B-Trees.
    3. Interval trees.
  - I believe that setting a monitoring service for the API is a good idea(Performance and Error).

TODO:
- Design Diagrams for API flow(Sequence Diagram for example)