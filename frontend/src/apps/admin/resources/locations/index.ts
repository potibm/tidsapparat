import LocationOnIcon from "@mui/icons-material/LocationOn";
import LocationsList from "./LocationsList";
import LocationsCreate from "./LocationsCreate";
import LocationsEdit from "./LocationsEdit";

export default {
  name: "locations",
  options: { label: "Locations" },
  list: LocationsList,
  create: LocationsCreate,
  edit: LocationsEdit,
  icon: LocationOnIcon,
};
