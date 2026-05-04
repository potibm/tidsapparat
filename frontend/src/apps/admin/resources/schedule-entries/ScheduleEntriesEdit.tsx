import { Edit, SimpleForm } from "react-admin";
import { ScheduleEntriesInputs } from "./ScheduleEntriesInputs";
import { transformScheduleToAPI } from "@admin/utils/time";

export const ScheduleEntriesEdit = () => {
  return (
    <Edit title="Edit Event" transform={transformScheduleToAPI}>
      <SimpleForm defaultValues={{ source: "*", type: "username" }}>
        <ScheduleEntriesInputs />
      </SimpleForm>
    </Edit>
  );
};

export default ScheduleEntriesEdit;
