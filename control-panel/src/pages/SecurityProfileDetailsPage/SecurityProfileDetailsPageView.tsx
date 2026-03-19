import { Link } from "react-router-dom";
import Badge from "react-bootstrap/Badge";
import Button from "react-bootstrap/Button";
import Card from "react-bootstrap/Card";
import Col from "react-bootstrap/Col";
import Form from "react-bootstrap/Form";
import Modal from "react-bootstrap/Modal";
import Row from "react-bootstrap/Row";
import Table from "react-bootstrap/Table";
import ErrorAlert from "../../components/ErrorAlert";
import Icon from "../../components/Icon";
import PageContainer from "../../components/PageContainer";
import SuccessToast from "../../components/SuccessToast";

type SecurityRuleRow = {
  id: number;
  createdAt: string;
  profileId: number;
  name: string;
  description?: string;
  priority: number;
  ruleType: "deterministic" | "ml";
  action: "allow" | "block";
  conditions?: {
    source_ip_cidr?: string[];
    uri_regex?: string[];
    host_regex?: string[];
    method_regex?: string[];
    headers?: { name: string; value_regex: string[] }[];
  } | null;
  conditionsSummary: string;
  mlModelId?: number | null;
  mlThreshold?: number | null;
  dryRun: boolean;
  isEnabled: boolean;
};

type MLModelOption = {
  id: number;
  name: string;
};

type RuleFormHeader = {
  key: string;
  name: string;
  valueRegexText: string;
};

type RuleFormState = {
  name: string;
  description: string;
  priority: number;
  ruleType: "deterministic" | "ml";
  action: "allow" | "block";
  matchAll: boolean;
  sourceIpCidrText: string;
  uriRegexText: string;
  hostRegexText: string;
  methodRegexText: string;
  headers: RuleFormHeader[];
  dryRun: boolean;
  isEnabled: boolean;
  mlModelId: number | null;
  mlThreshold: number | null;
};

type SecurityProfileDetailsPageViewProps = {
  title: string;
  name: string;
  description: string;
  baseAction: "allow" | "block";
  logEnabled: boolean;
  isEnabled: boolean;
  rules: SecurityRuleRow[];
  mlModels: MLModelOption[];
  ruleForm: RuleFormState;
  isLoading: boolean;
  error: string | null;
  onNameChange: (value: string) => void;
  onDescriptionChange: (value: string) => void;
  onBaseActionChange: (value: "allow" | "block") => void;
  onLogEnabledChange: (value: boolean) => void;
  onIsEnabledChange: (value: boolean) => void;
  onSave: () => void;
  onDeleteRule: (id: number) => void;
  onOpenCreateRule: () => void;
  showRuleModal: boolean;
  isCreateRuleMode: boolean;
  onRuleFormChange: (patch: Partial<RuleFormState>) => void;
  onAddHeaderCondition: () => void;
  onUpdateHeaderCondition: (
    key: string,
    patch: Partial<RuleFormHeader>,
  ) => void;
  onRemoveHeaderCondition: (key: string) => void;
  onSaveRule: () => void;
  onCloseRuleModal: () => void;
  showSuccessToast: boolean;
  onCloseSuccessToast: () => void;
};

