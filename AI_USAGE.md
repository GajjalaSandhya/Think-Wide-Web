# AI Usage Disclosure

## Tools Used
ChatGPT was used as an AI assistant during development.

## How AI Helped
- Helped understand the assignment requirements.
- Assisted with planning the Go API structure.
- Explained event validation and duplicate handling.
- Helped identify the four intentional bugs in the debugging exercise.
- Suggested verification commands, including the Go race detector.
- Assisted with drafting documentation.

## Suggestions Reviewed or Changed
The initial implementation used in-memory storage because it was suitable for the assignment's limited scope and time constraints. A production database was considered but intentionally deferred.

The debugging fixes were applied and verified against the supplied expected output.

## What I Tested
- Go formatting and compilation.
- Event ingestion using curl.
- Campaign statistics endpoint.
- Duplicate event handling.
- Invalid timestamp rejection.
- Debugging output comparison.
- Race detector execution.

## Limitations
The implementation is an assignment-level solution. It does not provide production-grade persistence, authentication, distributed processing, or deployment infrastructure.
