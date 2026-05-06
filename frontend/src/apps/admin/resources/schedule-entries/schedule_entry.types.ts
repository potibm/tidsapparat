import { RaRecord } from "react-admin";
import { CategoryRecord } from "../categories/category.types";

export interface ScheduleEntryRecord extends RaRecord {
  id: string | number;
  title: string;
  start_time: string; // ISO-String aus der DB
  end_time: string; // ISO-String aus der DB
  category_id?: number;
  category?: CategoryRecord;
  location?: string;
}

export interface ScheduleFormData {
  start_time_only: string | Date | number;
  party_day: string;
  duration_mins: number;
  title?: string;
  category_id?: number;
  location?: string;
  [key: string]: unknown;
}
