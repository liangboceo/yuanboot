import { defineStore } from "pinia";
import {
  type userType,
  store,
  router,
  resetRouter,
  routerArrays,
  storageLocal
} from "../utils";
import {
  type UserResult,
  type RefreshTokenResult,
  getLogin,
  refreshTokenApi
} from "@/api/user";
import { getUserInfoApi } from "@/api/system";
import { useMultiTagsStoreHook } from "./multiTags";
import { type DataInfo, setToken, removeToken, userKey } from "@/utils/auth";

export const useUserStore = defineStore("pure-user", {
  state: (): userType => ({
    // 头像
    avatar: storageLocal().getItem<DataInfo<number>>(userKey)?.avatar ?? "",
    // 用户名
    username: storageLocal().getItem<DataInfo<number>>(userKey)?.username ?? "",
    // 昵称
    nickname: storageLocal().getItem<DataInfo<number>>(userKey)?.nickname ?? "",
    // 页面级别权限
    roles: storageLocal().getItem<DataInfo<number>>(userKey)?.roles ?? [],
    // 按钮级别权限
    permissions:
      storageLocal().getItem<DataInfo<number>>(userKey)?.permissions ?? [],
    // 前端生成的验证码（按实际需求替换）
    verifyCode: "",
    // 判断登录页面显示哪个组件（0：登录（默认）、1：手机登录、2：二维码登录、3：注册、4：忘记密码）
    currentPage: 0,
    // 是否勾选了登录页的免登录
    isRemembered: false,
    // 登录页的免登录存储几天，默认7天
    loginDay: 7
  }),
  actions: {
    /** 存储头像 */
    SET_AVATAR(avatar: string) {
      this.avatar = avatar;
    },
    /** 存储用户名 */
    SET_USERNAME(username: string) {
      this.username = username;
    },
    /** 存储昵称 */
    SET_NICKNAME(nickname: string) {
      this.nickname = nickname;
    },
    /** 存储角色 */
    SET_ROLES(roles: Array<string>) {
      this.roles = roles;
    },
    /** 存储按钮级别权限 */
    SET_PERMS(permissions: Array<string>) {
      this.permissions = permissions;
    },
    /** 存储前端生成的验证码 */
    SET_VERIFYCODE(verifyCode: string) {
      this.verifyCode = verifyCode;
    },
    /** 存储登录页面显示哪个组件 */
    SET_CURRENTPAGE(value: number) {
      this.currentPage = value;
    },
    /** 存储是否勾选了登录页的免登录 */
    SET_ISREMEMBERED(bool: boolean) {
      this.isRemembered = bool;
    },
    /** 设置登录页的免登录存储几天 */
    SET_LOGINDAY(value: number) {
      this.loginDay = Number(value);
    },
    /** 登入 */
    async loginByUsername(data?: object) {
      return new Promise<UserResult>((resolve, reject) => {
        getLogin(data)
          .then(userData => {
            const tokenData = {
              accessToken: userData.accessToken,
              expires: userData.expires,
              refreshToken: userData.refreshToken,
              avatar: userData.avatar,
              username: userData.username,
              nickname: userData.nickname,
              roles: userData.roles,
              permissions: userData.permissions
            };
            setToken(tokenData as unknown as DataInfo<number>);
            resolve(userData);
          })
          .catch(error => {
            reject(error);
          });
      });
    },
    /** 获取当前登录用户信息 */
    async getUserInfo() {
      const data = await getUserInfoApi();
      this.SET_AVATAR(data?.avatar ?? "");
      this.SET_USERNAME(data?.username ?? "");
      this.SET_NICKNAME(data?.nickname ?? "");
      this.SET_ROLES(data?.roles ?? []);
      this.SET_PERMS(data?.permissions ?? []);
      const userInfo = storageLocal().getItem<DataInfo<number>>(userKey) ?? {};
      storageLocal().setItem(userKey, {
        ...userInfo,
        avatar: data?.avatar ?? "",
        username: data?.username ?? "",
        nickname: data?.nickname ?? "",
        roles: data?.roles ?? [],
        permissions: data?.permissions ?? []
      });
      return data;
    },
    /** 前端登出（不调用接口） */
    logOut() {
      this.username = "";
      this.roles = [];
      this.permissions = [];
      removeToken();
      useMultiTagsStoreHook().handleTags("equal", [...routerArrays]);
      resetRouter();
      const currentPath = router.currentRoute.value.fullPath;
      router.push(
        currentPath && currentPath !== "/"
          ? `/login?return_url=${encodeURIComponent(currentPath)}`
          : "/login"
      );
    },
    /** 刷新`token` */
    async handRefreshToken(data?: object) {
      return new Promise<RefreshTokenResult>((resolve, reject) => {
        refreshTokenApi(data)
          .then(response => {
            setToken(response as unknown as DataInfo<number>);
            resolve(response);
          })
          .catch(error => {
            reject(error);
          });
      });
    }
  }
});

export function useUserStoreHook() {
  return useUserStore(store);
}
