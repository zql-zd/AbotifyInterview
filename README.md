# Advertising System PoC

## Introduction

This project demonstrates a basic digital advertising system, similar to those used by online platforms to show ads to users. For those new to advertising systems:

- **How it works**: When a user visits a website/app, the system quickly chooses and displays relevant advertisements. It then tracks two main types of interactions:
  - **Impressions**: When an ad is shown to a user (CPM - Cost Per Mille/thousand impressions)
  - **Clicks**: When a user clicks on an ad (CPC - Cost Per Click)

- **Billing**: Advertisers are charged based on these interactions:
  - They pay a certain rate per 1000 impressions (CPM rate)
  - They pay a certain rate per click (CPC rate)
  - The system tracks their remaining budget to ensure they don't overspend

This proof-of-concept implements these core features in a distributed system that can handle multiple requests simultaneously through load balancing.

## Tech Stack

- Go 1.16+
- MySQL 8.0
- Nginx
- Docker & Docker Compose

## Setup

1. Ensure you have Docker and Docker Compose installed
2. Clone the repository
3. Start the system:

```bash
docker compose up --build
```

The system will start with:

- Nginx load balancer on port 8080
- Three Go server instances (ports 8081, 8082, 8083)
- MySQL database on port 3306
- Demo data automatically loaded

## Environment Variables

### Go Servers

- `DB_USER`: MySQL user (default: "aduser")
- `DB_PASS`: MySQL password (default: "adpass")
- `DB_HOST`: MySQL host (default: "mysql:3306")
- `DB_NAME`: MySQL database name (default: "adserver")
- `PORT`: Server port (different for each instance)

### MySQL

- `MYSQL_ROOT_PASSWORD`: Root password (default: "rootpass")
- `MYSQL_DATABASE`: Database name (default: "adserver")
- `MYSQL_USER`: Application user (default: "aduser")
- `MYSQL_PASSWORD`: Application password (default: "adpass")

## API Endpoints

All endpoints are accessible through the load balancer (port 8080):

- `GET /api/advertisers/{id}` - Get advertiser by ID
- `GET /api/advertisements/{id}` - Get advertisement by ID
- `GET /api/advertisements/{id}/impression` - Record an impression (CPM)
- `GET /api/advertisements/{id}/click` - Record a click (CPC)
- `GET /api/ad/request` - Request for an advertisement

## Database Schema

### Advertiser

- ID (int64)
- Name (string)
- Budget (float64)
- CreatedAt (timestamp)
- UpdatedAt (timestamp)

### Advertisement

- ID (int64)
- AdvertiserID (int64)
- Title (string)
- Content (string)
- CPMCount (int64)
- CPCCount (int64)
- CPMRate (float64)
- CPCRate (float64)
- CreatedAt (timestamp)
- UpdatedAt (timestamp)

## Demo Data

The system comes with pre-loaded demo data including:

- 3 advertisers with different budgets
- 6 advertisements (2 per advertiser)
- Various CPM and CPC rates

## Reset Demo Data

To reset the demo data to its initial state:

```bash
docker compose down -v
docker compose up --build
```

## List of Technical Solutions

Suggested Strategies 
1. Batch Processing: Instead of updating the counts in real-time, accumulate clicks and impressions in a temporary storage (like Redis or in-memory data structures) and periodically update the database in batches. This reduces the number of write operations on the database.

2. Data Aggregation Service: Implement a dedicated microservice for aggregating clicks and impressions. This service can handle in-memory counts and periodically flush these counts to the database. This approach centralizes the counting logic and reduces database load.

3. Asynchronous Processing: Utilize message queues (such as Kafka, RabbitMQ) to handle clicks and impressions asynchronously. The message consumers can then aggregate and update the counts in the database, minimizing the direct impact on the primary database operations.

4. Sharding: Distribute the data across multiple databases or tables (sharding) so that the load is balanced and not all write operations hit the same database or table.

5. Use of NoSQL Databases: Switch to a NoSQL database like Cassandra or MongoDB, which are better suited for write-heavy applications and can handle high volumes of data with eventual consistency.

