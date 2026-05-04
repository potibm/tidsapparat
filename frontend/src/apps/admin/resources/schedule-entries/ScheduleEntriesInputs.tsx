import { TimeSelector } from "@admin/components/fields/TimeSelector";
import { TextInput } from "react-admin";

const MOCK_CONFIG_FROM_BACKEND = {
  days: [
    { id: "2026-05-08", name: "Friday" },
    { id: "2026-05-09", name: "Saturday" },
  ],
  durations: [0, 15, 30, 45, 60, 90, 120],
};

export const ScheduleEntriesInputs = () => {
  return (
    <>
      <TextInput source="title" required />

      <TextInput source="description" multiline />

      <TimeSelector
        partyDays={MOCK_CONFIG_FROM_BACKEND.days}
        presetDurations={MOCK_CONFIG_FROM_BACKEND.durations}
      />

      <TextInput source="location" />

      <TextInput source="category" />
    </>
  );
};
