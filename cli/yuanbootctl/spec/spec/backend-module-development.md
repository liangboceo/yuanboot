# 后端模块开发规范

本文用于指导 `{{.ModelName}}` 后续业务模块开发，目标是让新模块在 Controller、DTO、Model、Service、缓存/锁、响应结构、查询分页、软删除、审计字段等方面保持一致。

## 1. 目录与分层约定

后端模块默认遵循现有目录结构：

```text
{{.ModelName}}/
├── internal/
│   ├── controller/   # HTTP Controller，负责参数校验、幂等锁、响应组装
│   ├── dto/          # 请求/响应 DTO，负责路由声明、入参字段、出参字段
│   ├── model/        # GORM 数据模型，负责表结构、JSON 字段、TableName
│   ├── service/      # 业务服务，负责数据库查询、事务、关联关系维护
│   ├── response/     # 统一响应
│   └── middleware/   # 中间件
├── pkg/utils/        # 登录态、通用工具
└── spec/             # 开发规范文档
```

开发一个新模块时至少包含：

- `internal/model/*_model.go`：新增模型结构体。
- `internal/dto/*_dto.go`：新增请求和响应 DTO。
- `internal/service/*_service.go` 或复用已有 service 文件：新增业务方法。
- `internal/controller/*_controller.go` 或复用已有 controller 文件：新增接口方法。

## 2. 命名规范

以模块 `Example` 为例：

| 类型 | 命名示例 | 说明 |
|---|---|---|
| Model | `Example` / `SysExample` | 与表名对应 |
| TableName | `example` / `sys_example` | 使用 snake_case |
| ListReq | `GetExampleListReq` | 必须与 Controller 方法 `GetExampleList` 一致并追加 `Req` |
| DetailReq | `GetExampleDetailReq` | 必须与 Controller 方法 `GetExampleDetail` 一致并追加 `Req` |
| SaveReq | `SaveExampleReq` | 必须与 Controller 方法 `SaveExample` 一致并追加 `Req` |
| DelReq | `DelExampleReq` | 必须与 Controller 方法 `DelExample` 一致并追加 `Req` |
| Resp | `ExampleResp` | 单条响应 |
| ListResp | `ExampleListResp` | 分页列表响应 |
| Controller方法 | `GetExampleList` / `SaveExample` / `DelExample` | 动宾结构 |
| Service方法 | `GetExampleList` / `GetExampleByID` / `CreateExample` / `UpdateExample` / `DeleteExample` | 与 Controller 对应 |

## 3. Model 标准

### 3.1 基础字段

业务模型优先嵌入现有基础模型：

```go
type Example struct {
	model.BaseModel
	ID     uint   `gorm:"type:int(11);primaryKey;column:id;comment:主键;AUTO_INCREMENT" json:"id"`
	Name   string `gorm:"type:varchar(100);column:name;comment:名称;not null" json:"name"`
	Status int    `gorm:"type:tinyint(2);column:status;comment:状态【0：禁用 1：正常】;default:1" json:"status"`
	Remark string `gorm:"type:varchar(500);column:remark;comment:备注;default:''" json:"remark"`
	model.CommonModel
}

func (Example) TableName() string {
	return "example"
}
```

### 3.2 字段要求

- 主键统一使用 `uint` 类型 `ID`。
- 软删除字段使用基础模型中的 `is_deleted`，查询默认加 `is_deleted=0`。
- 状态字段统一：`0` 禁用，`1` 正常。
- 排序字段优先命名为 `sort` 或现有业务约定字段，如菜单使用 `order_num`。
- 创建/更新审计字段由 `CommonModel` 承载，Controller 写入当前用户信息。
- 外键字段使用 `XxxID`，JSON 使用 `xxxId`。
- 关联子节点非持久化字段加 `gorm:"-"`。

## 4. DTO 标准

### 4.1 路由声明

本项目使用 Yuanboot MVC 请求对象声明路由。Controller 方法入参必须封装为 DTO 请求结构体，禁止直接使用匿名 `struct`、散装基础类型或从 `ctx` 手工取业务参数。

路由必须包含在请求结构体的 `mvc.RequestBody` tag 中，并且 `route` 必须与 Controller 结构体名和方法名保持一致：

```go
type GetRoutesReq struct {
	mvc.RequestBody `route:"/sys/GetRoutes" doc:"获取用户路由权限"`
}

func (c *SysController) GetRoutes(req *dto.GetRoutesReq, ctx *context.HttpContext) actionresult.IActionResult {
	// ...
}
```

