import { Navigate, Route, Routes } from "react-router-dom";
import ResourcesListPage from "./pages/ResourcesListPage";
import ResourceDetailsPage from "./pages/ResourceDetailsPage";
import SecurityProfilesListPage from "./pages/SecurityProfilesListPage";
import SecurityProfileDetailsPage from "./pages/SecurityProfileDetailsPage";
import TrafficProfilesListPage from "./pages/TrafficProfilesListPage";
import TrafficProfileDetailsPage from "./pages/TrafficProfileDetailsPage";
import RequestLogsPage from "./pages/RequestLogsPage";
import EventLogsPage from "./pages/EventLogsPage";
import MLModelsListPage from "./pages/MLModelsListPage";
import MLModelDetailsPage from "./pages/MLModelDetailsPage";

function AppRoutes() {
  return (
    <Routes>
      <Route path="/" element={<Navigate to="/resources" replace />} />
      <Route path="/resources" element={<ResourcesListPage />} />
      <Route path="/resources/:id" element={<ResourceDetailsPage />} />
      <Route path="/security-profiles" element={<SecurityProfilesListPage />} />
      <Route
        path="/security-profiles/:id"
        element={<SecurityProfileDetailsPage />}
      />
      <Route path="/traffic-profiles" element={<TrafficProfilesListPage />} />
      <Route
        path="/traffic-profiles/:id"
        element={<TrafficProfileDetailsPage />}
      />
      <Route path="/request-logs" element={<RequestLogsPage />} />
      <Route path="/event-logs" element={<EventLogsPage />} />
      <Route path="/ml-models" element={<MLModelsListPage />} />
      <Route path="/ml-models/:id" element={<MLModelDetailsPage />} />
      <Route path="*" element={<Navigate to="/resources" replace />} />
    </Routes>
  );
}

export default AppRoutes;
