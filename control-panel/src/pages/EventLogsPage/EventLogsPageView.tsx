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

type EventLogRow = {
  id: number;
  occurredAt: string;
  resourceId: number;
  eventType: string;
  severity: string;
  message: string;
  ruleId?: number;
  profileId?: number;
  requestId?: string;
  clientIp?: string;
  method?: string;
  path?: string;
  metadata?: Record<string, unknown>;
};

type EventLogsPageViewProps = {
  logs: EventLogRow[];
  isLoading: boolean;
  error: string | null;
  filters: {
    resourceId: string;
    eventType: string;
    severity: string;
    message: string;
  };
  onFilterChange: (
    field: keyof EventLogsPageViewProps["filters"],
    value: string,
  ) => void;
  showDetails: boolean;
  detailsLog: EventLogRow | null;
  onCloseDetails: () => void;
  currentPage: number;
  totalPages: number;
  onPageChange: (page: number) => void;
  pageSize: number;
  onPageSizeChange: (size: number) => void;
};

function EventLogsPageView({
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
}: EventLogsPageViewProps) {
  const severityVariant = (severity: string) => {
    const normalized = severity.toLowerCase();
    if (normalized === "high" || normalized === "critical") return "danger";
    if (normalized === "warning" || normalized === "medium") return "warning";
    return "success";
  };

  return (
    <PageContainer title="Логи событий">
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
            <Form.Label>Тип события</Form.Label>
            <Form.Control
              value={filters.eventType}
              onChange={(event) =>
                onFilterChange("eventType", event.target.value)
              }
            />
          </Col>
          <Col md={3}>
            <Form.Label>Severity</Form.Label>
            <Form.Control
              value={filters.severity}
              onChange={(event) =>
                onFilterChange("severity", event.target.value)
              }
            />
          </Col>
          <Col md={3}>
            <Form.Label>Сообщение</Form.Label>
            <Form.Control
              value={filters.message}
              onChange={(event) =>
                onFilterChange("message", event.target.value)
              }
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
            <th>Тип</th>
            <th>Severity</th>
            <th>Сообщение</th>
            <th>Request ID</th>
            <th>IP</th>
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
                  {log.resourceId}
                </Link>
              </td>
              <td>
                <Link className="table-link" to={`?logId=${log.id}`}>
                  {log.eventType}
                </Link>
              </td>
              <td>
                <Link className="table-link" to={`?logId=${log.id}`}>
                  <Badge bg={severityVariant(log.severity)}>
                    {log.severity}
                  </Badge>
                </Link>
              </td>
              <td>
                <Link className="table-link" to={`?logId=${log.id}`}>
                  {log.message}
                </Link>
              </td>
              <td>
                <Link className="table-link" to={`?logId=${log.id}`}>
                  {log.requestId ?? "-"}
                </Link>
              </td>
              <td>
                <Link className="table-link" to={`?logId=${log.id}`}>
                  {log.clientIp ?? "-"}
                </Link>
              </td>
            </tr>
          ))}
          {isLoading && logs.length === 0 ? (
            <tr>
              <td colSpan={8} className="text-center text-muted">
                Загрузка...
              </td>
            </tr>
          ) : logs.length === 0 ? (
            <tr>
              <td colSpan={8} className="text-center text-muted">
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
          <Modal.Title>Детали события</Modal.Title>
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
              <dd className="col-sm-8">{detailsLog.resourceId}</dd>
              <dt className="col-sm-4">Тип</dt>
              <dd className="col-sm-8">{detailsLog.eventType}</dd>
              <dt className="col-sm-4">Severity</dt>
              <dd className="col-sm-8">{detailsLog.severity}</dd>
              <dt className="col-sm-4">Сообщение</dt>
              <dd className="col-sm-8">{detailsLog.message}</dd>
              <dt className="col-sm-4">Rule ID</dt>
              <dd className="col-sm-8">{detailsLog.ruleId ?? "-"}</dd>
              <dt className="col-sm-4">Profile ID</dt>
              <dd className="col-sm-8">{detailsLog.profileId ?? "-"}</dd>
              <dt className="col-sm-4">Request ID</dt>
              <dd className="col-sm-8">{detailsLog.requestId ?? "-"}</dd>
              <dt className="col-sm-4">IP</dt>
              <dd className="col-sm-8">{detailsLog.clientIp ?? "-"}</dd>
              <dt className="col-sm-4">Метод</dt>
              <dd className="col-sm-8">{detailsLog.method ?? "-"}</dd>
              <dt className="col-sm-4">Путь</dt>
              <dd className="col-sm-8">{detailsLog.path ?? "-"}</dd>
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

export type { EventLogRow };
export default EventLogsPageView;