路由命名规则：

```text
/<ControllerName去掉Controller后缀>/<MethodName>
```

示例：

| Controller | Method | Req Struct | route |
|---|---|---|---|
| `SysController` | `GetRoutes` | `GetRoutesReq` | `/sys/GetRoutes` |
| `SysController` | `GetUserInfo` | `GetUserInfoReq` | `/sys/GetUserInfo` |
| `ExampleController` | `GetExampleList` | `GetExampleListReq` | `/example/GetExampleList` |
| `ExampleController` | `SaveExample` | `SaveExampleReq` | `/example/SaveExample` |
| `ExampleController` | `DelExample` | `DelExampleReq` | `/example/DelExample` |

请求结构体示例：

```go
type GetExampleListReq struct {
	mvc.RequestBody `route:"/example/GetExampleList" doc:"获取示例列表"`
	PageNum         int    `json:"pageNum" doc:"分页页码" default:"1"`
	PageSize        int    `json:"pageSize" doc:"每页数量" default:"20"`
	Name            string `json:"name" doc:"名称"`
	Status          int    `json:"status" doc:"状态"`
}
```

已有历史接口如不符合该规则，后续修改时必须逐步迁移到该标准。新增接口必须按该标准执行。

### 4.2 请求 DTO

```go
type SaveExampleReq struct {
	mvc.RequestBody `route:"/example/SaveExample" doc:"保存示例"`
	ID              uint   `json:"id" doc:"主键（更新时必填）"`
	Name            string `json:"name" doc:"名称" binding:"required"`
	Sort            int    `json:"sort" doc:"排序" default:"1"`
	Status          int    `json:"status" doc:"状态" default:"1"`
	Remark          string `json:"remark" doc:"备注"`
}

type GetExampleDetailReq struct {
	mvc.RequestBody `route:"/example/GetExampleDetail" doc:"获取示例详情"`
	ID              uint `json:"id" doc:"主键" binding:"required"`
}

type DelExampleReq struct {
	mvc.RequestBody `route:"/example/DelExample" doc:"删除示例"`
	ID              uint `json:"id" doc:"主键" binding:"required"`
}
```

要求：

- Controller 方法入参必须是 DTO 请求结构体指针，例如 `req *dto.SaveExampleReq`。
- 请求结构体必须嵌入 `mvc.RequestBody`，并声明 `route` 和 `doc`。
- 请求结构体名称必须与 Controller 方法名一致并追加 `Req`，例如方法 `SaveExample` 对应 `SaveExampleReq`。
- `route` 必须与 Controller 名和方法名一致，例如 `ExampleController.SaveExample` 对应 `/example/SaveExample`。
- 必填字段加 `binding:"required"`。
- 所有业务字段必须有 `json` 和 `doc` tag。
- 新增和修改优先共用 `SaveXxxReq`，通过 `ID > 0` 判断更新。
- 删除请求必须包含 `ID`，批量删除可定义 `IDs []uint`。
- 分页字段统一为 `pageNum`、`pageSize`。

### 4.3 响应 DTO

```go
type ExampleResp struct {
	dto.BaseResp
	Name   string `json:"name" doc:"名称"`
	Sort   int    `json:"sort" doc:"排序"`
	Status int    `json:"status" doc:"状态"`
	Remark string `json:"remark" doc:"备注"`
}

type ExampleListResp struct {
	dto.BasePageResp
	Data []ExampleResp `json:"list" doc:"示例列表"`
}
```

要求：

- 列表响应统一返回 `list` 字段。
- 分页响应嵌入 `BasePageResp`，设置 `Page` 和 `Total`。
- 不直接返回数据库模型给前端，统一转换为 Resp DTO。
- 密码、密钥、Token、内部配置等敏感字段禁止出现在响应 DTO。

## 5. Controller 标准

Controller 负责：参数校验、默认分页参数、幂等性锁、调用 Service、DTO 转换、统一响应。

### 5.1 Controller 结构

```go
type ExampleController struct {
	mvc.ApiController `doc:"示例管理"`
	ExampleService    *service.ExampleService
	CacheService      *service.CacheService
}

func NewExampleController(exampleService *service.ExampleService, cacheService *service.CacheService) *ExampleController {
	return &ExampleController{ExampleService: exampleService, CacheService: cacheService}
}
```

如果模块放在现有 `SysController` 中，则沿用已有结构和 `lockSysOperation`。

