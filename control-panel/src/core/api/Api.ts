/* eslint-disable */
/* tslint:disable */
// @ts-nocheck
/*
 * ---------------------------------------------------------------
 * ## THIS FILE WAS GENERATED VIA SWAGGER-TYPESCRIPT-API        ##
 * ##                                                           ##
 * ## AUTHOR: acacode                                           ##
 * ## SOURCE: https://github.com/acacode/swagger-typescript-api ##
 * ---------------------------------------------------------------
 */

export interface Resource {
  /** @format int64 */
  id: number;
  name: string;
  url_pattern: string;
  /** @format int64 */
  security_profile_id?: number | null;
  /** @format int64 */
  traffic_profile_id?: number | null;
  /** @format date-time */
  created_at: string;
  /** @format date-time */
  updated_at: string;
}

export interface ResourceCreate {
  name: string;
  url_pattern: string;
  /** @format int64 */
  security_profile_id?: number | null;
  /** @format int64 */
  traffic_profile_id?: number | null;
}

export interface ResourceUpdate {
  name: string;
  url_pattern: string;
  /** @format int64 */
  security_profile_id?: number | null;
  /** @format int64 */
  traffic_profile_id?: number | null;
}

export interface SecurityProfile {
  /** @format int64 */
  id: number;
  name: string;
  description?: string;
  base_action: "allow" | "block";
  log_enabled: boolean;
  is_enabled: boolean;
  /** @format date-time */
  created_at: string;
  /** @format date-time */
  updated_at: string;
}

export interface SecurityProfileCreate {
  name: string;
  description?: string;
  base_action: "allow" | "block";
  log_enabled: boolean;
  is_enabled: boolean;
}

export interface SecurityProfileUpdate {
  name: string;
  description?: string;
  base_action: "allow" | "block";
  log_enabled: boolean;
  is_enabled: boolean;
}

export interface SecurityRule {
  /** @format int64 */
  id: number;
  /** @format int64 */
  profile_id: number;
  name: string;
  description?: string;
  /** @format int32 */
  priority: number;
  rule_type: "deterministic" | "ml";
  action: "allow" | "block";
  /**
   * Условия для security rules. Блоки разных типов объединяются по AND.
   * Внутри одного блока значения объединяются по OR.
   */
  conditions?: SecurityRuleConditions | null;
  /** @format int64 */
  ml_model_id?: number | null;
  /**
   * @format int32
   * @min 0
   * @max 100
   */
  ml_threshold?: number | null;
  dry_run: boolean;
  is_enabled: boolean;
  /** @format date-time */
  created_at: string;
  /** @format date-time */
  updated_at: string;
}

export interface SecurityRuleCreate {
  /** @format int64 */
  profile_id: number;
  name: string;
  description?: string;
  /** @format int32 */
  priority: number;
  rule_type: "deterministic" | "ml";
  action: "allow" | "block";
  /**
   * Условия для security rules. Блоки разных типов объединяются по AND.
   * Внутри одного блока значения объединяются по OR.
   */
  conditions?: SecurityRuleConditions | null;
  /** @format int64 */
  ml_model_id?: number | null;
  /**
   * @format int32
   * @min 0
   * @max 100
   */
  ml_threshold?: number | null;
  dry_run: boolean;
  is_enabled: boolean;
}

export interface SecurityRuleUpdate {
  /** @format int64 */
  profile_id: number;
  name: string;
  description?: string;
  /** @format int32 */
  priority: number;
  rule_type: "deterministic" | "ml";
  action: "allow" | "block";
  /**
   * Условия для security rules. Блоки разных типов объединяются по AND.
   * Внутри одного блока значения объединяются по OR.
   */
  conditions?: SecurityRuleConditions | null;
  /** @format int64 */
  ml_model_id?: number | null;
  /**
   * @format int32
   * @min 0
   * @max 100
   */
  ml_threshold?: number | null;
  dry_run: boolean;
  is_enabled: boolean;
}

export interface TrafficProfile {
  /** @format int64 */
  id: number;
  name: string;
  description?: string;
  is_enabled: boolean;
  /** @format date-time */
  created_at: string;
  /** @format date-time */
  updated_at: string;
}

export interface TrafficProfileCreate {
  name: string;
  description?: string;
  is_enabled: boolean;
}

export interface TrafficProfileUpdate {
  name: string;
  description?: string;
  is_enabled: boolean;
}

