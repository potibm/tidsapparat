import { RaRecord } from "react-admin";
import { CategoryRecord } from "../categories/category.types";
import { LocationRecord } from "../locations/location.types";

export interface ScheduleEntryRecord extends RaRecord {
  id: string | number;
  title: string;
  start_time: string; // ISO string from the DB
  end_time: string; // ISO string from the DB
  category_id?: number;
  category?: CategoryRecord;
  location_id?: number;
  location?: LocationRecord;
}

export interface ScheduleFormData {
  start_time_only: string | Date | number;
  party_day: string;
  duration_mins: number;
  title?: string;
  category_id?: number;
  location_id?: number;
  [key: string]: unknown;
}
