# Review System Microservice

A scalable, event-driven microservice for ingesting, processing, and managing hotel review data. Built in Go with a serverless AWS architecture, this service offers a robust foundation for data ingestion pipelines and a full CRUD REST API.


## 📚 Table of Contents 

1. [Tech Stack](#tech-stack)
2. [Architecture Overview](#architecture-overview)
3. [Features](#features)
4. [Prerequisites](#prerequisites)
5. [Installation & Configuration](#installation--configuration)
6. [Local Development](#local-development)
7. [Usage Examples](#usage-examples)
8. [API Reference](#api-reference)
9. [Testing](#testing)
10. [Deployment](#deployment)
11. [Contributing](#contributing)
12. [License](#license)



## 💻 Tech Stack

* **Language:** Go 1.23+

* **Frameworks & Libraries:**

  * AWS Lambda Go SDK (`github.com/aws/aws-lambda-go`)
  * AWS SDK v2 (`github.com/aws/aws-sdk-go-v2`)
  * GORM (`gorm.io/gorm`, `gorm.io/driver/postgres`)
  * Gorilla Mux (`github.com/gorilla/mux`)
  * Viper & Dotenv (`github.com/spf13/viper`, `github.com/joho/godotenv`)
  * ZeroLog (`github.com/rs/zerolog`)
  * Swagger (`github.com/swaggo/swag`, `github.com/swaggo/http-swagger`)

* **Infrastructure:** AWS SAM, CloudFormation, S3, SQS, Secrets Manager, RDS (PostgreSQL), API Gateway, CloudWatch

* **Local Dev:** Docker, Docker Compose, PostgreSQL


## 🏛️ Architecture Overview

This service follows **Clean Architecture** principles, separating core business logic from external dependencies.

- **Green**: Components that are already implemented
- **Blue**: Componenets that are yet to be implemented

### Production setup:
```mermaid
graph TD
  classDef current fill:#7FC77F,stroke:#333,stroke-width:1px;
  classDef ideal fill:#6A9FD1,stroke:#333,stroke-width:1px;

  subgraph AWS_Cloud[AWS Cloud]
    direction LR

    subgraph Application_Stack[Application Stack SAM]
      direction LR
      APIGateway[API Gateway]:::current -->|"/api/v1/* - REST APIs"| Lambda
      APIGateway -->|"/swagger/*"| Lambda
      SQS[ReviewsQueue SQS]:::current -->|Trigger| Lambda
      Lambda -->|Fetch Secrets| SecretsManager(AWS Secrets Manager):::current
      Lambda -->|Logs| CloudWatch(CloudWatch Logs):::current
      Kinesis[Kinesis Stream]:::ideal -->|Real-time ingest| Lambda
      Tracing[OpenTelemetry Tracing]:::ideal
      Monitoring[Metrics & Dashboards]:::ideal
      CodePipeline[CI/CD Pipeline]:::ideal

    subgraph VPC[VPC: vpc-093a3a9e6470e3619]
    direction LR
    Lambda[Go Lambda Handler<br>SG: sg-00758af46e2ff1ae7]:::current
    RDS[(PostgreSQL in RDS)]:::current
    Redis[(Redis Cache)]:::ideal
    Lambda -->|CRUD via GORM| RDS
    Redis <-->|Cache hot reads| RDS
    end

    end

    S3[S3 Bucket: review-data-bucket]:::current -->|s3:ObjectCreated:*| SQS
  end
```

### Local Development:

```mermaid
%% Local Development Architecture
graph TD
  subgraph Local Development
    direction LR
    CLI[Importer CLI] -->|"go run cmd/importer/main.go <file>"| Processor[Processing Service]
    Processor -->|CRUD via GORM| LocalDB[(Local PostgreSQL)]
    Server[API Server]-->LocalDB
    class CLI,Processor,LocalDB,Server current;

    class DockerCompose,IDE ideal;
  end
```


## ⚡️ Features

* **Event-Driven:** Ingest reviews via S3 → SQS → Lambda pipeline.
* **Full CRUD API:** Manage providers, hotels, and reviews through REST endpoints.
* **Clean Architecture:** Ensures maintainable, testable code.
* **Secure:** Database credentials stored in AWS Secrets Manager.
* **IaC:** Resources defined with SAM & CloudFormation.
* **CI/CD**: CI/CD using Github Actions
* **Zero-Downtime Deployments:** Blue-green releases with automated rollback.
* **Local Development:** Dockerized PostgreSQL & easy setup.
* **Auto-Generated Docs:** Swagger UI for API exploration.



## Prerequisites

* Go 1.23+ installed
* Docker & Docker Compose
* AWS CLI & AWS SAM CLI
* AWS credentials configured (`~/.aws/credentials`)

---

## Installation & Configuration

1. **Clone the repo**

   ```bash
   git clone https://github.com/kirananto/review-system.git
   cd review-system
   ```

2. **Environment Variables**

   * Copy example and update as needed:

     ```bash
     cp .env.example .env
     ```
   * For local dev, default `.env` targets Docker Compose PostgreSQL.

3. **Docker Setup (Local DB)**

   ```bash
   docker-compose up -d postgres
   ```

---

## Local Development

### Run Importer CLI

```bash
go run cmd/importer/main.go /path/to/reviews.jl
```

### Start API Server

```bash
# Ensure RUN_MODE=local in .env
go run cmd/server/main.go
```

* Server listens on `http://localhost:8000`

### Invoke Lambda Locally (SAM)

```bash
sam build
sam local invoke ReviewImporterFunction \
  --event test/data/events/event.json \
  --env-vars env.json
```

---

## Usage Examples

### Import Reviews

```bash
# Process a local file
go run cmd/importer/main.go test/data/reviews.jl
```

### CRUD via cURL

```bash
# Create a Provider
curl -X POST http://localhost:8000/api/v1/providers \
  -H "Content-Type: application/json" \
  -d '{"name":"Agoda"}'

# Get Reviews
curl http://localhost:8000/api/v1/reviews
```

---

## API Reference

### Swagger UI

* Generate docs:

  ```bash
  go run -mod=mod github.com/swaggo/swag/cmd/swag init --generalInfo cmd/server/main.go
  ```
* View docs at `http://localhost:8000/swagger/index.html`

### Key Endpoints

| Resource     | Method | Path                   | Description          |
| ------------ | ------ | ---------------------- | -------------------- |
| Health Check | GET    | `/api/v1/health`       | Server health status |
| Providers    | GET    | `/api/v1/providers`    | List providers       |
| Hotels       | GET    | `/api/v1/hotels`       | Read hotel list      |
|              | POST   | `/api/v1/hotels`       | Create a hotel       |
|              | PUT    | `/api/v1/hotels/{id}`  | Update a hotel       |
| Provider Hotel| GET    | `/api/v1/provider-hotels`  | Get list of associations between Provider & Hotel       |
| Reviews      | GET    | `/api/v1/reviews`      | List reviews         |
|              | GET    | `/api/v1/reviews/{id}` | Get review by ID     |

---

## Testing

* **Unit Tests**

  ```bash
   go test -v -cover -coverprofile=coverage.out ./...
  ```
* **Integration Tests**

  * Use Docker Compose for DB.


## 🕸️ Deployment

### 1. Database Stack (One-Time)

```bash
aws cloudformation deploy \
  --template-file deploy/database.yaml \
  --stack-name review-system-database \
  --parameter-overrides DBPassword=YOUR_DB_PASSWORD \
  --capabilities CAPABILITY_IAM
```

* **Retrieve** `DBSecretArn` from stack outputs.

### 2. Application Stack (Blue-Green)

```bash
sam build --template-file deploy/sam/template.yaml
sam deploy --guided \
  --stack-name review-system-app \
  --parameter-overrides DBSecretArn=YOUR_DB_SECRET_ARN \
  --capabilities CAPABILITY_IAM
```

* **Traffic shifting**: 10% per minute
* **Rollback**: Triggered on `ApiGateway5XXErrorAlarm`

Yes, that's a solid explanation — you're showing that the system **fails gracefully**, **has retries**, and **prevents data loss** via SQS and DLQ.

Here’s a **cleaned-up version** of your write-up that you can directly put into your README under a section like **"Failure Recovery & Reliability"**:


### 🛡 Failure Recovery & Reliability

In case of downstream failures such as database throttling or service health issues, the system is designed to recover gracefully:

* Review ingestion is triggered via AWS Lambda consuming messages from an SQS queue.
* If a Lambda fails (e.g., due to a DB error), the message is automatically retried up to 5 times (default behavior).
* After the maximum retries, the message is moved to a **Dead Letter Queue (DLQ)** to avoid data loss.
* DLQ can be monitored via alerts (e.g., CloudWatch Alarms), and messages can be **redriven** for reprocessing after the root cause is resolved.
* This design ensures **at-least-once processing semantics** with **no data loss**.

A redrive policy has not been configured yet, but can be easily added based on the business use case and SLA requirements.



## 🤔 Assumptions and Design Decisions led by it

1. **Hotel IDs are not always reliable.**  
   There's no guarantee that `HotelID` will remain consistent across multiple files. We need to account for potential inconsistencies.  
   **Current Decision**: For now, we are assuming that `HotelID` is consistent across different providers and files, and treating it as a trustworthy identifier.

2. **Provider ID & Review ID is assumed to be internal.**  
   We assume that `ProviderID` & `ReviewID` is an internal identifier specific to each provider, and we will proceed based on this assumption.

3. **Ambiguity in Overall Score placement.**  
   The `overallScore` field conflicts with individual review lines — it's unclear when it should appear (before or after reviews), and the insertion order may affect interpretation.  
   **Assumption**: We store only the latest overall score, determined by the review entry with the highest count.

4. **One provider per file.**  
   Each file contains data from a single provider. A provider can upload multiple files over time, but no single file will contain reviews from multiple providers.

5. **No unique identifier for reviewers.**  
   Due to the absence of a unique reviewer ID, it's not feasible to use a relational model for reviewers. Reviewer details will instead be stored as a field in the `comment` table.

6. **Idempotency gotcha**
   The current design ensures idempotency by processing only newly added files from the S3 bucket. While checksum-based duplication checks were considered to prevent reprocessing if the same file is re-uploaded, they were intentionally avoided to reduce complexity and prevent ambiguity—especially during redrive scenarios in case of failures. This needs to be properly designed taking into consideration the re-drive policies and other failure mechanisms.




## 🐊 Gochas & Current Limitations

- Swagger Documentation works only on localhost
- Currently only supports upto 25K Reviews in a single file upto ~40MB. To support more than that we need to do either of the above approaches:
  - Fan out: Chunk into smaller groups of data by invoking lambdas or using Lambda Step functions.
  - Improve performance of a single lambda


## 🚧 TODO - Work in Progress

[Link to TODO.md file](./TODO.md)