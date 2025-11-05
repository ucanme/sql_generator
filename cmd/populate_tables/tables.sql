--[{"description": "用户ID", "is_primary": true, "is_required": true, "name": "id", "type": "BIGINT"}, {"description": "用户名", "is_primary": false, "is_required": true, "name": "username", "type": "VARCHAR(50)"}, {"description": "邮箱", "is_primary": false, "is_required": true, "name": "email", "type": "VARCHAR(100)"}]
CREATE TABLE IF NOT EXISTS users (
    id BIGINT PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    email VARCHAR(100) NOT NULL
);

--[{"description": "项目ID", "is_primary": true, "is_required": true, "name": "id", "type": "BIGINT"}, {"description": "项目名", "is_primary": false, "is_required": true, "name": "name", "type": "VARCHAR(100)"}, {"description": "项目描述", "is_primary": false, "is_required": false, "name": "description", "type": "TEXT"}]
CREATE TABLE IF NOT EXISTS projects (
    id BIGINT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT
);

--[{"description": "任务ID", "is_primary": true, "is_required": true, "name": "id", "type": "BIGINT"}, {"description": "项目ID", "is_primary": false, "is_required": true, "name": "project_id", "type": "BIGINT"}, {"description": "任务标题", "is_primary": false, "is_required": true, "name": "title", "type": "VARCHAR(200)"}, {"description": "任务状态", "is_primary": false, "is_required": true, "name": "status", "type": "VARCHAR(20)"}]
CREATE TABLE IF NOT EXISTS tasks (
    id BIGINT PRIMARY KEY,
    project_id BIGINT NOT NULL,
    title VARCHAR(200) NOT NULL,
    status VARCHAR(20) NOT NULL
);