6. Caching Layer: Implement a caching layer using technologies like Memcached or Redis to handle frequent read operations, thereby reducing the load on the primary database.

## Technical Solution Design

To alleviate frequent database I/O issues in high-traffic environments, I’ve implemented a strategy that leverages batch processing and caching mechanisms. This approach significantly optimizes performance and ensures data consistency across distributed systems. Here’s a detailed explanation of the workflow and technologies involved:

Initial Request Processing and Redis Caching
Upon receiving a request, the system first writes the data to Redis, a high-performance in-memory data store. Redis serves as a temporary cache, absorbing the immediate data write load. This caching layer is crucial for two reasons:

Reduced Latency: Writing to Redis is significantly faster than direct writes to a relational database. It helps in quickly responding to client requests, thereby enhancing user experience.
Load Distribution: By accumulating data in Redis, we distribute the write load more evenly over time, avoiding spikes that could overwhelm the database.
Scheduled Data Synchronization with Cron Jobs
To synchronize the cached data with the primary database, I employ Cron Jobs that run at specified intervals. These jobs execute batch processes that transfer data from Redis to the database in bulk. This method offers several advantages:

Efficiency: Batch processing minimizes the number of database transactions, reducing overhead and improving overall efficiency.
Reliability: By scheduling these operations, we can choose off-peak hours for data synchronization, ensuring minimal impact on the application’s performance and availability.
Transactional Integrity
Both the Redis writes and database synchronization processes are wrapped in transactions. This transactional approach is critical for maintaining consistency across the distributed system for several reasons:

Atomicity: Transactions ensure that each data operation is atomic. Either the entire operation succeeds, or it fails completely, leaving the system in a consistent state.
Isolation: By isolating each transaction, we prevent concurrent operations from interfering with each other, which is particularly important during the batch synchronization process.
Durability: Once a transaction is committed, the changes are permanently applied, ensuring data is not lost even in the event of a system failure.
Conclusion
By combining Redis caching with batch processing and employing transactional mechanisms, this architecture effectively mitigates frequent database I/O issues. It provides a scalable solution that balances performance with consistency, ensuring that high traffic volumes do not compromise the system’s responsiveness or reliability. This approach exemplifies how modern applications can leverage distributed systems and caching strategies to meet the demands of large-scale data management.

## Optimization Strategies

Optimizing Redis Write Operations:
  Utilizing pipelining can reduce network round trips and enhance write efficiency.
  For batch operations, consider using MSET or HMSET for updating multiple fields of a single key.

Batch Processing Data Synchronization:
  Ensure that batch operations do not negatively impact database performance due to the size of a single batch of data. Adjustments may be necessary based on actual conditions.
  Use transactions or batch insert statements (e.g., INSERT INTO ... VALUES (), (), ... ON DUPLICATE KEY UPDATE) to decrease the number of database operations.

Cold and Hot Data Classification:
  Classify data into cold and hot based on the frequency of access, storing only hot data (frequently accessed data) in the cache to optimize the use of cache resources.
  Regular analysis of data access logs or patterns can be used to dynamically adjust which data should be considered hot.

Asynchronous Processing:
  Synchronizing data to the database can be performed asynchronously to avoid blocking the main process and improve response speed.
  Use a message queue to manage the data needing synchronization. Backend services consuming data from the queue to perform database synchronization can further decouple system components, enhancing scalability and maintainability.
  These are basic ideas and examples. Actual implementation may need adjustments based on specific scenarios and performance testing results.

## Sorry

I’ve been deeply interested in the topics of this interview, specifically in the architecture design and implementation for distributed systems and high-traffic projects, which has been the focus of my work over the past year. I’ve been responsible for the iteration and development of the online album at Baidu Maps, dealing with data on the scale of over 1 billion entries. The interview process did not meet my expectations, and I recognize that this was due to my personal circumstances. Due to my current employment and the restrictions of Baidu’s internal network, I was unable to complete all the test questions. This unfortunate situation is also a reflection of my lack of preparation, for which I extend my apologies to you.