function SecurityProfileDetailsPageView({
  title,
  name,
  description,
  baseAction,
  logEnabled,
  isEnabled,
  rules,
  mlModels,
  ruleForm,
  isLoading,
  error,
  onNameChange,
  onDescriptionChange,
  onBaseActionChange,
  onLogEnabledChange,
  onIsEnabledChange,
  onSave,
  onDeleteRule,
  onOpenCreateRule,
  showRuleModal,
  isCreateRuleMode,
  onRuleFormChange,
  onAddHeaderCondition,
  onUpdateHeaderCondition,
  onRemoveHeaderCondition,
  onSaveRule,
  onCloseRuleModal,
  showSuccessToast,
  onCloseSuccessToast,
}: SecurityProfileDetailsPageViewProps) {
  const actionVariant = (action: string) =>
    action.toLowerCase() === "block" ? "danger" : "success";

  return (
    <PageContainer title={title}>
      <SuccessToast show={showSuccessToast} onClose={onCloseSuccessToast} />
      <div className="d-flex flex-wrap justify-content-between align-items-center gap-3 mb-4">
        <Button onClick={onSave} disabled={isLoading}>
          Сохранить
        </Button>
      </div>

      <Card className="mb-4">
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
              <Form.Label>Описание</Form.Label>
              <Form.Control
                as="textarea"
                rows={3}
                value={description}
                onChange={(event) => onDescriptionChange(event.target.value)}
                disabled={isLoading}
              />
            </div>
            <div>
              <Form.Label>Базовое действие</Form.Label>
              <Form.Select
                value={baseAction}
                onChange={(event) =>
                  onBaseActionChange(event.target.value as "allow" | "block")
                }
                disabled={isLoading}
              >
                <option value="allow">allow</option>
                <option value="block">block</option>
              </Form.Select>
            </div>
            <div>
              <Form.Check
                type="switch"
                label="Логирование включено"
                checked={logEnabled}
                onChange={(event) => onLogEnabledChange(event.target.checked)}
                disabled={isLoading}
              />
            </div>
            <div>
              <Form.Check
                type="switch"
                label="Профиль активен"
                checked={isEnabled}
                onChange={(event) => onIsEnabledChange(event.target.checked)}
                disabled={isLoading}
              />
            </div>
          </Form>
        </Card.Body>
      </Card>

      <Card>
        <Card.Header className="d-flex justify-content-between align-items-center">
          <div>
            <strong>Правила безопасности</strong>
          </div>
          <div className="d-flex align-items-center gap-3">
            <span className="text-muted">{rules.length} правил</span>
            <Button size="sm" onClick={onOpenCreateRule} disabled={isLoading}>
              Добавить правило
            </Button>
          </div>
        </Card.Header>
        <Card.Body>
          <Table striped bordered hover responsive>
            <thead>
              <tr>
                <th>Название</th>
                <th>Приоритет</th>
                <th>Тип</th>
                <th>Действие</th>
                <th>Условия</th>
                <th>Статус</th>
                <th></th>
              </tr>
            </thead>
            <tbody>
              {rules.map((rule) => (
                <tr key={rule.id}>
                  <td>
                    <Link className="table-link" to={`?ruleId=${rule.id}`}>
                      {rule.name}
                    </Link>
                  </td>
                  <td>
                    <Link className="table-link" to={`?ruleId=${rule.id}`}>
                      {rule.priority}
                    </Link>
                  </td>
                  <td>
                    <Link className="table-link" to={`?ruleId=${rule.id}`}>
                      <Badge bg={rule.ruleType === "ml" ? "primary" : "dark"}>
                        {rule.ruleType}
                      </Badge>
                    </Link>
                  </td>
                  <td>
                    <Link className="table-link" to={`?ruleId=${rule.id}`}>
                      <Badge bg={actionVariant(rule.action)}>
                        {rule.action}
                      </Badge>
                    </Link>
                  </td>
                  <td className="small text-muted">{rule.conditionsSummary}</td>
                  <td>
                    <div className="d-flex flex-column gap-1">
                      <Badge bg={rule.isEnabled ? "success" : "secondary"}>
                        {rule.isEnabled ? "Включено" : "Выключено"}
                      </Badge>
                      {rule.dryRun ? <Badge bg="warning">dry run</Badge> : null}
                    </div>
                  </td>
                  <td className="text-end">
                    <Button
                      size="sm"
                      variant="outline-danger"
                      onClick={() => onDeleteRule(rule.id)}
                      disabled={isLoading}
                      title="Удалить правило"
                      aria-label="Удалить правило"
                    >
                      <Icon name="x-lg" />
                    </Button>
                  </td>
                </tr>
              ))}
              {isLoading && rules.length === 0 ? (
                <tr>
                  <td colSpan={7} className="text-center text-muted">
                    Загрузка...
                  </td>
                </tr>
              ) : rules.length === 0 ? (
                <tr>
                  <td colSpan={7} className="text-center text-muted">
                    Нет правил
                  </td>
                </tr>
              ) : null}
            </tbody>
          </Table>
        </Card.Body>
      </Card>

      <Modal show={showRuleModal} onHide={onCloseRuleModal} centered size="lg">
        <Modal.Header closeButton>
          <Modal.Title>
            {isCreateRuleMode ? "Создание правила" : "Редактирование правила"}
          </Modal.Title>
        </Modal.Header>
        <Modal.Body>
          <Form className="d-grid gap-4">
            <Row className="g-3">
              <Col md={8}>
                <Form.Label>Название</Form.Label>
                <Form.Control
                  value={ruleForm.name}
                  onChange={(event) =>
                    onRuleFormChange({ name: event.target.value })
                  }
                />
              </Col>
              <Col md={4}>
                <Form.Label>Приоритет</Form.Label>
                <Form.Control
                  type="number"
                  value={ruleForm.priority}
                  onChange={(event) =>
                    onRuleFormChange({ priority: Number(event.target.value) })
                  }
                />
              </Col>
            </Row>

            <div>
              <Form.Label>Описание</Form.Label>
              <Form.Control
                as="textarea"
                rows={2}
                value={ruleForm.description}
                onChange={(event) =>
                  onRuleFormChange({ description: event.target.value })
                }
              />
            </div>

            <Row className="g-3">
              <Col md={6}>
                <Form.Label>Тип правила</Form.Label>
                <Form.Select
                  value={ruleForm.ruleType}
                  onChange={(event) =>
                    onRuleFormChange({
                      ruleType: event.target.value as "deterministic" | "ml",
                      mlModelId:
                        event.target.value === "ml" ? ruleForm.mlModelId : null,
                      mlThreshold:
                        event.target.value === "ml"
                          ? ruleForm.mlThreshold
                          : null,
                    })
                  }
                >
                  <option value="deterministic">deterministic</option>
                  <option value="ml">ml</option>
                </Form.Select>
              </Col>
              <Col md={6}>
                <Form.Label>Действие</Form.Label>
                <Form.Select
                  value={ruleForm.action}
                  onChange={(event) =>
                    onRuleFormChange({
                      action: event.target.value as "allow" | "block",
                    })
                  }
                >
                  <option value="allow">allow</option>
                  <option value="block">block</option>
                </Form.Select>
              </Col>
            </Row>

            {ruleForm.ruleType === "ml" ? (
              <Row className="g-3">
                <Col md={8}>
                  <Form.Label>ML модель</Form.Label>
                  <Form.Select
                    value={ruleForm.mlModelId ?? ""}
                    onChange={(event) =>
                      onRuleFormChange({
                        mlModelId: event.target.value
                          ? Number(event.target.value)
                          : null,
                      })
                    }
                  >
                    <option value="">Выберите модель</option>
                    {mlModels.map((model) => (
                      <option key={model.id} value={model.id}>
                        {model.name} (ID: {model.id})
                      </option>
                    ))}
                  </Form.Select>
                </Col>
                <Col md={4}>
                  <Form.Label>Порог ML</Form.Label>
                  <Form.Control
                    type="number"
                    min={0}
                    max={100}
                    value={ruleForm.mlThreshold ?? ""}
                    onChange={(event) =>
                      onRuleFormChange({
                        mlThreshold: event.target.value
                          ? Number(event.target.value)
                          : null,
                      })
                    }
                    placeholder="0..100"
                  />
                </Col>
              </Row>
            ) : null}

            <Form.Check
              type="switch"
              label="Применять ко всему трафику"
              checked={ruleForm.matchAll}
              onChange={(event) =>
                onRuleFormChange({ matchAll: event.target.checked })
              }
            />

            {!ruleForm.matchAll ? (
              <Card bg="light" border="secondary-subtle">
                <Card.Body className="d-grid gap-3">
                  <div className="d-flex justify-content-between align-items-center">
                    <div>
                      <strong>Условия</strong>
                      <div className="small text-muted">
                        Используется тот же набор фильтров, что и у rate
                        limiter.
                      </div>
                    </div>
                  </div>

                  <Row className="g-3">
                    <Col md={6}>
                      <Form.Label>
                        IP источника (CIDR, по одному на строку)
                      </Form.Label>
                      <Form.Control
                        as="textarea"
                        rows={4}
                        value={ruleForm.sourceIpCidrText}
                        onChange={(event) =>
                          onRuleFormChange({
                            sourceIpCidrText: event.target.value,
                          })
                        }
                        placeholder={"10.0.0.0/8\n192.168.1.0/24"}
                      />
                    </Col>
                    <Col md={6}>
                      <Form.Label>
                        Regexp по URI (по одному на строку)
                      </Form.Label>
                      <Form.Control
                        as="textarea"
                        rows={4}
                        value={ruleForm.uriRegexText}
                        onChange={(event) =>
                          onRuleFormChange({ uriRegexText: event.target.value })
                        }
                        placeholder={"^/api/.*\n^/admin"}
                      />
                    </Col>
                    <Col md={6}>
                      <Form.Label>
                        Regexp по Host (по одному на строку)
                      </Form.Label>
                      <Form.Control
                        as="textarea"
                        rows={4}
                        value={ruleForm.hostRegexText}
                        onChange={(event) =>
                          onRuleFormChange({
                            hostRegexText: event.target.value,
                          })
                        }
                        placeholder={"^example\\.com$\n^api\\."}
                      />
                    </Col>
                    <Col md={6}>
                      <Form.Label>
                        Regexp по HTTP методу (по одному на строку)
                      </Form.Label>
                      <Form.Control
                        as="textarea"
                        rows={4}
                        value={ruleForm.methodRegexText}
                        onChange={(event) =>
                          onRuleFormChange({
                            methodRegexText: event.target.value,
                          })
                        }
                        placeholder={"^GET$\n^POST$"}
                      />
                    </Col>
                  </Row>

                  <div className="d-grid gap-3">
                    <div className="d-flex justify-content-between align-items-center">
                      <div>
                        <Form.Label className="mb-0">
                          Условия по заголовкам
                        </Form.Label>
                        <div className="small text-muted">
                          Для каждого заголовка можно задать несколько regexp по
                          значению, по одному на строку.
                        </div>
                      </div>
                      <Button
                        size="sm"
                        variant="outline-secondary"
                        onClick={onAddHeaderCondition}
                      >
                        Добавить заголовок
                      </Button>
                    </div>

                    {ruleForm.headers.map((header) => (
                      <Card key={header.key} border="secondary-subtle">
                        <Card.Body>
                          <Row className="g-3 align-items-start">
                            <Col md={4}>
                              <Form.Label>Имя заголовка</Form.Label>
                              <Form.Control
                                value={header.name}
                                onChange={(event) =>
                                  onUpdateHeaderCondition(header.key, {
                                    name: event.target.value,
                                  })
                                }
                                placeholder="X-Forwarded-For"
                              />
                            </Col>
                            <Col md={7}>
                              <Form.Label>Regexp значения</Form.Label>
                              <Form.Control
                                as="textarea"
                                rows={3}
                                value={header.valueRegexText}
                                onChange={(event) =>
                                  onUpdateHeaderCondition(header.key, {
                                    valueRegexText: event.target.value,
                                  })
                                }
                                placeholder={"^Mozilla/.*\n^curl/.*"}
                              />
                            </Col>
                            <Col md={1} className="d-flex justify-content-end">
                              <Button
                                size="sm"
                                variant="outline-danger"
                                className="mt-4"
                                onClick={() =>
                                  onRemoveHeaderCondition(header.key)
                                }
                                aria-label="Удалить условие по заголовку"
                              >
                                <Icon name="x-lg" />
                              </Button>
                            </Col>
                          </Row>
                        </Card.Body>
                      </Card>
                    ))}
                  </div>
                </Card.Body>
              </Card>
            ) : null}

            <Row className="g-3">
              <Col md={6}>
                <Form.Check
                  type="switch"
                  label="Dry run"
                  checked={ruleForm.dryRun}
                  onChange={(event) =>
                    onRuleFormChange({ dryRun: event.target.checked })
                  }
                />
              </Col>
              <Col md={6}>
                <Form.Check
                  type="switch"
                  label="Правило активно"
                  checked={ruleForm.isEnabled}
                  onChange={(event) =>
                    onRuleFormChange({ isEnabled: event.target.checked })
                  }
                />
              </Col>
            </Row>
          </Form>
        </Modal.Body>
        <Modal.Footer>
          <Button variant="secondary" onClick={onCloseRuleModal}>
            Закрыть
          </Button>
          <Button onClick={onSaveRule} disabled={isLoading}>
            {isCreateRuleMode ? "Создать" : "Сохранить"}
          </Button>
        </Modal.Footer>
      </Modal>
    </PageContainer>
  );
}

export type { MLModelOption, RuleFormHeader, RuleFormState, SecurityRuleRow };
export default SecurityProfileDetailsPageView;
