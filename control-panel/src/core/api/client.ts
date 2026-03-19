import { Api } from "./Api";

const apiBaseUrl = import.meta.env.VITE_API_BASE_URL ?? "http://localhost:8001";

export const apiClient = new Api({
  baseURL: apiBaseUrl,
});

apiClient.instance.interceptors.request.use((config) => {
  const method = (config.method ?? "get").toLowerCase();
  if (method === "get") {
    if (config.headers && "set" in config.headers) {
      config.headers.set("Cache-Control", "no-store, no-cache, max-age=0");
      config.headers.set("Pragma", "no-cache");
      config.headers.set("Expires", "0");
    } else {
      config.headers = {
        "Cache-Control": "no-store, no-cache, max-age=0",
        Pragma: "no-cache",
        Expires: "0",
      } as unknown as typeof config.headers;
    }

    const params =
      typeof config.params === "object" && config.params !== null
        ? { ...config.params }
        : {};
    config.params = {
      ...params,
      _ts: Date.now(),
    };
  }

  return config;
});