### 5.2 统一响应

Controller 必须使用：

```go
return response.SuccessMessage(data, "操作成功")
return response.FailureMessage(nil, "错误原因")
```

不直接写 `ctx.JSON`，除非是中间件或特殊流式响应。

### 5.3 参数校验

保存/删除/关键动作入口必须校验：

```go
if err := validator.New().Struct(req); err != nil {
	return response.FailureMessage(nil, "参数错误")
}
```

简单 ID 查询也必须显式判断：

```go
if req.ID == 0 {
	return response.FailureMessage(nil, "ID不能为空")
}
```

### 5.4 分页列表 Controller

```go
func (c *ExampleController) GetExampleList(req *dto.GetExampleListReq) actionresult.IActionResult {
	if req.PageNum <= 0 {
		req.PageNum = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	items, total, err := c.ExampleService.GetExampleList(req.Name, req.Status, req.PageNum, req.PageSize)
	if err != nil {
		return response.FailureMessage(nil, "获取示例列表失败")
	}

	result := &dto.ExampleListResp{
		BasePageResp: dto.BasePageResp{Page: req.PageNum, Total: total},
	}
	result.Data = make([]dto.ExampleResp, len(items))
	for i, item := range items {
		result.Data[i] = dto.ExampleResp{
			BaseResp: dto.BaseResp{ID: item.ID, CreateTime: item.CreateTime, UpdateTime: item.UpdateTime},
			Name:     item.Name,
			Sort:     item.Sort,
			Status:   item.Status,
			Remark:   item.Remark,
		}
	}

	return response.SuccessMessage(result, "获取示例列表成功")
}
```

### 5.5 详情 Controller

```go
func (c *ExampleController) GetExampleDetail(req *dto.GetExampleDetailReq) actionresult.IActionResult {
	if req.ID == 0 {
		return response.FailureMessage(nil, "示例ID不能为空")
	}

	item, err := c.ExampleService.GetExampleByID(req.ID)
	if err != nil || item.ID == 0 {
		return response.FailureMessage(nil, "示例不存在")
	}

	resp := dto.ExampleResp{
		BaseResp: dto.BaseResp{ID: item.ID, CreateTime: item.CreateTime, UpdateTime: item.UpdateTime},
		Name:     item.Name,
		Sort:     item.Sort,
		Status:   item.Status,
		Remark:   item.Remark,
	}

	return response.SuccessMessage(resp, "获取示例详情成功")
}
```

### 5.6 新增/修改 Controller

新增、修改必须加幂等性锁，避免重复提交。

```go
func (c *ExampleController) SaveExample(saveReq *dto.SaveExampleReq, ctx *context.HttpContext) actionresult.IActionResult {
	if err := validator.New().Struct(saveReq); err != nil {
		return response.FailureMessage(nil, "参数错误")
	}

	lockKey, lock := c.lockExampleOperation("SaveExample", saveReq.Name, saveReq.ID)
	if !lock {
		return response.FailureMessage(nil, "操作频繁")
	}
	defer c.CacheService.Unlock(lockKey)

	if saveReq.ID > 0 {
		item, err := c.ExampleService.GetExampleByID(saveReq.ID)
		if err != nil || item.ID == 0 {
			return response.FailureMessage(nil, "示例不存在")
		}

		item.Name = saveReq.Name
		item.Sort = saveReq.Sort
		item.Status = saveReq.Status
		item.Remark = saveReq.Remark
		item.UpdateBy = utils.GetUserId(ctx)
		item.UpdateName = utils.GetUserName(ctx)
		item.UpdateTime = time.Now()

		if err := c.ExampleService.UpdateExample(item); err != nil {
			return response.FailureMessage(nil, "更新示例失败")
		}
		return response.SuccessMessage(nil, "更新示例成功")
	}

	item := &model.Example{
		Name:   saveReq.Name,
		Sort:   saveReq.Sort,
		Status: saveReq.Status,
		Remark: saveReq.Remark,
		CommonModel: model.CommonModel{
			CreateBy:   utils.GetUserId(ctx),
			CreateName: utils.GetUserName(ctx),
		},
	}

	if err := c.ExampleService.CreateExample(item); err != nil {
		return response.FailureMessage(nil, "创建示例失败")
	}
	return response.SuccessMessage(nil, "创建示例成功")
}
```

### 5.7 删除 Controller

删除也必须加幂等性锁，并统一软删除。

