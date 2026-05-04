import { Create, SimpleForm } from "react-admin";
import { ScheduleEntriesInputs } from "./ScheduleEntriesInputs";
import { transformScheduleToAPI } from "@admin/utils/time";

export const ScheduleEntriesCreate = () => {
  return (
    <Create title="Add Event" transform={transformScheduleToAPI}>
      <SimpleForm defaultValues={{ source: "*", type: "username" }}>
        <ScheduleEntriesInputs />
      </SimpleForm>
    </Create>
  );
};

export default ScheduleEntriesCreate;
