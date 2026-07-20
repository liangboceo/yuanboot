package service

import "sendex-server/internal/model"

type systemMenuSeed struct {
	ID         uint
	Name       string
	ParentID   uint
	OrderNum   int
	Path       string
	Component  string
	MenuType   int
	Icon       string
	Status     int
	Visible    int
	IsFrame    int
	IsCache    int
	Permission string
	Query      string
	Perms      string
}

func systemMenuSeeds() []systemMenuSeed {
	return []systemMenuSeed{
		menuDirectory(1, "系统管理", 0, 1, "/system", "ri/settings-3-line"),
		menuPage(100, "用户管理", 1, 1, "user", "system/user/index", "ep:avatar"),
		menuPage(101, "角色管理", 1, 2, "role", "system/role/index", "ri/admin-line"),
		menuPage(102, "菜单管理", 1, 3, "menu", "system/menu/index", "ep:menu"),
		menuPage(107, "系统配置", 1, 4, "system-config", "system/system-config/index", "ri/settings-2-line"),
	}
}

func menuDirectory(id uint, name string, parentId uint, orderNum int, path string, icon string) systemMenuSeed {
	return systemMenuSeed{
		ID:        id,
		Name:      name,
		ParentID:  parentId,
		OrderNum:  orderNum,
		Path:      path,
		Component: "Layout",
		MenuType:  0,
		Icon:      icon,
		Status:    1,
	}
}

func menuPage(id uint, name string, parentId uint, orderNum int, path string, component string, icon string) systemMenuSeed {
	return systemMenuSeed{
		ID:        id,
		Name:      name,
		ParentID:  parentId,
		OrderNum:  orderNum,
		Path:      path,
		Component: component,
		MenuType:  1,
		Icon:      icon,
		Status:    1,
	}
}

func innerPage(id uint, name string, parentId uint, orderNum int, path string, component string, icon string) systemMenuSeed {
	seed := menuPage(id, name, parentId, orderNum, path, component, icon)
	seed.MenuType = 3
	seed.Visible = 1
	return seed
}

func (seed systemMenuSeed) toModel() model.SysMenu {
	return model.SysMenu{
		ID:         seed.ID,
		Name:       seed.Name,
		ParentID:   seed.ParentID,
		OrderNum:   seed.OrderNum,
		Path:       seed.Path,
		Component:  seed.Component,
		MenuType:   seed.MenuType,
		Visible:    seed.Visible,
		Status:     seed.Status,
		Icon:       seed.Icon,
		IsFrame:    seed.IsFrame,
		IsCache:    seed.IsCache,
		Permission: seed.Permission,
		Query:      seed.Query,
		Perms:      seed.Perms,
	}
}
