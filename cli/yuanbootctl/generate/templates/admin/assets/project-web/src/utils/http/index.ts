import Axios, {
  type AxiosInstance,
  type AxiosRequestConfig,
  type CustomParamsSerializer
} from "axios";
import type {
  PureHttpError,
  RequestMethods,
  PureHttpResponse,
  PureHttpRequestConfig
} from "./types.d";
import { stringify } from "qs";
import { message } from "@/utils/message";
import { $t, transformI18n } from "@/plugins/i18n";
import { getToken, formatToken, removeToken } from "@/utils/auth";
import { useUserStoreHook } from "@/store/modules/user";
import router from "@/router";

// 相关配置请参考：www.axios-js.com/zh-cn/docs/#axios-request-config-1
const defaultConfig: AxiosRequestConfig = {
  // 后端接口请求前缀
  baseURL: import.meta.env.VITE_API_PREFIX,
  // 取消请求超时时间限制，避免大文件下载耗时较长时被前端中断
  timeout: 0,
  headers: {
    Accept: "application/json, text/plain, */*",
    "Content-Type": "application/json",
    "X-Requested-With": "XMLHttpRequest"
  },
  // 数组格式参数序列化（https://github.com/axios/axios/issues/5142）
  paramsSerializer: {
    serialize: stringify as unknown as CustomParamsSerializer
  }
};

class PureHttp {
  constructor() {
    this.httpInterceptorsRequest();
    this.httpInterceptorsResponse();
  }

  /** `token`过期后，暂存待执行的请求 */
  private static requests: Array<{
    resolve: (config: PureHttpRequestConfig) => void;
    reject: (error: any) => void;
    config: PureHttpRequestConfig;
  }> = [];

  /** 防止重复刷新`token` */
  private static isRefreshing = false;

  /** 初始化配置对象 */
  private static initConfig: PureHttpRequestConfig = {};

  /** 保存当前`Axios`实例对象 */
  private static axiosInstance: AxiosInstance = Axios.create(defaultConfig);

  /** 重连原始请求 */
  private static retryOriginalRequest(config: PureHttpRequestConfig) {
    return new Promise((resolve, reject) => {
      PureHttp.requests.push({ resolve, reject, config });
    });
  }

  /** 请求拦截 */
  private httpInterceptorsRequest(): void {
    PureHttp.axiosInstance.interceptors.request.use(
      async (config: PureHttpRequestConfig): Promise<any> => {
        // 优先判断post/get等方法是否传入回调，否则执行初始化设置等回调
        if (typeof config.beforeRequestCallback === "function") {
          config.beforeRequestCallback(config);
          return config;
        }
        if (PureHttp.initConfig.beforeRequestCallback) {
          PureHttp.initConfig.beforeRequestCallback(config);
          return config;
        }
        /** 请求白名单，放置一些不需要`token`的接口（通过设置请求白名单，防止`token`过期后再请求造成的死循环问题） */
        const whiteList = [
          "/refresh-token",
          "/login",
          "refresh-token",
          "sys/Login",
          "sys/captcha",
          "sys/GetSystemConfig",
          "user/Login",
          "user/Register"
        ];
        return whiteList.some(url => config.url.endsWith(url))
          ? config
          : new Promise((resolve, reject) => {
              const data = getToken();
              if (data) {
                const now = new Date().getTime();
                const expired = parseInt(data.expires) - now <= 0;
                if (expired) {
                  if (!PureHttp.isRefreshing) {
                    PureHttp.isRefreshing = true;
                    // token过期刷新
                    useUserStoreHook()
                      .handRefreshToken({ refreshToken: data.refreshToken })
                      .then(res => {
                        const token = res.accessToken;
                        PureHttp.requests.forEach(
                          ({ resolve, config: reqConfig }) => {
                            reqConfig.headers["Authorization"] =
                              formatToken(token);
                            resolve(reqConfig);
                          }
                        );
                        PureHttp.requests = [];
                      })
                      .catch(_err => {
                        reject(_err);
                        PureHttp.requests.forEach(({ reject }) => {
                          reject(_err);
                        });
                        PureHttp.requests = [];
                        useUserStoreHook().logOut();
                        message(transformI18n($t("login.pureLoginExpired")), {
                          type: "warning"
                        });
                      })
                      .finally(() => {
                        PureHttp.isRefreshing = false;
                      });
                  }
                  resolve(PureHttp.retryOriginalRequest(config));
                } else {
                  config.headers["Authorization"] = formatToken(
                    data.accessToken
                  );
                  resolve(config);
                }
              } else {
                resolve(config);
              }
            });
      },
      error => {
        return Promise.reject(error);
      }
    );
  }

