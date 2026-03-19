import { useEffect, useMemo, useState } from "react";
import { useSearchParams } from "react-router-dom";
import { apiClient } from "../../core/api/client";
import EventLogsPageView, { type EventLogRow } from "./EventLogsPageView";

function EventLogsPage() {
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [searchParams, setSearchParams] = useSearchParams();
  const [logs, setLogs] = useState<EventLogRow[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [filters, setFilters] = useState({
    resourceId: "",
    eventType: "",
    severity: "",
    message: "",
  });

  const loadLogs = async () => {
    setIsLoading(true);
    setError(null);
    try {
      const response = await apiClient.api.listEventLogs();
      setLogs(
        response.data.map((log) => {
          const logAny = log as Record<string, any>;
          return {
            id: log.id,
            occurredAt: log.occurred_at,
            resourceId: log.resource_id,
            eventType: log.event_type,
            severity: log.severity,
            message: log.message,
            ruleId: log.rule_id,
            profileId: log.profile_id,
            requestId: logAny.request_id,
            clientIp: logAny.client_ip,
            method: logAny.method,
            path: logAny.path,
            metadata: log.metadata,
          };
        }),
      );
    } catch (requestError) {
      setError("Не удалось загрузить логи событий.");
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    void loadLogs();
  }, []);

  const handleFilterChange = (field: keyof typeof filters, value: string) => {
    setFilters((prevFilters) => ({
      ...prevFilters,
      [field]: value,
    }));
  };

  const filteredLogs = useMemo(() => {
    return logs
      .filter((log) => {
        const resourceMatch = filters.resourceId
          ? log.resourceId.toString().includes(filters.resourceId)
          : true;
        const typeMatch = filters.eventType
          ? log.eventType
              .toLowerCase()
              .includes(filters.eventType.toLowerCase())
          : true;
        const severityMatch = filters.severity
          ? log.severity.toLowerCase().includes(filters.severity.toLowerCase())
          : true;
        const messageMatch = filters.message
          ? log.message.toLowerCase().includes(filters.message.toLowerCase())
          : true;
        return resourceMatch && typeMatch && severityMatch && messageMatch;
      })
      .sort(
        (a, b) =>
          new Date(b.occurredAt).getTime() - new Date(a.occurredAt).getTime(),
      );
  }, [filters, logs]);

  useEffect(() => {
    setCurrentPage(1);
  }, [filters]);

  const totalPages = Math.max(1, Math.ceil(filteredLogs.length / pageSize));
  const normalizedCurrentPage = Math.min(currentPage, totalPages);
  const paginatedLogs = useMemo(() => {
    const start = (normalizedCurrentPage - 1) * pageSize;
    return filteredLogs.slice(start, start + pageSize);
  }, [filteredLogs, normalizedCurrentPage, pageSize]);

  const logIdParam = searchParams.get("logId");
  const detailsLog = useMemo(() => {
    if (!logIdParam) {
      return null;
    }
    const logId = Number(logIdParam);
    return logs.find((log) => log.id === logId) ?? null;
  }, [logIdParam, logs]);

  useEffect(() => {
    if (logIdParam && !detailsLog && !isLoading) {
      setSearchParams({});
    }
  }, [logIdParam, detailsLog, isLoading, setSearchParams]);

  const handleCloseDetails = () => {
    setSearchParams({});
  };

  return (
    <EventLogsPageView
      logs={paginatedLogs}
      isLoading={isLoading}
      error={error}
      filters={filters}
      onFilterChange={handleFilterChange}
      showDetails={Boolean(logIdParam)}
      detailsLog={detailsLog}
      onCloseDetails={handleCloseDetails}
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

export default EventLogsPage;
