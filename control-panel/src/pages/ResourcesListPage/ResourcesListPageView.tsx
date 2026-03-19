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
import Icon from "../../components/Icon";
import TablePagination from "../../components/TablePagination";
import { formatDateTime } from "../../utils/formatDateTime";

type ResourceRow = {
  id: number;
  createdAt: string;
  name: string;
  urlPattern: string;
  securityProfileId: string;
  trafficProfileId: string;
  securityProfile: string;
  trafficProfile: string;
};

type ResourcesListPageViewProps = {
  resources: ResourceRow[];
  securityProfileOptions: SelectOption[];
  trafficProfileOptions: SelectOption[];
  nameFilter: string;
  urlFilter: string;
  securityFilter: string;
  trafficFilter: string;
  onNameFilterChange: (value: string) => void;
  onUrlFilterChange: (value: string) => void;
  onSecurityFilterChange: (value: string) => void;
  onTrafficFilterChange: (value: string) => void;
  isLoading: boolean;
  error: string | null;
  showCreateModal: boolean;
  onOpenCreateModal: () => void;
  onCloseCreateModal: () => void;
  createName: string;
  createUrlPattern: string;
  createSecurityProfileId: string;
  createTrafficProfileId: string;
  onCreateNameChange: (value: string) => void;
  onCreateUrlPatternChange: (value: string) => void;
  onCreateSecurityProfileIdChange: (value: string) => void;
  onCreateTrafficProfileIdChange: (value: string) => void;
  onCreateSubmit: () => void;
  createAttempted: boolean;
  currentPage: number;
  totalPages: number;
  onPageChange: (page: number) => void;
  pageSize: number;
  onPageSizeChange: (size: number) => void;
};