  /** 响应拦截 */
  private httpInterceptorsResponse(): void {
    const instance = PureHttp.axiosInstance;
    instance.interceptors.response.use(
      async (response: PureHttpResponse) => {
        const $config = response.config;
        const { data, status } = response;

        // 文件下载：blob 响应直接返回，跳过 JSON 业务码校验
        if ($config.responseType === "blob") {
          // 支持 beforeResponseCallback 捕获 headers 等
          if (typeof $config.beforeResponseCallback === "function") {
            $config.beforeResponseCallback(response);
          }
          if (status !== 200) {
            message(`下载失败(${status})`, { type: "error" });
            return Promise.reject(data);
          }
          // 服务端可能返回 JSON 错误（Content-Type: application/json），需要解析提示
          if (
            data instanceof Blob &&
            (data.type.includes("application/json") ||
              data.type.includes("text/plain"))
          ) {
            const text = await data.text();
            try {
              const err = JSON.parse(text);
              message(err?.msg || err?.message || `下载失败(${status})`, {
                type: "error"
              });
              return Promise.reject(err);
            } catch {
              // 非 JSON 文本，正常下载
            }
          }
          return data; // 直接返回 Blob
        }

        // 非 200 状态码提示
        if (status !== 200) {
          message(data?.message || `请求失败(${status})`, { type: "error" });
          return Promise.reject(data);
        }

        // 业务层错误码处理（后端返回 code 字段）
        if (data?.code !== undefined && data.code !== 0 && data.code !== 200) {
          // 401 未授权，清除 token 并跳转登录页
          if (data.code === 401) {
            if (($config as any).evaluationUserRequest) {
              return Promise.reject(data);
            }
            removeToken();
            useUserStoreHook().logOut();
            message(transformI18n($t("login.pureLoginExpired")), {
              type: "warning"
            });
            const currentPath = router.currentRoute.value.fullPath;
            router.push(
              currentPath && currentPath !== "/"
                ? `/login?return_url=${encodeURIComponent(currentPath)}`
                : "/login"
            );
            return Promise.reject(data);
          }
          // 其他业务错误提示
          message(data?.message || `操作失败(${data.code})`, {
            type: "error"
          });
          return Promise.reject(data);
        }

        // 优先判断post/get等方法是否传入回调，否则执行初始化设置等回调
        if (typeof $config.beforeResponseCallback === "function") {
          $config.beforeResponseCallback(response);
          return data.data;
        }
        if (PureHttp.initConfig.beforeResponseCallback) {
          PureHttp.initConfig.beforeResponseCallback(response);
          return data.data;
        }
        return data.data;
      },
      (error: PureHttpError) => {
        const $error = error;
        $error.isCancelRequest = Axios.isCancel($error);

        // HTTP 状态码级别的错误处理
        const status = $error?.response?.status;
        if (status === 401) {
          if (($error?.config as any)?.evaluationUserRequest) {
            return Promise.reject($error);
          }
          removeToken();
          useUserStoreHook().logOut();
          message(transformI18n($t("login.pureLoginExpired")), {
            type: "warning"
          });
          const currentPath = router.currentRoute.value.fullPath;
          router.push(
            currentPath && currentPath !== "/"
              ? `/login?return_url=${encodeURIComponent(currentPath)}`
              : "/login"
          );
        } else if (status) {
          const errData = $error?.response?.data as
            | { message?: string }
            | undefined;
          const msg = errData?.message || `请求失败(${status})`;
          message(msg, { type: "error" });
        } else {
          const isLoggedOut =
            !getToken() && router.currentRoute.value.path === "/login";
          if (!isLoggedOut) {
            message($error?.message || "网络异常，请稍后重试", {
              type: "error"
            });
          }
        }

        return Promise.reject($error);
      }
    );
  }

  /** 通用请求工具函数 */
  public request<T>(
    method: RequestMethods,
    url: string,
    param?: AxiosRequestConfig,
    axiosConfig?: PureHttpRequestConfig
  ): Promise<T> {
    const config = {
      method,
      url,
      ...param,
      ...axiosConfig
    } as PureHttpRequestConfig;

    // 单独处理自定义请求/响应回调
    return new Promise((resolve, reject) => {
      PureHttp.axiosInstance
        .request(config)
        .then((response: undefined) => {
          resolve(response);
        })
        .catch(error => {
          reject(error);
        });
    });
  }

  /** 单独抽离的`post`工具函数 */
  public post<T, P>(
    url: string,
    params?: AxiosRequestConfig<P>,
    config?: PureHttpRequestConfig
  ): Promise<T> {
    return this.request<T>("post", url, params, config);
  }

  /** 单独抽离的`get`工具函数 */
  public get<T, P>(
    url: string,
    params?: AxiosRequestConfig<P>,
    config?: PureHttpRequestConfig
  ): Promise<T> {
    return this.request<T>("get", url, params, config);
  }

  /** 文件下载：以 blob 形式请求，返回 Blob */
  public download(
    url: string,
    data?: object,
    config?: PureHttpRequestConfig
  ): Promise<Blob> {
    return this.request<Blob>(
      "post",
      url,
      {
        data,
        responseType: "blob"
      } as any,
      config
    );
  }

  /** 文件上传：POST FormData，显式移除 Content-Type 让浏览器自动设置 multipart/form-data + boundary */
  public upload<T>(
    url: string,
    data: FormData,
    config?: PureHttpRequestConfig
  ): Promise<T> {
    return this.request<T>(
      "post",
      url,
      {
        data,
        headers: { "Content-Type": undefined as unknown as string }
      },
      config
    );
  }

