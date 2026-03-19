import { Link } from "react-router-dom";
import { useRef } from "react";
import Button from "react-bootstrap/Button";
import Col from "react-bootstrap/Col";
import Form from "react-bootstrap/Form";
import Modal from "react-bootstrap/Modal";
import Row from "react-bootstrap/Row";
import Table from "react-bootstrap/Table";
import PageContainer from "../../components/PageContainer";
import ErrorAlert from "../../components/ErrorAlert";
import Icon from "../../components/Icon";
import TablePagination from "../../components/TablePagination";

type MLModelRow = {
  id: number;
  name: string;
  sizeBytes: number;
};

type MLModelsListPageViewProps = {
  models: MLModelRow[];
  formatBytes: (value: number) => string;
  isLoading: boolean;
  error: string | null;
  nameFilter: string;
  onNameFilterChange: (value: string) => void;
  showCreateModal: boolean;
  onOpenCreateModal: () => void;
  onCloseCreateModal: () => void;
  createName: string;
  onCreateNameChange: (value: string) => void;
  createFileName: string;
  createModelData: string;
  onCreateSubmit: () => void;
  createAttempted: boolean;
  isDragOver: boolean;
  onDragOver: () => void;
  onDragLeave: () => void;
  onDrop: (file: File) => Promise<void>;
  onBrowseFile: (file: File) => Promise<void>;
  currentPage: number;
  totalPages: number;
  onPageChange: (page: number) => void;
  pageSize: number;
  onPageSizeChange: (size: number) => void;
};

function MLModelsListPageView({
  models,
  formatBytes,
  isLoading,
  error,
  nameFilter,
  onNameFilterChange,
  showCreateModal,
  onOpenCreateModal,
  onCloseCreateModal,
  createName,
  onCreateNameChange,
  createFileName,
  createModelData,
  onCreateSubmit,
  createAttempted,
  isDragOver,
  onDragOver,
  onDragLeave,
  onDrop,
  onBrowseFile,
  currentPage,
  totalPages,
  onPageChange,
  pageSize,
  onPageSizeChange,
}: MLModelsListPageViewProps) {
  const fileInputRef = useRef<HTMLInputElement>(null);

  return (
    <PageContainer title="ML-модели">
      <div className="d-flex flex-wrap justify-content-between align-items-center gap-3 mb-4">
        <Button onClick={onOpenCreateModal}>
          <Icon name="plus-lg" className="me-1" />
          Загрузить модель
        </Button>
      </div>

      <Form className="mb-4">
        <Row className="g-3">
          <Col md={6}>
            <Form.Label>Поиск по названию</Form.Label>
            <Form.Control
              placeholder="Например, fraud-v1"
              value={nameFilter}
              onChange={(event) => onNameFilterChange(event.target.value)}
            />
          </Col>
        </Row>
      </Form>

      <ErrorAlert error={error} />

      <Table striped bordered hover responsive>
        <thead>
          <tr>
            <th>ID</th>
            <th>Название</th>
            <th>Размер</th>
          </tr>
        </thead>
        <tbody>
          {models.map((model) => (
            <tr key={model.id}>
              <td>
                <Link className="table-link" to={`/ml-models/${model.id}`}>
                  {model.id}
                </Link>
              </td>
              <td>
                <Link className="table-link" to={`/ml-models/${model.id}`}>
                  {model.name}
                </Link>
              </td>
              <td>
                <Link className="table-link" to={`/ml-models/${model.id}`}>
                  {formatBytes(model.sizeBytes)}
                </Link>
              </td>
            </tr>
          ))}
          {isLoading && models.length === 0 ? (
            <tr>
              <td colSpan={3} className="text-center text-muted">
                Загрузка...
              </td>
            </tr>
          ) : models.length === 0 ? (
            <tr>
              <td colSpan={3} className="text-center text-muted">
                Нет моделей по выбранным фильтрам
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
          <Modal.Title>Новая ML-модель</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          <Form className="d-grid gap-3">
            <Form.Group>
              <Form.Label>Название</Form.Label>
              <Form.Control
                placeholder="fraud-model-v1"
                value={createName}
                onChange={(event) => onCreateNameChange(event.target.value)}
                isInvalid={createAttempted && !createName.trim()}
              />
              <Form.Control.Feedback type="invalid">
                Укажите название модели.
              </Form.Control.Feedback>
            </Form.Group>

            <Form.Group>
              <Form.Label>Файл ONNX</Form.Label>
              <div
                className={`model-dropzone ${isDragOver ? "is-drag-over" : ""}`}
                onDragOver={(event) => {
                  event.preventDefault();
                  onDragOver();
                }}
                onDragLeave={(event) => {
                  event.preventDefault();
                  onDragLeave();
                }}
                onDrop={(event) => {
                  event.preventDefault();
                  const file = event.dataTransfer.files?.[0];
                  if (file) {
                    void onDrop(file);
                  }
                }}
              >
                <p className="mb-2">
                  Перетащите `.onnx` файл сюда или выберите вручную
                </p>
                <Button
                  variant="outline-secondary"
                  onClick={() => fileInputRef.current?.click()}
                >
                  Выбрать файл
                </Button>
                <Form.Control
                  ref={fileInputRef}
                  type="file"
                  accept=".onnx"
                  className="d-none"
                  onChange={(event) => {
                    const input = event.target as HTMLInputElement;
                    const file = input.files?.[0];
                    if (file) {
                      void onBrowseFile(file);
                    }
                  }}
                />
                <div className="text-muted mt-2 small">
                  {createFileName || "Файл не выбран"}
                </div>
              </div>
              {createAttempted && !createModelData ? (
                <div className="invalid-feedback d-block mt-1">
                  Загрузите ONNX-файл.
                </div>
              ) : null}
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

export type { MLModelRow };
export default MLModelsListPageView;
