# Debugging Exercise – Bug Report

## Bug 1: Duplicate event detection
- Problem: The seen map was declared inside processBatch().
- Why it happened: A new map was created for every batch, so duplicate events in different batches were counted again.
- Fix: Moved event ID tracking to a global map shared across batches.
- Verification: Compared actual output with expected_output.txt.

## Bug 2: Incorrect unique opens
- Problem: The openedBy map used only contact_id.
- Why it happened: The same contact opening two different campaigns was counted only once globally.
- Fix: Used campaign_id + contact_id as the map key.
- Verification: All unique_opens values matched expected_output.txt.

## Bug 3: Incorrect daily delivery date
- Problem: Timestamp.Local() was used.
- Why it happened: Local time depends on the machine timezone.
- Fix: Used Timestamp.UTC() before formatting the date.
- Verification: Daily delivered buckets matched expected_output.txt.

## Bug 4: Concurrent counter updates
- Problem: Multiple worker goroutines updated shared counters without synchronization.
- Why it happened: Concurrent writes to shared integer fields can cause data races.
- Fix: Added a mutex around counter updates in apply().
- Verification: Ran go run -race . events.jsonl without race warnings.
