import { Admin, Resource } from "react-admin";
import { MyTheme, MyDarkTheme } from "./theme/MyTheme";
import { MyLayout } from "./theme/MyLayout";
import { dataProvider } from "./providers/dataProvider";
import scheduleEntries from "./resources/schedule-entries";

export const AdminApp = () => (
  <Admin
    dataProvider={dataProvider}
    theme={MyTheme}
    darkTheme={MyDarkTheme}
    layout={MyLayout}
    title="Tidsapparat Admin"
  >
    <Resource {...scheduleEntries} />
  </Admin>
);
