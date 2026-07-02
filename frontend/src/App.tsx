import * as Sentry from "@sentry/react";
import { AdminApp } from "@admin/Admin";

function App() {
  return (
    <div className="App">
      <Sentry.ErrorBoundary
        fallback={
          <p>
            A serious error has occurred. Please restart the beamer application.
          </p>
        }
      >
        <AdminApp />
      </Sentry.ErrorBoundary>
    </div>
  );
}

export default App;
