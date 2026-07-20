import { http } from "@/utils/http";

export type UserResult = {
  /** 访问令牌 */
  accessToken: string;
  /** 过期时间戳 */
  expires: number;
  /** 刷新令牌 */
  refreshToken: string;
  /** 头像 */
  avatar: string;
  /** 用户名 */
  username: string;
  /** 昵称 */
  nickname: string;
  /** 当前登录用户的角色 */
  roles: Array<string>;
  /** 按钮级别权限 */
  permissions: Array<string>;
};

export type CaptchaResult = {
  captchaId: string;
  image: string;
};

export type RefreshTokenResult = {
  /** token */
  accessToken: string;
  /** 用于调用刷新`accessToken`的接口时所需的`token` */
  refreshToken: string;
  /** `accessToken`的过期时间（秒级时间戳） */
  expires: number;
};

export type UserInfo = {
  /** 头像 */
  avatar: string;
  /** 用户名 */
  username: string;
  /** 昵称 */
  nickname: string;
  /** 邮箱 */
  email: string;
  /** 联系电话 */
  phone: string;
  /** 简介 */
  description: string;
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
export const getLogin = (data?: object) => {
  return http.request<UserResult>("post", "sys/Login", { data });
};

/** 获取图形验证码 */
export const getCaptcha = () => {
  return http.request<CaptchaResult>("post", "sys/captcha");
};

/** 刷新`token` */
export const refreshTokenApi = (data?: object) => {
  return http.request<RefreshTokenResult>("post", "refresh-token", { data });
};

/** 账户设置-个人信息 */
export const getMine = (data?: object) => {
  return http.request<UserInfo>("get", "mine", { data });
};

/** 账户设置-个人安全日志 */
export const getMineLogs = (data?: object) => {
  return http.request<ResultTable>("get", "mine-logs", { data });
};
