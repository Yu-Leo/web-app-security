import { useEffect, useMemo, useState } from "react";
import { useParams, useSearchParams } from "react-router-dom";
import { apiClient } from "../../core/api/client";
import type {
  SecurityHeaderCondition,
  SecurityRuleConditions,
} from "../../core/api/Api";
import SecurityProfileDetailsPageView, {
  type MLModelOption,
  type RuleFormHeader,
  type RuleFormState,
  type SecurityRuleRow,
} from "./SecurityProfileDetailsPageView";

const createEmptyRuleForm = (): RuleFormState => ({
  name: "",
  description: "",
  priority: 1,
  ruleType: "deterministic",
  action: "block",
  matchAll: true,
  sourceIpCidrText: "",
  uriRegexText: "",
  hostRegexText: "",
  methodRegexText: "",
  headers: [{ key: crypto.randomUUID(), name: "", valueRegexText: "" }],
  dryRun: false,
  isEnabled: true,
  mlModelId: null,
  mlThreshold: null,
});

function SecurityProfileDetailsPage() {
  const { id } = useParams();
  const [searchParams, setSearchParams] = useSearchParams();
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [baseAction, setBaseAction] = useState<"allow" | "block">("allow");
  const [logEnabled, setLogEnabled] = useState(true);
  const [isEnabled, setIsEnabled] = useState(true);
  const [rules, setRules] = useState<SecurityRuleRow[]>([]);
  const [mlModels, setMLModels] = useState<MLModelOption[]>([]);
  const [ruleForm, setRuleForm] = useState<RuleFormState>(createEmptyRuleForm);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [showSuccessToast, setShowSuccessToast] = useState(false);

  const profileId = Number(id);
  const ruleIdParam = searchParams.get("ruleId");
  const isCreateRuleMode = ruleIdParam === "new";
  const selectedRule = useMemo(() => {
    if (!ruleIdParam || ruleIdParam === "new") {
      return null;
    }
    const ruleId = Number(ruleIdParam);
    return rules.find((rule) => rule.id === ruleId) ?? null;
  }, [ruleIdParam, rules]);

  const loadProfile = async () => {
    if (!profileId) {
      setError("Некорректный ID профиля.");
      return;
    }
    setIsLoading(true);
    setError(null);
    try {
      const [profileResponse, rulesResponse, mlModelsResponse] =
        await Promise.all([
          apiClient.api.getSecurityProfile(profileId),
          apiClient.api.listSecurityRules(),
          apiClient.api.listMlModels().catch(() => ({ data: [] })),
        ]);

      const profile = profileResponse.data;
      setName(profile.name);
      setDescription(profile.description ?? "");
      setBaseAction(profile.base_action);
      setLogEnabled(profile.log_enabled);
      setIsEnabled(profile.is_enabled);
      setRules(
        rulesResponse.data
          .filter(
            (rule: (typeof rulesResponse.data)[number]) =>
              rule.profile_id === profileId,
          )
          .map((rule: (typeof rulesResponse.data)[number]) => ({
            id: rule.id,
            createdAt: rule.created_at,
            profileId: rule.profile_id,
            name: rule.name,
            description: rule.description ?? "",
            priority: rule.priority,
            ruleType: rule.rule_type,
            action: rule.action,
            conditions: rule.conditions ?? null,
            conditionsSummary: formatConditionsSummary(rule.conditions ?? null),
            mlModelId: rule.ml_model_id ?? null,
            mlThreshold: rule.ml_threshold ?? null,
            dryRun: rule.dry_run,
            isEnabled: rule.is_enabled,
          }))
          .sort(
            (a: SecurityRuleRow, b: SecurityRuleRow) => a.priority - b.priority,
          ),
      );
      setMLModels(
        mlModelsResponse.data.map(
          (model: (typeof mlModelsResponse.data)[number]) => ({
            id: model.id,
            name: model.name,
          }),
        ),
      );
    } catch {
      setError("Не удалось загрузить профиль безопасности.");
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    void loadProfile();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [id]);

  useEffect(() => {
    if (isCreateRuleMode) {
      setRuleForm(createEmptyRuleForm());
      return;
    }
    if (!selectedRule) {
      return;
    }
    setRuleForm(ruleFormFromRule(selectedRule));
  }, [isCreateRuleMode, selectedRule]);

  useEffect(() => {
    if (ruleIdParam && !isCreateRuleMode && !selectedRule && !isLoading) {
      setSearchParams({});
    }
  }, [ruleIdParam, isCreateRuleMode, selectedRule, isLoading, setSearchParams]);

  const updateRuleForm = (patch: Partial<RuleFormState>) => {
    setRuleForm((current) => ({ ...current, ...patch }));
  };

  const handleSave = async () => {
    if (!profileId) {
      return;
    }
    setError(null);
    setIsLoading(true);
    try {
      await apiClient.api.updateSecurityProfile(profileId, {
        name,
        description: description || undefined,
        base_action: baseAction,
        log_enabled: logEnabled,
        is_enabled: isEnabled,
      });
      setShowSuccessToast(true);
    } catch {
      setError("Не удалось сохранить профиль безопасности.");
    } finally {
      setIsLoading(false);
    }
  };

  const handleDeleteRule = async (ruleId: number) => {
    setError(null);
    setIsLoading(true);
    try {
      await apiClient.api.deleteSecurityRule(ruleId);
      setRules((prevRules) => prevRules.filter((rule) => rule.id !== ruleId));
    } catch {
      setError("Не удалось удалить правило безопасности.");
    } finally {
      setIsLoading(false);
    }
  };

  const handleCloseRuleModal = () => {
    setSearchParams({});
  };

  const handleOpenCreateRuleModal = () => {
    setSearchParams({ ruleId: "new" });
  };

  const handleAddHeaderCondition = () => {
    updateRuleForm({
      headers: [
        ...ruleForm.headers,
        { key: crypto.randomUUID(), name: "", valueRegexText: "" },
      ],
    });
  };

  const handleUpdateHeaderCondition = (
    key: string,
    patch: Partial<RuleFormHeader>,
  ) => {
    updateRuleForm({
      headers: ruleForm.headers.map((header) =>
        header.key === key ? { ...header, ...patch } : header,
      ),
    });
  };

  const handleRemoveHeaderCondition = (key: string) => {
    const nextHeaders = ruleForm.headers.filter((header) => header.key !== key);
    updateRuleForm({
      headers:
        nextHeaders.length > 0
          ? nextHeaders
          : [{ key: crypto.randomUUID(), name: "", valueRegexText: "" }],
    });
  };

  const handleSaveRule = async () => {
    if (!profileId) {
      return;
    }

    const normalizedName = ruleForm.name.trim();
    if (!normalizedName) {
      setError("У правила должно быть название.");
      return;
    }

    if (Number.isNaN(ruleForm.priority)) {
      setError("Приоритет должен быть числом.");
      return;
    }

    if (ruleForm.ruleType === "ml") {
      if (ruleForm.mlModelId === null || ruleForm.mlThreshold === null) {
        setError("Для ML-правила нужно указать модель и порог.");
        return;
      }
      if (ruleForm.mlThreshold < 0 || ruleForm.mlThreshold > 100) {
        setError("ML порог должен быть в диапазоне 0..100.");
        return;
      }
    }

    let conditions: SecurityRuleConditions | undefined;
    try {
      conditions = ruleForm.matchAll
        ? undefined
        : buildSecurityRuleConditions(ruleForm);
    } catch (buildError) {
      setError(
        buildError instanceof Error
          ? buildError.message
          : "Не удалось обработать условия правила.",
      );
      return;
    }

    setError(null);
    setIsLoading(true);

    const payload = {
      profile_id: isCreateRuleMode
        ? profileId
        : (selectedRule?.profileId ?? profileId),
      name: normalizedName,
      description: ruleForm.description.trim() || undefined,
      priority: ruleForm.priority,
      rule_type: ruleForm.ruleType,
      action: ruleForm.action,
      conditions,
      ml_model_id: ruleForm.ruleType === "ml" ? ruleForm.mlModelId : null,
      ml_threshold: ruleForm.ruleType === "ml" ? ruleForm.mlThreshold : null,
      dry_run: ruleForm.dryRun,
      is_enabled: ruleForm.isEnabled,
    } as const;

    try {
      if (isCreateRuleMode) {
        await apiClient.api.createSecurityRule(payload);
      } else if (selectedRule) {
        await apiClient.api.updateSecurityRule(selectedRule.id, payload);
      } else {
        setError("Не найдено правило для обновления.");
        return;
      }

      await loadProfile();
      handleCloseRuleModal();
    } catch {
      setError("Не удалось сохранить правило безопасности.");
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <SecurityProfileDetailsPageView
      title={
        name.trim() ? `Профиль безопасности: ${name}` : "Профиль безопасности"
      }
      name={name}
      description={description}
      baseAction={baseAction}
      logEnabled={logEnabled}
      isEnabled={isEnabled}
      rules={rules}
      mlModels={mlModels}
      ruleForm={ruleForm}
      isLoading={isLoading}
      error={error}
      onNameChange={setName}
      onDescriptionChange={setDescription}
      onBaseActionChange={setBaseAction}
      onLogEnabledChange={setLogEnabled}
      onIsEnabledChange={setIsEnabled}
      onSave={handleSave}
      onDeleteRule={handleDeleteRule}
      onOpenCreateRule={handleOpenCreateRuleModal}
      showRuleModal={Boolean(ruleIdParam)}
      isCreateRuleMode={isCreateRuleMode}
      onRuleFormChange={updateRuleForm}
      onAddHeaderCondition={handleAddHeaderCondition}
      onUpdateHeaderCondition={handleUpdateHeaderCondition}
      onRemoveHeaderCondition={handleRemoveHeaderCondition}
      onSaveRule={handleSaveRule}
      onCloseRuleModal={handleCloseRuleModal}
      showSuccessToast={showSuccessToast}
      onCloseSuccessToast={() => setShowSuccessToast(false)}
    />
  );
}

function ruleFormFromRule(rule: SecurityRuleRow): RuleFormState {
  const conditions = rule.conditions ?? {};
  return {
    name: rule.name,
    description: rule.description ?? "",
    priority: rule.priority,
    ruleType: rule.ruleType,
    action: rule.action,
    matchAll:
      !conditions ||
      ((!conditions.source_ip_cidr || conditions.source_ip_cidr.length === 0) &&
        (!conditions.uri_regex || conditions.uri_regex.length === 0) &&
        (!conditions.host_regex || conditions.host_regex.length === 0) &&
        (!conditions.method_regex || conditions.method_regex.length === 0) &&
        (!conditions.headers || conditions.headers.length === 0)),
    sourceIpCidrText: linesFromList(conditions.source_ip_cidr),
    uriRegexText: linesFromList(conditions.uri_regex),
    hostRegexText: linesFromList(conditions.host_regex),
    methodRegexText: linesFromList(conditions.method_regex),
    headers:
      conditions.headers && conditions.headers.length > 0
        ? conditions.headers.map((header) => ({
            key: crypto.randomUUID(),
            name: header.name,
            valueRegexText: linesFromList(header.value_regex),
          }))
        : [{ key: crypto.randomUUID(), name: "", valueRegexText: "" }],
    dryRun: rule.dryRun,
    isEnabled: rule.isEnabled,
    mlModelId: rule.mlModelId ?? null,
    mlThreshold: rule.mlThreshold ?? null,
  };
}

function buildSecurityRuleConditions(
  ruleForm: RuleFormState,
): SecurityRuleConditions | undefined {
  const sourceIpCidr = listFromLines(ruleForm.sourceIpCidrText);
  const uriRegex = listFromLines(ruleForm.uriRegexText);
  const hostRegex = listFromLines(ruleForm.hostRegexText);
  const methodRegex = listFromLines(ruleForm.methodRegexText);
  const headers = buildHeaderConditions(ruleForm.headers);

  const conditions: SecurityRuleConditions = {};
  if (sourceIpCidr.length > 0) {
    conditions.source_ip_cidr = sourceIpCidr;
  }
  if (uriRegex.length > 0) {
    conditions.uri_regex = uriRegex;
  }
  if (hostRegex.length > 0) {
    conditions.host_regex = hostRegex;
  }
  if (methodRegex.length > 0) {
    conditions.method_regex = methodRegex;
  }
  if (headers.length > 0) {
    conditions.headers = headers;
  }

  return Object.keys(conditions).length > 0 ? conditions : undefined;
}

function buildHeaderConditions(
  headers: RuleFormHeader[],
): SecurityHeaderCondition[] {
  return headers.flatMap((header) => {
    const name = header.name.trim();
    const valueRegex = listFromLines(header.valueRegexText);

    if (!name && valueRegex.length === 0) {
      return [];
    }
    if (!name) {
      throw new Error("Для условия по заголовку нужно указать имя заголовка.");
    }
    if (valueRegex.length === 0) {
      throw new Error(
        `Для заголовка "${name}" нужно указать хотя бы один regexp по значению.`,
      );
    }

    return [{ name, value_regex: valueRegex }];
  });
}

function listFromLines(value: string): string[] {
  return value
    .split("\n")
    .map((item) => item.trim())
    .filter(Boolean);
}

function linesFromList(items?: string[]): string {
  return items?.join("\n") ?? "";
}

function formatConditionsSummary(
  conditions?: SecurityRuleConditions | null,
): string {
  if (!conditions) {
    return "Весь трафик";
  }

  const parts: string[] = [];
  if (conditions.source_ip_cidr?.length) {
    parts.push(`IP: ${conditions.source_ip_cidr.length}`);
  }
  if (conditions.uri_regex?.length) {
    parts.push(`URI: ${conditions.uri_regex.length}`);
  }
  if (conditions.host_regex?.length) {
    parts.push(`Host: ${conditions.host_regex.length}`);
  }
  if (conditions.method_regex?.length) {
    parts.push(`Method: ${conditions.method_regex.length}`);
  }
  if (conditions.headers?.length) {
    parts.push(`Headers: ${conditions.headers.length}`);
  }

  return parts.length > 0 ? parts.join(" · ") : "Весь трафик";
}

export default SecurityProfileDetailsPage;
