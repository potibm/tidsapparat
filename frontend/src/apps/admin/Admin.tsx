import { Admin, Resource } from "react-admin";
import { BrowserRouter } from "react-router";
import { MyTheme, MyDarkTheme } from "./theme/MyTheme";
import { MyLayout } from "./theme/MyLayout";
import { dataProvider } from "./providers/dataProvider";
import scheduleEntries from "./resources/schedule-entries";
import categories from "./resources/categories";
import locations from "./resources/locations";
import { authProvider, configureOidc } from "./providers/authProvider";
import { OidcLogin } from "./components/auth/OidcLogin";
import { useAppConfig } from "@core/config/useConfig";

export const AdminApp = () => (
  <BrowserRouter>
    <AppBootstrapper />
  </BrowserRouter>
);

export const AppBootstrapper = () => {
  const appConfig = useAppConfig();

  const isOidcActive = appConfig.auth?.type === "oidc";
  if (isOidcActive) {
    configureOidc(appConfig.auth.authority, appConfig.auth.client_id);
  }

  return (
    <Admin
      loginPage={isOidcActive ? OidcLogin : false}
      dataProvider={dataProvider}
      authProvider={isOidcActive ? authProvider : undefined}
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
};
