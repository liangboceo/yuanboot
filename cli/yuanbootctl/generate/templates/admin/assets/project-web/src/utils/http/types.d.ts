import type {
  Method,
  AxiosError,
  AxiosResponse,
  AxiosRequestConfig
} from "axios";

export type resultType = {
  accessToken?: string;
};

export type RequestMethods = Extract<
  Method,
  "get" | "post" | "put" | "delete" | "patch" | "option" | "head"
>;

export interface PureHttpError extends AxiosError {
  isCancelRequest?: boolean;
}

export interface PureHttpResponse extends AxiosResponse {
  config: PureHttpRequestConfig;
}

export interface PureHttpRequestConfig extends AxiosRequestConfig {
  beforeRequestCallback?: (request: PureHttpRequestConfig) => void;
  beforeResponseCallback?: (response: PureHttpResponse) => void;
}

export default class PureHttp {
  request<T>(
    method: RequestMethods,
    url: string,
    param?: AxiosRequestConfig,
    axiosConfig?: PureHttpRequestConfig
  ): Promise<T>;
  post<T, P>(
    url: string,
    params?: P,
    config?: PureHttpRequestConfig
  ): Promise<T>;
  get<T, P>(
    url: string,
    params?: P,
    config?: PureHttpRequestConfig
  ): Promise<T>;
  download(
    url: string,
    data?: object,
    config?: PureHttpRequestConfig
  ): Promise<Blob>;
  upload<T>(
    url: string,
    data: FormData,
    config?: PureHttpRequestConfig
  ): Promise<T>;
  downloadWithFilename(
    url: string,
    data?: object,
    config?: PureHttpRequestConfig
  ): Promise<{ blob: Blob; filename: string }>;
  stream(
    url: string,
    data: object,
    handlers: {
      onMessage?: (event: string, payload: any) => void;
      onError?: (message: string) => void;
      onDone?: () => void;
    },
    config?: PureHttpRequestConfig
  ): () => void;
}
