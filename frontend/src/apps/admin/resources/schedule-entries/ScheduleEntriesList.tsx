import {
  List,
  Datagrid,
  TextField,
  EditButton,
  DeleteButton,
  FunctionField,
  DateField,
} from "react-admin";

export const ScheduleEntriesList = () => (
  <List title="Events" sort={{ field: "id", order: "DESC" }}>
    <Datagrid rowClick="edit" bulkActionButtons={false}>
      <DateField
        source="start_time"
        label="Start"
        showTime
        options={{
          weekday: "short",
          hour: "2-digit",
          minute: "2-digit",
          hour12: false,
        }}
      />

      <FunctionField
        label="Duration"
        sortBy="end_time"
        render={(record: {
          start_time?: string;
          end_time?: string;
          title?: string;
        }) => {
          if (!record || !record.start_time || !record.end_time) return "-";

          const start = new Date(record.start_time);
          const end = new Date(record.end_time);

          const diffMins = Math.round(
            (end.getTime() - start.getTime()) / 60000,
          );

          const endTimeString = end.toLocaleString("en-US", {
            weekday: "short",
            hour: "2-digit",
            minute: "2-digit",
            hour12: false,
          });

          return <span title={`Ends at: ${endTimeString}`}>{diffMins}m</span>;
        }}
      />

      <TextField source="title" label="Title" />

      <EditButton />
      <DeleteButton />
    </Datagrid>
  </List>
);

export default ScheduleEntriesList;
