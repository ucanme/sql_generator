# SQL Query Bot with RAG Enhancement / SQL查询机器人与RAG增强

An intelligent SQL query generator based on large language models and RAG technology that can automatically generate SQL queries based on natural language descriptions.
一个基于大语言模型和RAG技术的智能SQL查询生成器，可以根据自然语言描述自动生成SQL查询语句。

## Features / 功能特性

- Store and manage millions of table structure definitions
- Intelligent SQL query generation based on natural language descriptions (supports multi-table joins)
- RESTful API interface
- MongoDB storage support for handling massive data
- Integration with large language models (such as DeepSeek, OpenAI GPT series)
- RAG enhancement features to improve query accuracy

- 存储和管理数百万张表结构定义
- 基于自然语言描述智能生成SQL查询（支持多表关联）
- RESTful API接口
- 支持MongoDB存储，可处理海量数据
- 集成大语言模型（如DeepSeek、OpenAI GPT系列）
- RAG增强功能，提高查询准确性

*(For Chinese version, please refer to the end of this document or see [README_zh.md](README_zh.md))*
*(有关中文版本，请参阅本文档末尾或查看 [README_zh.md](README_zh.md))*

## Technical Architecture / 技术架构

### Core Components / 核心组件

1. **Web Service Layer** - HTTP service based on Gin framework
2. **Storage Layer** - MongoDB for persistent storage, Pinecone for vector storage
3. **Large Language Model Layer** - Supports multiple LLMs (DeepSeek, OpenAI, etc.)
4. **RAG Enhancement Layer** - Retrieval-Augmented Generation technology based on vector search

1. **Web服务层** - 基于Gin框架的HTTP服务
2. **存储层** - MongoDB用于持久化存储，Pinecone用于向量存储
3. **大语言模型层** - 支持多种大语言模型（DeepSeek、OpenAI等）
4. **RAG增强层** - 基于向量检索的增强生成技术

### Data Flow / 数据流向

```
User natural language query                         用户自然语言查询
       ↓                                        ↓
[Query understanding and keyword extraction]    [查询理解与关键词提取]
       ↓                                        ↓
[Retrieve relevant table structures from vector database]  [向量数据库中检索相关表结构]
       ↓                                        ↓
[Re-rank most relevant table structures]        [重排序最相关表结构]
       ↓                                        ↓
[Construct prompt and send to LLM]              [构造Prompt发送给大语言模型]
       ↓                                        ↓
[Generate and return SQL query]                 [生成并返回SQL查询]
```

## Quick Start / 快速开始

### Requirements / 环境要求

- Go 1.21+
- MySQL 8.0+ (for test scripts)
- MongoDB 4.0+ (for application storage)
- Pinecone account (or other vector database)
- LLM API Key (DeepSeek, OpenAI, etc.)

- Go 1.21+
- MySQL 8.0+ (用于测试脚本)
- MongoDB 4.0+ (用于应用存储)
- Pinecone账户（或其他向量数据库）
- 大语言模型API Key（DeepSeek、OpenAI等）

### Starting Dependencies with Docker Compose / 使用 Docker Compose 启动依赖服务

