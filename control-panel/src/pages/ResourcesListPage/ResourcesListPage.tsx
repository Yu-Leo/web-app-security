import { useEffect, useMemo, useState } from "react";
import { apiClient } from "../../core/api/client";
import ResourcesListPageView, {
  type ResourceRow,
  type SelectOption,
} from "./ResourcesListPageView";

function ResourcesListPage() {
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [nameFilter, setNameFilter] = useState("");
  const [urlFilter, setUrlFilter] = useState("");
  const [securityFilter, setSecurityFilter] = useState("");
  const [trafficFilter, setTrafficFilter] = useState("");
  const [resources, setResources] = useState<ResourceRow[]>([]);
  const [securityProfileOptions, setSecurityProfileOptions] = useState<
    SelectOption[]
  >([]);
  const [trafficProfileOptions, setTrafficProfileOptions] = useState<
    SelectOption[]
  >([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [createName, setCreateName] = useState("");
  const [createUrlPattern, setCreateUrlPattern] = useState("");
  const [createSecurityProfileId, setCreateSecurityProfileId] = useState("");
  const [createTrafficProfileId, setCreateTrafficProfileId] = useState("");
  const [createAttempted, setCreateAttempted] = useState(false);
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);

  const loadResources = async () => {
    setIsLoading(true);
    setError(null);
    try {
      const [
        resourcesResponse,
        securityProfilesResponse,
        trafficProfilesResponse,
      ] = await Promise.all([
        apiClient.api.listResources(),
        apiClient.api.listSecurityProfiles(),
        apiClient.api.listTrafficProfiles(),
      ]);

      const securityOptions = securityProfilesResponse.data.map((profile) => ({
        value: profile.id.toString(),
        label: profile.name,
      }));
      const trafficOptions = trafficProfilesResponse.data.map((profile) => ({
        value: profile.id.toString(),
        label: profile.name,
      }));
      const securityMap = new Map(
        securityProfilesResponse.data.map((profile) => [
          profile.id,
          profile.name,
        ]),
      );
      const trafficMap = new Map(
        trafficProfilesResponse.data.map((profile) => [
          profile.id,
          profile.name,
        ]),
      );
      const resourceRows = resourcesResponse.data.map((resource) => ({
        id: resource.id,
        name: resource.name,
        createdAt: resource.created_at,
        urlPattern: resource.url_pattern,
        securityProfileId:
          resource.security_profile_id === null ||
          resource.security_profile_id === undefined
            ? ""
            : resource.security_profile_id.toString(),
        trafficProfileId:
          resource.traffic_profile_id === null ||
          resource.traffic_profile_id === undefined
            ? ""
            : resource.traffic_profile_id.toString(),
        securityProfile:
          resource.security_profile_id === null ||
          resource.security_profile_id === undefined
            ? "Не задан"
            : (securityMap.get(resource.security_profile_id) ?? "—"),
        trafficProfile:
          resource.traffic_profile_id === null ||
          resource.traffic_profile_id === undefined
            ? "Не задан"
            : (trafficMap.get(resource.traffic_profile_id) ?? "—"),
      }));

      setResources(resourceRows);
      setSecurityProfileOptions(securityOptions);
      setTrafficProfileOptions(trafficOptions);
    } catch (requestError) {
      setError("Не удалось загрузить ресурсы.");
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    void loadResources();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const filteredResources = useMemo(() => {
    return resources
      .filter((resource) => {
        const nameMatch = resource.name
          .toLowerCase()
          .includes(nameFilter.toLowerCase());
        const urlMatch = resource.urlPattern
          .toLowerCase()
          .includes(urlFilter.toLowerCase());
        const securityMatch = securityFilter
          ? resource.securityProfileId === securityFilter
          : true;
        const trafficMatch = trafficFilter
          ? resource.trafficProfileId === trafficFilter
          : true;
        return nameMatch && urlMatch && securityMatch && trafficMatch;
      })
      .sort(
        (a, b) =>
          new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime(),
      );
  }, [
    nameFilter,
    urlFilter,
    securityFilter,
    trafficFilter,
    resources,
    securityProfileOptions,
    trafficProfileOptions,
  ]);

  useEffect(() => {
    setCurrentPage(1);
  }, [nameFilter, urlFilter, securityFilter, trafficFilter]);

  const totalPages = Math.max(
    1,
    Math.ceil(filteredResources.length / pageSize),
  );
  const normalizedCurrentPage = Math.min(currentPage, totalPages);

  const paginatedResources = useMemo(() => {
    const start = (normalizedCurrentPage - 1) * pageSize;
    return filteredResources.slice(start, start + pageSize);
  }, [filteredResources, normalizedCurrentPage, pageSize]);

  const handleCreateSubmit = async () => {
    setCreateAttempted(true);
    if (!createName.trim() || !createUrlPattern.trim()) {
      return;
    }
    setError(null);
    try {
      await apiClient.api.createResource({
        name: createName,
        url_pattern: createUrlPattern,
        security_profile_id: createSecurityProfileId
          ? Number(createSecurityProfileId)
          : null,
        traffic_profile_id: createTrafficProfileId
          ? Number(createTrafficProfileId)
          : null,
      });
      setShowCreateModal(false);
      setCreateAttempted(false);
      setCreateName("");
      setCreateUrlPattern("");
      await loadResources();
    } catch (requestError) {
      setError("Не удалось создать ресурс.");
    }
  };

  return (
    <ResourcesListPageView
      resources={paginatedResources}
      securityProfileOptions={securityProfileOptions}
      trafficProfileOptions={trafficProfileOptions}
      nameFilter={nameFilter}
      urlFilter={urlFilter}
      securityFilter={securityFilter}
      trafficFilter={trafficFilter}
      onNameFilterChange={setNameFilter}
      onUrlFilterChange={setUrlFilter}
      onSecurityFilterChange={setSecurityFilter}
      onTrafficFilterChange={setTrafficFilter}
      isLoading={isLoading}
      error={error}
      showCreateModal={showCreateModal}
      onOpenCreateModal={() => setShowCreateModal(true)}
      onCloseCreateModal={() => {
        setShowCreateModal(false);
        setCreateAttempted(false);
      }}
      createName={createName}
      createUrlPattern={createUrlPattern}
      createSecurityProfileId={createSecurityProfileId}
      createTrafficProfileId={createTrafficProfileId}
      onCreateNameChange={setCreateName}
      onCreateUrlPatternChange={setCreateUrlPattern}
      onCreateSecurityProfileIdChange={setCreateSecurityProfileId}
      onCreateTrafficProfileIdChange={setCreateTrafficProfileId}
      onCreateSubmit={handleCreateSubmit}
      createAttempted={createAttempted}
      currentPage={normalizedCurrentPage}
      totalPages={totalPages}
      onPageChange={setCurrentPage}
      pageSize={pageSize}
      onPageSizeChange={(value) => {
        setPageSize(value);
        setCurrentPage(1);
      }}
    />
  );
}

export default ResourcesListPage;
