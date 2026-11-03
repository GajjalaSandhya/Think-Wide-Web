# Campaign Events API

A simple Go HTTP API for ingesting campaign events and retrieving campaign statistics.

## Requirements

- Go 1.22 or later.
- No external Go dependencies.

## Run the application

From the starter directory:

```bash
cd starter
go run .

## API Endpoints

### 1. Ingest Events

POST `/events`

Accepts a JSON array of events.

Example:

```json
[
  {
    "event_id": "evt_001",
    "campaign_id": "cmp_summer_sale",
    "contact_id": "ct_001",
    "type": "opened",
    "timestamp": "2026-08-10T09:15:04Z",
    "metadata": {
      "url": "https://example.com/offer"
    }
  }
]
### 2. Retrieve Campaign Statistics

GET `/campaigns/{campaignID}/stats`

Example:

    curl http://localhost:8080/campaigns/cmp_summer_sale/stats

## Validation and Duplicate Handling

- Required fields: event_id, campaign_id, contact_id, type, timestamp.
- Supported types: sent, delivered, opened, clicked.
- Timestamps must use RFC3339 format.
- Exact duplicates are ignored.
- Reusing an event ID with different data is rejected.
- Invalid events are reported individually while valid events can still be accepted.

## Storage

Events are stored in an in-memory map protected by a mutex.

Data is lost when the application stops or restarts.

## Limitations

- No persistent database.
- No authentication.
- No optional pagination endpoint.
- No frontend.
- Not designed for distributed production workloads.