The project provides a [docker-compose.yml](file:///Users/ji.wu/work/project/awesomeProject2/docker-compose.yml) file that can be used to quickly start the required MySQL and MongoDB services:

项目提供了 [docker-compose.yml](file:///Users/ji.wu/work/project/awesomeProject2/docker-compose.yml) 文件，可用于快速启动所需的 MySQL 和 MongoDB 服务：

```bash
docker-compose up -d
```

This will start the following services:
- MySQL: Port 3306, root user password is password
- MongoDB: Port 27017, root user password is password

这将启动以下服务：
- MySQL: 端口 3306，root 用户密码为 password
- MongoDB: 端口 27017，root 用户密码为 password

### Installation Steps / 安装步骤

1. Clone the project:
```bash
git clone <repository-url>
cd awesomeProject2
```

1. 克隆项目：
```bash
git clone <repository-url>
cd awesomeProject2
```

2. Install dependencies:
```bash
go mod tidy
```

2. 安装依赖：
```bash
go mod tidy
```

3. Configure environment variables (copy `.example_env` to `.env` and fill in the configuration):
```bash
cp .example_env .env
```

3. 配置环境变量（复制 `.example_env` 为 `.env` 并填写相应配置）：
```bash
cp .example_env .env
```

4. Initialize database (optional, for testing):
```bash
# Initialize database with default configuration
./scripts/init_db.sh

# Or with custom database configuration
MYSQL_HOST=localhost MYSQL_PORT=3306 MYSQL_USER=root MYSQL_PASS=password MYSQL_DB=test ./scripts/init_db.sh
```

4. 初始化数据库（可选，用于测试）：
```bash
# 使用默认配置初始化数据库
./scripts/init_db.sh

# 或者自定义数据库配置
MYSQL_HOST=localhost MYSQL_PORT=3306 MYSQL_USER=root MYSQL_PASS=password MYSQL_DB=test ./scripts/init_db.sh
```

5. Start the service:
```bash
go run main.go
```

5. 启动服务：
```bash
go run main.go
```

The service will start on the configured port (default 8080).

服务将在配置的端口（默认8080）上启动。

## Configuration / 配置说明

Configure the application via environment variables or `.env` file:

通过环境变量或 `.env` 文件配置应用：

| Variable | Default | Description | 变量名 | 默认值 | 描述 |
|----------|---------|-------------|--------|--------|------|
| SERVER_PORT | 8080 | Service port | SERVER_PORT | 8080 | 服务端口 |
| SERVER_READ_TIMEOUT | 30 | Read timeout (seconds) | SERVER_READ_TIMEOUT | 30 | 读取超时时间（秒） |
| SERVER_WRITE_TIMEOUT | 30 | Write timeout (seconds) | SERVER_WRITE_TIMEOUT | 30 | 写入超时时间（秒） |
| MONGO_URI | mongodb://localhost:27017 | MongoDB connection string | MONGO_URI | mongodb://localhost:27017 | MongoDB连接字符串 |
| MONGO_DATABASE | sqlbot | MongoDB database name | MONGO_DATABASE | sqlbot | MongoDB数据库名称 |
| MYSQL_DSN | - | MySQL connection string (for some test functions) | MYSQL_DSN | - | MySQL连接字符串（用于某些测试功能） |
| MYSQL_DATABASE | test | MySQL database name | MYSQL_DATABASE | test | MySQL数据库名称 |
| LLM_API_KEY | - | LLM API key | LLM_API_KEY | - | 大语言模型API密钥 |
| LLM_MODEL | gpt-3.5-turbo | LLM model name | LLM_MODEL | gpt-3.5-turbo | 大语言模型名称 |
| LLM_MAX_TOKENS | 2000 | Maximum tokens | LLM_MAX_TOKENS | 2000 | 最大token数 |
| LLM_TEMPERATURE | 0.3 | Temperature parameter | LLM_TEMPERATURE | 0.3 | 温度参数 |
| EMBEDDING_API_KEY | - | Embedding model API key | EMBEDDING_API_KEY | - | 嵌入模型API密钥 |
| EMBEDDING_MODEL | sentence-transformers/all-MiniLM-L6-v2 | Embedding model name | EMBEDDING_MODEL | sentence-transformers/all-MiniLM-L6-v2 | 嵌入模型名称 |
| EMBEDDING_PROVIDER | sbert | Embedding provider (supports: openai, deepseek, huggingface, sbert) | EMBEDDING_PROVIDER | sbert | 嵌入服务提供商 (支持: openai, deepseek, huggingface, sbert) |
| HF_ENDPOINT | https://api-inference.huggingface.co/models/ | Hugging Face API endpoint | HF_ENDPOINT | https://api-inference.huggingface.co/models/ | Hugging Face API端点 |
| HF_MODEL | sentence-transformers/all-MiniLM-L6-v2 | Hugging Face model name | HF_MODEL | sentence-transformers/all-MiniLM-L6-v2 | Hugging Face模型名称 |
| QWEN_API_KEY | - | Qwen API key | QWEN_API_KEY | - | 阿里云千问API密钥 |
| QWEN_MODEL | text-embedding-v1 | Qwen embedding model name | QWEN_MODEL | text-embedding-v1 | 阿里云千问嵌入模型名称 |
| VECTOR_DB_API_KEY | - | Vector database API key | VECTOR_DB_API_KEY | - | 向量数据库API密钥 |
| VECTOR_DB_INDEX_NAME | sqlbot-tables | Vector database index name | VECTOR_DB_INDEX_NAME | sqlbot-tables | 向量数据库索引名称 |
| VECTOR_DB_ENVIRONMENT | us-west1-gcp | Vector database environment | VECTOR_DB_ENVIRONMENT | us-west1-gcp | 向量数据库环境 |

## Database Initialization Scripts / 数据库初始化脚本

The project includes scripts for initializing test data:

项目包含了用于初始化测试数据的脚本：

1. **[init_db.sh](file:///Users/ji.wu/work/project/awesomeProject2/scripts/init_db.sh)** - Complete database initialization script
   - Check database connection
   - Create database and table structures
   - Run Go program to populate test data

1. **[init_db.sh](file:///Users/ji.wu/work/project/awesomeProject2/scripts/init_db.sh)** - 完整的数据库初始化脚本
   - 检查数据库连接
   - 创建数据库和表结构
   - 运行 Go 程序填充测试数据

2. **[populate_tables.sh](file:///Users/ji.wu/work/project/awesomeProject2/scripts/populate_tables.sh)** - Script to only populate table structure data

2. **[populate_tables.sh](file:///Users/ji.wu/work/project/awesomeProject2/scripts/populate_tables.sh)** - 仅填充表结构数据的脚本

### init_db.sh Usage / init_db.sh 使用方法

```bash
# Use default configuration
./scripts/init_db.sh

# Show verbose output
./scripts/init_db.sh -v

# Show help information
./scripts/init_db.sh -h

# Use custom configuration
MYSQL_HOST=localhost MYSQL_PORT=3306 MYSQL_USER=user MYSQL_PASS=pass MYSQL_DB=mydb ./scripts/init_db.sh
```

```bash
# 使用默认配置
./scripts/init_db.sh

# 显示详细输出
./scripts/init_db.sh -v

# 显示帮助信息
./scripts/init_db.sh -h

# 使用自定义配置
MYSQL_HOST=localhost MYSQL_PORT=3306 MYSQL_USER=user MYSQL_PASS=pass MYSQL_DB=mydb ./scripts/init_db.sh
```

## Embedding Service Providers / 嵌入服务提供商说明

The system supports multiple embedding service providers:

系统支持多种嵌入服务提供商：

1. **Hugging Face** (`huggingface`)
   - Uses Hugging Face model API
   - Requires API key setup
   - Default model: `sentence-transformers/all-MiniLM-L6-v2`
   - Custom API endpoint can be configured via `HF_ENDPOINT`
   - Specific model can be configured via `HF_MODEL`

1. **Hugging Face** (`huggingface`)
   - 使用Hugging Face的模型API
   - 需要设置API密钥
   - 默认模型: `sentence-transformers/all-MiniLM-L6-v2`
   - 可通过 `HF_ENDPOINT` 配置自定义API端点
   - 可通过 `HF_MODEL` 配置特定模型

2. **OpenAI** (`openai`)
   - Uses OpenAI embedding API
   - Requires API key setup
   - Default model: `text-embedding-ada-002`

2. **OpenAI** (`openai`)
   - 使用OpenAI的嵌入API
   - 需要设置API密钥
   - 默认模型: `text-embedding-ada-002`

3. **DeepSeek** (`deepseek`)
   - Uses DeepSeek embedding API (if supported)
   - Requires API key setup
   - Default model: `text-embedding-3-small`

3. **DeepSeek** (`deepseek`)
   - 使用DeepSeek的嵌入API（如果支持）
   - 需要设置API密钥
   - 默认模型: `text-embedding-3-small`

4. **Qwen** (`qwen`)
   - Uses Qwen embedding API
   - Requires API key setup
   - Default model: `text-embedding-v1`
   - Accessed via DashScope OpenAI-compatible API endpoint

4. **阿里云千问** (`qwen`)
   - 使用阿里云千问的嵌入API
   - 需要设置API密钥
   - 默认模型: `text-embedding-v1`
   - 通过DashScope兼容OpenAI的API端点访问

## API Endpoints / API 接口

### Health Check / 健康检查
- `GET /health` - Check service status
- `GET /health` - 检查服务状态

### Table Structure Management / 表结构管理
- `POST /tables` - Create table structure definition
- `GET /tables` - List all table structures with pagination
- `GET /tables/:name` - Get specified table structure
- `PUT /tables/:name` - Update specified table structure
- `DELETE /tables/:name` - Delete specified table structure
- `GET /tables/search/:keyword` - Search related table structures

- `POST /tables` - 创建表结构定义
- `GET /tables` - 分页列出所有表结构
- `GET /tables/:name` - 获取指定表结构
- `PUT /tables/:name` - 更新指定表结构
- `DELETE /tables/:name` - 删除指定表结构
- `GET /tables/search/:keyword` - 搜索相关表结构

### Query Generation / 查询生成
- `POST /queries/generate` - Generate SQL query based on description
- `GET /queries` - List all generated queries with pagination
- `GET /queries/:id` - Get specified query

- `POST /queries/generate` - 根据描述生成SQL查询
- `GET /queries` - 分页列出所有已生成的查询
- `GET /queries/:id` - 获取指定查询

## Usage Examples / 使用示例

### 1. Define Table Structure / 定义表结构

```bash
curl -X POST http://localhost:8080/tables \
  -H "Content-Type: application/json" \
  -d '{
    "name": "users",
    "description": "User table",
    "columns": [
      {
        "name": "id",
        "type": "int",
        "description": "User ID",
        "is_primary": true,
        "is_required": true
      },
      {
        "name": "name",
        "type": "varchar(100)",
        "description": "Username",
        "is_required": true
      },
      {
        "name": "email",
        "type": "varchar(255)",
        "description": "User email"
      }
    ]
  }'
```

### 2. Generate SQL Query / 生成SQL查询

```bash
curl -X POST http://localhost:8080/queries/generate \
  -H "Content-Type: application/json" \
  -d '{
    "description": "Get names and emails of all users"
  }'
```

### 3. Generate Query for Specific Tables / 指定特定表生成查询

```bash
curl -X POST http://localhost:8080/queries/generate \
  -H "Content-Type: application/json" \
  -d '{
    "description": "Get user information with ID greater than 100",
    "table_names": ["users"]
  }'
```

## RAG Enhancement Features / RAG增强功能

The system enhances query accuracy through RAG (Retrieval-Augmented Generation) technology:

本系统通过RAG（Retrieval-Augmented Generation）技术增强查询准确性：

1. **Vector Storage**: Convert table structure information to vectors and store in vector database
2. **Semantic Retrieval**: Retrieve relevant table structures based on semantic similarity of user queries
3. **Context Enhancement**: Provide retrieved table structures as context to the large language model
4. **Precise Generation**: The LLM generates more accurate SQL queries based on enhanced context

1. **向量化存储**：将表结构信息转换为向量并存储在向量数据库中
2. **语义检索**：根据用户查询的语义相似度检索相关表结构
3. **上下文增强**：将检索到的相关表结构作为上下文提供给大语言模型
4. **精准生成**：大语言模型基于增强的上下文生成更准确的SQL查询

## Performance Optimization / 性能优化

1. **Intelligent Table Structure Retrieval**: Fast retrieval of relevant table structures using vector database
2. **Prompt Optimization**: Control context length sent to large language models
3. **Pagination Support**: All list interfaces support pagination
4. **Connection Pooling**: Database connection pool optimization

1. **智能表结构检索**：使用向量数据库快速检索相关表结构
2. **提示词优化**：控制发送给大语言模型的上下文长度
3. **分页支持**：所有列表接口均支持分页
4. **连接池**：数据库连接池优化

## Extensibility / 扩展性

- Easy to integrate with other LLMs (such as Claude, ERNIE Bot, etc.)
- Modular design for easy feature expansion
- Supports horizontal scaling to handle larger data volumes

- 易于集成其他大语言模型（如Claude、文心一言等）
- 模块化设计便于功能扩展
- 支持水平扩展以处理更大规模数据

---

## 中文版本 / Chinese Version

请参阅 [README_zh.md](README_zh.md) 获取完整的中文文档。