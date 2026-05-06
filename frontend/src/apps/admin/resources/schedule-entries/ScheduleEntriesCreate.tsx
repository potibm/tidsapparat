import { Create, SimpleForm } from "react-admin";
import { ScheduleEntriesInputs } from "./ScheduleEntriesInputs";
import { transformScheduleToAPI } from "@admin/utils/time";
import { useAppConfig } from "@core/config/useConfig";
import { ScheduleFormData } from "./schedule_entry.types";

export const ScheduleEntriesCreate = () => {
  const { timezone } = useAppConfig();

  const transform = (data: ScheduleFormData) =>
    transformScheduleToAPI(data, timezone);

  return (
    <Create title="Add Event" transform={transform}>
      <SimpleForm defaultValues={{ source: "*", type: "username" }}>
        <ScheduleEntriesInputs />
      </SimpleForm>
    </Create>
  );
};

export default ScheduleEntriesCreate;
