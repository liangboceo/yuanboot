import "@/utils/sso";
import Cookies from "js-cookie";
import { getConfig } from "@/config";
import NProgress from "@/utils/progress";
import { transformI18n } from "@/plugins/i18n";
import remainingRouter from "./modules/remaining";
import { useSystemConfigStoreHook } from "@/store/modules/systemConfig";
import { useUserStoreHook } from "@/store/modules/user";
import { useMultiTagsStoreHook } from "@/store/modules/multiTags";
import { usePermissionStoreHook } from "@/store/modules/permission";
import {
  isUrl,
  openLink,
  cloneDeep,
  isAllEmpty,
  storageLocal
} from "@pureadmin/utils";
import {
  getTopMenu,
  initRouter,
  isOneOfArray,
  getHistoryMode,
  findRouteByPath,
  handleAliveRoute
} from "./utils";
import { type Router, type RouteRecordRaw, createRouter } from "vue-router";
import {
  type DataInfo,
  userKey,
  removeToken,
  multipleTabsKey
} from "@/utils/auth";

const Layout = () => import("@/layout/index.vue");

/** 只保留动态路由挂载所需的根路由，其余菜单路由由后端动态加载 */
export const constantRoutes: Array<RouteRecordRaw> = [
  {
    path: "/",
    name: "Home",
    component: Layout,
    meta: {
      title: "menus.pureHome",
      rank: 0
    },
    redirect: "/home",
    children: [
      {
        path: "/home",
        name: "HomeDashboard",
        component: () => import("@/views/home/index.vue"),
        meta: {
          title: "首页",
          icon: "ri:dashboard-line",
          rank: 0
        }
      }
    ]
  }
];

/** 初始的静态路由，用于退出登录时重置路由 */
const initConstantRoutes: Array<RouteRecordRaw> = cloneDeep(constantRoutes);

/** 用于渲染菜单，静态菜单已移除，仅展示动态路由 */
export const constantMenus: Array<RouteRecordRaw> = [];

/** 不参与菜单的路由 */
export const remainingPaths = remainingRouter.map(route => route.path);

/** 创建路由实例 */
export const router: Router = createRouter({
  history: getHistoryMode(import.meta.env.VITE_ROUTER_HISTORY),
  routes: constantRoutes.concat(...(remainingRouter as any)),
  strict: true,
  scrollBehavior(to, from, savedPosition) {
    return new Promise(resolve => {
      if (savedPosition) {
        return savedPosition;
      } else {
        if (from.meta.saveSrollTop) {
          const top: number =
            document.documentElement.scrollTop || document.body.scrollTop;
          resolve({ left: 0, top });
        }
      }
    });
  }
});

/** 记录已经加载的页面路径 */
const loadedPaths = new Set<string>();

/** 重置已加载页面记录 */
export function resetLoadedPaths() {
  loadedPaths.clear();
}

/** 重置路由 */
export function resetRouter() {
  router.clearRoutes();
  for (const route of initConstantRoutes.concat(...(remainingRouter as any))) {
    router.addRoute(route);
  }
  router.options.routes = cloneDeep(initConstantRoutes);
  usePermissionStoreHook().clearAllCachePage();
  resetLoadedPaths();
}

/** 路由白名单 */
const whiteList = ["/login"];

router.beforeEach((to: ToRouteType, _from) => {
  to.meta.loaded = loadedPaths.has(to.path);

  if (!to.meta.loaded) {
    NProgress.start();
  }

  if (to.meta?.keepAlive) {
    handleAliveRoute(to, "add");
    // 页面整体刷新和点击标签页刷新
    if (_from.name === undefined || _from.name === "Redirect") {
      handleAliveRoute(to);
    }
  }
  if (!to.path.includes("/redirect")) {
    useMultiTagsStoreHook().handleTags("push", {
      path: to.path,
      name: to.name,
      meta: to.meta,
      query: isAllEmpty(to.query) ? undefined : to.query,
      params: isAllEmpty(to.params) ? undefined : to.params
    });
  }
  const userInfo = storageLocal().getItem<DataInfo<number>>(userKey);
  const externalLink = isUrl(to?.name as string);
  if (!externalLink) {
    to.matched.some(item => {
      if (!item.meta.title) return "";
      const Title = useSystemConfigStoreHook().getName || getConfig().Title;
      if (Title)
        document.title = `${transformI18n(item.meta.title)} | ${Title}`;
      else document.title = transformI18n(item.meta.title);
    });
  }
  /** 如果已经登录并存在登录信息后不能跳转到路由白名单，而是继续保持在当前页面 */
  function toCorrectRoute() {
    if (whiteList.includes(to.path)) {
      // 已登录访问登录页时，重定向到 return_url 或默认首页
      const returnUrl = Array.isArray(to.query.return_url)
        ? to.query.return_url[0]
        : to.query.return_url;
      return typeof returnUrl === "string" &&
        returnUrl &&
        returnUrl !== "/login"
        ? returnUrl
        : getTopMenu()?.path || "/";
    }
    return undefined;
  }
  if (Cookies.get(multipleTabsKey) && userInfo) {
    // 无权限跳转403页面
    if (to.meta?.roles && !isOneOfArray(to.meta?.roles, userInfo?.roles)) {
      return { path: "/error/403" };
    }
    if (_from?.name) {
      // name为超链接
      if (externalLink) {
        openLink(to?.name as string);
        NProgress.done();
        return false;
      } else {
        return toCorrectRoute();
      }
    } else {
      // 刷新
      if (
        usePermissionStoreHook().wholeMenus.length === 0 &&
        to.path !== "/login"
      ) {
        useUserStoreHook()
          .getUserInfo()
          .then(() => initRouter())
          .then((router: Router) => {
            if (!useMultiTagsStoreHook().getMultiTagsCache) {
              const { path } = to;
              const rootRoute = router.options.routes.find(
                route => route.path === "/"
              );
              const route = findRouteByPath(path, rootRoute?.children ?? []);
              getTopMenu(true);
              // query、params模式路由传参数的标签页不在此处处理
              if (route && route.meta?.title) {
                if (
                  isAllEmpty(route.parentId) &&
                  route.meta?.backstage &&
                  route.children?.length
                ) {
                  // 此处为动态顶级路由（目录）
                  const { path, name, meta } = route.children[0];
                  useMultiTagsStoreHook().handleTags("push", {
                    path,
                    name,
                    meta
                  });
                } else {
                  const { path, name, meta } = route;
                  useMultiTagsStoreHook().handleTags("push", {
                    path,
                    name,
                    meta
                  });
                }
              }
            }
            // 确保动态路由完全加入路由列表并且不影响静态路由（注意：动态路由刷新时router.beforeEach可能会触发两次，第一次触发动态路由还未完全添加，第二次动态路由才完全添加到路由列表，如果需要在router.beforeEach做一些判断可以在to.name存在的条件下去判断，这样就只会触发一次）
            if (isAllEmpty(to.name)) router.push(to.fullPath);
          });
      }
      return toCorrectRoute();
    }
  } else {
    if (to.path !== "/login") {
      if (whiteList.indexOf(to.path) !== -1) {
        return true;
      } else {
        removeToken();
        return { path: "/login", query: { return_url: to.fullPath } };
      }
    } else {
      return true;
    }
  }
});

router.afterEach(to => {
  loadedPaths.add(to.path);
  NProgress.done();
});

export default router;
