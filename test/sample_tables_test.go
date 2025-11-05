package test

import (
	"testing"
	"time"

	"awesomeProject2/internal/models"
	"github.com/google/uuid"
)

// GenerateSampleTables creates 15 interconnected sample tables
func GenerateSampleTables() []*models.Table {
	tables := make([]*models.Table, 0, 15)

	now := time.Now()

	// 1. Users table
	users := &models.Table{
		ID:          uuid.New().String(),
		Name:        "users",
		Description: "用户基本信息表，存储系统用户的核心信息",
		Columns: []models.Column{
			{
				Name:        "id",
				Type:        "BIGINT",
				Description: "用户唯一标识符，主键",
				IsPrimary:   true,
				IsRequired:  true,
			},
			{
				Name:        "username",
				Type:        "VARCHAR(50)",
				Description: "用户名，用于登录系统",
				IsPrimary:   false,
				IsRequired:  true,
			},
			{
				Name:        "email",
				Type:        "VARCHAR(100)",
				Description: "用户邮箱地址",
				IsPrimary:   false,
				IsRequired:  true,
			},
			{
				Name:        "full_name",
				Type:        "VARCHAR(100)",
				Description: "用户全名",
				IsPrimary:   false,
				IsRequired:  false,
			},
			{
				Name:        "created_at",
				Type:        "TIMESTAMP",
				Description: "记录创建时间",
				IsPrimary:   false,
				IsRequired:  true,
			},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
	tables = append(tables, users)

	// 2. Departments table
	departments := &models.Table{
		ID:          uuid.New().String(),
		Name:        "departments",
		Description: "部门信息表，存储组织架构中的部门信息",
		Columns: []models.Column{
			{
				Name:        "id",
				Type:        "BIGINT",
				Description: "部门唯一标识符，主键",
				IsPrimary:   true,
				IsRequired:  true,
			},
			{
				Name:        "name",
				Type:        "VARCHAR(100)",
				Description: "部门名称",
				IsPrimary:   false,
				IsRequired:  true,
			},
			{
				Name:        "description",
				Type:        "TEXT",
				Description: "部门描述信息",
				IsPrimary:   false,
				IsRequired:  false,
			},
			{
				Name:        "manager_id",
				Type:        "BIGINT",
				Description: "部门经理的用户ID，关联users表",
				IsPrimary:   false,
				IsRequired:  false,
			},
			{
				Name:        "created_at",
				Type:        "TIMESTAMP",
				Description: "记录创建时间",
				IsPrimary:   false,
				IsRequired:  true,
			},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
	tables = append(tables, departments)

	// 3. User-Department mapping table
	userDepartments := &models.Table{
		ID:          uuid.New().String(),
		Name:        "user_departments",
		Description: "用户部门关联表，记录用户所属的部门",
		Columns: []models.Column{
			{
				Name:        "id",
				Type:        "BIGINT",
				Description: "关联唯一标识符，主键",
				IsPrimary:   true,
				IsRequired:  true,
			},
			{
				Name:        "user_id",
				Type:        "BIGINT",
				Description: "用户ID，关联users表",
				IsPrimary:   false,
				IsRequired:  true,
			},
			{
				Name:        "department_id",
				Type:        "BIGINT",
				Description: "部门ID，关联departments表",
				IsPrimary:   false,
				IsRequired:  true,
			},
			{
				Name:        "is_primary",
				Type:        "BOOLEAN",
				Description: "是否为主要部门",
				IsPrimary:   false,
				IsRequired:  true,
			},
			{
				Name:        "assigned_at",
				Type:        "TIMESTAMP",
				Description: "分配时间",
				IsPrimary:   false,
				IsRequired:  true,
			},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
	tables = append(tables, userDepartments)

	// 4. Roles table
	roles := &models.Table{
		ID:          uuid.New().String(),
		Name:        "roles",
		Description: "角色信息表，定义系统中的各种角色",
		Columns: []models.Column{
			{
				Name:        "id",
				Type:        "BIGINT",
				Description: "角色唯一标识符，主键",
				IsPrimary:   true,
				IsRequired:  true,
			},
			{
				Name:        "name",
				Type:        "VARCHAR(50)",
				Description: "角色名称",
				IsPrimary:   false,
				IsRequired:  true,
			},
			{
				Name:        "description",
				Type:        "TEXT",
				Description: "角色详细描述",
				IsPrimary:   false,
				IsRequired:  false,
			},
			{
				Name:        "permissions",
				Type:        "JSON",
				Description: "角色权限列表，以JSON格式存储",
				IsPrimary:   false,
				IsRequired:  false,
			},
			{
				Name:        "created_at",
				Type:        "TIMESTAMP",
				Description: "记录创建时间",
				IsPrimary:   false,
				IsRequired:  true,
			},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
	tables = append(tables, roles)

	// 5. User-Role mapping table
	userRoles := &models.Table{
		ID:          uuid.New().String(),
		Name:        "user_roles",
		Description: "用户角色关联表，记录用户拥有的角色",
		Columns: []models.Column{
			{
				Name:        "id",
				Type:        "BIGINT",
				Description: "关联唯一标识符，主键",
				IsPrimary:   true,
				IsRequired:  true,
			},
			{
				Name:        "user_id",
				Type:        "BIGINT",
				Description: "用户ID，关联users表",
				IsPrimary:   false,
				IsRequired:  true,
			},
			{
				Name:        "role_id",
				Type:        "BIGINT",
				Description: "角色ID，关联roles表",
				IsPrimary:   false,
				IsRequired:  true,
			},
			{
				Name:        "assigned_at",
				Type:        "TIMESTAMP",
				Description: "分配时间",
				IsPrimary:   false,
				IsRequired:  true,
			},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
	tables = append(tables, userRoles)

	// 6. Projects table
	projects := &models.Table{
		ID:          uuid.New().String(),
		Name:        "projects",
		Description: "项目信息表，存储项目的基本信息",
		Columns: []models.Column{
			{
				Name:        "id",
				Type:        "BIGINT",
				Description: "项目唯一标识符，主键",
				IsPrimary:   true,
				IsRequired:  true,
			},
			{
				Name:        "name",
				Type:        "VARCHAR(100)",
				Description: "项目名称",
				IsPrimary:   false,
				IsRequired:  true,
			},
			{
				Name:        "description",
				Type:        "TEXT",
				Description: "项目详细描述",
				IsPrimary:   false,
				IsRequired:  false,
			},
			{
				Name:        "start_date",
				Type:        "DATE",
				Description: "项目开始日期",
				IsPrimary:   false,
				IsRequired:  false,
			},
			{
				Name:        "end_date",
				Type:        "DATE",
				Description: "项目结束日期",
				IsPrimary:   false,
				IsRequired:  false,
			},
			{
				Name:        "status",
				Type:        "VARCHAR(20)",
				Description: "项目状态 (active, completed, cancelled)",
				IsPrimary:   false,
				IsRequired:  true,
			},
			{
				Name:        "created_at",
				Type:        "TIMESTAMP",
				Description: "记录创建时间",
				IsPrimary:   false,
				IsRequired:  true,
			},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
	tables = append(tables, projects)

	// 7. Project members table
	projectMembers := &models.Table{
		ID:          uuid.New().String(),
		Name:        "project_members",
		Description: "项目成员表，记录项目中的成员及其角色",
		Columns: []models.Column{
			{
				Name:        "id",
				Type:        "BIGINT",
				Description: "关联唯一标识符，主键",
				IsPrimary:   true,
				IsRequired:  true,
			},
			{
				Name:        "project_id",
				Type:        "BIGINT",
				Description: "项目ID，关联projects表",
				IsPrimary:   false,
				IsRequired:  true,
			},
			{
				Name:        "user_id",
				Type:        "BIGINT",
				Description: "用户ID，关联users表",
				IsPrimary:   false,
				IsRequired:  true,
			},
			{
				Name:        "role",
				Type:        "VARCHAR(50)",
				Description: "在项目中的角色 (manager, developer, tester等)",
				IsPrimary:   false,
				IsRequired:  true,
			},
			{
				Name:        "joined_at",
				Type:        "TIMESTAMP",
				Description: "加入项目的时间",
				IsPrimary:   false,
				IsRequired:  true,
			},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
	tables = append(tables, projectMembers)

	// 8. Tasks table
	tasks := &models.Table{
		ID:          uuid.New().String(),
		Name:        "tasks",
		Description: "任务信息表，存储项目中的具体任务",
		Columns: []models.Column{
			{
				Name:        "id",
				Type:        "BIGINT",
				Description: "任务唯一标识符，主键",
				IsPrimary:   true,
				IsRequired:  true,
			},
			{
				Name:        "project_id",
				Type:        "BIGINT",
				Description: "项目ID，关联projects表",
				IsPrimary:   false,
				IsRequired:  true,
			},
			{
				Name:        "title",
				Type:        "VARCHAR(200)",
				Description: "任务标题",
				IsPrimary:   false,
				IsRequired:  true,
			},
			{
				Name:        "description",
				Type:        "TEXT",
				Description: "任务详细描述",
				IsPrimary:   false,
				IsRequired:  false,
			},
			{
				Name:        "assignee_id",
				Type:        "BIGINT",
				Description: "任务负责人ID，关联users表",
				IsPrimary:   false,
				IsRequired:  false,
			},
			{
				Name:        "status",
				Type:        "VARCHAR(20)",
				Description: "任务状态 (todo, in_progress, review, done)",
				IsPrimary:   false,
				IsRequired:  true,
			},
			{
				Name:        "priority",
				Type:        "VARCHAR(20)",
				Description: "任务优先级 (low, medium, high, urgent)",
				IsPrimary:   false,
				IsRequired:  true,
			},
			{
				Name:        "due_date",
				Type:        "DATE",
				Description: "任务截止日期",
				IsPrimary:   false,
				IsRequired:  false,
			},
			{
				Name:        "created_at",
				Type:        "TIMESTAMP",
				Description: "记录创建时间",
				IsPrimary:   false,
				IsRequired:  true,
			},
			{
				Name:        "updated_at",
				Type:        "TIMESTAMP",
				Description: "记录更新时间",
				IsPrimary:   false,
				IsRequired:  true,
			},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
	tables = append(tables, tasks)

	// 9. Task comments table
	taskComments := &models.Table{
		ID:          uuid.New().String(),
		Name:        "task_comments",
		Description: "任务评论表，存储任务相关的评论信息",
		Columns: []models.Column{
			{
				Name:        "id",
				Type:        "BIGINT",
				Description: "评论唯一标识符，主键",
				IsPrimary:   true,
				IsRequired:  true,
			},
			{
				Name:        "task_id",
				Type:        "BIGINT",
				Description: "任务ID，关联tasks表",
				IsPrimary:   false,
				IsRequired:  true,
			},
			{
				Name:        "user_id",
				Type:        "BIGINT",
				Description: "评论者ID，关联users表",
				IsPrimary:   false,
				IsRequired:  true,
			},
			{
				Name:        "content",
				Type:        "TEXT",
				Description: "评论内容",
				IsPrimary:   false,
				IsRequired:  true,
			},
			{
				Name:        "created_at",
				Type:        "TIMESTAMP",
				Description: "评论创建时间",
				IsPrimary:   false,
				IsRequired:  true,
			},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
	tables = append(tables, taskComments)

	// 10. Documents table
	documents := &models.Table{
		ID:          uuid.New().String(),
		Name:        "documents",
		Description: "文档信息表，存储项目相关的文档",
		Columns: []models.Column{
			{
				Name:        "id",
				Type:        "BIGINT",
				Description: "文档唯一标识符，主键",
				IsPrimary:   true,
				IsRequired:  true,
			},
			{
				Name:        "project_id",
				Type:        "BIGINT",
				Description: "项目ID，关联projects表",
				IsPrimary:   false,
				IsRequired:  true,
			},
			{
				Name:        "title",
				Type:        "VARCHAR(200)",
				Description: "文档标题",
				IsPrimary:   false,
				IsRequired:  true,
			},
			{
				Name:        "file_path",
				Type:        "VARCHAR(500)",
				Description: "文档文件路径",
				IsPrimary:   false,
				IsRequired:  true,
			},
			{
				Name:        "uploaded_by",
				Type:        "BIGINT",
				Description: "上传者ID，关联users表",
				IsPrimary:   false,
				IsRequired:  true,
			},
			{
				Name:        "version",
				Type:        "VARCHAR(20)",
				Description: "文档版本号",
				IsPrimary:   false,
				IsRequired:  true,
			},
			{
				Name:        "created_at",
				Type:        "TIMESTAMP",
				Description: "文档上传时间",
				IsPrimary:   false,
				IsRequired:  true,
			},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
	tables = append(tables, documents)

	// 11. Issues table
	issues := &models.Table{
		ID:          uuid.New().String(),
		Name:        "issues",
		Description: "问题跟踪表，记录项目中发现的问题",
		Columns: []models.Column{
			{
				Name:        "id",
				Type:        "BIGINT",
				Description: "问题唯一标识符，主键",
				IsPrimary:   true,
				IsRequired:  true,
			},
			{
				Name:        "project_id",
				Type:        "BIGINT",
				Description: "项目ID，关联projects表",
				IsPrimary:   false,
				IsRequired:  true,
			},
			{
				Name:        "title",
				Type:        "VARCHAR(200)",
				Description: "问题标题",
				IsPrimary:   false,
				IsRequired:  true,
			},
			{
				Name:        "description",
				Type:        "TEXT",
				Description: "问题详细描述",
				IsPrimary:   false,
				IsRequired:  true,
			},
			{
				Name:        "reporter_id",
				Type:        "BIGINT",
				Description: "报告者ID，关联users表",
				IsPrimary:   false,
				IsRequired:  true,
			},
			{
				Name:        "assignee_id",
				Type:        "BIGINT",
				Description: "处理者ID，关联users表",
				IsPrimary:   false,
				IsRequired:  false,
			},
			{
				Name:        "status",
				Type:        "VARCHAR(20)",
				Description: "问题状态 (open, in_progress, resolved, closed)",
				IsPrimary:   false,
				IsRequired:  true,
			},
			{
				Name:        "priority",
				Type:        "VARCHAR(20)",
				Description: "问题优先级 (low, medium, high, urgent)",
				IsPrimary:   false,
				IsRequired:  true,
			},
			{
				Name:        "created_at",
				Type:        "TIMESTAMP",
				Description: "问题报告时间",
				IsPrimary:   false,
				IsRequired:  true,
			},
			{
				Name:        "updated_at",
				Type:        "TIMESTAMP",
				Description: "问题更新时间",
				IsPrimary:   false,
				IsRequired:  true,
			},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
	tables = append(tables, issues)

	// 12. Issue comments table
	issueComments := &models.Table{
		ID:          uuid.New().String(),
		Name:        "issue_comments",
		Description: "问题评论表，存储问题相关的评论信息",
		Columns: []models.Column{
			{
				Name:        "id",
				Type:        "BIGINT",
				Description: "评论唯一标识符，主键",
				IsPrimary:   true,
				IsRequired:  true,
			},
			{
				Name:        "issue_id",
				Type:        "BIGINT",
				Description: "问题ID，关联issues表",
				IsPrimary:   false,
				IsRequired:  true,
			},
			{
				Name:        "user_id",
				Type:        "BIGINT",
				Description: "评论者ID，关联users表",
				IsPrimary:   false,
				IsRequired:  true,
			},
			{
				Name:        "content",
				Type:        "TEXT",
				Description: "评论内容",
				IsPrimary:   false,
				IsRequired:  true,
			},
			{
				Name:        "created_at",
				Type:        "TIMESTAMP",
				Description: "评论创建时间",
				IsPrimary:   false,
				IsRequired:  true,
			},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
	tables = append(tables, issueComments)

	// 13. Time tracking table
	timeTracking := &models.Table{
		ID:          uuid.New().String(),
		Name:        "time_tracking",
		Description: "时间跟踪表，记录用户在任务上花费的时间",
		Columns: []models.Column{
			{
				Name:        "id",
				Type:        "BIGINT",
				Description: "时间记录唯一标识符，主键",
				IsPrimary:   true,
				IsRequired:  true,
			},
			{
				Name:        "task_id",
				Type:        "BIGINT",
				Description: "任务ID，关联tasks表",
				IsPrimary:   false,
				IsRequired:  true,
			},
			{
				Name:        "user_id",
				Type:        "BIGINT",
				Description: "用户ID，关联users表",
				IsPrimary:   false,
				IsRequired:  true,
			},
			{
				Name:        "hours_spent",
				Type:        "DECIMAL(5,2)",
				Description: "花费的小时数",
				IsPrimary:   false,
				IsRequired:  true,
			},
			{
				Name:        "description",
				Type:        "TEXT",
				Description: "工作内容描述",
				IsPrimary:   false,
				IsRequired:  false,
			},
			{
				Name:        "tracked_at",
				Type:        "DATE",
				Description: "时间记录日期",
				IsPrimary:   false,
				IsRequired:  true,
			},
			{
				Name:        "created_at",
				Type:        "TIMESTAMP",
				Description: "记录创建时间",
				IsPrimary:   false,
				IsRequired:  true,
			},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
	tables = append(tables, timeTracking)

	// 14. Notifications table
	notifications := &models.Table{
		ID:          uuid.New().String(),
		Name:        "notifications",
		Description: "通知表，存储系统发送给用户的通知",
		Columns: []models.Column{
			{
				Name:        "id",
				Type:        "BIGINT",
				Description: "通知唯一标识符，主键",
				IsPrimary:   true,
				IsRequired:  true,
			},
			{
				Name:        "user_id",
				Type:        "BIGINT",
				Description: "接收用户ID，关联users表",
				IsPrimary:   false,
				IsRequired:  true,
			},
			{
				Name:        "title",
				Type:        "VARCHAR(200)",
				Description: "通知标题",
				IsPrimary:   false,
				IsRequired:  true,
			},
			{
				Name:        "content",
				Type:        "TEXT",
				Description: "通知内容",
				IsPrimary:   false,
				IsRequired:  true,
			},
			{
				Name:        "is_read",
				Type:        "BOOLEAN",
				Description: "是否已读",
				IsPrimary:   false,
				IsRequired:  true,
			},
			{
				Name:        "created_at",
				Type:        "TIMESTAMP",
				Description: "通知创建时间",
				IsPrimary:   false,
				IsRequired:  true,
			},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
	tables = append(tables, notifications)

	// 15. Audit logs table
	auditLogs := &models.Table{
		ID:          uuid.New().String(),
		Name:        "audit_logs",
		Description: "审计日志表，记录系统中的重要操作",
		Columns: []models.Column{
			{
				Name:        "id",
				Type:        "BIGINT",
				Description: "日志唯一标识符，主键",
				IsPrimary:   true,
				IsRequired:  true,
			},
			{
				Name:        "user_id",
				Type:        "BIGINT",
				Description: "操作用户ID，关联users表",
				IsPrimary:   false,
				IsRequired:  true,
			},
			{
				Name:        "action",
				Type:        "VARCHAR(100)",
				Description: "操作类型 (create, update, delete等)",
				IsPrimary:   false,
				IsRequired:  true,
			},
			{
				Name:        "table_name",
				Type:        "VARCHAR(100)",
				Description: "涉及的表名",
				IsPrimary:   false,
				IsRequired:  true,
			},
			{
				Name:        "record_id",
				Type:        "VARCHAR(100)",
				Description: "涉及的记录ID",
				IsPrimary:   false,
				IsRequired:  true,
			},
			{
				Name:        "old_values",
				Type:        "JSON",
				Description: "修改前的值，以JSON格式存储",
				IsPrimary:   false,
				IsRequired:  false,
			},
			{
				Name:        "new_values",
				Type:        "JSON",
				Description: "修改后的值，以JSON格式存储",
				IsPrimary:   false,
				IsRequired:  false,
			},
			{
				Name:        "created_at",
				Type:        "TIMESTAMP",
				Description: "操作时间",
				IsPrimary:   false,
				IsRequired:  true,
			},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
	tables = append(tables, auditLogs)

	return tables
}

// TestGenerateSampleTables tests the sample tables generation
func TestGenerateSampleTables(t *testing.T) {
	tables := GenerateSampleTables()

	// 检查生成的表数量
	if len(tables) != 15 {
		t.Errorf("Expected 15 tables, got %d", len(tables))
	}

	// 检查特定的关键表是否存在
	requiredTables := map[string]bool{
		"users":        false,
		"departments":  false,
		"user_roles":   false,
		"projects":     false,
		"tasks":        false,
		"issues":       false,
		"audit_logs":   false,
		"notifications":false,
	}

	for _, table := range tables {
		if _, required := requiredTables[table.Name]; required {
			requiredTables[table.Name] = true
		}
	}

	for tableName, found := range requiredTables {
		if !found {
			t.Errorf("Required table '%s' not found in generated tables", tableName)
		}
	}

	// 检查 users 表的列定义
	usersTable := findTableByName(tables, "users")
	if usersTable == nil {
		t.Fatal("users table not found")
	}

	// 验证 users 表有正确的列
	expectedUserColumns := []string{"id", "username", "email", "full_name", "created_at"}
	validateTableColumns(t, usersTable, expectedUserColumns)

	// 检查每个表都有创建和更新时间
	for _, table := range tables {
		if table.CreatedAt.IsZero() {
			t.Errorf("Table %s has zero CreatedAt", table.Name)
		}
		if table.UpdatedAt.IsZero() {
			t.Errorf("Table %s has zero UpdatedAt", table.Name)
		}
	}

	// 检查每个表都有ID
	for _, table := range tables {
		if table.ID == "" {
			t.Errorf("Table %s has empty ID", table.Name)
		}
	}
}

// 辅助函数：根据表名查找表
func findTableByName(tables []*models.Table, name string) *models.Table {
	for _, table := range tables {
		if table.Name == name {
			return table
		}
	}
	return nil
}

// 辅助函数：验证表是否包含预期的列
func validateTableColumns(t *testing.T, table *models.Table, expectedColumns []string) {
	t.Helper()
	columnMap := make(map[string]bool)
	for _, col := range table.Columns {
		columnMap[col.Name] = true
	}

	for _, expectedCol := range expectedColumns {
		if !columnMap[expectedCol] {
			t.Errorf("Expected column %s not found in table %s", expectedCol, table.Name)
		}
	}
}

// TestTableStructureValidation 测试表结构验证
func TestTableStructureValidation(t *testing.T) {
	tables := GenerateSampleTables()

	// 检查是否有重复的表名
	tableNames := make(map[string]bool)
	for _, table := range tables {
		if tableNames[table.Name] {
			t.Errorf("Duplicate table name found: %s", table.Name)
		}
		tableNames[table.Name] = true
	}

	// 检查每个表至少有一个主键列
	for _, table := range tables {
		hasPrimaryKey := false
		for _, col := range table.Columns {
			if col.IsPrimary {
				hasPrimaryKey = true
				break
			}
		}
		if !hasPrimaryKey {
			t.Errorf("Table %s has no primary key", table.Name)
		}
	}

	// 检查每个表的必需字段是否正确设置
	for _, table := range tables {
		for _, col := range table.Columns {
			if col.Name == "" {
				t.Errorf("Table %s has column with empty name", table.Name)
			}
			if col.Type == "" {
				t.Errorf("Table %s has column %s with empty type", table.Name, col.Name)
			}
		}
	}
}