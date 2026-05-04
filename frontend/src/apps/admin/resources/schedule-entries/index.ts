import EventIcon from "@mui/icons-material/Event";
import ScheduleEntriesList from "./ScheduleEntriesList";
import ScheduleEntriesCreate from "./ScheduleEntriesCreate";
import ScheduleEntriesEdit from "./ScheduleEntriesEdit";

export default {
  name: "schedule-entries",
  options: { label: "Events" },
  list: ScheduleEntriesList,
  create: ScheduleEntriesCreate,
  edit: ScheduleEntriesEdit,
  icon: EventIcon,
};
