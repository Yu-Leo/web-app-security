import { Link } from "react-router-dom";
import Button from "react-bootstrap/Button";
import Card from "react-bootstrap/Card";
import Form from "react-bootstrap/Form";
import PageContainer from "../../components/PageContainer";
import ErrorAlert from "../../components/ErrorAlert";
import Icon from "../../components/Icon";
import SuccessToast from "../../components/SuccessToast";

type SelectOption = {
  value: string;
  label: string;
};

type ResourceDetailsPageViewProps = {
  title: string;
  name: string;
  urlPattern: string;
  securityProfileId: string;
  trafficProfileId: string;
  securityProfiles: SelectOption[];
  trafficProfiles: SelectOption[];
  isLoading: boolean;
  error: string | null;
  onNameChange: (value: string) => void;
  onUrlPatternChange: (value: string) => void;
  onSecurityProfileChange: (value: string) => void;
  onTrafficProfileChange: (value: string) => void;
  onSave: () => void;
  onDelete: () => void;
  showSuccessToast: boolean;
  onCloseSuccessToast: () => void;
};

function ResourceDetailsPageView({
  title,
  name,
  urlPattern,
  securityProfileId,
  trafficProfileId,
  securityProfiles,
  trafficProfiles,
  isLoading,
  error,
  onNameChange,
  onUrlPatternChange,
  onSecurityProfileChange,
  onTrafficProfileChange,
  onSave,
  onDelete,
  showSuccessToast,
  onCloseSuccessToast,
}: ResourceDetailsPageViewProps) {
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
              <Form.Label>URL-паттерн</Form.Label>
              <Form.Control
                value={urlPattern}
                onChange={(event) => onUrlPatternChange(event.target.value)}
                disabled={isLoading}
              />
            </div>
            <div>
              <div className="d-flex justify-content-between align-items-center mb-1">
                <Form.Label className="mb-0">Профиль безопасности</Form.Label>
                <div className="d-flex align-items-center gap-3">
                  <Link to="/security-profiles">Создать профиль</Link>
                  {securityProfileId ? (
                    <Link to={`/security-profiles/${securityProfileId}`}>
                      Открыть профиль
                    </Link>
                  ) : null}
                </div>
              </div>
              <Form.Select
                value={securityProfileId}
                onChange={(event) =>
                  onSecurityProfileChange(event.target.value)
                }
                disabled={isLoading}
              >
                <option value="">Не выбран</option>
                {securityProfiles.map((profile) => (
                  <option key={profile.value} value={profile.value}>
                    {profile.label}
                  </option>
                ))}
              </Form.Select>
            </div>
            <div>
              <div className="d-flex justify-content-between align-items-center mb-1">
                <Form.Label className="mb-0">
                  Профиль ограничителя трафика
                </Form.Label>
                <div className="d-flex align-items-center gap-3">
                  <Link to="/traffic-profiles">Создать профиль</Link>
                  {trafficProfileId ? (
                    <Link to={`/traffic-profiles/${trafficProfileId}`}>
                      Открыть профиль
                    </Link>
                  ) : null}
                </div>
              </div>
              <Form.Select
                value={trafficProfileId}
                onChange={(event) => onTrafficProfileChange(event.target.value)}
                disabled={isLoading}
              >
                <option value="">Не выбран</option>
                {trafficProfiles.map((profile) => (
                  <option key={profile.value} value={profile.value}>
                    {profile.label}
                  </option>
                ))}
              </Form.Select>
            </div>
          </Form>
        </Card.Body>
      </Card>
    </PageContainer>
  );
}

export type { SelectOption };
export default ResourceDetailsPageView;
