import { http } from "@/utils/http";

type Result = {
  data: any;
  code: number;
  status: number;
  message: string;
};

type UserInfoResult = {
  id: number;
  username: string;
  nickname: string;
  avatar: string;
  email: string;
  phone: string;
  deptId: number;
  deptName: string;
  roles: Array<string>;
  permissions: Array<string>;
};

type ResultTable = {
  /** 列表数据 */
  list: Array<any>;
  /** 总条目数 */
  total?: number;
  /** 每页显示条目个数 */
  pageSize?: number;
  /** 当前页数 */
  currentPage?: number;
};

/** 登录 */
export const loginApi = (data?: object) => {
  return http.request<Result>("post", "sys/Login", { data });
};

/** 获取当前用户信息 */
export const getUserInfoApi = () => {
  return http.request<UserInfoResult>("post", "sys/GetUserInfo");
};

/** 获取路由权限 */
export const getRoutesApi = () => {
  return http.request<Result>("post", "sys/GetRoutes");
};

// ==================== 用户管理 ====================

/** 获取用户管理列表 */
export const getUserList = (data?: object) => {
  return http.request<ResultTable>("post", "sys/getUserList", { data });
};

/** 获取用户详情 */
export const getUserDetail = (data?: object) => {
  return http.request<Result>("post", "sys/getUserDetail", { data });
};

/** 保存用户 */
export const saveUser = (data?: object) => {
  return http.request<Result>("post", "sys/saveUser", { data });
};

/** 删除用户 */
export const delUser = (data?: object) => {
  return http.request<Result>("post", "sys/delUser", { data });
};

/** 重置密码 */
export const resetPwd = (data?: object) => {
  return http.request<Result>("post", "sys/resetPwd", { data });
};

/** 上传用户头像 */
export const uploadUserAvatar = (file: File) => {
  const formData = new FormData();
  formData.append("file", file);
  return http.upload<{ avatar: string }>("sys/UploadUserAvatar", formData);
};

/** 修改当前用户密码 */
export const changePwd = (data?: object) => {
  return http.request<Result>("post", "sys/changePwd", { data });
};

// ==================== 角色管理 ====================

/** 获取角色管理列表 */
export const getRoleList = (data?: object) => {
  return http.request<ResultTable>("post", "sys/getRoleList", { data });
};

/** 获取所有角色列表 */
export const getAllRoleList = () => {
  return http.request<Result>("post", "sys/getAllRoleList");
};

/** 获取角色详情 */
export const getRoleDetail = (data?: object) => {
  return http.request<Result>("post", "sys/getRoleDetail", { data });
};

/** 保存角色 */
export const saveRole = (data?: object) => {
  return http.request<Result>("post", "sys/saveRole", { data });
};

/** 删除角色 */
export const delRole = (data?: object) => {
  return http.request<Result>("post", "sys/delRole", { data });
};

/** 获取角色菜单权限 */
export const getRoleMenu = (data?: object) => {
  return http.request<Result>("post", "sys/getRoleMenu", { data });
};

// ==================== 菜单管理 ====================

/** 获取菜单管理列表 */
export const getMenuList = (data?: object) => {
  return http.request<Result>("post", "sys/getMenuList", { data });
};

/** 获取菜单树 */
export const getMenuTree = () => {
  return http.request<Result>("post", "sys/getMenuTree");
};

/** 获取菜单详情 */
export const getMenuDetail = (data?: object) => {
  return http.request<Result>("post", "sys/getMenuDetail", { data });
};

/** 保存菜单 */
export const saveMenu = (data?: object) => {
  return http.request<Result>("post", "sys/saveMenu", { data });
};

/** 删除菜单 */
export const delMenu = (data?: object) => {
  return http.request<Result>("post", "sys/delMenu", { data });
};

// ==================== 部门管理 ====================

/** 获取部门管理列表 */
export const getDeptList = (data?: object) => {
  return http.request<Result>("post", "sys/getDeptList", { data });
};

/** 获取部门树 */
export const getDeptTree = () => {
  return http.request<Result>("post", "sys/getDeptTree");
};

/** 获取部门详情 */
export const getDeptDetail = (data?: object) => {
  return http.request<Result>("post", "sys/getDeptDetail", { data });
};

/** 保存部门 */
export const saveDept = (data?: object) => {
  return http.request<Result>("post", "sys/saveDept", { data });
};

/** 删除部门 */
export const delDept = (data?: object) => {
  return http.request<Result>("post", "sys/delDept", { data });
};
