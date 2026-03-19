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

type SecurityProfileRow = {
  id: number;
  createdAt: string;
  name: string;
  baseAction: "allow" | "block";
  logEnabled: boolean;
  isEnabled: boolean;
};

type SecurityProfilesListPageViewProps = {
  profiles: SecurityProfileRow[];
  nameFilter: string;
  baseActionFilter: string;
  statusFilter: string;
  onNameFilterChange: (value: string) => void;
  onBaseActionFilterChange: (value: string) => void;
  onStatusFilterChange: (value: string) => void;
  isLoading: boolean;
  error: string | null;
  showCreateModal: boolean;
  onOpenCreateModal: () => void;
  onCloseCreateModal: () => void;
  createName: string;
  createDescription: string;
  createBaseAction: "allow" | "block";
  createLogEnabled: boolean;
  createIsEnabled: boolean;
  onCreateNameChange: (value: string) => void;
  onCreateDescriptionChange: (value: string) => void;
  onCreateBaseActionChange: (value: "allow" | "block") => void;
  onCreateLogEnabledChange: (value: boolean) => void;
  onCreateIsEnabledChange: (value: boolean) => void;
  onCreateSubmit: () => void;
  createAttempted: boolean;
  currentPage: number;
  totalPages: number;
  onPageChange: (page: number) => void;
  pageSize: number;
  onPageSizeChange: (size: number) => void;
};

function SecurityProfilesListPageView({
  profiles,
  nameFilter,
  baseActionFilter,
  statusFilter,
  onNameFilterChange,
  onBaseActionFilterChange,
  onStatusFilterChange,
  isLoading,
  error,
  showCreateModal,
  onOpenCreateModal,
  onCloseCreateModal,
  createName,
  createDescription,
  createBaseAction,
  createLogEnabled,
  createIsEnabled,
  onCreateNameChange,
  onCreateDescriptionChange,
  onCreateBaseActionChange,
  onCreateLogEnabledChange,
  onCreateIsEnabledChange,
  onCreateSubmit,
  createAttempted,
  currentPage,
  totalPages,
  onPageChange,
  pageSize,
  onPageSizeChange,
}: SecurityProfilesListPageViewProps) {
  const baseActionVariant = (action: string) => {
    const normalized = action.toLowerCase();
    if (normalized === "block" || normalized === "deny") return "danger";
    return "success";
  };

  return (
    <PageContainer title="Профили безопасности">
      <div className="d-flex flex-wrap justify-content-between align-items-center gap-3 mb-4">
        <Button onClick={onOpenCreateModal}>
          <Icon name="plus-lg" className="me-1" />
          Добавить профиль
        </Button>
      </div>

      <Form className="mb-4">
        <Row className="g-3">
          <Col md={4}>
            <Form.Label>Название</Form.Label>
            <Form.Control
              placeholder="Например, Default Security"
              value={nameFilter}
              onChange={(event) => onNameFilterChange(event.target.value)}
            />
          </Col>
          <Col md={4}>
            <Form.Label>Базовое действие</Form.Label>
            <Form.Select
              value={baseActionFilter}
              onChange={(event) => onBaseActionFilterChange(event.target.value)}
            >
              <option value="">Все</option>
              <option value="allow">allow</option>
              <option value="block">block</option>
            </Form.Select>
          </Col>
          <Col md={4}>
            <Form.Label>Статус</Form.Label>
            <Form.Select
              value={statusFilter}
              onChange={(event) => onStatusFilterChange(event.target.value)}
            >
              <option value="">Все</option>
              <option value="enabled">Включен</option>
              <option value="disabled">Выключен</option>
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
            <th>Базовое действие</th>
            <th>Логи</th>
            <th>Статус</th>
          </tr>
        </thead>
        <tbody>
          {profiles.map((profile) => (
            <tr key={profile.id}>
              <td>
                <Link
                  className="table-link"
                  to={`/security-profiles/${profile.id}`}
                >
                  {profile.name}
                </Link>
              </td>
              <td>
                <Link
                  className="table-link"
                  to={`/security-profiles/${profile.id}`}
                >
                  {formatDateTime(profile.createdAt)}
                </Link>
              </td>
              <td>
                <Link
                  className="table-link"
                  to={`/security-profiles/${profile.id}`}
                >
                  <Badge bg={baseActionVariant(profile.baseAction)}>
                    {profile.baseAction}
                  </Badge>
                </Link>
              </td>
              <td>
                <Link
                  className="table-link"
                  to={`/security-profiles/${profile.id}`}
                >
                  <Badge bg={profile.logEnabled ? "success" : "secondary"}>
                    {profile.logEnabled ? "Включены" : "Выключены"}
                  </Badge>
                </Link>
              </td>
              <td>
                <Link
                  className="table-link"
                  to={`/security-profiles/${profile.id}`}
                >
                  <Badge bg={profile.isEnabled ? "success" : "secondary"}>
                    {profile.isEnabled ? "Включен" : "Выключен"}
                  </Badge>
                </Link>
              </td>
            </tr>
          ))}
          {isLoading && profiles.length === 0 ? (
            <tr>
              <td colSpan={5} className="text-center text-muted">
                Загрузка...
              </td>
            </tr>
          ) : profiles.length === 0 ? (
            <tr>
              <td colSpan={5} className="text-center text-muted">
                Нет профилей по выбранным фильтрам
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
          <Modal.Title>Новый профиль безопасности</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          <Form className="d-grid gap-3">
            <div>
              <Form.Label>Название</Form.Label>
              <Form.Control
                placeholder="Security profile name"
                value={createName}
                onChange={(event) => onCreateNameChange(event.target.value)}
                isInvalid={createAttempted && !createName.trim()}
              />
              <Form.Control.Feedback type="invalid">
                Укажите название профиля безопасности.
              </Form.Control.Feedback>
            </div>
            <div>
              <Form.Label>Описание</Form.Label>
              <Form.Control
                as="textarea"
                rows={3}
                placeholder="Описание профиля"
                value={createDescription}
                onChange={(event) =>
                  onCreateDescriptionChange(event.target.value)
                }
              />
            </div>
            <div>
              <Form.Label>Базовое действие</Form.Label>
              <Form.Select
                value={createBaseAction}
                onChange={(event) =>
                  onCreateBaseActionChange(
                    event.target.value as "allow" | "block",
                  )
                }
              >
                <option value="allow">allow</option>
                <option value="block">block</option>
              </Form.Select>
            </div>
            <div>
              <Form.Check
                type="switch"
                label="Логирование включено"
                checked={createLogEnabled}
                onChange={(event) =>
                  onCreateLogEnabledChange(event.target.checked)
                }
              />
            </div>
            <div>
              <Form.Check
                type="switch"
                label="Профиль активен"
                checked={createIsEnabled}
                onChange={(event) =>
                  onCreateIsEnabledChange(event.target.checked)
                }
              />
            </div>
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

export type { SecurityProfileRow };
export default SecurityProfilesListPageView;