export interface TrafficRule {
  /** @format int64 */
  id: number;
  /** @format int64 */
  profile_id: number;
  name: string;
  description?: string;
  /** @format int32 */
  priority: number;
  dry_run: boolean;
  match_all: boolean;
  /** @format int32 */
  requests_limit: number;
  /** @format int32 */
  period_seconds: number;
  /**
   * Условия для security rules. Блоки разных типов объединяются по AND.
   * Внутри одного блока значения объединяются по OR.
   */
  conditions?: SecurityRuleConditions | null;
  is_enabled: boolean;
  /** @format date-time */
  created_at: string;
  /** @format date-time */
  updated_at: string;
}

export interface TrafficRuleCreate {
  /** @format int64 */
  profile_id: number;
  name: string;
  description?: string;
  /** @format int32 */
  priority: number;
  dry_run: boolean;
  match_all: boolean;
  /** @format int32 */
  requests_limit: number;
  /** @format int32 */
  period_seconds: number;
  /**
   * Условия для security rules. Блоки разных типов объединяются по AND.
   * Внутри одного блока значения объединяются по OR.
   */
  conditions?: SecurityRuleConditions | null;
  is_enabled: boolean;
}

export interface TrafficRuleUpdate {
  /** @format int64 */
  profile_id: number;
  name: string;
  description?: string;
  /** @format int32 */
  priority: number;
  dry_run: boolean;
  match_all: boolean;
  /** @format int32 */
  requests_limit: number;
  /** @format int32 */
  period_seconds: number;
  /**
   * Условия для security rules. Блоки разных типов объединяются по AND.
   * Внутри одного блока значения объединяются по OR.
   */
  conditions?: SecurityRuleConditions | null;
  is_enabled: boolean;
}

export interface SecurityHeaderCondition {
  name: string;
  value_regex: string[];
}

/**
 * Условия для security rules. Блоки разных типов объединяются по AND.
 * Внутри одного блока значения объединяются по OR.
 */
export interface SecurityRuleConditions {
  /** CIDR диапазоны клиентских IP */
  source_ip_cidr?: string[];
  /** Регулярные выражения для HTTP URI */
  uri_regex?: string[];
  /** Регулярные выражения для host */
  host_regex?: string[];
  /** Регулярные выражения для HTTP метода */
  method_regex?: string[];
  /** Условия по имени заголовка и regexp значения */
  headers?: SecurityHeaderCondition[];
}

/**
 * Условия для правил безопасности и лимитов трафика. Все блоки объединяются по AND.
 * Внутри каждого блока значения объединяются по OR.
 */
export interface RuleConditions {
  /** HTTP методы (GET, POST и т.п.) */
  methods?: string[];
  /** Регулярные выражения для пути */
  path_regex?: string[];
  /** Префиксы пути */
  path_prefix?: string[];
  /** Регулярные выражения для query-string */
  query_regex?: string[];
  /** Точные значения host */
  hosts?: string[];
  /** Регулярные выражения для host */
  host_regex?: string[];
  /** Карта header -> список допустимых значений */
  headers?: Record<string, string[]>;
  /** Карта header -> список regex для значения */
  headers_regex?: Record<string, string[]>;
  /** CIDR-диапазоны клиентских IP */
  ip_cidr?: string[];
  /** Регулярные выражения для User-Agent */
  user_agent_regex?: string[];
}

export interface MLModel {
  /** @format int64 */
  id: number;
  name: string;
  /** @format byte */
  model_data: string;
}

export interface MLModelCreate {
  name: string;
  /** @format byte */
  model_data: string;
}

export interface MLModelUpdate {
  name: string;
  /** @format byte */
  model_data: string;
}

export interface RequestLog {
  /** @format int64 */
  id: number;
  /** @format int64 */
  resource_id?: number | null;
  /** @format date-time */
  occurred_at: string;
  client_ip: string;
  method: string;
  path: string;
  /** @format int32 */
  status_code: number;
  action: string;
  /** @format int64 */
  rule_id?: number;
  /** @format int64 */
  profile_id?: number;
  user_agent?: string;
  country?: string;
  /** @format int32 */
  latency_ms?: number;
  request_id?: string;
  metadata?: Record<string, any>;
  host?: string;
  scheme?: string;
  protocol?: string;
  authority?: string;
  query?: string;
  /** @format int32 */
  source_port?: number;
  destination_ip?: string;
  /** @format int32 */
  destination_port?: number;
  source_principal?: string;
  source_service?: string;
  source_labels?: Record<string, any>;
  destination_service?: string;
  destination_labels?: Record<string, any>;
  request_http_id?: string;
  fragment?: string;
  request_headers?: Record<string, any>;
  /** @format int32 */
  request_body_size?: number;
  request_body?: string;
  context_extensions?: Record<string, any>;
  metadata_context?: Record<string, any>;
  route_metadata_context?: Record<string, any>;
}

