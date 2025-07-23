# ASSUMPTIONS

1. Id's are not trustable. HotelID, no guarantee that it'll be same from multiple files. We need to handle that. So what is the unique identifier?
2. Is providerId an internal Specific id? It's best to assume that the IDs are internal and we proceed with this.
3. Overall score: This is conflicting with the review lines, as at what point would the overall score be? Before or after the score? Doesn't the order or insertion largely affect this?
Assumption: We store only the latest overall score (Based on the highest review Count)
4. Each file would be by each provider, meaning no overlap between providers and files.
5. No Unique identifier for Reviewer - Hence not possible to move it to relational approach, just another field in comment table.

# GOCHAS

- [x] Swagger Documentation works only on localhost
