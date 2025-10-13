curl http://localhost:8080/analytics/abc123?group=day

curl http://localhost:8080/analytics/abc123?group=month

curl http://localhost:8080/analytics/abc123?group=detailed

{
  "total_visits": 150,
  "unique_ips": 45,
  "daily_activity": {
    "2025-10-13": 23,
    "2025-10-12": 18,
    ...
  },
  "monthly_activity": {
    "2025-10": 87,
    "2025-09": 63
  },
  "device_stats": {
    "mobile": 85,
    "desktop": 55,
    "tablet": 10
  },
  "visits": [
    {
      "id": 150,
      "visited_at": "2025-10-13T14:30:00Z",
      ...
    }
  ]
}