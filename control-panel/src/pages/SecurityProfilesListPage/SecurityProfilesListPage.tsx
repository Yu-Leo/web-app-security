import { useEffect, useMemo, useState } from "react";
import { apiClient } from "../../core/api/client";
import SecurityProfilesListPageView, {
  type SecurityProfileRow,
} from "./SecurityProfilesListPageView";

function SecurityProfilesListPage() {
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [nameFilter, setNameFilter] = useState("");
  const [baseActionFilter, setBaseActionFilter] = useState("");
  const [statusFilter, setStatusFilter] = useState("");
  const [profiles, setProfiles] = useState<SecurityProfileRow[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [createName, setCreateName] = useState("");
  const [createDescription, setCreateDescription] = useState("");
  const [createBaseAction, setCreateBaseAction] = useState<"allow" | "block">(
    "allow",
  );
  const [createLogEnabled, setCreateLogEnabled] = useState(true);
  const [createIsEnabled, setCreateIsEnabled] = useState(true);
  const [createAttempted, setCreateAttempted] = useState(false);
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);

  const loadProfiles = async () => {
    setIsLoading(true);
    setError(null);
    try {
      const response = await apiClient.api.listSecurityProfiles();
      setProfiles(
        response.data.map((profile) => ({
          id: profile.id,
          createdAt: profile.created_at,
          name: profile.name,
          baseAction: profile.base_action,
          logEnabled: profile.log_enabled,
          isEnabled: profile.is_enabled,
        })),
      );
    } catch (requestError) {
      setError("Не удалось загрузить профили безопасности.");
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    void loadProfiles();
  }, []);

  const filteredProfiles = useMemo(() => {
    return profiles
      .filter((profile) => {
        const nameMatch = profile.name
          .toLowerCase()
          .includes(nameFilter.toLowerCase());
        const actionMatch = baseActionFilter
          ? profile.baseAction === baseActionFilter
          : true;
        const statusMatch = statusFilter
          ? statusFilter === "enabled"
            ? profile.isEnabled
            : !profile.isEnabled
          : true;
        return nameMatch && actionMatch && statusMatch;
      })
      .sort(
        (a, b) =>
          new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime(),
      );
  }, [nameFilter, baseActionFilter, statusFilter, profiles]);

  useEffect(() => {
    setCurrentPage(1);
  }, [nameFilter, baseActionFilter, statusFilter]);

  const totalPages = Math.max(1, Math.ceil(filteredProfiles.length / pageSize));
  const normalizedCurrentPage = Math.min(currentPage, totalPages);

  const paginatedProfiles = useMemo(() => {
    const start = (normalizedCurrentPage - 1) * pageSize;
    return filteredProfiles.slice(start, start + pageSize);
  }, [filteredProfiles, normalizedCurrentPage, pageSize]);

  const handleCreateSubmit = async () => {
    setCreateAttempted(true);
    if (!createName.trim()) {
      return;
    }
    setError(null);
    try {
      await apiClient.api.createSecurityProfile({
        name: createName,
        description: createDescription || undefined,
        base_action: createBaseAction,
        log_enabled: createLogEnabled,
        is_enabled: createIsEnabled,
      });
      setShowCreateModal(false);
      setCreateAttempted(false);
      setCreateName("");
      setCreateDescription("");
      await loadProfiles();
    } catch (requestError) {
      setError("Не удалось создать профиль безопасности.");
    }
  };

  return (
    <SecurityProfilesListPageView
      profiles={paginatedProfiles}
      nameFilter={nameFilter}
      baseActionFilter={baseActionFilter}
      statusFilter={statusFilter}
      onNameFilterChange={setNameFilter}
      onBaseActionFilterChange={setBaseActionFilter}
      onStatusFilterChange={setStatusFilter}
      isLoading={isLoading}
      error={error}
      showCreateModal={showCreateModal}
      onOpenCreateModal={() => setShowCreateModal(true)}
      onCloseCreateModal={() => {
        setShowCreateModal(false);
        setCreateAttempted(false);
      }}
      createName={createName}
      createDescription={createDescription}
      createBaseAction={createBaseAction}
      createLogEnabled={createLogEnabled}
      createIsEnabled={createIsEnabled}
      onCreateNameChange={setCreateName}
      onCreateDescriptionChange={setCreateDescription}
      onCreateBaseActionChange={setCreateBaseAction}
      onCreateLogEnabledChange={setCreateLogEnabled}
      onCreateIsEnabledChange={setCreateIsEnabled}
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

export default SecurityProfilesListPage;
