import jsonServerProvider from "ra-data-json-server";
import { DataProvider } from "react-admin";

export const dataProvider: DataProvider = jsonServerProvider("/api/admin");