function ResourcesListPageView({
  resources,
  securityProfileOptions,
  trafficProfileOptions,
  nameFilter,
  urlFilter,
  securityFilter,
  trafficFilter,
  onNameFilterChange,
  onUrlFilterChange,
  onSecurityFilterChange,
  onTrafficFilterChange,
  isLoading,
  error,
  showCreateModal,
  onOpenCreateModal,
  onCloseCreateModal,
  createName,
  createUrlPattern,
  createSecurityProfileId,
  createTrafficProfileId,
  onCreateNameChange,
  onCreateUrlPatternChange,
  onCreateSecurityProfileIdChange,
  onCreateTrafficProfileIdChange,
  onCreateSubmit,
  createAttempted,
  currentPage,
  totalPages,
  onPageChange,
  pageSize,
  onPageSizeChange,
}: ResourcesListPageViewProps) {
  return (
    <PageContainer title="Ресурсы">
      <div className="d-flex flex-wrap justify-content-between align-items-center gap-3 mb-4">
        <Button onClick={onOpenCreateModal}>
          <Icon name="plus-lg" className="me-1" />
          Добавить ресурс
        </Button>
      </div>

      <Form className="mb-4">
        <Row className="g-3">
          <Col md={4}>
            <Form.Label>Название</Form.Label>
            <Form.Control
              placeholder="Например, Public API"
              value={nameFilter}
              onChange={(event) => onNameFilterChange(event.target.value)}
            />
          </Col>
          <Col md={4}>
            <Form.Label>URL-паттерн</Form.Label>
            <Form.Control
              placeholder="https://api.example.com/*"
              value={urlFilter}
              onChange={(event) => onUrlFilterChange(event.target.value)}
            />
          </Col>
          <Col md={4}>
            <Form.Label>Профиль безопасности</Form.Label>
            <Form.Select
              value={securityFilter}
              onChange={(event) => onSecurityFilterChange(event.target.value)}
            >
              <option value="">Все</option>
              {securityProfileOptions.map((option) => (
                <option key={option.value} value={option.value}>
                  {option.label}
                </option>
              ))}
            </Form.Select>
          </Col>
          <Col md={4}>
            <Form.Label>Профиль ограничителя трафика</Form.Label>
            <Form.Select
              value={trafficFilter}
              onChange={(event) => onTrafficFilterChange(event.target.value)}
            >
              <option value="">Все</option>
              {trafficProfileOptions.map((option) => (
                <option key={option.value} value={option.value}>
                  {option.label}
                </option>
              ))}
            </Form.Select>
          </Col>
        </Row>
      </Form>

      <ErrorAlert error={error} />

      <Table striped bordered hover responsive>
        <thead>
          <tr>
            <th>Название</th>
            <th>Создан</th>
            <th>URL-паттерн</th>
            <th>Профиль безопасности</th>
            <th>Профиль трафика</th>
          </tr>
        </thead>
        <tbody>
          {resources.map((resource) => (
            <tr key={resource.id}>
              <td>
                <Link className="table-link" to={`/resources/${resource.id}`}>
                  {resource.name}
                </Link>
              </td>
              <td>
                <Link className="table-link" to={`/resources/${resource.id}`}>
                  {formatDateTime(resource.createdAt)}
                </Link>
              </td>
              <td>
                <Link className="table-link" to={`/resources/${resource.id}`}>
                  {resource.urlPattern}
                </Link>
              </td>
              <td>
                <Link className="table-link" to={`/resources/${resource.id}`}>
                  <Badge
                    bg={
                      resource.securityProfile === "Не задан"
                        ? "secondary"
                        : "success"
                    }
                  >
                    {resource.securityProfile}
                  </Badge>
                </Link>
              </td>
              <td>
                <Link className="table-link" to={`/resources/${resource.id}`}>
                  <Badge
                    bg={
                      resource.trafficProfile === "Не задан"
                        ? "secondary"
                        : "success"
                    }
                  >
                    {resource.trafficProfile}
                  </Badge>
                </Link>
              </td>
            </tr>
          ))}
          {isLoading && resources.length === 0 ? (
            <tr>
              <td colSpan={5} className="text-center text-muted">
                Загрузка...
              </td>
            </tr>
          ) : resources.length === 0 ? (
            <tr>
              <td colSpan={5} className="text-center text-muted">
                Нет ресурсов по выбранным фильтрам
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

      <Modal show={showCreateModal} onHide={onCloseCreateModal} centered>
        <Modal.Header closeButton>
          <Modal.Title>Новый ресурс</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          <Form className="d-grid gap-3">
            <Form.Group className="mb-3">
              <Form.Label>Название</Form.Label>
              <Form.Control
                placeholder="Resource name"
                value={createName}
                onChange={(event) => onCreateNameChange(event.target.value)}
                isInvalid={createAttempted && !createName.trim()}
              />
              <Form.Control.Feedback type="invalid">
                Укажите название ресурса.
              </Form.Control.Feedback>
            </Form.Group>
            <Form.Group className="mb-3">
              <Form.Label>URL-паттерн</Form.Label>
              <Form.Control
                placeholder="https://app.example.com/*"
                value={createUrlPattern}
                onChange={(event) =>
                  onCreateUrlPatternChange(event.target.value)
                }
                isInvalid={createAttempted && !createUrlPattern.trim()}
              />
              <Form.Control.Feedback type="invalid">
                Укажите URL-паттерн ресурса.
              </Form.Control.Feedback>
            </Form.Group>
            <Form.Group className="mb-3">
              <Form.Label>Профиль безопасности</Form.Label>
              <Form.Select
                value={createSecurityProfileId}
                onChange={(event) =>
                  onCreateSecurityProfileIdChange(event.target.value)
                }
              >
                <option value="">Не выбран</option>
                {securityProfileOptions.map((option) => (
                  <option key={option.value} value={option.value}>
                    {option.label}
                  </option>
                ))}
              </Form.Select>
            </Form.Group>
            <Form.Group className="mb-3">
              <Form.Label>Профиль ограничителя трафика</Form.Label>
              <Form.Select
                value={createTrafficProfileId}
                onChange={(event) =>
                  onCreateTrafficProfileIdChange(event.target.value)
                }
              >
                <option value="">Не выбран</option>
                {trafficProfileOptions.map((option) => (
                  <option key={option.value} value={option.value}>
                    {option.label}
                  </option>
                ))}
              </Form.Select>
            </Form.Group>
          </Form>
        </Modal.Body>
        <Modal.Footer>
          <Button variant="secondary" onClick={onCloseCreateModal}>
            Отмена
          </Button>
          <Button onClick={onCreateSubmit}>Создать</Button>
        </Modal.Footer>
      </Modal>
    </PageContainer>
  );
}

type SelectOption = {
  value: string;
  label: string;
};

export type { ResourceRow, SelectOption };
export default ResourcesListPageView;
