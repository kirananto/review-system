# ✅ TODO & Technical Debt – Review System Microservice


## 🔴 Immediate Priorities (Must-have for POC)

- [ ] Cache Upsert Logic for overall data
- [x] Large files, notify on success/100/100 inserted, how many failed etc. in logs/slack etc & Store in db, audit log
- [x] Support for adding middlewares
- [x] Document Deployment Strategy - Without delay
- [x] Explain Key Design Pattern in `README.md`
- [x] Consistent Error Handling across all endpoints
- [x] Verify working of CLI (go run cmd/importer/main.go test/data/reviews.jl)
- [x] Fix CRUD Operations CRUD APIs (Just bare minimum due to time constraints)
- [x] Add Pagination to CRUD APIs
- [x] Filtering & Sorting in CRUD APIs
- [x] Write Unit Tests for core flows (CRUD, error handling, validation)


## 🔵 Known Gaps – Not Required for POC, but Must for Production

These are best practices or production concerns that are out of scope for now:

- [ ] Hardcoded VPC configuration in CloudFormation templates
- [ ] Add redrive policy from the DLQ (Dead letter queue) 
- [ ] Move Auto-Migration to CI/CD instead of running on every start
- [ ] Missing S3 Bucket Encryption 
- [ ] Authentication Mechanism | Move to actual one (Most of the time we use a Third Party Service)
- [ ] Use Production-grade IAM Policies for Lambda  
- [ ] Restrict CORS to known domains only  
- [ ] Add Rate Limiting to API Gateway or via middleware  
- [ ] Set up CloudWatch Alarms  
- [ ] Tweak Database Indexes based on access patterns (Currently few indexes are already added) 
- [ ] Tune Lambda Concurrency settings  
- [ ] Tune RDS or DB Configurations for optimal use case


## 🟢 Nice to Have (Future Enhancements)

- [ ] Push logs of object creation/failure to external logging/analytics  
- [ ] Add alerts for unprocessed lines – helps improve ingestion accuracy over time
