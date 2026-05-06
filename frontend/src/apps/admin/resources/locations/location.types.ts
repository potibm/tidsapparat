import { RaRecord } from "react-admin";

export interface LocationRecord extends RaRecord {
  id: number;
  name: string;
  address?: string;
}
