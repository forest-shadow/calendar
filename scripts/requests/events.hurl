# Test events endpoints
# Test event creation
POST http://{{host}}/api/v1/events

{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "title": "new_event_title",
  "start_time": "2024-10-10T10:10:10Z",
  "end_time": "2024-10-10T10:10:10Z",
  "description": "new_event_description",
  "user_id": "123e4567-e89b-12d3-a456-426614174001",
  "notify_before": "PT10S"
}

HTTP/1.1 201

# Test event update
PUT http://{{host}}/api/v1/events/123e4567-e89b-12d3-a456-426614174000

{
  "title": "updated_event_title",
  "description": "updated_event_description",
  "notify_before": "PT20S"
}

HTTP/1.1 204

# Test get event
GET http://{{host}}/api/v1/events/123e4567-e89b-12d3-a456-426614174000

HTTP/1.1 200
[Asserts]
body jsonpath "$.title" == "updated_event_title"
body jsonpath "$.description" == "updated_event_description"

# Test event deletion
DELETE http://{{host}}/api/v1/events/123e4567-e89b-12d3-a456-426614174000

HTTP/1.1 204