```go
func (c *ExampleController) DelExample(req *dto.DelExampleReq, ctx *context.HttpContext) actionresult.IActionResult {
	if req.ID == 0 {
		return response.FailureMessage(nil, "示例ID不能为空")
	}

	lockKey, lock := c.lockExampleOperation("DelExample", "example", req.ID)
	if !lock {
		return response.FailureMessage(nil, "操作频繁")
	}
	defer c.CacheService.Unlock(lockKey)

	item, err := c.ExampleService.GetExampleByID(req.ID)
	if err != nil || item.ID == 0 {
		return response.FailureMessage(nil, "示例不存在")
	}

	item.IsDeleted = 1
	item.UpdateBy = utils.GetUserId(ctx)
	item.UpdateName = utils.GetUserName(ctx)
	item.UpdateTime = time.Now()

	if err := c.ExampleService.DeleteExample(req.ID); err != nil {
		return response.FailureMessage(nil, err.Error())
	}

	return response.SuccessMessage(nil, "删除示例成功")
}
```

## 6. 幂等性锁标准

### 6.1 锁 Key 标准

沿用现有缓存锁：

```go
const DATALOCK = "{{.ModelName}}:bridge:data:lock:%s:%s:%s"
```

Controller 中封装：

```go
func (c *ExampleController) lockExampleOperation(module, key string, id uint) (string, bool) {
	lockKey := fmt.Sprintf(service.DATALOCK, module, key, strconv.FormatUint(uint64(id), 10))
	return lockKey, c.CacheService.Lock(lockKey, 10)
}
```

现有 `SysController` 已有：

```go
func (c *SysController) lockSysOperation(module, key string, id uint) (string, bool)
```

### 6.2 必须加锁的接口

- 新增接口。
- 修改接口。
- 删除接口。
- 重置密码、启停用、授权、绑定关系变更等有副作用动作。

### 6.3 锁使用要求

- 获取锁失败统一返回：`操作频繁`。
- 获取锁成功后必须 `defer Unlock`。
- 锁超时时间默认 10 秒。
- 锁 Key 至少包含：动作名、业务唯一键、ID。
- 查询类接口不加锁。

## 7. Service 标准

Service 负责数据库访问、事务、业务关联、删除前约束检查，不负责 HTTP 响应。

### 7.1 获取详情

```go
func (s *ExampleService) GetExampleByID(id uint) (*model.Example, error) {
	var item model.Example
	if err := s.DbService.Db.Where("id=? and is_deleted=?", id, 0).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}
```

### 7.2 分页查询

```go
func (s *ExampleService) GetExampleList(name string, status int, pageNum, pageSize int) ([]*model.Example, int64, error) {
	var items []*model.Example
	var total int64

	query := s.DbService.Db.Model(&model.Example{}).Where("is_deleted=?", 0)
	if name != "" {
		query = query.Where("name like ?", "%"+name+"%")
	}
	if status > 0 {
		query = query.Where("status=?", status)
	}

	query.Count(&total)
	err := query.Offset((pageNum - 1) * pageSize).Order("sort asc, id desc").Limit(pageSize).Find(&items).Error
	return items, total, err
}
```

要求：

- 查询默认加 `is_deleted=0`。
- 模糊查询统一使用 `like` + `%keyword%`。
- 状态查询沿用现有习惯：`status > 0` 才追加条件。如需查询禁用状态，应改为指针或约定 `-1` 表示全部。
- 分页先 `Count` 再 `Offset/Limit/Find`。
- 排序明确写出，不能依赖数据库默认排序。

### 7.3 全量启用列表

```go
func (s *ExampleService) GetAllExamples() ([]*model.Example, error) {
	var items []*model.Example
	err := s.DbService.Db.Where("status=1 and is_deleted=0").Order("sort asc, id desc").Find(&items).Error
	return items, err
}
```

### 7.4 新增

```go
func (s *ExampleService) CreateExample(item *model.Example) error {
	return s.DbService.Db.Create(item).Error
}
```

有子表/关联表时必须使用事务：

```go
func (s *ExampleService) CreateExample(item *model.Example, relationIds []uint) error {
	return s.DbService.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(item).Error; err != nil {
			return err
		}
		for _, relationId := range relationIds {
			if err := tx.Create(&model.ExampleRelation{ExampleID: item.ID, RelationID: relationId}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
```

### 7.5 修改

```go
func (s *ExampleService) UpdateExample(item *model.Example) error {
	return s.DbService.Db.Save(item).Error
}
```

