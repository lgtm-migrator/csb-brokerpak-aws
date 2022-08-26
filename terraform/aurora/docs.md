### What is Amazon Aurora?

Amazon Aurora (Aurora) is a fully managed relational database engine that's compatible with MySQL and PostgreSQL.
You already know how MySQL and PostgreSQL combine the speed and reliability of high-end commercial databases with
the simplicity and cost-effectiveness of open-source databases. The code, tools, and applications you use today with
your existing MySQL and PostgreSQL databases can be used with Aurora. With some workloads, Aurora can deliver up to five
times the throughput of MySQL and up to three times the throughput of PostgreSQL without requiring changes to most
of your existing applications.

Aurora includes a high-performance storage subsystem. Its MySQL- and PostgreSQL-compatible database engines are
customized to take advantage of that fast distributed storage. The underlying storage grows automatically as needed.
An Aurora cluster volume can grow to a maximum size of 128 tebibytes (TiB). Aurora also automates and standardizes
database clustering and replication, which are typically among the most challenging aspects of database configuration
and administration.

Aurora is part of the managed database service Amazon Relational Database Service (Amazon RDS). Amazon RDS is a web
service that makes it easier to set up, operate, and scale a relational database in the cloud. If you are not already
familiar with Amazon RDS, see the Amazon Relational Database Service User Guide.

The following points illustrate how Aurora relates to the standard MySQL and PostgreSQL engines available in Amazon RDS:

* You choose Aurora as the DB engine option when setting up new database servers through Amazon RDS.

* Aurora takes advantage of the familiar Amazon Relational Database Service (Amazon RDS) features for management and
administration. Aurora uses the Amazon RDS AWS Management Console interface, AWS CLI commands, and API operations
to handle routine database tasks such as provisioning, patching, backup, recovery, failure detection, and repair.
* Aurora management operations typically involve entire clusters of database servers that are synchronized through
replication, instead of individual database instances. The automatic clustering, replication, and storage allocation
make it simple and cost-effective to set up, operate, and scale your largest MySQL and PostgreSQL deployments.

You can bring data from Amazon RDS for MySQL and Amazon RDS for PostgreSQL into Aurora by creating and restoring
snapshots, or by setting up one-way replication. You can use push-button migration tools to convert your existing
Amazon RDS for MySQL and Amazon RDS for PostgreSQL applications to Aurora.

### Amazon Aurora DB clusters

An Amazon Aurora DB cluster consists of one or more DB instances and a cluster volume that manages the data for those
DB instances. An Aurora cluster volume is a virtual database storage volume that spans multiple Availability Zones,
with each Availability Zone having a copy of the DB cluster data. Two types of DB instances make up an Aurora DB cluster:

* Primary DB instance – Supports read and write operations, and performs all of the data modifications to the cluster
volume. Each Aurora DB cluster has one primary DB instance.
* Aurora Replica – Connects to the same storage volume as the primary DB instance and supports only read operations.
Each Aurora DB cluster can have up to 15 Aurora Replicas in addition to the primary DB instance. Maintain high
availability by locating Aurora Replicas in separate Availability Zones. Aurora automatically fails over to an
Aurora Replica in case the primary DB instance becomes unavailable. You can specify the failover priority for
Aurora Replicas. Aurora Replicas can also offload read workloads from the primary DB instance.


Note
The preceding information applies to all the Aurora clusters that use single-master replication. These include
provisioned clusters, parallel query clusters, global database clusters, serverless clusters, and all MySQL
8.0-compatible, 5.7-compatible, and PostgreSQL-compatible clusters.

Aurora clusters that use multi-master replication have a different arrangement of read/write and read-only DB instances.
All DB instances in a multi-master cluster can perform write operations. There isn't a single DB instance that performs
all the write operations, and there aren't any read-only DB instances. Therefore, the terms primary instance and Aurora
Replica don't apply to multi-master clusters. When we discuss clusters that might use multi-master replication,
we refer to writer DB instances and reader DB instances.

The Aurora cluster illustrates the separation of compute capacity and storage. For example, an Aurora configuration with
only a single DB instance is still a cluster, because the underlying storage volume involves multiple storage nodes
distributed across multiple Availability Zones (AZs).


### Amazon Aurora versions

Amazon Aurora reuses code and maintains compatibility with the underlying MySQL and PostgreSQL DB engines. However, 
Aurora has its own version numbers, release cycle, time line for version deprecation, and so on. 


### Amazon Serverless v1 vs Amazon Serverless v2

#### What is serverless v1?

At the time of becoming generally available in 2018, it was called serverless v1 and had some really exciting new
features for a relational database, like

* Data API - It offers an HTTPS endpoint to insert, update, delete, or query data directly from your web applications
* Autoscaling - It scales up or down automatically. You don’t have to provision in advance and you do not have to rely
on guesswork for those unknown spikes or downtimes
* Sleep - It can be configured to go to sleep after a certain period of inactivity, meaning you only pay for what you
use (serverless at its best!). However, this also means that your application will have to deal with some kind of cold
* start situation. At worst, it can take more than 30 seconds to warm up and this could lead to request timing out errors.

 
However, there were some pain points with serverless v1. It does not come with any kind of multi-A-Z support.
Hence, if something goes wrong in the region or zone of your database, then your application will face downtime.
AWS does guarantee no data loss due to failure in your zone because it uses distributed, fault-tolerant, self-healing
storage with 6-way replication that will automatically grow as you add more data.
It also had no support for read-only replicas and you are not able to export data to S3.


#### What is serverless v2?

In Re: Invent 2020, AWS announced Aurora v2 with some major updates to v1

* Multi A-Z support: Aurora Serverless v2 comes with support for multi-AZ. No need to worry about downtime due to a
failure in one particular zone
* Read Replicas: With v2, you can now access data from the read replicas of your DB cluster.

However, there are some things missing from v1. You do not have access to the data API and you can no longer put your DB
cluster to sleep during the inactive period. The pricing model for v2 is also very high for smaller applications and if
you combine that with no sleep mode, then you are looking at a heavy bill at the end of your monthly cycle.