import { useEffect, useMemo, useState } from "react";
import { useParams, useSearchParams } from "react-router-dom";
import { apiClient } from "../../core/api/client";
import type {
  SecurityHeaderCondition,
  SecurityRuleConditions,
} from "../../core/api/Api";
import TrafficProfileDetailsPageView, {
  type RuleFormHeader,
  type TrafficRuleFormState,
  type TrafficRuleRow,
} from "./TrafficProfileDetailsPageView";

const createEmptyRuleForm = (): TrafficRuleFormState => ({
  name: "",
  description: "",
  priority: 1,
  requestsLimit: 100,
  periodSeconds: 60,
  matchAll: true,
  sourceIpCidrText: "",
  uriRegexText: "",
  hostRegexText: "",
  methodRegexText: "",
  headers: [{ key: crypto.randomUUID(), name: "", valueRegexText: "" }],
  dryRun: false,
  isEnabled: true,
});

function TrafficProfileDetailsPage() {
  const { id } = useParams();
  const [searchParams, setSearchParams] = useSearchParams();
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [isEnabled, setIsEnabled] = useState(true);
  const [rules, setRules] = useState<TrafficRuleRow[]>([]);
  const [ruleForm, setRuleForm] = useState<TrafficRuleFormState>(
    createEmptyRuleForm,
  );
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
      const [profileResponse, rulesResponse] = await Promise.all([
        apiClient.api.getTrafficProfile(profileId),
        apiClient.api.listTrafficRules(),
      ]);

      const profile = profileResponse.data;
      setName(profile.name);
      setDescription(profile.description ?? "");
      setIsEnabled(profile.is_enabled);
      setRules(
        rulesResponse.data
          .filter((rule) => rule.profile_id === profileId)
          .map((rule) => ({
            id: rule.id,
            createdAt: rule.created_at,
            profileId: rule.profile_id,
            name: rule.name,
            description: rule.description ?? "",
            priority: rule.priority,
            requestsLimit: rule.requests_limit,
            periodSeconds: rule.period_seconds,
            matchAll: rule.match_all,
            conditions: rule.conditions ?? null,
            conditionsSummary: formatConditionsSummary(
              rule.match_all,
              rule.conditions ?? null,
            ),
            dryRun: rule.dry_run,
            isEnabled: rule.is_enabled,
          }))
          .sort((a, b) => a.priority - b.priority),
      );
    } catch {
      setError("Не удалось загрузить профиль трафика.");
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

  const updateRuleForm = (patch: Partial<TrafficRuleFormState>) => {
    setRuleForm((current) => ({ ...current, ...patch }));
  };

  const handleSave = async () => {
    if (!profileId) {
      return;
    }
    setError(null);
    setIsLoading(true);
    try {
      await apiClient.api.updateTrafficProfile(profileId, {
        name,
        description: description || undefined,
        is_enabled: isEnabled,
      });
      setShowSuccessToast(true);
    } catch {
      setError("Не удалось сохранить профиль трафика.");
    } finally {
      setIsLoading(false);
    }
  };

  const handleDeleteRule = async (ruleId: number) => {
    setError(null);
    setIsLoading(true);
    try {
      await apiClient.api.deleteTrafficRule(ruleId);
      setRules((prevRules) => prevRules.filter((rule) => rule.id !== ruleId));
    } catch {
      setError("Не удалось удалить правило трафика.");
    } finally {
      setIsLoading(false);
    }
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

  const handleCloseRuleModal = () => {
    setSearchParams({});
  };

  const handleOpenCreateRuleModal = () => {
    setSearchParams({ ruleId: "new" });
  };

  const handleSaveRule = async () => {
    if (!profileId) {
      return;
    }

    if (!ruleForm.name.trim()) {
      setError("У правила должно быть название.");
      return;
    }
    if (Number.isNaN(ruleForm.priority)) {
      setError("Приоритет должен быть числом.");
      return;
    }
    if (Number.isNaN(ruleForm.requestsLimit) || ruleForm.requestsLimit <= 0) {
      setError("Лимит запросов должен быть положительным числом.");
      return;
    }
    if (Number.isNaN(ruleForm.periodSeconds) || ruleForm.periodSeconds <= 0) {
      setError("Период должен быть положительным числом.");
      return;
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
      profile_id: isCreateRuleMode ? profileId : (selectedRule?.profileId ?? profileId),
      name: ruleForm.name.trim(),
      description: ruleForm.description.trim() || undefined,
      priority: ruleForm.priority,
      dry_run: ruleForm.dryRun,
      match_all: ruleForm.matchAll,
      requests_limit: ruleForm.requestsLimit,
      period_seconds: ruleForm.periodSeconds,
      conditions,
      is_enabled: ruleForm.isEnabled,
    } as const;

    try {
      if (isCreateRuleMode) {
        await apiClient.api.createTrafficRule(payload);
      } else if (selectedRule) {
        await apiClient.api.updateTrafficRule(selectedRule.id, payload);
      } else {
        setError("Не найдено правило для обновления.");
        return;
      }

      await loadProfile();
      handleCloseRuleModal();
    } catch {
      setError("Не удалось сохранить правило трафика.");
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <TrafficProfileDetailsPageView
      title={
        name.trim()
          ? `Профиль ограничителя трафика: ${name}`
          : "Профиль ограничителя трафика"
      }
      name={name}
      description={description}
      isEnabled={isEnabled}
      rules={rules}
      ruleForm={ruleForm}
      isLoading={isLoading}
      error={error}
      onNameChange={setName}
      onDescriptionChange={setDescription}
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

function ruleFormFromRule(rule: TrafficRuleRow): TrafficRuleFormState {
  const conditions = rule.conditions ?? {};
  return {
    name: rule.name,
    description: rule.description ?? "",
    priority: rule.priority,
    requestsLimit: rule.requestsLimit,
    periodSeconds: rule.periodSeconds,
    matchAll: rule.matchAll,
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
  };
}

function buildSecurityRuleConditions(
  ruleForm: TrafficRuleFormState,
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

function buildHeaderConditions(headers: RuleFormHeader[]): SecurityHeaderCondition[] {
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
  matchAll: boolean,
  conditions?: SecurityRuleConditions | null,
): string {
  if (matchAll) {
    return "Весь трафик";
  }
  if (!conditions) {
    return "Без условий";
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

  return parts.length > 0 ? parts.join(" · ") : "Без условий";
}

export default TrafficProfileDetailsPage;
