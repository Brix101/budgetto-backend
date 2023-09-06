import App from "@/App.tsx";
import { ThemeProvider } from "@/components/theme-provider";
import "@/index.css";
import queryClient from "@/lib/queryClient";
import { Auth0Provider } from "@auth0/auth0-react";
import { QueryClientProvider } from "@tanstack/react-query";
import React from "react";
import ReactDOM from "react-dom/client";

ReactDOM.createRoot(document.getElementById("root")!).render(
  <Auth0Provider
    domain={import.meta.env.VITE_AUTH0_DOMAIN}
    clientId={import.meta.env.VITE_AUTH0_CLIENT_ID}
    authorizationParams={{
      audience: import.meta.env.VITE_AUTH0_AUDIENCE,
      redirect_uri: window.location.origin,
    }}
  >
    <React.StrictMode>
      <ThemeProvider defaultTheme="light" storageKey="budgetto-ui-theme">
        <QueryClientProvider client={queryClient}>
          <App />
        </QueryClientProvider>
      </ThemeProvider>
    </React.StrictMode>
  </Auth0Provider>
);
