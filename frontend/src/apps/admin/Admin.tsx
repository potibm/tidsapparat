import { Admin, Resource } from "react-admin";
import { MyTheme, MyDarkTheme } from "./theme/MyTheme";
import { MyLayout } from "./theme/MyLayout";
import { dataProvider } from "./providers/dataProvider";
import scheduleEntries from "./resources/schedule-entries";
import categories from "./resources/categories";
import locations from "./resources/locations";

export const AdminApp = () => (
  <Admin
    dataProvider={dataProvider}
    theme={MyTheme}
    darkTheme={MyDarkTheme}
    layout={MyLayout}
    title="Tidsapparat Admin"
  >
    <Resource {...scheduleEntries} />
    <Resource {...categories} />
    <Resource {...locations} />
  </Admin>
);
