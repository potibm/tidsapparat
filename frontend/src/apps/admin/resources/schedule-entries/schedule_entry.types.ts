import { RaRecord } from "react-admin";

export interface ScheduleEntryRecord extends RaRecord {
  id: string | number;
  title: string;
  start_time: string; // ISO-String aus der DB
  end_time: string; // ISO-String aus der DB
  category?: string;
  location?: string;
}

export interface ScheduleFormData {
  start_time_only: string | Date | number;
  party_day: string;
  duration_mins: number;
  title?: string;
  category?: string;
  location?: string;
  [key: string]: unknown;
}
