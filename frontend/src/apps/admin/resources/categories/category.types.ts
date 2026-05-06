import { RaRecord } from "react-admin";

export interface CategoryRecord extends RaRecord {
  id: number;
  name: string;
  color: string;
}
