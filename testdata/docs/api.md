# Product API

## Authentication

Use a bearer token in the `Authorization` header. Tokens should be scoped to the minimum API permissions needed by the client.

## Rate Limits

The API allows 120 requests per minute for each workspace. Retry with exponential backoff when a `429` response is returned.
