import {
  List,
  DatagridConfigurable,
  TextField,
  SearchInput,
  FunctionField,
  TopToolbar,
  CreateButton,
  ExportButton,
  SelectColumnsButton,
  BooleanInput,
  ReferenceField,
  ChipField,
  ReferenceInput,
  SelectInput,
  FilterButton,
} from "react-admin";
import { Chip } from "@mui/material";
import dayjs from "dayjs";
import isBetween from "dayjs/plugin/isBetween";
import { ScheduleEntryRecord } from "./schedule_entry.types";
import { LocationWithIcon } from "@admin/components/fields/LocationWithIcon";
import { BooleanToggleField } from "@admin/components/fields/BooleanToggleField";

dayjs.extend(isBetween);

const scheduleFilters = [
  <SearchInput key="q" source="q" alwaysOn placeholder="Search by title..." />,

  <ReferenceInput key="category_id" source="category_id" reference="categories">
    <SelectInput label="Category" optionText="name" />
  </ReferenceInput>,

  <ReferenceInput key="location_id" source="location_id" reference="locations">
    <SelectInput label="Location" optionText="name" />
  </ReferenceInput>,

  <BooleanInput key="hide_past" source="hide_past" label="Hide past" />,

  <BooleanInput
    key="hide_hidden"
    source="hide_hidden"
    label="Show only visible"
  />,
];

const ListActions = () => (
  <TopToolbar>
    <FilterButton />
    <SelectColumnsButton />
    <CreateButton />
    <ExportButton />
  </TopToolbar>
);

export const ScheduleEntriesList = () => (
  <List
    title="Events"
    sort={{ field: "id", order: "DESC" }}
    actions={<ListActions />}
    filters={scheduleFilters}
    filterDefaultValues={{ hide_past: true, hide_hidden: true }}
  >
    <DatagridConfigurable
      rowClick="edit"
      bulkActionButtons={false}
      omit={["end_time_display", "category_id", "location_id"]}
    >
      {/* 1. STATUS FIELD */}
      <FunctionField
        label="Status"
        sortable={false}
        render={(record: ScheduleEntryRecord) => {
          if (!record.start_time || !record.end_time) return null;

          const now = dayjs();
          const start = dayjs(record.start_time);
          const end = dayjs(record.end_time);
          const hidden = record.hidden;

          if (hidden) {
            return <Chip label="Hidden" size="small" variant="outlined" />;
          } else if (now.isBetween(start, end)) {
            return (
              <Chip
                label="Live"
                color="success"
                size="small"
                variant="filled"
              />
            );
          } else if (now.isAfter(end)) {
            return <Chip label="Done" size="small" variant="outlined" />;
          } else {
            return (
              <Chip
                label="Upcoming"
                color="primary"
                size="small"
                variant="outlined"
              />
            );
          }
        }}
      />

      {/* 3. TITLE */}
      <TextField source="title" label="Event Name" />

      {/* 3. START TIME */}
      <FunctionField
        label="Start"
        sortBy="start_time"
        render={(record: ScheduleEntryRecord) => {
          if (!record.start_time) return "-";
          return dayjs(record.start_time).format("ddd, HH:mm");
        }}
      />

      {/* 3. END TIME (hidden) */}
      <FunctionField
        source="end_time_display"
        label="End"
        sortBy="end_time"
        render={(record: ScheduleEntryRecord) =>
          record.end_time ? dayjs(record.end_time).format("ddd, HH:mm") : "-"
        }
      />

      {/* 4. DURATION */}
      <FunctionField
        label="Duration"
        sortBy="end_time"
        sortable={false}
        render={(record: ScheduleEntryRecord) => {
          if (!record.start_time || !record.end_time) return "-";
          const start = dayjs(record.start_time);
          const end = dayjs(record.end_time);
          const diffMins = end.diff(start, "minute");
          const endTimeString = end.format("ddd, HH:mm");

          return (
            <span
              title={`Ends at: ${endTimeString}`}
              style={{ cursor: "help", borderBottom: "1px dotted #888" }}
            >
              {diffMins}m
            </span>
          );
        }}
      />

      {/* 5. LOCATION / CATEGORY */}
      <ReferenceField
        source="location_id"
        reference="locations"
        label="Location"
      >
        <LocationWithIcon />
      </ReferenceField>
      <ReferenceField
        source="category_id"
        reference="categories"
        label="Category"
      >
        <ChipField source="name" />
      </ReferenceField>

      <BooleanToggleField source="hidden" label="Hidden" />
    </DatagridConfigurable>
  </List>
);

export default ScheduleEntriesList;
