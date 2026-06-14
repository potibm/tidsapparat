import jsonServerProvider from "ra-data-json-server";
import { DataProvider, fetchUtils } from "react-admin";
import { getAccessToken } from "./authProvider"; // Pfad ggf. anpassen

const httpClient = async (url: string, options: fetchUtils.Options = {}) => {
  if (!options.headers) {
    options.headers = new Headers({ Accept: "application/json" });
  }

  const headers =
    options.headers instanceof Headers
      ? options.headers
      : new Headers(options.headers);

  const token = await getAccessToken();
  if (token) {
    headers.set("Authorization", `Bearer ${token}`);
  }

  options.headers = headers;

  return fetchUtils.fetchJson(url, options);
};

export const dataProvider: DataProvider = jsonServerProvider(
  "/api/admin", 
  httpClient
);