import { TimeSelector } from "@admin/components/fields/TimeSelector";
import { useAppConfig } from "@core/config/useConfig";
import { TextInput } from "react-admin";

export const ScheduleEntriesInputs = () => {
  const { party_days: partyDays, event_durations: eventDurations } =
    useAppConfig();

  return (
    <>
      <TextInput source="title" required />

      <TextInput source="description" multiline />

      <TimeSelector partyDays={partyDays} presetDurations={eventDurations} />

      <TextInput source="location" />

      <TextInput source="category" />
    </>
  );
};
