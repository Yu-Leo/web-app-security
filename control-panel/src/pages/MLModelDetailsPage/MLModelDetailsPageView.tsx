import { useRef } from "react";
import Button from "react-bootstrap/Button";
import Card from "react-bootstrap/Card";
import Form from "react-bootstrap/Form";
import PageContainer from "../../components/PageContainer";
import ErrorAlert from "../../components/ErrorAlert";
import Icon from "../../components/Icon";
import SuccessToast from "../../components/SuccessToast";

type MLModelDetailsPageViewProps = {
  title: string;
  name: string;
  modelSize: string;
  fileName: string;
  isLoading: boolean;
  error: string | null;
  onNameChange: (value: string) => void;
  onSave: () => void;
  onDelete: () => void;
  isDragOver: boolean;
  onDragOver: () => void;
  onDragLeave: () => void;
  onDrop: (file: File) => Promise<void>;
  onBrowseFile: (file: File) => Promise<void>;
  showSuccessToast: boolean;
  onCloseSuccessToast: () => void;
};

function MLModelDetailsPageView({
  title,
  name,
  modelSize,
  fileName,
  isLoading,
  error,
  onNameChange,
  onSave,
  onDelete,
  isDragOver,
  onDragOver,
  onDragLeave,
  onDrop,
  onBrowseFile,
  showSuccessToast,
  onCloseSuccessToast,
}: MLModelDetailsPageViewProps) {
  const fileInputRef = useRef<HTMLInputElement>(null);

  return (
    <PageContainer title={title}>
      <SuccessToast show={showSuccessToast} onClose={onCloseSuccessToast} />
      <div className="d-flex flex-wrap justify-content-between align-items-center gap-3 mb-4">
        <div className="d-flex flex-wrap gap-2">
          <Button
            variant="outline-danger"
            onClick={onDelete}
            disabled={isLoading}
          >
            <Icon name="trash" className="me-1" />
            Удалить
          </Button>
          <Button onClick={onSave} disabled={isLoading}>
            Сохранить
          </Button>
        </div>
      </div>

      <Card>
        <Card.Body>
          <ErrorAlert error={error} />
          <Form className="d-grid gap-3">
            <div>
              <Form.Label>Название</Form.Label>
              <Form.Control
                value={name}
                onChange={(event) => onNameChange(event.target.value)}
                disabled={isLoading}
              />
            </div>

            <div>
              <Form.Label>Текущий файл</Form.Label>
              <Card className="bg-light border">
                <Card.Body className="py-2">
                  <div className="small">
                    <div>
                      <strong>Имя:</strong> {fileName || "—"}
                    </div>
                    <div>
                      <strong>Размер:</strong> {modelSize}
                    </div>
                  </div>
                </Card.Body>
              </Card>
            </div>

            <div>
              <Form.Label>Заменить ONNX-файл</Form.Label>
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
                  Перетащите новый `.onnx` файл сюда или выберите вручную
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
                  {fileName || "Файл не выбран"}
                </div>
              </div>
            </div>
          </Form>
        </Card.Body>
      </Card>
    </PageContainer>
  );
}

export default MLModelDetailsPageView;
