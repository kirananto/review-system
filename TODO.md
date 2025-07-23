# ✅ TODO & Technical Debt – Review System Microservice


## 🔴 Immediate Priorities (Must-have for POC)

- [ ] Fix CRUD Operations CRUD APIs
- [ ] Add Pagination to CRUD APIs
- [ ] Consistent Error Handling across all endpoints
- [x] Write Unit Tests for core flows (CRUD, error handling, validation)
- [ ] Support for adding middlewares
- [ ] Document Deployment Strategy - Without delay
- [ ] Explain Key Design Pattern in `README.md`
- [ ] Large files, notify on success/100/100 inserted, how many failed etc. in logs/slack etc
- [.] Verify working of CLI (go run cmd/importer/main.go test/data/reviews.jl)
---

## 🟡 Mid Priority (If time permits - Helpful for UX & Dev Experience)

- [ ] Filtering & Sorting in CRUD APIs  
- [ ] Authentication Mechanism (Token/OAuth - even a stub is okay)  
- [ ] Move Auto-Migration to CI/CD instead of running on every start  
- [ ] Cache Upsert Logic for overall data

---

## 🔵 Known Gaps – Not Required for POC, but Must for Production

These are best practices or production concerns that are out of scope for now:

- [ ] Hardcoded VPC configuration in CloudFormation templates  
- [ ] Missing S3 Bucket Encryption 
- [ ] Use Production-grade IAM Policies for Lambda  
- [ ] Restrict CORS to known domains only  
- [ ] Add Rate Limiting to API Gateway or via middleware  
- [ ] Use Structured Logging (e.g., `zap`, `logrus`)  
- [ ] Set up CloudWatch Alarms  
- [ ] Add Database Indexes based on access patterns  
- [ ] Tune Lambda Concurrency settings  
- [ ] Tune RDS or DB Configurations for optimal use case

---

## 🟢 Nice to Have (Future Enhancements)

- [ ] Push logs of object creation/failure to external logging/analytics  
- [ ] Add alerts for unprocessed lines – helps improve ingestion accuracy over time
