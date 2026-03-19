import { Link } from "react-router-dom";
import Button from "react-bootstrap/Button";
import Badge from "react-bootstrap/Badge";
import Col from "react-bootstrap/Col";
import Form from "react-bootstrap/Form";
import Modal from "react-bootstrap/Modal";
import Row from "react-bootstrap/Row";
import Table from "react-bootstrap/Table";
import PageContainer from "../../components/PageContainer";
import ErrorAlert from "../../components/ErrorAlert";
import TablePagination from "../../components/TablePagination";
import { formatDateTime } from "../../utils/formatDateTime";

type RequestLogRow = {
  id: number;
  occurredAt: string;
  resourceId: number | null;
  method: string;
  path: string;
  statusCode: number;
  action: string;
  clientIp: string;
  ruleId?: number;
  profileId?: number;
  userAgent?: string;
  country?: string;
  latencyMs?: number;
  requestId?: string;
  metadata?: Record<string, unknown>;
  host?: string;
  scheme?: string;
  protocol?: string;
  authority?: string;
  query?: string;
  sourcePort?: number;
  destinationIp?: string;
  destinationPort?: number;
  sourcePrincipal?: string;
  sourceService?: string;
  sourceLabels?: Record<string, unknown>;
  destinationService?: string;
  destinationLabels?: Record<string, unknown>;
  requestHttpId?: string;
  fragment?: string;
  requestHeaders?: Record<string, unknown>;
  requestBodySize?: number;
  requestBody?: string;
  contextExtensions?: Record<string, unknown>;
  metadataContext?: Record<string, unknown>;
  routeMetadataContext?: Record<string, unknown>;
};

type RequestLogsPageViewProps = {
  logs: RequestLogRow[];
  isLoading: boolean;
  error: string | null;
  filters: {
    resourceId: string;
    method: string;
    statusCode: string;
    action: string;
    clientIp: string;
    path: string;
  };
  onFilterChange: (
    field: keyof RequestLogsPageViewProps["filters"],
    value: string,
  ) => void;
  showDetails: boolean;
  detailsLog: RequestLogRow | null;
  onCloseDetails: () => void;
  currentPage: number;
  totalPages: number;
  onPageChange: (page: number) => void;
  pageSize: number;
  onPageSizeChange: (size: number) => void;
};

