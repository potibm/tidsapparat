import { TimeSelector } from "@admin/components/fields/TimeSelector";
import { useAppConfig } from "@core/config/useConfig";
import { TextInput, ReferenceInput, SelectInput } from "react-admin";

export const ScheduleEntriesInputs = () => {
  const { party_days: partyDays, event_durations: eventDurations } =
    useAppConfig();

  return (
    <>
      <TextInput source="title" required />

      <TextInput source="description" multiline />

      <TimeSelector partyDays={partyDays} presetDurations={eventDurations} />

      <ReferenceInput source="location_id" reference="locations">
        <SelectInput label="Location" optionText="name" />
      </ReferenceInput>

      <ReferenceInput source="category_id" reference="categories">
        <SelectInput label="Category" optionText="name" />
      </ReferenceInput>
    </>
  );
};