  /** 文件下载并提取文件名：返回 blob + 解析自 Content-Disposition 的 filename */
  public downloadWithFilename(
    url: string,
    data?: object,
    config?: PureHttpRequestConfig
  ): Promise<{ blob: Blob; filename: string }> {
    let respHeaders: Record<string, string> = {};
    return this.request<Blob>(
      "post",
      url,
      { data, responseType: "blob" } as any,
      {
        ...config,
        beforeResponseCallback: response => {
          respHeaders = { ...response.headers } as any;
          config?.beforeResponseCallback?.(response);
        }
      }
    ).then(blob => {
      const disposition =
        respHeaders["content-disposition"] ||
        respHeaders["Content-Disposition"] ||
        "";
      const encodedMatch = disposition.match(/filename\*=UTF-8''([^;]+)/i);
      const plainMatch = disposition.match(/filename="?([^";]+)"?/i);
      let filename = encodedMatch?.[1] || plainMatch?.[1] || "download";
      try {
        filename = decodeURIComponent(filename);
      } catch {
        // 保留服务端返回的原始文件名
      }
      return { blob, filename };
    });
  }

  /** SSE 流式请求：封装 fetch + SSE 解析，自动处理 baseURL / token / Content-Type
   *  支持心跳维持连接、超时检测 */
  public stream(
    url: string,
    data: object,
    handlers: {
      onMessage?: (event: string, payload: any) => void;
      onError?: (message: string) => void;
      onDone?: () => void;
    },
    _config?: PureHttpRequestConfig
  ): () => void {
    const controller = new AbortController();
    const base = defaultConfig.baseURL || "";
    const fullUrl = `${base}${url}`;
    const token = getToken()?.accessToken;

    // 连接超时检测：如果超过 60s 没有收到任何数据，认为连接已断开
    const READ_TIMEOUT_MS = 60000;
    let readTimeoutId: ReturnType<typeof setTimeout> | null = null;
    let isAborted = false;

    const resetReadTimeout = () => {
      if (readTimeoutId) clearTimeout(readTimeoutId);
      readTimeoutId = setTimeout(() => {
        if (!isAborted) {
          isAborted = true;
          controller.abort();
          handlers.onError?.("连接超时，请重试");
        }
      }, READ_TIMEOUT_MS);
    };

    (async () => {
      try {
        resetReadTimeout();

        const response = await fetch(fullUrl, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            ...(token ? { Authorization: formatToken(token) } : {})
          },
          body: JSON.stringify(data),
          signal: controller.signal
        });

        if (!response.ok || !response.body) {
          handlers.onError?.("请求失败");
          isAborted = true;
          return;
        }

        const reader = response.body.getReader();
        const decoder = new TextDecoder("utf-8");
        let buffer = "";
        let hasReceivedDone = false;

        const nextFrame = () =>
          new Promise<void>(resolve => setTimeout(resolve, 0));

        const dispatchEvent = async (raw: string) => {
          const lines = raw.split("\n");
          const event = lines
            .find(line => line.startsWith("event:"))
            ?.replace("event:", "")
            .trim();
          const dataLine = lines
            .filter(line => line.startsWith("data:"))
            .map(line => line.replace("data:", "").trim())
            .join("\n");
          if (!event || !dataLine) return;

          let payload: any;
          try {
            payload = JSON.parse(dataLine);
          } catch {
            return; // 忽略解析失败的数据
          }

          const normalizedEvent = event.toLowerCase();

          // 心跳事件：只重置超时，不传递给业务层
          if (normalizedEvent === "heartbeat") {
            resetReadTimeout();
            return;
          }

          if (normalizedEvent === "error") {
            hasReceivedDone = true;
            handlers.onError?.(payload.message || payload.msg || "请求失败");
            return;
          }
          if (normalizedEvent === "done") {
            hasReceivedDone = true;
            handlers.onDone?.();
            return;
          }
          handlers.onMessage?.(normalizedEvent, payload);

          // 仅对高频 message chunk 等待帧渲染，status/tool_call/file_output 立即分发
          if (normalizedEvent === "message") {
            await nextFrame();
          }
        };

        while (true) {
          const { value, done } = await reader.read();
          // Reset timeout immediately on data arrival, before potentially slow chunk processing
          resetReadTimeout();
          if (done) break;
          buffer += decoder.decode(value, { stream: true });
          const chunks = buffer.split("\n\n");
          buffer = chunks.pop() || "";
          for (const chunk of chunks) {
            await dispatchEvent(chunk);
          }
        }
        if (buffer.trim()) await dispatchEvent(buffer);
        if (!hasReceivedDone) {
          handlers.onDone?.();
        }
      } catch (err: any) {
        if (err?.name !== "AbortError") {
          handlers.onError?.(err?.message || "网络异常");
        }
      } finally {
        isAborted = true;
        if (readTimeoutId) clearTimeout(readTimeoutId);
      }
    })();

    return () => {
      isAborted = true;
      if (readTimeoutId) clearTimeout(readTimeoutId);
      controller.abort();
    };
  }
}

export const http = new PureHttp();