function RequestLogsPageView({
  logs,
  isLoading,
  error,
  filters,
  onFilterChange,
  showDetails,
  detailsLog,
  onCloseDetails,
  currentPage,
  totalPages,
  onPageChange,
  pageSize,
  onPageSizeChange,
}: RequestLogsPageViewProps) {
  const formatResourceId = (resourceId: number | null) =>
    resourceId === null ? "-" : resourceId;
  const actionVariant = (action: string) => {
    const normalized = action.toLowerCase();
    if (normalized === "block" || normalized === "deny") return "danger";
    if (normalized === "challenge" || normalized === "monitor")
      return "warning";
    return "success";
  };

  return (
    <PageContainer title="Логи запросов">
      <Form className="mb-4">
        <Row className="g-3">
          <Col md={3}>
            <Form.Label>Resource ID</Form.Label>
            <Form.Control
              value={filters.resourceId}
              onChange={(event) =>
                onFilterChange("resourceId", event.target.value)
              }
            />
          </Col>
          <Col md={3}>
            <Form.Label>Метод</Form.Label>
            <Form.Control
              value={filters.method}
              onChange={(event) => onFilterChange("method", event.target.value)}
            />
          </Col>
          <Col md={3}>
            <Form.Label>Статус</Form.Label>
            <Form.Control
              value={filters.statusCode}
              onChange={(event) =>
                onFilterChange("statusCode", event.target.value)
              }
            />
          </Col>
          <Col md={3}>
            <Form.Label>Действие</Form.Label>
            <Form.Select
              value={filters.action}
              onChange={(event) => onFilterChange("action", event.target.value)}
            >
              <option value="">Все</option>
              <option value="allow">allow</option>
              <option value="block">block</option>
            </Form.Select>
          </Col>
          <Col md={4}>
            <Form.Label>IP</Form.Label>
            <Form.Control
              value={filters.clientIp}
              onChange={(event) =>
                onFilterChange("clientIp", event.target.value)
              }
            />
          </Col>
          <Col md={8}>
            <Form.Label>Путь</Form.Label>
            <Form.Control
              value={filters.path}
              onChange={(event) => onFilterChange("path", event.target.value)}
            />
          </Col>
        </Row>
      </Form>

      <ErrorAlert error={error} />

      <Table striped bordered hover responsive>
        <thead>
          <tr>
            <th>ID</th>
            <th>Время</th>
            <th>Resource</th>
            <th>IP</th>
            <th>Метод</th>
            <th>Путь</th>
            <th>Статус</th>
            <th>Действие</th>
            <th>Rule</th>
            <th>Request ID</th>
          </tr>
        </thead>
        <tbody>
          {logs.map((log) => (
            <tr key={log.id}>
              <td>
                <Link className="table-link" to={`?logId=${log.id}`}>
                  {log.id}
                </Link>
              </td>
              <td>
                <Link className="table-link" to={`?logId=${log.id}`}>
                  {formatDateTime(log.occurredAt)}
                </Link>
              </td>
              <td>
                <Link className="table-link" to={`?logId=${log.id}`}>
                  {formatResourceId(log.resourceId)}
                </Link>
              </td>
              <td>
                <Link className="table-link" to={`?logId=${log.id}`}>
                  {log.clientIp}
                </Link>
              </td>
              <td>
                <Link className="table-link" to={`?logId=${log.id}`}>
                  {log.method}
                </Link>
              </td>
              <td>
                <Link className="table-link" to={`?logId=${log.id}`}>
                  {log.path}
                </Link>
              </td>
              <td>
                <Link className="table-link" to={`?logId=${log.id}`}>
                  <Badge bg={log.statusCode >= 400 ? "danger" : "success"}>
                    {log.statusCode}
                  </Badge>
                </Link>
              </td>
              <td>
                <Link className="table-link" to={`?logId=${log.id}`}>
                  <Badge bg={actionVariant(log.action)}>{log.action}</Badge>
                </Link>
              </td>
              <td>
                <Link className="table-link" to={`?logId=${log.id}`}>
                  {log.ruleId ?? "-"}
                </Link>
              </td>
              <td>
                <Link className="table-link" to={`?logId=${log.id}`}>
                  {log.requestId ?? "-"}
                </Link>
              </td>
            </tr>
          ))}
          {isLoading && logs.length === 0 ? (
            <tr>
              <td colSpan={10} className="text-center text-muted">
                Загрузка...
              </td>
            </tr>
          ) : logs.length === 0 ? (
            <tr>
              <td colSpan={10} className="text-center text-muted">
                Нет логов по выбранным фильтрам
              </td>
            </tr>
          ) : null}
        </tbody>
      </Table>
      <TablePagination
        currentPage={currentPage}
        totalPages={totalPages}
        onPageChange={onPageChange}
        pageSize={pageSize}
        onPageSizeChange={onPageSizeChange}
      />

      <Modal show={showDetails} onHide={onCloseDetails} centered size="xl">
        <Modal.Header closeButton>
          <Modal.Title>Детали запроса</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          {detailsLog ? (
            <dl className="row mb-0">
              <dt className="col-sm-4">ID</dt>
              <dd className="col-sm-8">{detailsLog.id}</dd>
              <dt className="col-sm-4">Время</dt>
              <dd className="col-sm-8">
                {formatDateTime(detailsLog.occurredAt)}
              </dd>
              <dt className="col-sm-4">Resource</dt>
              <dd className="col-sm-8">
                {formatResourceId(detailsLog.resourceId)}
              </dd>
              <dt className="col-sm-4">IP</dt>
              <dd className="col-sm-8">{detailsLog.clientIp}</dd>
              <dt className="col-sm-4">Метод</dt>
              <dd className="col-sm-8">{detailsLog.method}</dd>
              <dt className="col-sm-4">Путь</dt>
              <dd className="col-sm-8">{detailsLog.path}</dd>
              <dt className="col-sm-4">Статус</dt>
              <dd className="col-sm-8">{detailsLog.statusCode}</dd>
              <dt className="col-sm-4">Действие</dt>
              <dd className="col-sm-8">{detailsLog.action}</dd>
              <dt className="col-sm-4">Rule ID</dt>
              <dd className="col-sm-8">{detailsLog.ruleId ?? "-"}</dd>
              <dt className="col-sm-4">Profile ID</dt>
              <dd className="col-sm-8">{detailsLog.profileId ?? "-"}</dd>
              <dt className="col-sm-4">Request ID</dt>
              <dd className="col-sm-8">{detailsLog.requestId ?? "-"}</dd>
              <dt className="col-sm-4">User-Agent</dt>
              <dd className="col-sm-8">{detailsLog.userAgent ?? "-"}</dd>
              <dt className="col-sm-4">Country</dt>
              <dd className="col-sm-8">{detailsLog.country ?? "-"}</dd>
              <dt className="col-sm-4">Latency (ms)</dt>
              <dd className="col-sm-8">{detailsLog.latencyMs ?? "-"}</dd>
              <dt className="col-sm-4">Host</dt>
              <dd className="col-sm-8">{detailsLog.host ?? "-"}</dd>
              <dt className="col-sm-4">Scheme</dt>
              <dd className="col-sm-8">{detailsLog.scheme ?? "-"}</dd>
              <dt className="col-sm-4">Protocol</dt>
              <dd className="col-sm-8">{detailsLog.protocol ?? "-"}</dd>
              <dt className="col-sm-4">Authority</dt>
              <dd className="col-sm-8">{detailsLog.authority ?? "-"}</dd>
              <dt className="col-sm-4">Query</dt>
              <dd className="col-sm-8">{detailsLog.query ?? "-"}</dd>
              <dt className="col-sm-4">Source Port</dt>
              <dd className="col-sm-8">{detailsLog.sourcePort ?? "-"}</dd>
              <dt className="col-sm-4">Destination IP</dt>
              <dd className="col-sm-8">{detailsLog.destinationIp ?? "-"}</dd>
              <dt className="col-sm-4">Destination Port</dt>
              <dd className="col-sm-8">{detailsLog.destinationPort ?? "-"}</dd>
              <dt className="col-sm-4">Source Principal</dt>
              <dd className="col-sm-8">{detailsLog.sourcePrincipal ?? "-"}</dd>
              <dt className="col-sm-4">Source Service</dt>
              <dd className="col-sm-8">{detailsLog.sourceService ?? "-"}</dd>
              <dt className="col-sm-4">Source Labels</dt>
              <dd className="col-sm-8">
                <pre className="mb-0">
                  {detailsLog.sourceLabels
                    ? JSON.stringify(detailsLog.sourceLabels, null, 2)
                    : "-"}
                </pre>
              </dd>
              <dt className="col-sm-4">Destination Service</dt>
              <dd className="col-sm-8">
                {detailsLog.destinationService ?? "-"}
              </dd>
              <dt className="col-sm-4">Destination Labels</dt>
              <dd className="col-sm-8">
                <pre className="mb-0">
                  {detailsLog.destinationLabels
                    ? JSON.stringify(detailsLog.destinationLabels, null, 2)
                    : "-"}
                </pre>
              </dd>
              <dt className="col-sm-4">HTTP Request ID</dt>
              <dd className="col-sm-8">{detailsLog.requestHttpId ?? "-"}</dd>
              <dt className="col-sm-4">Fragment</dt>
              <dd className="col-sm-8">{detailsLog.fragment ?? "-"}</dd>
              <dt className="col-sm-4">Request Body Size</dt>
              <dd className="col-sm-8">{detailsLog.requestBodySize ?? "-"}</dd>
              <dt className="col-sm-4">Request Body</dt>
              <dd className="col-sm-8">{detailsLog.requestBody ?? "-"}</dd>
              <dt className="col-sm-4">Request Headers</dt>
              <dd className="col-sm-8">
                <pre className="mb-0">
                  {detailsLog.requestHeaders
                    ? JSON.stringify(detailsLog.requestHeaders, null, 2)
                    : "-"}
                </pre>
              </dd>
              <dt className="col-sm-4">Context Extensions</dt>
              <dd className="col-sm-8">
                <pre className="mb-0">
                  {detailsLog.contextExtensions
                    ? JSON.stringify(detailsLog.contextExtensions, null, 2)
                    : "-"}
                </pre>
              </dd>
              <dt className="col-sm-4">Metadata Context</dt>
              <dd className="col-sm-8">
                <pre className="mb-0">
                  {detailsLog.metadataContext
                    ? JSON.stringify(detailsLog.metadataContext, null, 2)
                    : "-"}
                </pre>
              </dd>
              <dt className="col-sm-4">Route Metadata Context</dt>
              <dd className="col-sm-8">
                <pre className="mb-0">
                  {detailsLog.routeMetadataContext
                    ? JSON.stringify(detailsLog.routeMetadataContext, null, 2)
                    : "-"}
                </pre>
              </dd>
              <dt className="col-sm-4">Metadata</dt>
              <dd className="col-sm-8">
                <pre className="mb-0">
                  {detailsLog.metadata
                    ? JSON.stringify(detailsLog.metadata, null, 2)
                    : "-"}
                </pre>
              </dd>
            </dl>
          ) : (
            <div className="text-muted">Нет данных</div>
          )}
        </Modal.Body>
        <Modal.Footer>
          <Button variant="secondary" onClick={onCloseDetails}>
            Закрыть
          </Button>
        </Modal.Footer>
      </Modal>
    </PageContainer>
  );
}

export type { RequestLogRow };
export default RequestLogsPageView;
