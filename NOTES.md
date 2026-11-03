# Campaign Events Assignment – Design Notes

## 1. Problem Interpretation

The application receives campaign events from providers and stores them.
It exposes an API to retrieve aggregate statistics for a campaign.

Supported event types:
- sent
- delivered
- opened
- clicked

## 2. Assumptions

1. event_id uniquely identifies an event across providers.
2. An exact duplicate event is ignored.
3. Reusing an event_id with different data is rejected.
4. Each valid event is counted once.
5. Statistics count event records, not distinct contacts.
6. Timestamps must follow RFC3339 format.
7. Events may arrive late or out of order.
8. Invalid events are rejected individually; valid events in the same batch are still accepted.
9. Data is stored in memory for this implementation.
10. A mutex protects shared data during concurrent API access.

## 3. Questions for the Product Manager

1. Is event_id globally unique or unique only within a provider?
2. Should opened statistics count total open events or unique contacts?
3. Can providers update or delete previously submitted events?
4. Should batch processing be atomic or allow partial acceptance?
5. How long should event data be retained?
6. Should statistics support custom date ranges?
7. What metadata fields should be supported?

## 4. Priorities

- Correct event ingestion and validation.
- Duplicate prevention.
- Accurate campaign statistics.
- Safe concurrent access.
- Clear API responses and error handling.
- Reliable persistence for production use.

## 5. Scaling to 100 Million Events Per Day

The current in-memory implementation would not be suitable for this volume.

### What would break?
- Memory usage would grow continuously.
- Data would be lost if the application restarted.
- A single application instance would become a bottleneck.
- Repeatedly scanning all events for statistics would become expensive.

### Proposed improvements

1. Use a durable database such as PostgreSQL.
2. Add a unique index on event_id to prevent duplicate insertion.
3. Add indexes on campaign_id and timestamp.
4. Introduce a message queue such as Kafka to absorb traffic spikes.
5. Use multiple ingestion workers for parallel processing.
6. Maintain pre-aggregated campaign statistics instead of scanning every event.
7. Partition large event tables by date or another suitable key.
8. Use idempotent processing so retries do not increase counts twice.

### Queue and reliability considerations

- Apply backpressure when consumers cannot keep up.
- Retry temporary failures.
- Send repeatedly failing messages to a dead-letter queue.
- Monitor queue lag and processing failures.
- Make database writes and aggregate updates safe against duplicate processing.

### Monitoring

Monitor:
- Event ingestion rate.
- Duplicate and rejected event counts.
- API latency and error rate.
- Queue lag.
- Database performance.
- Aggregation delays.
- Differences between raw events and aggregate counters.

## 6. Marketer Scenario

At 10:00:
- Sent: 1,000,000
- Delivered: 970,000
- Opened: 250,000
- Clicked: 20,000

At 10:30:
- Sent: 1,000,000
- Delivered: 975,000
- Opened: 248,000
- Clicked: 20,000

### Initial interpretation

Delivered increased by 5,000, while opened decreased by 2,000.
Clicked remained unchanged.

I would not immediately conclude that the system is broken.

### Possible explanations

1. Late delivery events may have arrived.
2. The reporting window or timezone may have changed.
3. Open events may have been corrected or removed.
4. The dashboard may use unique opens rather than total open events.
5. Different filters or campaign scopes may have been applied.
6. Cached or replicated statistics may be temporarily inconsistent.

### Investigation

- Compare raw provider events at both times.
- Check event timestamps and ingestion timestamps.
- Verify duplicate handling.
- Check dashboard filters and reporting windows.
- Compare total opens with unique opens.
- Inspect aggregation and cache update logs.

### Decision

The increase in delivered events can be normal because delivery receipts may arrive later.

However, a decrease in cumulative opened counts requires investigation. If the metric is supposed to be a cumulative count of immutable events, it should not decrease. I would verify the metric definition and underlying records before deciding whether there is a bug.

## 7. What I Completed / What I Intentionally Skipped

### Completed
- Go HTTP API for event ingestion.
- Campaign statistics endpoint.
- Event validation.
- Duplicate event handling.
- In-memory storage with mutex protection.
- Debugging exercise with four fixes.
- Race detector and expected-output verification.

### Intentionally skipped
- Optional paginated events endpoint.
- Frontend application.
- Production database integration.
- Authentication and deployment infrastructure.

## 8. If I Had Another Day

I would add persistent database storage, automated unit and integration tests, pagination, better observability, and performance testing with larger event volumes.