有关联表时：

- 事务内先保存主表。
- 删除原有关联。
- 写入新关联。
- 任一失败回滚。

```go
func (s *ExampleService) UpdateExample(item *model.Example, relationIds []uint) error {
	return s.DbService.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(item).Error; err != nil {
			return err
		}
		if err := tx.Where("example_id=?", item.ID).Delete(&model.ExampleRelation{}).Error; err != nil {
			return err
		}
		for _, relationId := range relationIds {
			if err := tx.Create(&model.ExampleRelation{ExampleID: item.ID, RelationID: relationId}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
```

### 7.6 删除

删除必须软删除：

```go
func (s *ExampleService) DeleteExample(id uint) error {
	return s.DbService.Db.Model(&model.Example{}).Where("id=?", id).Update("is_deleted", 1).Error
}
```

删除前如存在依赖关系必须检查：

```go
func (s *ExampleService) DeleteExample(id uint) error {
	var count int64
	s.DbService.Db.Model(&model.Child{}).Where("example_id=? and is_deleted=0", id).Count(&count)
	if count > 0 {
		return errors.New("存在子数据，无法删除")
	}
	return s.DbService.Db.Model(&model.Example{}).Where("id=?", id).Update("is_deleted", 1).Error
}
```

## 8. 树结构模块标准

适用于菜单、部门、分类等父子结构。

### 8.1 Model 字段

```go
ParentID uint       `gorm:"type:int(11);column:parent_id;comment:父级ID;default:0" json:"parentId"`
Children []*Example `gorm:"-" json:"children"`
```

### 8.2 查询排序

```go
err := query.Order("parent_id asc, sort asc, id asc").Find(&items).Error
```

### 8.3 构树

Service 或 Controller 中统一构树：

```go
func (s *ExampleService) BuildExampleTree(items []*model.Example) []*model.Example {
	itemMap := make(map[uint]*model.Example)
	var roots []*model.Example
	for _, item := range items {
		itemMap[item.ID] = item
	}
	for _, item := range items {
		if item.ParentID == 0 {
			roots = append(roots, item)
			continue
		}
		if parent, ok := itemMap[item.ParentID]; ok {
			parent.Children = append(parent.Children, item)
		} else {
			roots = append(roots, item)
		}
	}
	return roots
}
```

### 8.4 删除约束

树节点删除前必须检查是否存在未删除子节点：

```go
var count int64
s.DbService.Db.Model(&model.Example{}).Where("parent_id=? and is_deleted=0", id).Count(&count)
if count > 0 {
	return errors.New("存在子节点，无法删除")
}
```

## 9. 审计字段标准

新增时写入：

```go
CommonModel: model.CommonModel{
	CreateBy:   utils.GetUserId(ctx),
	CreateName: utils.GetUserName(ctx),
}
```

修改/删除时写入：

```go
item.UpdateBy = utils.GetUserId(ctx)
item.UpdateName = utils.GetUserName(ctx)
item.UpdateTime = time.Now()
```

要求：

- 所有有副作用操作都必须带 `ctx *context.HttpContext`。
- 不允许在 Service 中直接读取 HTTP 上下文。
- Controller 负责将审计信息赋值到 Model。

## 10. 缓存标准

### 10.1 使用场景

- 登录用户信息。
- 高频读取且变更较少的数据。
- 分布式锁。

### 10.2 缓存失效

凡是修改了缓存对应数据，必须清理缓存：

```go
s.DelLoginUserCache(user.Username)
```

要求：

- 用户启停用、修改密码、删除用户后必须清理登录缓存。
- 菜单/权限/角色变更后，如存在权限缓存，必须清理对应用户或全局权限缓存。
- 登录、权限等安全敏感场景优先读取数据库最新状态，再按需刷新缓存。

## 11. 错误处理标准

### 11.1 Controller 错误消息

- 参数错误：`参数错误`
- ID 为空：`xxxID不能为空`
- 数据不存在：`xxx不存在`
- 列表失败：`获取xxx列表失败`
- 详情失败：`获取xxx详情失败`
- 新增失败：`创建xxx失败`
- 修改失败：`更新xxx失败`
- 删除失败：优先返回 Service 的业务错误 `err.Error()`
- 操作频繁：`操作频繁`

### 11.2 Service 错误

Service 只返回 `error`，不引用 `response`。

```go
return errors.New("存在子数据，无法删除")
```