export interface EventLog {
  /** @format int64 */
  id: number;
  /** @format int64 */
  resource_id: number;
  /** @format date-time */
  occurred_at: string;
  event_type: string;
  severity: string;
  message: string;
  /** @format int64 */
  rule_id?: number;
  /** @format int64 */
  profile_id?: number;
  metadata?: Record<string, any>;
}

import type {
  AxiosInstance,
  AxiosRequestConfig,
  AxiosResponse,
  HeadersDefaults,
  ResponseType,
} from "axios";
import axios from "axios";

export type QueryParamsType = Record<string | number, any>;

export interface FullRequestParams
  extends Omit<AxiosRequestConfig, "data" | "params" | "url" | "responseType"> {
  /** set parameter to `true` for call `securityWorker` for this request */
  secure?: boolean;
  /** request path */
  path: string;
  /** content type of request body */
  type?: ContentType;
  /** query params */
  query?: QueryParamsType;
  /** format of response (i.e. response.json() -> format: "json") */
  format?: ResponseType;
  /** request body */
  body?: unknown;
}

export type RequestParams = Omit<
  FullRequestParams,
  "body" | "method" | "query" | "path"
>;

export interface ApiConfig<SecurityDataType = unknown>
  extends Omit<AxiosRequestConfig, "data" | "cancelToken"> {
  securityWorker?: (
    securityData: SecurityDataType | null,
  ) => Promise<AxiosRequestConfig | void> | AxiosRequestConfig | void;
  secure?: boolean;
  format?: ResponseType;
}

export enum ContentType {
  Json = "application/json",
  JsonApi = "application/vnd.api+json",
  FormData = "multipart/form-data",
  UrlEncoded = "application/x-www-form-urlencoded",
  Text = "text/plain",
}

export class HttpClient<SecurityDataType = unknown> {
  public instance: AxiosInstance;
  private securityData: SecurityDataType | null = null;
  private securityWorker?: ApiConfig<SecurityDataType>["securityWorker"];
  private secure?: boolean;
  private format?: ResponseType;

  constructor({
    securityWorker,
    secure,
    format,
    ...axiosConfig
  }: ApiConfig<SecurityDataType> = {}) {
    this.instance = axios.create({
      ...axiosConfig,
      baseURL: axiosConfig.baseURL || "",
    });
    this.secure = secure;
    this.format = format;
    this.securityWorker = securityWorker;
  }

  public setSecurityData = (data: SecurityDataType | null) => {
    this.securityData = data;
  };

  protected mergeRequestParams(
    params1: AxiosRequestConfig,
    params2?: AxiosRequestConfig,
  ): AxiosRequestConfig {
    const method = params1.method || (params2 && params2.method);

    return {
      ...this.instance.defaults,
      ...params1,
      ...(params2 || {}),
      headers: {
        ...((method &&
          this.instance.defaults.headers[
            method.toLowerCase() as keyof HeadersDefaults
          ]) ||
          {}),
        ...(params1.headers || {}),
        ...((params2 && params2.headers) || {}),
      },
    };
  }

  protected stringifyFormItem(formItem: unknown) {
    if (typeof formItem === "object" && formItem !== null) {
      return JSON.stringify(formItem);
    } else {
      return `${formItem}`;
    }
  }

  protected createFormData(input: Record<string, unknown>): FormData {
    if (input instanceof FormData) {
      return input;
    }
    return Object.keys(input || {}).reduce((formData, key) => {
      const property = input[key];
      const propertyContent: any[] =
        property instanceof Array ? property : [property];

      for (const formItem of propertyContent) {
        const isFileType = formItem instanceof Blob || formItem instanceof File;
        formData.append(
          key,
          isFileType ? formItem : this.stringifyFormItem(formItem),
        );
      }

      return formData;
    }, new FormData());
  }

