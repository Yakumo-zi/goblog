-- PostgreSQL 数据库初始化脚本
-- 这个脚本会在数据库容器首次启动时自动执行

-- 设置编码和本地化
ALTER DATABASE goblog SET datestyle TO "ISO, YMD";
ALTER DATABASE goblog SET timezone TO 'Asia/Shanghai';
ALTER DATABASE goblog SET default_text_search_config TO 'pg_catalog.simple';

-- 创建扩展(如果需要)
-- CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 设置用户权限
GRANT ALL PRIVILEGES ON DATABASE goblog TO goblog;

-- 记录初始化完成
INSERT INTO pg_catalog.pg_description (objoid, classoid, objsubid, description)
VALUES (
    (SELECT oid FROM pg_database WHERE datname = 'goblog'),
    (SELECT oid FROM pg_class WHERE relname = 'pg_database'),
    0,
    'GoBlog database initialized at ' || CURRENT_TIMESTAMP
) ON CONFLICT DO NOTHING; 