## 12. 安全标准

- 密码必须加密存储，使用现有 `EncryptPassword` / `VerifyPassword`。
- 密码、密钥、Token 不允许返回前端。
- 登录、修改密码、启停用用户必须读取数据库最新状态，避免缓存绕过。
- 删除、授权、重置密码等敏感动作必须加锁。
- 不记录明文密码、Token、密钥到日志。
- 后端必须做关键规则校验，不能只依赖前端校验。

## 13. 开发检查清单

新增模块完成前逐项检查：

- [ ] Model 包含 `TableName()`。
- [ ] 查询默认过滤 `is_deleted=0`。
- [ ] DTO 字段包含 `json`、`doc`，必填字段包含 `binding:"required"`。
- [ ] 列表请求包含 `pageNum`、`pageSize`（树/全量列表除外）。
- [ ] Controller 列表设置默认分页。
- [ ] Controller 保存/删除/副作用动作加幂等性锁。
- [ ] Controller 使用 `response.SuccessMessage` / `response.FailureMessage`。
- [ ] 新增写入 `CreateBy`、`CreateName`。
- [ ] 修改/删除写入 `UpdateBy`、`UpdateName`、`UpdateTime`。
- [ ] 删除使用软删除，不物理删除主表。
- [ ] 删除前检查子数据或关联数据。
- [ ] 多表写入使用事务。
- [ ] 响应不返回敏感字段。
- [ ] 修改缓存关联数据后清理缓存。
- [ ] 运行 `gofmt`。
- [ ] 运行 `go test ./...`。

## 14. 推荐实现顺序

1. 定义 Model 和 `TableName()`。
2. 定义 DTO：ListReq、DetailReq、SaveReq、DelReq、Resp、ListResp。
3. 编写 Service：ByID、List、All、Create、Update、Delete。
4. 编写 Controller：List、Detail、Save、Delete。
5. 为 Save/Delete 增加幂等性锁。
6. 补充树结构、关联表、缓存失效等特殊逻辑。
7. `gofmt` 格式化。
8. `go test ./...` 验证。

## 15. 最小模块模板

```go
// DTO
type GetExampleListReq struct {
	mvc.RequestBody `route:"/example/GetExampleList" doc:"获取示例列表"`
	PageNum         int    `json:"pageNum" doc:"分页页码" default:"1"`
	PageSize        int    `json:"pageSize" doc:"每页数量" default:"20"`
	Name            string `json:"name" doc:"名称"`
	Status          int    `json:"status" doc:"状态"`
}

type SaveExampleReq struct {
	mvc.RequestBody `route:"/example/SaveExample" doc:"保存示例"`
	ID              uint   `json:"id" doc:"主键（更新时必填）"`
	Name            string `json:"name" doc:"名称" binding:"required"`
	Sort            int    `json:"sort" doc:"排序" default:"1"`
	Status          int    `json:"status" doc:"状态" default:"1"`
	Remark          string `json:"remark" doc:"备注"`
}

type GetExampleDetailReq struct {
	mvc.RequestBody `route:"/example/GetExampleDetail" doc:"获取示例详情"`
	ID              uint `json:"id" doc:"主键" binding:"required"`
}

type DelExampleReq struct {
	mvc.RequestBody `route:"/example/DelExample" doc:"删除示例"`
	ID              uint `json:"id" doc:"主键" binding:"required"`
}
```

```go
// Service
func (s *ExampleService) GetExampleByID(id uint) (*model.Example, error) {
	var item model.Example
	if err := s.DbService.Db.Where("id=? and is_deleted=?", id, 0).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *ExampleService) GetExampleList(name string, status int, pageNum, pageSize int) ([]*model.Example, int64, error) {
	var items []*model.Example
	var total int64
	query := s.DbService.Db.Model(&model.Example{}).Where("is_deleted=?", 0)
	if name != "" {
		query = query.Where("name like ?", "%"+name+"%")
	}
	if status > 0 {
		query = query.Where("status=?", status)
	}
	query.Count(&total)
	err := query.Offset((pageNum - 1) * pageSize).Order("sort asc, id desc").Limit(pageSize).Find(&items).Error
	return items, total, err
}
```

```go
// Controller save/delete 必须包含锁
lockKey, lock := c.lockExampleOperation("SaveExample", saveReq.Name, saveReq.ID)
if !lock {
	return response.FailureMessage(nil, "操作频繁")
}
defer c.CacheService.Unlock(lockKey)
```