  public request = async <T = any, _E = any>({
    secure,
    path,
    type,
    query,
    format,
    body,
    ...params
  }: FullRequestParams): Promise<AxiosResponse<T>> => {
    const secureParams =
      ((typeof secure === "boolean" ? secure : this.secure) &&
        this.securityWorker &&
        (await this.securityWorker(this.securityData))) ||
      {};
    const requestParams = this.mergeRequestParams(params, secureParams);
    const responseFormat = format || this.format || undefined;

    if (
      type === ContentType.FormData &&
      body &&
      body !== null &&
      typeof body === "object"
    ) {
      body = this.createFormData(body as Record<string, unknown>);
    }

    if (
      type === ContentType.Text &&
      body &&
      body !== null &&
      typeof body !== "string"
    ) {
      body = JSON.stringify(body);
    }

    return this.instance.request({
      ...requestParams,
      headers: {
        ...(requestParams.headers || {}),
        ...(type ? { "Content-Type": type } : {}),
      },
      params: query,
      responseType: responseFormat,
      data: body,
      url: path,
    });
  };
}

/**
 * @title WAS API
 * @version 1.0.0
 */
export class Api<
  SecurityDataType extends unknown,
> extends HttpClient<SecurityDataType> {
  info = {
    /**
     * No description
     *
     * @name InfoList
     * @summary Service info
     * @request GET:/_info
     */
    infoList: (params: RequestParams = {}) =>
      this.request<void, any>({
        path: `/_info`,
        method: "GET",
        ...params,
      }),
  };
  api = {
    /**
     * No description
     *
     * @tags resources
     * @name ListResources
     * @summary List resources
     * @request GET:/api/v1/resources
     */
    listResources: (params: RequestParams = {}) =>
      this.request<Resource[], any>({
        path: `/api/v1/resources`,
        method: "GET",
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags resources
     * @name CreateResource
     * @summary Create resource
     * @request POST:/api/v1/resources
     */
    createResource: (data: ResourceCreate, params: RequestParams = {}) =>
      this.request<Resource, any>({
        path: `/api/v1/resources`,
        method: "POST",
        body: data,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags resources
     * @name GetResource
     * @summary Get resource by id
     * @request GET:/api/v1/resources/{id}
     */
    getResource: (id: number, params: RequestParams = {}) =>
      this.request<Resource, void>({
        path: `/api/v1/resources/${id}`,
        method: "GET",
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags resources
     * @name UpdateResource
     * @summary Update resource by id
     * @request PUT:/api/v1/resources/{id}
     */
    updateResource: (
      id: number,
      data: ResourceUpdate,
      params: RequestParams = {},
    ) =>
      this.request<Resource, void>({
        path: `/api/v1/resources/${id}`,
        method: "PUT",
        body: data,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags resources
     * @name DeleteResource
     * @summary Delete resource by id
     * @request DELETE:/api/v1/resources/{id}
     */
    deleteResource: (id: number, params: RequestParams = {}) =>
      this.request<void, void>({
        path: `/api/v1/resources/${id}`,
        method: "DELETE",
        ...params,
      }),

    /**
     * No description
     *
     * @tags security-profiles
     * @name ListSecurityProfiles
     * @summary List security profiles
     * @request GET:/api/v1/security-profiles
     */
    listSecurityProfiles: (params: RequestParams = {}) =>
      this.request<SecurityProfile[], any>({
        path: `/api/v1/security-profiles`,
        method: "GET",
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags security-profiles
     * @name CreateSecurityProfile
     * @summary Create security profile
     * @request POST:/api/v1/security-profiles
     */
    createSecurityProfile: (
      data: SecurityProfileCreate,
      params: RequestParams = {},
    ) =>
      this.request<SecurityProfile, any>({
        path: `/api/v1/security-profiles`,
        method: "POST",
        body: data,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags security-profiles
     * @name GetSecurityProfile
     * @summary Get security profile by id
     * @request GET:/api/v1/security-profiles/{id}
     */
    getSecurityProfile: (id: number, params: RequestParams = {}) =>
      this.request<SecurityProfile, void>({
        path: `/api/v1/security-profiles/${id}`,
        method: "GET",
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags security-profiles
     * @name UpdateSecurityProfile
     * @summary Update security profile by id
     * @request PUT:/api/v1/security-profiles/{id}
     */
    updateSecurityProfile: (
      id: number,
      data: SecurityProfileUpdate,
      params: RequestParams = {},
    ) =>
      this.request<SecurityProfile, void>({
        path: `/api/v1/security-profiles/${id}`,
        method: "PUT",
        body: data,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags security-profiles
     * @name DeleteSecurityProfile
     * @summary Delete security profile by id
     * @request DELETE:/api/v1/security-profiles/{id}
     */
    deleteSecurityProfile: (id: number, params: RequestParams = {}) =>
      this.request<void, void>({
        path: `/api/v1/security-profiles/${id}`,
        method: "DELETE",
        ...params,
      }),

    /**
     * No description
     *
     * @tags security-rules
     * @name ListSecurityRules
     * @summary List security rules
     * @request GET:/api/v1/security-rules
     */
    listSecurityRules: (params: RequestParams = {}) =>
      this.request<SecurityRule[], any>({
        path: `/api/v1/security-rules`,
        method: "GET",
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags security-rules
     * @name CreateSecurityRule
     * @summary Create security rule
     * @request POST:/api/v1/security-rules
     */
    createSecurityRule: (
      data: SecurityRuleCreate,
      params: RequestParams = {},
    ) =>
      this.request<SecurityRule, any>({
        path: `/api/v1/security-rules`,
        method: "POST",
        body: data,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags security-rules
     * @name GetSecurityRule
     * @summary Get security rule by id
     * @request GET:/api/v1/security-rules/{id}
     */
    getSecurityRule: (id: number, params: RequestParams = {}) =>
      this.request<SecurityRule, void>({
        path: `/api/v1/security-rules/${id}`,
        method: "GET",
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags security-rules
     * @name UpdateSecurityRule
     * @summary Update security rule by id
     * @request PUT:/api/v1/security-rules/{id}
     */
    updateSecurityRule: (
      id: number,
      data: SecurityRuleUpdate,
      params: RequestParams = {},
    ) =>
      this.request<SecurityRule, void>({
        path: `/api/v1/security-rules/${id}`,
        method: "PUT",
        body: data,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags security-rules
     * @name DeleteSecurityRule
     * @summary Delete security rule by id
     * @request DELETE:/api/v1/security-rules/{id}
     */
    deleteSecurityRule: (id: number, params: RequestParams = {}) =>
      this.request<void, void>({
        path: `/api/v1/security-rules/${id}`,
        method: "DELETE",
        ...params,
      }),

    /**
     * No description
     *
     * @tags traffic-profiles
     * @name ListTrafficProfiles
     * @summary List traffic profiles
     * @request GET:/api/v1/traffic-profiles
     */
    listTrafficProfiles: (params: RequestParams = {}) =>
      this.request<TrafficProfile[], any>({
        path: `/api/v1/traffic-profiles`,
        method: "GET",
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags traffic-profiles
     * @name CreateTrafficProfile
     * @summary Create traffic profile
     * @request POST:/api/v1/traffic-profiles
     */
    createTrafficProfile: (
      data: TrafficProfileCreate,
      params: RequestParams = {},
    ) =>
      this.request<TrafficProfile, any>({
        path: `/api/v1/traffic-profiles`,
        method: "POST",
        body: data,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags traffic-profiles
     * @name GetTrafficProfile
     * @summary Get traffic profile by id
     * @request GET:/api/v1/traffic-profiles/{id}
     */
    getTrafficProfile: (id: number, params: RequestParams = {}) =>
      this.request<TrafficProfile, void>({
        path: `/api/v1/traffic-profiles/${id}`,
        method: "GET",
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags traffic-profiles
     * @name UpdateTrafficProfile
     * @summary Update traffic profile by id
     * @request PUT:/api/v1/traffic-profiles/{id}
     */
    updateTrafficProfile: (
      id: number,
      data: TrafficProfileUpdate,
      params: RequestParams = {},
    ) =>
      this.request<TrafficProfile, void>({
        path: `/api/v1/traffic-profiles/${id}`,
        method: "PUT",
        body: data,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags traffic-profiles
     * @name DeleteTrafficProfile
     * @summary Delete traffic profile by id
     * @request DELETE:/api/v1/traffic-profiles/{id}
     */
    deleteTrafficProfile: (id: number, params: RequestParams = {}) =>
      this.request<void, void>({
        path: `/api/v1/traffic-profiles/${id}`,
        method: "DELETE",
        ...params,
      }),

    /**
     * No description
     *
     * @tags traffic-rules
     * @name ListTrafficRules
     * @summary List traffic rules
     * @request GET:/api/v1/traffic-rules
     */
    listTrafficRules: (params: RequestParams = {}) =>
      this.request<TrafficRule[], any>({
        path: `/api/v1/traffic-rules`,
        method: "GET",
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags traffic-rules
     * @name CreateTrafficRule
     * @summary Create traffic rule
     * @request POST:/api/v1/traffic-rules
     */
    createTrafficRule: (data: TrafficRuleCreate, params: RequestParams = {}) =>
      this.request<TrafficRule, any>({
        path: `/api/v1/traffic-rules`,
        method: "POST",
        body: data,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags traffic-rules
     * @name GetTrafficRule
     * @summary Get traffic rule by id
     * @request GET:/api/v1/traffic-rules/{id}
     */
    getTrafficRule: (id: number, params: RequestParams = {}) =>
      this.request<TrafficRule, void>({
        path: `/api/v1/traffic-rules/${id}`,
        method: "GET",
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags traffic-rules
     * @name UpdateTrafficRule
     * @summary Update traffic rule by id
     * @request PUT:/api/v1/traffic-rules/{id}
     */
    updateTrafficRule: (
      id: number,
      data: TrafficRuleUpdate,
      params: RequestParams = {},
    ) =>
      this.request<TrafficRule, void>({
        path: `/api/v1/traffic-rules/${id}`,
        method: "PUT",
        body: data,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags traffic-rules
     * @name DeleteTrafficRule
     * @summary Delete traffic rule by id
     * @request DELETE:/api/v1/traffic-rules/{id}
     */
    deleteTrafficRule: (id: number, params: RequestParams = {}) =>
      this.request<void, void>({
        path: `/api/v1/traffic-rules/${id}`,
        method: "DELETE",
        ...params,
      }),

    /**
     * No description
     *
     * @tags ml-models
     * @name ListMlModels
     * @summary List ML models
     * @request GET:/api/v1/ml-models
     */
    listMlModels: (params: RequestParams = {}) =>
      this.request<MLModel[], any>({
        path: `/api/v1/ml-models`,
        method: "GET",
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags ml-models
     * @name CreateMlModel
     * @summary Create ML model
     * @request POST:/api/v1/ml-models
     */
    createMlModel: (data: MLModelCreate, params: RequestParams = {}) =>
      this.request<MLModel, any>({
        path: `/api/v1/ml-models`,
        method: "POST",
        body: data,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags ml-models
     * @name GetMlModel
     * @summary Get ML model by id
     * @request GET:/api/v1/ml-models/{id}
     */
    getMlModel: (id: number, params: RequestParams = {}) =>
      this.request<MLModel, void>({
        path: `/api/v1/ml-models/${id}`,
        method: "GET",
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags ml-models
     * @name UpdateMlModel
     * @summary Update ML model by id
     * @request PUT:/api/v1/ml-models/{id}
     */
    updateMlModel: (
      id: number,
      data: MLModelUpdate,
      params: RequestParams = {},
    ) =>
      this.request<MLModel, void>({
        path: `/api/v1/ml-models/${id}`,
        method: "PUT",
        body: data,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags ml-models
     * @name DeleteMlModel
     * @summary Delete ML model by id
     * @request DELETE:/api/v1/ml-models/{id}
     */
    deleteMlModel: (id: number, params: RequestParams = {}) =>
      this.request<void, void>({
        path: `/api/v1/ml-models/${id}`,
        method: "DELETE",
        ...params,
      }),

    /**
     * No description
     *
     * @tags request-logs
     * @name ListRequestLogs
     * @summary List request logs
     * @request GET:/api/v1/request-logs
     */
    listRequestLogs: (params: RequestParams = {}) =>
      this.request<RequestLog[], any>({
        path: `/api/v1/request-logs`,
        method: "GET",
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags request-logs
     * @name GetRequestLog
     * @summary Get request log by id
     * @request GET:/api/v1/request-logs/{id}
     */
    getRequestLog: (id: number, params: RequestParams = {}) =>
      this.request<RequestLog, void>({
        path: `/api/v1/request-logs/${id}`,
        method: "GET",
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags event-logs
     * @name ListEventLogs
     * @summary List event logs
     * @request GET:/api/v1/event-logs
     */
    listEventLogs: (params: RequestParams = {}) =>
      this.request<EventLog[], any>({
        path: `/api/v1/event-logs`,
        method: "GET",
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags event-logs
     * @name GetEventLog
     * @summary Get event log by id
     * @request GET:/api/v1/event-logs/{id}
     */
    getEventLog: (id: number, params: RequestParams = {}) =>
      this.request<EventLog, void>({
        path: `/api/v1/event-logs/${id}`,
        method: "GET",
        format: "json",
        ...params,
      }),
  };
}
