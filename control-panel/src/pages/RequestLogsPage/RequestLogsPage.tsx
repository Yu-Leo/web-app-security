import { useEffect, useMemo, useState } from "react";
import { useSearchParams } from "react-router-dom";
import { apiClient } from "../../core/api/client";
import RequestLogsPageView, { type RequestLogRow } from "./RequestLogsPageView";

function RequestLogsPage() {
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [searchParams, setSearchParams] = useSearchParams();
  const [logs, setLogs] = useState<RequestLogRow[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [filters, setFilters] = useState({
    resourceId: "",
    method: "",
    statusCode: "",
    action: "",
    clientIp: "",
    path: "",
  });

  const loadLogs = async () => {
    setIsLoading(true);
    setError(null);
    try {
      const response = await apiClient.api.listRequestLogs();
      setLogs(
        response.data.map((log) => {
          return {
            id: log.id,
            occurredAt: log.occurred_at,
            resourceId: log.resource_id ?? null,
            method: log.method,
            path: log.path,
            statusCode: log.status_code,
            action: log.action,
            clientIp: log.client_ip,
            ruleId: log.rule_id,
            profileId: log.profile_id,
            userAgent: log.user_agent,
            country: log.country,
            latencyMs: log.latency_ms,
            requestId: log.request_id,
            metadata: log.metadata,
            host: log.host,
            scheme: log.scheme,
            protocol: log.protocol,
            authority: log.authority,
            query: log.query,
            sourcePort: log.source_port,
            destinationIp: log.destination_ip,
            destinationPort: log.destination_port,
            sourcePrincipal: log.source_principal,
            sourceService: log.source_service,
            sourceLabels: log.source_labels,
            destinationService: log.destination_service,
            destinationLabels: log.destination_labels,
            requestHttpId: log.request_http_id,
            fragment: log.fragment,
            requestHeaders: log.request_headers,
            requestBodySize: log.request_body_size,
            requestBody: log.request_body,
            contextExtensions: log.context_extensions,
            metadataContext: log.metadata_context,
            routeMetadataContext: log.route_metadata_context,
          };
        }),
      );
    } catch (requestError) {
      setError("Не удалось загрузить логи запросов.");
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
          ? log.resourceId === null
            ? false
            : log.resourceId.toString().includes(filters.resourceId)
          : true;
        const methodMatch = filters.method
          ? log.method.toLowerCase().includes(filters.method.toLowerCase())
          : true;
        const statusMatch = filters.statusCode
          ? log.statusCode.toString().includes(filters.statusCode)
          : true;
        const actionMatch = filters.action
          ? log.action.toLowerCase() === filters.action.toLowerCase()
          : true;
        const ipMatch = filters.clientIp
          ? log.clientIp.toLowerCase().includes(filters.clientIp.toLowerCase())
          : true;
        const pathMatch = filters.path
          ? log.path.toLowerCase().includes(filters.path.toLowerCase())
          : true;
        return (
          resourceMatch &&
          methodMatch &&
          statusMatch &&
          actionMatch &&
          ipMatch &&
          pathMatch
        );
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
    <RequestLogsPageView
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

export default RequestLogsPage;
