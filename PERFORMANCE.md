# Performance Optimizations

این پروژه با تمرکز بر performance طراحی شده است.

## انتخاب‌های معماری

### ✅ Raw SQL + Migrations
- **بدون overhead ORM** - queries مستقیم و سریع
- **کنترل کامل** روی query execution
- **بهینه‌سازی دستی** queries پیچیده

### ✅ Connection Pooling
- Max connections: 25
- Idle connections: 5
- Connection lifetime: 5 minutes
- Idle timeout: 10 minutes

## Database Indexes

### Primary Indexes
- `idx_users_email` - برای login/authentication
- `idx_content_type` - فیلتر بر اساس type
- `idx_content_is_active` - فیلتر active content
- `idx_content_rating` - sort بر اساس rating

### Composite Indexes (برای queries ترکیبی)
- `idx_content_type_active_rating` - برای pagination با فیلتر
- `idx_watch_history_user_content` - برای continue watching
- `idx_watchlist_user_profile` - برای watchlist queries

### Full-Text Search
- `idx_content_title_trgm` - GIN index برای search
- `idx_content_description_trgm` - GIN index برای search

## Query Optimization Tips

### 1. استفاده از LIMIT و OFFSET
```go
// ✅ Good - با limit
SELECT * FROM content LIMIT 20 OFFSET 0

// ❌ Bad - بدون limit
SELECT * FROM content
```

### 2. Select فقط فیلدهای مورد نیاز
```go
// ✅ Good
SELECT id, title, poster_url FROM content

// ❌ Bad
SELECT * FROM content
```

### 3. استفاده از Indexes در WHERE clauses
```go
// ✅ Good - از index استفاده می‌کند
WHERE type = 'movie' AND is_active = true

// ❌ Bad - index استفاده نمی‌شود
WHERE LOWER(title) = 'movie'
```

### 4. Prepared Statements
تمام queries از prepared statements استفاده می‌کنند (مزیت `database/sql`)

## Monitoring Performance

### PostgreSQL Query Analysis
```sql
-- بررسی slow queries
SELECT query, mean_exec_time, calls 
FROM pg_stat_statements 
ORDER BY mean_exec_time DESC 
LIMIT 10;

-- بررسی index usage
SELECT schemaname, tablename, indexname, idx_scan 
FROM pg_stat_user_indexes 
ORDER BY idx_scan;
```

### Go Profiling
```powershell
# CPU profiling
go tool pprof http://localhost:8080/debug/pprof/profile

# Memory profiling
go tool pprof http://localhost:8080/debug/pprof/heap
```

## Best Practices

1. **Connection Pooling** - تنظیمات بهینه برای load
2. **Indexes** - برای تمام foreign keys و frequently queried fields
3. **Query Optimization** - استفاده از EXPLAIN ANALYZE
4. **Caching** - Redis برای frequently accessed data (آینده)
5. **Pagination** - همیشه از LIMIT/OFFSET استفاده کنید

## Future Optimizations

- [ ] Redis caching برای recommendations
- [ ] Database read replicas برای scaling
- [ ] Query result caching
- [ ] Materialized views برای complex aggregations
