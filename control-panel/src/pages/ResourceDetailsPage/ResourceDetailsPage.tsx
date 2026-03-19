import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { apiClient } from "../../core/api/client";
import ResourceDetailsPageView, {
  type SelectOption,
} from "./ResourceDetailsPageView";

function ResourceDetailsPage() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [name, setName] = useState("");
  const [urlPattern, setUrlPattern] = useState("");
  const [securityProfileId, setSecurityProfileId] = useState("");
  const [trafficProfileId, setTrafficProfileId] = useState("");
  const [securityProfiles, setSecurityProfiles] = useState<SelectOption[]>([]);
  const [trafficProfiles, setTrafficProfiles] = useState<SelectOption[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [showSuccessToast, setShowSuccessToast] = useState(false);

  const resourceId = Number(id);

  const loadResource = async () => {
    if (!resourceId) {
      setError("Некорректный ID ресурса.");
      return;
    }
    setIsLoading(true);
    setError(null);
    try {
      const [
        resourceResponse,
        securityProfilesResponse,
        trafficProfilesResponse,
      ] = await Promise.all([
        apiClient.api.getResource(resourceId),
        apiClient.api.listSecurityProfiles(),
        apiClient.api.listTrafficProfiles(),
      ]);

      const resource = resourceResponse.data;
      setName(resource.name);
      setUrlPattern(resource.url_pattern);
      setSecurityProfileId(
        resource.security_profile_id === null ||
          resource.security_profile_id === undefined
          ? ""
          : resource.security_profile_id.toString(),
      );
      setTrafficProfileId(
        resource.traffic_profile_id === null ||
          resource.traffic_profile_id === undefined
          ? ""
          : resource.traffic_profile_id.toString(),
      );
      setSecurityProfiles(
        securityProfilesResponse.data.map((profile) => ({
          value: profile.id.toString(),
          label: profile.name,
        })),
      );
      setTrafficProfiles(
        trafficProfilesResponse.data.map((profile) => ({
          value: profile.id.toString(),
          label: profile.name,
        })),
      );
    } catch (requestError) {
      setError("Не удалось загрузить ресурс.");
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    void loadResource();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [id]);

  const handleSave = async () => {
    if (!resourceId) {
      return;
    }
    setError(null);
    setIsLoading(true);
    try {
      await apiClient.api.updateResource(resourceId, {
        name,
        url_pattern: urlPattern,
        security_profile_id: securityProfileId
          ? Number(securityProfileId)
          : null,
        traffic_profile_id: trafficProfileId ? Number(trafficProfileId) : null,
      });
      setShowSuccessToast(true);
    } catch (requestError) {
      setError("Не удалось сохранить ресурс.");
    } finally {
      setIsLoading(false);
    }
  };

  const handleDelete = async () => {
    if (!resourceId) {
      return;
    }
    setError(null);
    setIsLoading(true);
    try {
      await apiClient.api.deleteResource(resourceId);
      navigate("/resources");
    } catch (requestError) {
      setError("Не удалось удалить ресурс.");
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <ResourceDetailsPageView
      title={name.trim() ? `Ресурс: ${name}` : "Ресурс"}
      name={name}
      urlPattern={urlPattern}
      securityProfileId={securityProfileId}
      trafficProfileId={trafficProfileId}
      securityProfiles={securityProfiles}
      trafficProfiles={trafficProfiles}
      isLoading={isLoading}
      error={error}
      onNameChange={setName}
      onUrlPatternChange={setUrlPattern}
      onSecurityProfileChange={setSecurityProfileId}
      onTrafficProfileChange={setTrafficProfileId}
      onSave={handleSave}
      onDelete={handleDelete}
      showSuccessToast={showSuccessToast}
      onCloseSuccessToast={() => setShowSuccessToast(false)}
    />
  );
}

export default ResourceDetailsPage;
