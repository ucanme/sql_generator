- [ ] 结合hive表的结构## Overview
帮我使用golang实现一个问答机器人，可以帮我完成查找表的沟通，自动生成hive sql功能(需要支持多表联查)，表可能有数百万张

## Task 1: Basic Project Structure
- [x] 设计表结构及其描述的存储（支持数百万张表）
- [x] 实现HTTP服务功能
- [x] 设置合理的包结构

## Task 2: Data Models and Storage
- [x] 定义数据模型（Table, Column, Query）
- [x] 使用mysql（替代内存存储以支持海量数据）
- [x] 添加CRUD操作

## Task 3: HTTP Server Implementation
- [x] 创建HTTP服务器
- [x] 实现基本路由
- [x] 添加健康检查端点

## Task 4: REST API Endpoints
- [x] 实现表结构管理的GET/POST/PUT/DELETE端点
- [x] 实现SQL查询生成和管理端点
- [x] 添加分页支持

## Task 5: Query Processing with Large Language Models
- [x] 集成DeepSeek大模型API
- [x] 结合Hive表结构生成SQL
- [x] 实现多表关联查询生成

## Task 6: RAG (Retrieval-Augmented Generation) Enhancement
- [x] 集成向量数据库（如Pinecone)，支持准实时更新表结构信息
- [x] 实现表结构向量化存储
- [x] 开发语义化表结构检索功能
- [x] 结合RAG和大模型生成更准确的SQL

## Task 7: Error Handling and Validation
- [x] 实现错误处理中间件
- [x] 添加请求验证
- [x] 创建一致的错误响应格式

## Task 8: Performance Optimization
- [ ] 实现数据库索引优化
- [ ] 添加查询结果缓存
- [ ] 优化大语言模型调用性能

## Task 9: Testing
- [ ] 添加单元测试
- [ ] 添加集成测试
- [ ] 配置测试覆盖率
- [ ] 生成一个包含15个关联的表的case，测试多表联查功能

## Task 10: Documentation
- [x] 记录API端点
- [x] 添加包含使用说明的README文件
- [ ] 记录部署过程
- [x] 添加RAG增强功能说明